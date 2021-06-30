package main

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"encoding/json"
	"net/http"

	mgmtauthorization "github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2018-09-01-preview/authorization"
	mgmtfeatures "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-07-01/features"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/to"

	"github.com/Azure/ARO-RP/pkg/api"
	"github.com/Azure/ARO-RP/pkg/util/arm"
	"github.com/Azure/ARO-RP/pkg/util/azureclient"
	"github.com/Azure/ARO-RP/pkg/util/azureerrors"
	"github.com/Azure/ARO-RP/pkg/util/stringutils"
)

// deploy is intended to test customer actions on the cluster
func (s *sshTool) deploy(ctx context.Context) error {

	resourceGroup := stringutils.LastTokenByte(s.oc.Properties.ClusterProfile.ResourceGroupID, '/')

	t := &arm.Template{
		Schema:         "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
		ContentVersion: "1.0.0.0",
		Resources: []*arm.Resource{
			s.denyAssignment(),
		},
	}

	return s.deployARMTemplate(ctx, resourceGroup, "storage", t, nil)

}

func (s *sshTool) deployARMTemplate(ctx context.Context, resourceGroupName string, deploymentName string, template *arm.Template, parameters map[string]interface{}) error {
	s.log.Printf("deploying %s template", deploymentName)
	err := s.deployments.CreateOrUpdateAndWait(ctx, resourceGroupName, deploymentName, mgmtfeatures.Deployment{
		Properties: &mgmtfeatures.DeploymentProperties{
			Template:   template,
			Parameters: parameters,
			Mode:       mgmtfeatures.Incremental,
		},
	})

	if azureerrors.IsDeploymentActiveError(err) {
		s.log.Printf("waiting for %s template to be deployed", deploymentName)
		err = s.deployments.Wait(ctx, resourceGroupName, deploymentName)
	}

	if azureerrors.HasAuthorizationFailedError(err) ||
		azureerrors.HasLinkedAuthorizationFailedError(err) {
		return err
	}

	serviceErr, _ := err.(*azure.ServiceError) // futures return *azure.ServiceError directly

	// CreateOrUpdate() returns a wrapped *azure.ServiceError
	if detailedErr, ok := err.(autorest.DetailedError); ok {
		serviceErr, _ = detailedErr.Original.(*azure.ServiceError)
	}

	if serviceErr != nil {
		b, _ := json.Marshal(serviceErr)

		return &api.CloudError{
			StatusCode: http.StatusBadRequest,
			CloudErrorBody: &api.CloudErrorBody{
				Code:    api.CloudErrorCodeDeploymentFailed,
				Message: "Deployment failed.",
				Details: []api.CloudErrorBody{
					{
						Message: string(b),
					},
				},
			},
		}
	}

	return err
}

func (s *sshTool) denyAssignment() *arm.Resource {
	return &arm.Resource{
		Resource: &mgmtauthorization.DenyAssignment{
			Name: to.StringPtr("[guid(resourceGroup().id, 'ARO cluster resource group deny assignment')]"),
			Type: to.StringPtr("Microsoft.Authorization/denyAssignments"),
			DenyAssignmentProperties: &mgmtauthorization.DenyAssignmentProperties{
				DenyAssignmentName: to.StringPtr("[guid(resourceGroup().id, 'ARO cluster resource group deny assignment')]"),
				Permissions: &[]mgmtauthorization.DenyAssignmentPermission{
					{
						Actions: &[]string{
							"*/action",
							"*/delete",
							"*/write",
						},
						NotActions: &[]string{
							"Microsoft.Compute/disks/beginGetAccess/action",
							"Microsoft.Compute/disks/endGetAccess/action",
							"Microsoft.Compute/disks/write",
							"Microsoft.Compute/snapshots/beginGetAccess/action",
							"Microsoft.Compute/snapshots/delete",
							"Microsoft.Compute/snapshots/endGetAccess/action",
							"Microsoft.Compute/snapshots/write",
							"Microsoft.Network/networkInterfaces/effectiveRouteTable/action",
							"Microsoft.Network/networkSecurityGroups/join/action",
							"Microsoft.Resources/tags/write",
						},
					},
				},
				Scope: &s.oc.Properties.ClusterProfile.ResourceGroupID,
				Principals: &[]mgmtauthorization.Principal{
					{
						ID:   to.StringPtr("00000000-0000-0000-0000-000000000000"),
						Type: to.StringPtr("SystemDefined"),
					},
				},
				ExcludePrincipals: &[]mgmtauthorization.Principal{
					{
						ID:   &s.oc.Properties.ServicePrincipalProfile.SPObjectID,
						Type: to.StringPtr("ServicePrincipal"),
					},
				},
				IsSystemProtected: to.BoolPtr(true),
			},
		},
		APIVersion: azureclient.APIVersion("Microsoft.Authorization/denyAssignments"),
	}
}
