package cluster

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"strings"
	"time"

	mgmtauthorization "github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2018-09-01-preview/authorization"
	"github.com/ghodss/yaml"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"

	"github.com/Azure/ARO-RP/pkg/util/arm"
	"github.com/Azure/ARO-RP/pkg/util/rbac"
	"github.com/Azure/ARO-RP/pkg/util/stringutils"
)

func (m *manager) createOrUpdateClusterServicePrincipalRBAC(ctx context.Context) error {
	resourceGroupID := m.doc.OpenShiftCluster.Properties.ClusterProfile.ResourceGroupID
	resourceGroup := stringutils.LastTokenByte(resourceGroupID, '/')
	clusterSPObjectID := m.doc.OpenShiftCluster.Properties.ServicePrincipalProfile.SPObjectID

	roleAssignments, err := m.roleAssignments.ListForResourceGroup(ctx, resourceGroup, "")
	if err != nil {
		return err
	}

	// If we have Contributor RBAC role for Cluster SP on the resource group in question
	// We are interested in Resource group scope only (inherited are returned too).
	var toDelete []mgmtauthorization.RoleAssignment
	var found bool
	for _, assignment := range roleAssignments {
		// Contributor assignments only!
		if strings.EqualFold(*assignment.Scope, resourceGroupID) && strings.HasSuffix(strings.ToLower(*assignment.RoleDefinitionID), rbac.RoleContributor) {
			if strings.EqualFold(*assignment.PrincipalID, clusterSPObjectID) {
				found = true
			} else {
				toDelete = append(toDelete, assignment)
			}
		}
	}

	for _, assignment := range toDelete {
		m.log.Infof("Deleting Contributor roleAssignment %s", *assignment.Name)
		_, err := m.roleAssignments.Delete(ctx, *assignment.Scope, *assignment.Name)
		if err != nil {
			return err
		}
	}

	if !found {
		m.log.Info("Contributor roleAssignment not found for cluster service principal. Creating")
		t := &arm.Template{
			Schema:         "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
			ContentVersion: "1.0.0.0",
			Resources:      []*arm.Resource{m.clusterServicePrincipalRBAC()},
		}
		err = m.deployARMTemplate(ctx, resourceGroup, "storage", t, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *manager) updateAROSecret(ctx context.Context) error {
	spp := m.doc.OpenShiftCluster.Properties.ServicePrincipalProfile
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		//data:
		// cloud-config: <base64 map[string]string with keys 'aadClientId' and 'aadClientSecret'>
		secret, err := m.kubernetescli.CoreV1().Secrets("kube-system").Get(ctx, "azure-cloud-provider", metav1.GetOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) { // we are not in control if secret is not present
				return nil
			}
			return err
		}

		var cf map[string]interface{}
		var changed bool
		if secret != nil && secret.Data != nil {
			err = yaml.Unmarshal(secret.Data["cloud-config"], &cf)
			if err != nil {
				return err
			}
			if val, ok := cf["aadClientId"].(string); ok {
				if val != spp.ClientID {
					cf["aadClientId"] = spp.ClientID
					changed = true
				}
			}
			if val, ok := cf["aadClientSecret"].(string); ok {
				if val != string(spp.ClientSecret) {
					cf["aadClientSecret"] = spp.ClientSecret
					changed = true
				}
			}
		}

		if changed {
			data, err := yaml.Marshal(cf)
			if err != nil {
				return err
			}
			secret.Data["cloud-config"] = data

			_, err = m.kubernetescli.CoreV1().Secrets("kube-system").Update(ctx, secret, metav1.UpdateOptions{})
			if err != nil {
				return err
			}

			// If secret change we need to trigger kube-api-server and kube-controller-manager restarts
			err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
				kAPIServer, err := m.operatorcli.OperatorV1().KubeAPIServers().Get(ctx, "cluster", metav1.GetOptions{})
				if err != nil {
					return err
				}
				kAPIServer.Spec.ForceRedeploymentReason = "Credential rotation " + time.Now().Format("2006-01-02 3:4:5")

				_, err = m.operatorcli.OperatorV1().KubeAPIServers().Update(ctx, kAPIServer, metav1.UpdateOptions{})
				if err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				return err
			}

			return retry.RetryOnConflict(retry.DefaultRetry, func() error {
				kManager, err := m.operatorcli.OperatorV1().KubeControllerManagers().Get(ctx, "cluster", metav1.GetOptions{})
				if err != nil {
					return err
				}
				kManager.Spec.ForceRedeploymentReason = "Credential rotation " + time.Now().Format("2006-01-02 3:4:5")

				_, err = m.operatorcli.OperatorV1().KubeControllerManagers().Update(ctx, kManager, metav1.UpdateOptions{})
				if err != nil {
					return err
				}
				return nil
			})
		}
		return nil
	})
}

func (m *manager) updateOpenShiftSecret(ctx context.Context) error {
	spp := m.doc.OpenShiftCluster.Properties.ServicePrincipalProfile
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		//data:
		// azure_client_id: secret_id
		// azure_client_secret: secret_value
		// azure_tenant_id: tenant_id
		secret, err := m.kubernetescli.CoreV1().Secrets("kube-system").Get(ctx, "azure-credentials", metav1.GetOptions{})
		if err != nil {
			return err
		}

		var changed bool
		if string(secret.Data["azure_client_id"]) != spp.ClientID {
			secret.Data["azure_client_id"] = []byte(spp.ClientID)
			changed = true
		}

		if string(secret.Data["azure_client_secret"]) != string(spp.ClientSecret) {
			secret.Data["azure_client_secret"] = []byte(spp.ClientSecret)
			changed = true
		}

		if string(secret.Data["azure_tenant_id"]) != m.subscriptionDoc.Subscription.Properties.TenantID {
			secret.Data["azure_tenant_id"] = []byte(m.subscriptionDoc.Subscription.Properties.TenantID)
			changed = true
		}

		if changed {
			_, err = m.kubernetescli.CoreV1().Secrets("kube-system").Update(ctx, secret, metav1.UpdateOptions{})
			if err != nil {
				return err
			}

			// restart cloud credentials operator to trigger rotation
			return retry.RetryOnConflict(retry.DefaultRetry, func() error {
				return m.kubernetescli.CoreV1().Pods("openshift-cloud-credential-operator").DeleteCollection(ctx, metav1.DeleteOptions{}, metav1.ListOptions{
					LabelSelector: "app=cloud-credential-operator",
				})
			})

		}

		return nil
	})
}