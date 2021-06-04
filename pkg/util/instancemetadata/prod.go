package instancemetadata

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/form3tech-oss/jwt-go"
	"github.com/sirupsen/logrus"

	"github.com/Azure/ARO-RP/pkg/util/azureclaim"
)

type ServicePrincipalToken interface {
	RefreshWithContext(context.Context) error
	OAuthToken() string
}

type prod struct {
	log *logrus.Entry
	instanceMetadata

	do                              func(*http.Request) (*http.Response, error)
	newServicePrincipalTokenFromMSI func(string, string) (ServicePrincipalToken, error)
}

func newProd(ctx context.Context, log *logrus.Entry) (InstanceMetadata, error) {
	p := &prod{
		log: log,
		do:  http.DefaultClient.Do,
		newServicePrincipalTokenFromMSI: func(msiEndpoint, resource string) (ServicePrincipalToken, error) {
			return adal.NewServicePrincipalTokenFromMSI(msiEndpoint, resource)
		},
	}

	err := p.populateInstanceMetadata()
	if err != nil {
		return nil, err
	}

	err = p.populateTenantIDFromMSI(ctx)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *prod) populateTenantIDFromMSI(ctx context.Context) error {
	msiEndpoint, err := adal.GetMSIVMEndpoint()
	if err != nil {
		return err
	}

	token, err := p.newServicePrincipalTokenFromMSI(msiEndpoint, p.Environment().ResourceManagerEndpoint)
	if err != nil {
		return err
	}

	err = token.RefreshWithContext(ctx)
	if err != nil {
		return err
	}

	parser := &jwt.Parser{}
	c := &azureclaim.AzureClaim{}
	_, _, err = parser.ParseUnverified(token.OAuthToken(), c)
	if err != nil {
		return err
	}

	p.tenantID = c.TenantID

	return nil
}

// populateInstanceMetadataFromEnvironment attempts to populate instance metadata from environment, it does not throw errors.
func (p *prod) populateInstanceMetadataFromEnvironment() bool {
	environment := azure.PublicCloud
	if value, found := os.LookupEnv("AZURE_ENVIRONMENT"); found {
		var err error
		environment, err = azure.EnvironmentFromName(value)
		if err == nil {
			p.environment = &environment
		}
	}

	// Getenv returns empty string if not set which indicates no value for the struct
	p.subscriptionID = os.Getenv("AZURE_SUBSCRIPTION_ID")
	p.tenantID = os.Getenv("AZURE_TENANT_ID")
	p.location = os.Getenv("LOCATION")
	p.resourceGroup = os.Getenv("RESOURCEGROUP")

	hostname, err := os.Hostname()
	if err == nil {
		p.hostname = hostname
	}

	return p.haveIntanceMetadata()
}

func (p *prod) haveIntanceMetadata() bool {
	return p.environment != nil &&
		p.subscriptionID != "" &&
		p.tenantID != "" &&
		p.location != "" &&
		p.resourceGroup != "" &&
		p.hostname != ""
}

func (p *prod) populateInstanceMetadata() error {
	if os.Getenv("AZURE_EV2") != "" {
		if p.populateInstanceMetadataFromEnvironment() {
			// everything was populated from environment, can stop.
			p.log.Info("created InstanceMetadata from Environment, not using http metadata endpoint")
			return nil
		}
	}

	// didn't find in env vars, build form metadata
	// note this doesn't work in ACI, used with Ev2

	req, err := http.NewRequest(http.MethodGet, "http://169.254.169.254/metadata/instance/compute?api-version=2019-03-11", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Metadata", "true")

	resp, err := p.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	if strings.SplitN(resp.Header.Get("Content-Type"), ";", 2)[0] != "application/json" {
		return fmt.Errorf("unexpected content type %q", resp.Header.Get("Content-Type"))
	}

	var m *struct {
		Location          string `json:"location,omitempty"`
		ResourceGroupName string `json:"resourceGroupName,omitempty"`
		SubscriptionID    string `json:"subscriptionId,omitempty"`
		AzEnvironment     string `json:"azEnvironment,omitempty"`
	}

	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return err
	}

	environment, err := azure.EnvironmentFromName(m.AzEnvironment)
	if err != nil {
		return err
	}
	p.environment = &environment
	p.subscriptionID = m.SubscriptionID
	p.location = m.Location
	p.resourceGroup = m.ResourceGroupName

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	p.hostname = hostname

	return nil
}
