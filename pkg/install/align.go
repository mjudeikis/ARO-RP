package install

import (
	"context"

	mgmtauthorization "github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2018-09-01-preview/authorization"
	mgmtfeatures "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-07-01/features"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/Azure/ARO-RP/pkg/api"
	"github.com/Azure/ARO-RP/pkg/util/arm"
	"github.com/Azure/ARO-RP/pkg/util/azureclient"
	"github.com/Azure/ARO-RP/pkg/util/azureclient/graphrbac"
	"github.com/Azure/ARO-RP/pkg/util/stringutils"
)

func (i *Installer) CreateOrUpdateDenyAssignment(ctx context.Context, doc *api.OpenShiftClusterDocument) error {
	spp := doc.OpenShiftCluster.Properties.ServicePrincipalProfile

	conf := auth.NewClientCredentialsConfig(spp.ClientID, string(spp.ClientSecret), spp.TenantID)
	conf.Resource = azure.PublicCloud.GraphEndpoint

	spGraphAuthorizer, err := conf.Authorizer()
	if err != nil {
		return err
	}

	applications := graphrbac.NewApplicationsClient(spp.TenantID, spGraphAuthorizer)

	res, err := applications.GetServicePrincipalsIDByAppID(ctx, spp.ClientID)
	if err != nil {
		return err
	}

	clusterSPObjectID := *res.Value

	t := &arm.Template{
		Schema:         "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
		ContentVersion: "1.0.0.0",
		Resources: []*arm.Resource{
			{
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
									"Microsoft.Network/networkSecurityGroups/join/action",
								},
							},
						},
						Scope: &doc.OpenShiftCluster.Properties.ClusterProfile.ResourceGroupID,
						Principals: &[]mgmtauthorization.Principal{
							{
								ID:   to.StringPtr("00000000-0000-0000-0000-000000000000"),
								Type: to.StringPtr("SystemDefined"),
							},
						},
						ExcludePrincipals: &[]mgmtauthorization.Principal{
							{
								ID:   &clusterSPObjectID,
								Type: to.StringPtr("ServicePrincipal"),
							},
						},
						IsSystemProtected: to.BoolPtr(true),
					},
				},
				APIVersion: azureclient.APIVersions["Microsoft.Authorization/denyAssignments"],
			},
		},
	}

	resourceGroup := stringutils.LastTokenByte(doc.OpenShiftCluster.Properties.ClusterProfile.ResourceGroupID, '/')

	i.log.Info("deploying")
	err = i.deployments.CreateOrUpdateAndWait(ctx, resourceGroup, "denyassignment", mgmtfeatures.Deployment{
		Properties: &mgmtfeatures.DeploymentProperties{
			Template: t,
			Mode:     mgmtfeatures.Incremental,
		},
	})
	if err != nil {
		return err
	}

	i.log.Info("deleting deployment")
	return i.deployments.DeleteAndWait(ctx, resourceGroup, "denyassignment")
}

func (i *Installer) InstallerFixups(ctx context.Context, doc *api.OpenShiftClusterDocument) error {
	err := i.initializeKubernetesClients(ctx)
	if err != nil {
		return err
	}

	i.log.Info("creating billing record")
	err = i.createBillingRecord(ctx)
	if err != nil {
		return err
	}

	i.log.Info("disable alertmanager warning")
	err = i.disableAlertManagerWarning(ctx)
	if err != nil {
		return err
	}

	i.log.Info("remove bootstrap ignition")
	err = i.removeBootstrapIgnition(ctx)
	if err != nil {
		return err
	}

	i.log.Info("ensure genevaLogging")
	return i.ensureGenevaLogging(ctx)
}

func (i *Installer) ConfigurationFixup(ctx context.Context, doc *api.OpenShiftClusterDocument) error {
	scc, err := i.securitycli.SecurityV1().SecurityContextConstraints().Get("privileged", metav1.GetOptions{})
	if err != nil {
		return err
	}

	var needsUpdate bool
	var users []string
	for _, u := range scc.Users {
		if u != "system:serviceaccount:openshift-azure-logging:geneva" {
			users = append(users, u)
		} else {
			needsUpdate = true
		}
	}

	if needsUpdate {
		i.log.Info("updating privileged scc")
		scc.Users = users

		_, err := i.securitycli.SecurityV1().SecurityContextConstraints().Update(scc)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Installer) KubeConfigFixup(ctx context.Context, doc *api.OpenShiftClusterDocument) error {
	g, err := i.loadGraph(ctx)
	if err != nil {
		return err
	}

	aroServiceInternalClient, err := i.generateAROServiceKubeconfig(g)
	if err != nil {
		return err
	}

	_, err = i.db.Patch(ctx, doc.Key, func(doc *api.OpenShiftClusterDocument) error {
		if len(doc.OpenShiftCluster.Properties.AROServiceKubeconfig) == 0 {
			i.log.Print("updating aro service kubeconfig")
			doc.OpenShiftCluster.Properties.AROServiceKubeconfig = aroServiceInternalClient.File.Data
		}
		return nil
	})
	return err
}
