package validate

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	mgmtnetwork "github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-07-01/network"
	mgmtfeatures "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-07-01/features"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/apparentlymart/go-cidr/cidr"
	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/sirupsen/logrus"

	"github.com/Azure/ARO-RP/pkg/api"
	"github.com/Azure/ARO-RP/pkg/api/validate/dynamic"
	"github.com/Azure/ARO-RP/pkg/env"
	"github.com/Azure/ARO-RP/pkg/util/aad"
	"github.com/Azure/ARO-RP/pkg/util/azureclaim"
	"github.com/Azure/ARO-RP/pkg/util/azureclient/mgmt/features"
	"github.com/Azure/ARO-RP/pkg/util/refreshable"
	"github.com/Azure/ARO-RP/pkg/util/subnet"
)

// OpenShiftClusterDynamicValidator is the dynamic validator interface
type OpenShiftClusterDynamicValidator interface {
	Dynamic(context.Context) error
}

// NewOpenShiftClusterDynamicValidator creates a new OpenShiftClusterDynamicValidator
func NewOpenShiftClusterDynamicValidator(log *logrus.Entry, env env.Core, oc *api.OpenShiftCluster, subscriptionDoc *api.SubscriptionDocument, fpAuthorizer refreshable.Authorizer) OpenShiftClusterDynamicValidator {
	return &openShiftClusterDynamicValidator{
		log: log,
		env: env,

		oc:              oc,
		subscriptionDoc: subscriptionDoc,
		fpAuthorizer:    fpAuthorizer,
		providers:       features.NewProvidersClient(env.Environment(), subscriptionDoc.ID, fpAuthorizer),
	}
}

type openShiftClusterDynamicValidator struct {
	log *logrus.Entry
	env env.Core

	oc              *api.OpenShiftCluster
	subscriptionDoc *api.SubscriptionDocument
	fpAuthorizer    refreshable.Authorizer
	providers       features.ProvidersClient
}

