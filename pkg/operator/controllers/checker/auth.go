package checker

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest/azure"
	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/Azure/ARO-RP/pkg/api"
	"github.com/Azure/ARO-RP/pkg/util/aad"
	"github.com/Azure/ARO-RP/pkg/util/azureclaim"
)

type credentials struct {
	clientID     string
	clientSecret string
	tenantID     string
}

//TODO - this function is duplicated in openshiftcluster_validatedynamic.go move this to a common location
func validateServicePrincipalProfile(ctx context.Context, log *logrus.Entry, env *azure.Environment, clientID string, clientSecret api.SecureString, tenantID string) error {
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
			return api.NewCloudError(http.StatusBadRequest, api.CloudErrorCodeInvalidServicePrincipalCredentials, "properties.servicePrincipalProfile", "The provided service principal must not have the Application.ReadWrite.OwnedBy permission.")
		}
	}

	return nil
}

func azCredentials(ctx context.Context, kubernetescli kubernetes.Interface) (*credentials, error) {
	var creds credentials

	mysec, err := kubernetescli.CoreV1().Secrets(azureCredentialSecretNamespace).Get(ctx, azureCredentialSecretName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	creds.clientID = string(mysec.Data["azure_client_id"])
	creds.clientSecret = string(mysec.Data["azure_client_secret"])
	creds.tenantID = string(mysec.Data["azure_tenant_id"])

	return &creds, nil
}
