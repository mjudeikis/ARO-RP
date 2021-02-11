package validate

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"errors"
	"net/http"
	"sort"
	"strings"
	"time"

	mgmtnetwork "github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-07-01/network"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/Azure/ARO-RP/pkg/util/azureclient/mgmt/authorization"
	"github.com/Azure/ARO-RP/pkg/util/azureclient/mgmt/network"
	utilpermissions "github.com/Azure/ARO-RP/pkg/util/permissions"
	"github.com/Azure/ARO-RP/pkg/util/refreshable"
	"github.com/Azure/ARO-RP/pkg/util/steps"
	"github.com/Azure/ARO-RP/pkg/util/subnet"
)

// DynamicValidator validates in the operator context.
type DynamicValidator interface {
	ValidateVnetPermissions(ctx context.Context) error
	ValidateRouteTablesPermissions(ctx context.Context) error
	ValidateVnetDNS(ctx context.Context) error
	// etc
	// does Quota code go in here too?
}

var _ DynamicValidator = &dynamic{}

type dynamic struct {
	log       *logrus.Entry
	vnetr     *azure.Resource
	subnetIDs []string

	permissions     authorization.PermissionsClient
	virtualNetworks virtualNetworksGetClient
}

func NewValidator(log *logrus.Entry, azEnv *azure.Environment, subnetIDs []string, subscriptionID string, authorizer refreshable.Authorizer) (*dynamic, error) {
	// TODO - there exists a possibility that the subnets passed into the validator are not part of the same vnet.
	// this logic should change to account for the fact that there may be multiple vnets that need to be validated along
	// with parsing through all associated route tables.
	if len(subnetIDs) < 1 {
		return nil, errors.New("No subnets found in the OpenShift cluster")
	}

	vnetID, _, err := subnet.Split(subnetIDs[0])
	if err != nil {
		return nil, err
	}

	vnetr, err := azure.ParseResourceID(vnetID)
	if err != nil {
		return nil, err
	}

	return &dynamic{
		log:       log,
		vnetr:     &vnetr,
		subnetIDs: subnetIDs,

		permissions:     authorization.NewPermissionsClient(azEnv, subscriptionID, authorizer),
		virtualNetworks: newVirtualNetworksCache(network.NewVirtualNetworksClient(azEnv, subscriptionID, authorizer)),
	}, nil
}

func (dv *dynamic) ValidateVnetPermissions(ctx context.Context) error {
	dv.log.Printf("ValidateVnetPermissions")

	err := dv.validateActions(ctx, dv.vnetr, []string{
		"Microsoft.Network/virtualNetworks/join/action",
		"Microsoft.Network/virtualNetworks/read",
		"Microsoft.Network/virtualNetworks/write",
		"Microsoft.Network/virtualNetworks/subnets/join/action",
		"Microsoft.Network/virtualNetworks/subnets/read",
		"Microsoft.Network/virtualNetworks/subnets/write",
	})

	if err == wait.ErrWaitTimeout {
		return &PermissionError{resourceID: dv.vnetr.String(), resourceType: vnetResource}
	}
	if detailedErr, ok := err.(autorest.DetailedError); ok &&
		detailedErr.StatusCode == http.StatusNotFound {
		return &NotFoundError{resourceID: dv.vnetr.String(), resourceType: vnetResource}
	}
	return err
}

func (dv *dynamic) ValidateRouteTablesPermissions(ctx context.Context) error {
	vnet, err := dv.virtualNetworks.Get(ctx, dv.vnetr.ResourceGroup, dv.vnetr.ResourceName, "")
	if err != nil {
		return err
	}

	m := make(map[string]struct{})

	for _, s := range dv.subnetIDs {
		rtID, err := getRouteTableID(&vnet, s)
		if err != nil {
			return err
		}

		if _, ok := m[strings.ToLower(rtID)]; ok || rtID == "" {
			continue
		}
		m[strings.ToLower(rtID)] = struct{}{}
	}

	rts := make([]string, 0, len(m))
	for rt := range m {
		rts = append(rts, rt)
	}

	sort.Slice(rts, func(i, j int) bool { return strings.Compare(rts[i], rts[j]) < 0 })

	for _, rt := range rts {
		err := dv.validateRouteTablePermissions(ctx, rt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (dv *dynamic) validateRouteTablePermissions(ctx context.Context, rtID string) error {
	dv.log.Printf("validateRouteTablePermissions(%s)", rtID)

	rtr, err := azure.ParseResourceID(rtID)
	if err != nil {
		return err
	}

	err = dv.validateActions(ctx, &rtr, []string{
		"Microsoft.Network/routeTables/join/action",
		"Microsoft.Network/routeTables/read",
		"Microsoft.Network/routeTables/write",
	})
	if err == wait.ErrWaitTimeout {
		return &PermissionError{resourceID: rtID, resourceType: routeTableResource}
	}
	if detailedErr, ok := err.(autorest.DetailedError); ok &&
		detailedErr.StatusCode == http.StatusNotFound {
		return &NotFoundError{resourceID: rtID, resourceType: routeTableResource}
	}
	return err
}

func (dv *dynamic) ValidateVnetDNS(ctx context.Context) error {
	dv.log.Print("validateVnetDNS")

	vnet, err := dv.virtualNetworks.Get(ctx, dv.vnetr.ResourceGroup, dv.vnetr.ResourceName, "")
	if err != nil {
		return err
	}

	if vnet.DhcpOptions != nil &&
		vnet.DhcpOptions.DNSServers != nil &&
		len(*vnet.DhcpOptions.DNSServers) > 0 {
		return &InvalidResourceError{resourceID: *vnet.ID, resourceType: vnetResource, message: "custom DNS servers are not supported"}
	}

	return nil
}

func (dv *dynamic) validateActions(ctx context.Context, r *azure.Resource, actions []string) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	return wait.PollImmediateUntil(20*time.Second, func() (bool, error) {
		dv.log.Debug("retry validateActions")
		permissions, err := dv.permissions.ListForResource(ctx, r.ResourceGroup, r.Provider, "", r.ResourceType, r.ResourceName)

		if detailedErr, ok := err.(autorest.DetailedError); ok &&
			detailedErr.StatusCode == http.StatusForbidden {
			return false, steps.ErrWantRefresh
		}
		if err != nil {
			return false, err
		}

		for _, action := range actions {
			ok, err := utilpermissions.CanDoAction(permissions, action)
			if !ok || err != nil {
				// TODO(jminter): I don't understand if there are genuinely
				// cases where CanDoAction can return false then true shortly
				// after. I'm a little skeptical; if it can't happen we can
				// simplify this code.  We should add a metric on this.
				return false, err
			}
		}

		return true, nil
	}, timeoutCtx.Done())
}

func getRouteTableID(vnet *mgmtnetwork.VirtualNetwork, subnetID string) (string, error) {
	s := findSubnet(vnet, subnetID)
	if s == nil {
		return "", &NotFoundError{resourceID: subnetID, resourceType: subnetResource}
	}

	if s.RouteTable == nil {
		return "", nil
	}

	return *s.RouteTable.ID, nil
}
