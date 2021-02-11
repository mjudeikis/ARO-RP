package checker

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type credentials struct {
	clientID     string
	clientSecret string
	tenantID     string
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
