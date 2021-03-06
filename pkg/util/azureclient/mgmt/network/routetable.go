package network

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"

	mgmtnetwork "github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-07-01/network"
	"github.com/Azure/go-autorest/autorest"
)

// RouteTablesClient is a minimal interface for azure RouteTablesClient
type RouteTablesClient interface {
	Get(ctx context.Context, resourceGroupName string, routeTableName string, expand string) (result mgmtnetwork.RouteTable, err error)
	RouteTablesClientAddons
}

type routeTablesClient struct {
	mgmtnetwork.RouteTablesClient
}

var _ RouteTablesClient = &routeTablesClient{}

// NewRouteTablesClient creates a new RouteTablesClient
func NewRouteTablesClient(subscriptionID string, authorizer autorest.Authorizer) RouteTablesClient {
	client := mgmtnetwork.NewRouteTablesClient(subscriptionID)
	client.Authorizer = authorizer

	return &routeTablesClient{
		RouteTablesClient: client,
	}
}
