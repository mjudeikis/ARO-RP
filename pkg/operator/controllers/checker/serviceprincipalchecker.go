package checker

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"

	"github.com/Azure/go-autorest/autorest/azure"
	maoclient "github.com/openshift/machine-api-operator/pkg/generated/clientset/versioned"
	"github.com/operator-framework/operator-sdk/pkg/status"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/Azure/ARO-RP/pkg/api"
	"github.com/Azure/ARO-RP/pkg/api/validate"
	arov1alpha1 "github.com/Azure/ARO-RP/pkg/operator/apis/aro.openshift.io/v1alpha1"
	aroclient "github.com/Azure/ARO-RP/pkg/operator/clientset/versioned"
	"github.com/Azure/ARO-RP/pkg/operator/controllers"
	"github.com/Azure/ARO-RP/pkg/util/aad"
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

func (r *ServicePrincipalChecker) servicePrincipalValid(ctx context.Context) error {
	cluster, err := r.arocli.AroV1alpha1().Clusters().Get(ctx, arov1alpha1.SingletonClusterName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	resource, err := azure.ParseResourceID(cluster.Spec.ResourceID)
	if err != nil {
		return err
	}

	azEnv, err := azure.EnvironmentFromName(cluster.Spec.AZEnvironment)
	if err != nil {
		return err
	}

	azCred, err := azCredentials(ctx, r.kubernetescli)
	if err != nil {
		return err
	}

	token, err := aad.GetToken(ctx, r.log, azCred.clientID, api.SecureString(azCred.clientSecret), azCred.tenantID, azEnv.ActiveDirectoryEndpoint, azEnv.ResourceManagerEndpoint)
	if err != nil {
		return err
	}

	authorizer, err := newAuthorizer(token)
	if err != nil {
		return err
	}

	masterSubnetID, workerSubnetIDs, err := getSubnetIDs(ctx, cluster.Spec.VnetID, r.clustercli)
	if err != nil {
		return err
	}

	validator, err := validate.NewValidator(r.log, &azEnv, masterSubnetID, workerSubnetIDs, resource.SubscriptionID, authorizer)
	if err != nil {
		return err
	}

	err = validator.ValidateVnetPermissions(ctx)
	if err != nil {
		if cloudErr, ok := err.(*api.CloudError); ok && cloudErr.Code == "" {
			cloudErr.Code = api.CloudErrorCodeInvalidServicePrincipalPermissions
		}
	}

	return err
}

func (r *ServicePrincipalChecker) Name() string {
	return "ServicePrincipalChecker"
}

func (r *ServicePrincipalChecker) Check(ctx context.Context) error {
	cond := &status.Condition{
		Type:    arov1alpha1.ServicePrincipalValid,
		Status:  corev1.ConditionTrue,
		Message: "service principal is valid",
		Reason:  "CheckDone",
	}

	err := r.servicePrincipalValid(ctx)
	if err != nil {
		//TODO - bvesel - check detailed error before setting message
		//TODO - bvesel: how should we form at error messages
		cond.Status = corev1.ConditionFalse
		cond.Message = err.Error()
	}

	return controllers.SetCondition(ctx, r.arocli, cond, r.role)
}
