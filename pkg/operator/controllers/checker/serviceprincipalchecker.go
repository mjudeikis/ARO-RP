package checker

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"

	"github.com/Azure/go-autorest/autorest/azure"
	maoclient "github.com/openshift/machine-api-operator/pkg/generated/clientset/versioned"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/Azure/ARO-RP/pkg/api"
	"github.com/Azure/ARO-RP/pkg/api/validate/dynamic"
	arov1alpha1 "github.com/Azure/ARO-RP/pkg/operator/apis/aro.openshift.io/v1alpha1"
	aroclient "github.com/Azure/ARO-RP/pkg/operator/clientset/versioned"
	"github.com/Azure/ARO-RP/pkg/operator/controllers"
	"github.com/Azure/ARO-RP/pkg/util/aad"
	"github.com/Azure/ARO-RP/pkg/util/azureclient"
	"github.com/Azure/ARO-RP/pkg/util/clusterauthorizer"
)

type ServicePrincipalChecker struct {
	log           *logrus.Entry
	clustercli    maoclient.Interface
	arocli        aroclient.Interface
	kubernetescli kubernetes.Interface
	role          string
}

func NewServicePrincipalChecker(log *logrus.Entry, maocli maoclient.Interface, arocli aroclient.Interface, kubernetescli kubernetes.Interface, role string) *ServicePrincipalChecker {
	return &ServicePrincipalChecker{
		log:           log,
		clustercli:    maocli,
		arocli:        arocli,
		kubernetescli: kubernetescli,
		role:          role,
	}
}

func (r *ServicePrincipalChecker) Name() string {
	return "ServicePrincipalChecker"
}

func (r *ServicePrincipalChecker) Check(ctx context.Context) error {
	cond := &arov1alpha1.Condition{
		Type:    arov1alpha1.ServicePrincipalValid,
		Status:  corev1.ConditionTrue,
		Message: "service principal is valid",
		Reason:  "CheckDone",
	}

	cluster, err := r.arocli.AroV1alpha1().Clusters().Get(ctx, arov1alpha1.SingletonClusterName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	resource, err := azure.ParseResourceID(cluster.Spec.ResourceID)
	if err != nil {
		return err
	}

	azEnv, err := azureclient.EnvironmentFromName(cluster.Spec.AZEnvironment)
	if err != nil {
		return err
	}

	azCred, err := clusterauthorizer.AzCredentials(ctx, r.kubernetescli)
	if err != nil {
		return err
	}

	_, err = aad.GetToken(ctx, r.log, string(azCred.ClientID), string(azCred.ClientSecret), string(azCred.TenantID), azEnv.ActiveDirectoryEndpoint, azEnv.ResourceManagerEndpoint)
	if err != nil {
		updateFailedCondition(cond, err)
	}

	spDynamic, err := dynamic.NewValidator(r.log, &azEnv, resource.SubscriptionID, nil, dynamic.AuthorizerClusterServicePrincipal)
	if err != nil {
		return err
	}

	err = spDynamic.ValidateServicePrincipal(ctx, string(azCred.ClientID), string(azCred.ClientSecret), string(azCred.TenantID))
	if err != nil {
		updateFailedCondition(cond, err)
	}

	return controllers.SetCondition(ctx, r.arocli, cond, r.role)
}

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

func updateFailedCondition(cond *arov1alpha1.Condition, err error) {
	cond.Status = corev1.ConditionFalse
	if tErr, ok := err.(*api.CloudError); ok {
		cond.Message = tErr.Message
	} else {
		cond.Message = err.Error()
	}
}