// Dynamic validates an OpenShift cluster
// this function should only be called from the context of the RP, as it utilizes document objects
// TODO - move Dynamic out of this package
func (dv *openShiftClusterDynamicValidator) Dynamic(ctx context.Context) error {
	subnetIDs := []string{dv.oc.Properties.MasterProfile.SubnetID}

	for _, s := range dv.oc.Properties.WorkerProfiles {
		subnetIDs = append(subnetIDs, s.SubnetID)
	}

	vnetID, _, err := subnet.Split(subnetIDs[0])
	if err != nil {
		return err
	}

	// FP validation
	fpDynamic, err := dynamic.NewValidator(dv.log, dv.env.Environment(), dv.subscriptionDoc.ID, dv.fpAuthorizer)
	if err != nil {
		return err
	}

	err = fpDynamic.ValidateVnetPermissions(ctx, vnetID)
	if err != nil {
		return translateError(err, api.CloudErrorCodeInvalidResourceProviderPermissions, "resource provider")
	}

	err = fpDynamic.ValidateSubnetsRouteTablesPermissions(ctx, subnetIDs)
	if err != nil {
		return translateError(err, api.CloudErrorCodeInvalidResourceProviderPermissions, "resource provider")
	}

	// SP validation
	err = ValidateServicePrincipalProfile(ctx, dv.log, dv.env.Environment(), dv.oc.Properties.ServicePrincipalProfile.ClientID, dv.oc.Properties.ServicePrincipalProfile.ClientSecret, dv.subscriptionDoc.Subscription.Properties.TenantID)
	if err != nil {
		return translateError(err, api.CloudErrorCodeInvalidServicePrincipalPermissions, "service principal")
	}

	token, err := aad.GetToken(ctx, dv.log, dv.oc.Properties.ServicePrincipalProfile.ClientID, dv.oc.Properties.ServicePrincipalProfile.ClientSecret, dv.subscriptionDoc.Subscription.Properties.TenantID, dv.env.Environment().ActiveDirectoryEndpoint, dv.env.Environment().ResourceManagerEndpoint)
	if err != nil {
		return translateError(err, api.CloudErrorCodeInvalidServicePrincipalPermissions, "service principal")
	}

	spAuthorizer := refreshable.NewAuthorizer(token)

	spDynamic, err := dynamic.NewValidator(dv.log, dv.env.Environment(), dv.subscriptionDoc.ID, spAuthorizer)
	if err != nil {
		return err
	}

	err = spDynamic.ValidateVnetPermissions(ctx, vnetID)
	if err != nil {
		return translateError(err, api.CloudErrorCodeInvalidServicePrincipalPermissions, "service principal")
	}

	err = spDynamic.ValidateSubnetsRouteTablesPermissions(ctx, subnetIDs)
	if err != nil {
		return translateError(err, api.CloudErrorCodeInvalidServicePrincipalPermissions, "service principal")
	}

	// Additional checks - use any dynamic because they both have the correct permissions
	err = spDynamic.ValidateVnetDNS(ctx, vnetID)
	if err != nil {
		return translateError(err, api.CloudErrorCodeInvalidServicePrincipalPermissions, "service principal")
	}

	err = spDynamic.ValidateSubnetsCIDRRanges(ctx, subnetIDs, dv.oc.Properties.NetworkProfile.PodCIDR, dv.oc.Properties.NetworkProfile.ServiceCIDR)
	if err != nil {
		return translateError(err, api.CloudErrorCodeInvalidServicePrincipalPermissions, "service principal")
	}

	err = dv.validateCIDRRanges(ctx, &vnet)
	if err != nil {
		return err
	}

	err = dv.validateVnetLocation(ctx, &vnet)
	if err != nil {
		return err
	}

	err = dv.validateProviders(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (dv *openShiftClusterDynamicValidator) validateVnetLocation(ctx context.Context, vnet *mgmtnetwork.VirtualNetwork) error {
	dv.log.Print("validateVnetLocation")

	if !strings.EqualFold(*vnet.Location, dv.oc.Location) {
		return api.NewCloudError(http.StatusBadRequest, api.CloudErrorCodeInvalidLinkedVNet, "", "The vnet location '%s' must match the cluster location '%s'.", *vnet.Location, dv.oc.Location)
	}

	return nil
}

func (dv *openShiftClusterDynamicValidator) validateCIDRRanges(ctx context.Context, vnet *mgmtnetwork.VirtualNetwork) error {
	dv.log.Print("validateCIDRRanges")

	var subnets []string
	var CIDRArray []*net.IPNet

	// unique names of subnets from all node pools
	for i, subnet := range dv.oc.Properties.WorkerProfiles {
		exists := false
		for _, s := range subnets {
			if strings.EqualFold(strings.ToLower(subnet.SubnetID), strings.ToLower(s)) {
				exists = true
				break
			}
		}

		if !exists {
			subnets = append(subnets, subnet.SubnetID)
			path := fmt.Sprintf("properties.workerProfiles[%d].subnetId", i)
			c, err := dv.validateSubnet(ctx, vnet, path, subnet.SubnetID)
			if err != nil {
				return err
			}

			CIDRArray = append(CIDRArray, c)
		}
	}

	masterCIDR, err := dv.validateSubnet(ctx, vnet, "properties.MasterProfile.subnetId", dv.oc.Properties.MasterProfile.SubnetID)
	if err != nil {
		return err
	}
	_, podCIDR, err := net.ParseCIDR(dv.oc.Properties.NetworkProfile.PodCIDR)
	if err != nil {
		return err
	}

	_, serviceCIDR, err := net.ParseCIDR(dv.oc.Properties.NetworkProfile.ServiceCIDR)
	if err != nil {
		return err
	}

	CIDRArray = append(CIDRArray, masterCIDR, podCIDR, serviceCIDR)

	err = cidr.VerifyNoOverlap(CIDRArray, &net.IPNet{IP: net.IPv4zero, Mask: net.IPMask(net.IPv4zero)})
	if err != nil {
		return api.NewCloudError(http.StatusBadRequest, api.CloudErrorCodeInvalidLinkedVNet, "", "The provided CIDRs must not overlap: '%s'.", err)
	}

	return nil
}

func (dv *openShiftClusterDynamicValidator) validateSubnet(ctx context.Context, vnet *mgmtnetwork.VirtualNetwork, path, subnetID string) (*net.IPNet, error) {
	dv.log.Printf("validateSubnet (%s)", path)

	s := findSubnet(vnet, subnetID)
	if s == nil {
		return nil, api.NewCloudError(http.StatusBadRequest, api.CloudErrorCodeInvalidLinkedVNet, path, "The provided subnet '%s' could not be found.", subnetID)
	}

	if strings.EqualFold(dv.oc.Properties.MasterProfile.SubnetID, subnetID) {
		if s.PrivateLinkServiceNetworkPolicies == nil ||
			!strings.EqualFold(*s.PrivateLinkServiceNetworkPolicies, "Disabled") {
			return nil, api.NewCloudError(http.StatusBadRequest, api.CloudErrorCodeInvalidLinkedVNet, path, "The provided subnet '%s' is invalid: must have privateLinkServiceNetworkPolicies disabled.", subnetID)
		}
	}

	var found bool
	if s.ServiceEndpoints != nil {
		for _, se := range *s.ServiceEndpoints {
			if strings.EqualFold(*se.Service, "Microsoft.ContainerRegistry") &&
				se.ProvisioningState == mgmtnetwork.Succeeded {
				found = true
				break
			}
		}
	}
	if !found {
		return nil, api.NewCloudError(http.StatusBadRequest, api.CloudErrorCodeInvalidLinkedVNet, path, "The provided subnet '%s' is invalid: must have Microsoft.ContainerRegistry serviceEndpoint.", subnetID)
	}

	if dv.oc.Properties.ProvisioningState == api.ProvisioningStateCreating {
		if s.SubnetPropertiesFormat != nil &&
			s.SubnetPropertiesFormat.NetworkSecurityGroup != nil {
			return nil, api.NewCloudError(http.StatusBadRequest, api.CloudErrorCodeInvalidLinkedVNet, path, "The provided subnet '%s' is invalid: must not have a network security group attached.", subnetID)
		}

	} else {
		nsgID, err := subnet.NetworkSecurityGroupID(dv.oc, *s.ID)
		if err != nil {
			return nil, err
		}

		if s.SubnetPropertiesFormat == nil ||
			s.SubnetPropertiesFormat.NetworkSecurityGroup == nil ||
			!strings.EqualFold(*s.SubnetPropertiesFormat.NetworkSecurityGroup.ID, nsgID) {
			return nil, api.NewCloudError(http.StatusBadRequest, api.CloudErrorCodeInvalidLinkedVNet, path, "The provided subnet '%s' is invalid: must have network security group '%s' attached.", subnetID, nsgID)
		}
	}

	_, net, err := net.ParseCIDR(*s.AddressPrefix)
	if err != nil {
		return nil, err
	}

	ones, _ := net.Mask.Size()
	if ones > 27 {
		return nil, api.NewCloudError(http.StatusBadRequest, api.CloudErrorCodeInvalidLinkedVNet, path, "The provided subnet '%s' is invalid: must be /27 or larger.", subnetID)
	}

	return net, nil
}

func (dv *openShiftClusterDynamicValidator) validateProviders(ctx context.Context) error {
	dv.log.Print("validateProviders")

	providers, err := dv.providers.List(ctx, nil, "")
	if err != nil {
		return err
	}

	providerMap := make(map[string]mgmtfeatures.Provider, len(providers))

	for _, provider := range providers {
		providerMap[*provider.Namespace] = provider
	}

	for _, provider := range []string{
		"Microsoft.Authorization",
		"Microsoft.Compute",
		"Microsoft.Network",
		"Microsoft.Storage",
	} {
		if providerMap[provider].RegistrationState == nil ||
			*providerMap[provider].RegistrationState != "Registered" {
			return api.NewCloudError(http.StatusBadRequest, api.CloudErrorResourceProviderNotRegistered, "", "The resource provider '%s' is not registered.", provider)
		}
	}

	return nil
}

// ValidateServicePrincipalProfile validates the cluster service principal against the Azure graph endpoint to ensure that the service principal doesn't have
// the Application.ReadWrite.OwnedBy permission
func ValidateServicePrincipalProfile(ctx context.Context, log *logrus.Entry, env *azure.Environment, clientID string, clientSecret api.SecureString, tenantID string) error {
	// TODO: once aad.GetToken is mockable, write a unit test for this function

	log.Print("validateServicePrincipalProfile")

	token, err := aad.GetToken(ctx, log, clientID, clientSecret, tenantID, env.ActiveDirectoryEndpoint, env.GraphEndpoint)
	if err != nil {
		return err
	}

	p := &jwt.Parser{}
	c := &azureclaim.AzureClaim{}
	_, _, err = p.ParseUnverified(token.OAuthToken(), c)
	if err != nil {
		return err
	}

	for _, role := range c.Roles {
		if role == "Application.ReadWrite.OwnedBy" {
			return &PermissionError{ResourceType: servicePrincipalResource, Message: "must not have the Application.ReadWrite.OwnedBy permission"}
		}
	}

	return nil
}

// translate an error from validate package into a CloudError type
func translateError(err error, code string, typ string) error {
	switch err {
	case err.(*PermissionError):
		tErr := err.(*PermissionError)

		switch tErr.ResourceType {
		case vnetResource, subnetResource, routeTableResource:
			return api.NewCloudError(http.StatusBadRequest, code, "", "The %s does not have Network Contributor permission on %s '%s'", typ, tErr.ResourceType, tErr.ResourceID)
		case servicePrincipalResource:
			return api.NewCloudError(http.StatusBadRequest, code, "", "The provided service principal must not have the Application.ReadWrite.OwnedBy permission.")
		}

	case err.(*NotFoundError):
		tErr := err.(*NotFoundError)

		switch tErr.ResourceType {
		case vnetResource, subnetResource:
			return api.NewCloudError(http.StatusBadRequest, api.CloudErrorCodeInvalidLinkedVNet, "", "The %s '%s' could not be found.", vnetResource, tErr.ResourceID)
		case routeTableResource:
			return api.NewCloudError(http.StatusBadRequest, api.CloudErrorCodeInvalidLinkedRouteTable, "", "The %s '%s' could not be found.", routeTableResource, tErr.ResourceID)
		}

	case err.(*InvalidResourceError):
		tErr := err.(*InvalidResourceError)
		return api.NewCloudError(http.StatusBadRequest, api.CloudErrorCodeInvalidLinkedVNet, "", "The provided %s '%s' is invalid: %s", tErr.ResourceType, tErr.ResourceID, tErr.Message)

	case err.(*InvalidCredentialsError):
		return api.NewCloudError(http.StatusBadRequest, api.CloudErrorCodeInvalidServicePrincipalCredentials, "properties.servicePrincipalProfile", "The provided service principal credentials are invalid")

	case err.(*InvalidTokenClaims):
		return api.NewCloudError(http.StatusBadRequest, api.CloudErrorCodeInvalidServicePrincipalClaims, "properties.servicePrincipalProfile", "The provided service principal does not give an access token with at least one of the claims 'altsecid', 'oid', or 'puid'.")
	}

	return err
}
