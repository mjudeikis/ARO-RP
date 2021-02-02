package checker

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"errors"
	"net/http"

	"github.com/Azure/go-autorest/autorest/adal"
	jwt "github.com/form3tech-oss/jwt-go"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/Azure/ARO-RP/pkg/api"
	"github.com/Azure/ARO-RP/pkg/util/azureclaim"
	"github.com/Azure/ARO-RP/pkg/util/refreshable"
)

type credentials struct {
	clientID     []byte
	clientSecret []byte
	tenantID     []byte
}

func newAuthorizer(token *adal.ServicePrincipalToken) (refreshable.Authorizer, error) {
	p := &jwt.Parser{}
	c := &azureclaim.AzureClaim{}
	_, _, err := p.ParseUnverified(token.OAuthToken(), c)
	if err != nil {
		return nil, err
	}

	for _, role := range c.Roles {
		if role == "Application.ReadWrite.OwnedBy" {
			return nil, api.NewCloudError(http.StatusBadRequest, api.CloudErrorCodeInvalidServicePrincipalCredentials, "properties.servicePrincipalProfile", "The provided service principal must not have the Application.ReadWrite.OwnedBy permission.")
		}
	}

	return refreshable.NewAuthorizer(token), nil
}

func azCredentials(ctx context.Context, kubernetescli kubernetes.Interface) (*credentials, error) {
	var creds credentials
	var ok bool

	mysec, err := kubernetescli.CoreV1().Secrets(azureCredentialSecretNamespace).Get(ctx, azureCredentialSecretName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if creds.clientID, ok = mysec.Data["azure_client_id"]; !ok {
		return nil, errors.New("azure_client_id does not exist")
	}

	if creds.clientSecret, ok = mysec.Data["azure_client_secret"]; !ok {
		return nil, errors.New("azure_client_secret does not exist")
	}

	if creds.tenantID, ok = mysec.Data["azure_tenant_id"]; !ok {
		return nil, errors.New("azure_tenant_id does not exist")
	}

	return &creds, nil
}
