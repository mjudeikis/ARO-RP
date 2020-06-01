package network

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	mgmtnetwork "github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-07-01/network"
	"github.com/Azure/go-autorest/autorest"
)

// SecurityGroupsClient is a minimal interface for azure SecurityGroupsClient
type SecurityGroupsClient interface {
	SecurityGroupsClientAddons
}

type securityGroupsClient struct {
	mgmtnetwork.SecurityGroupsClient
}

var _ SecurityGroupsClient = &securityGroupsClient{}

// NewSecurityGroupsClient creates a new SecurityGroupsClient
func NewSecurityGroupsClient(subscriptionID string, authorizer autorest.Authorizer) SecurityGroupsClient {
	client := mgmtnetwork.NewSecurityGroupsClient(subscriptionID)
	client.Authorizer = authorizer

	return &securityGroupsClient{
		SecurityGroupsClient: client,
	}
}
