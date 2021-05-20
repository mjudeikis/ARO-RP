package controllers

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeclock "k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/client-go/util/retry"

	"github.com/Azure/ARO-RP/pkg/operator"
	arov1alpha1 "github.com/Azure/ARO-RP/pkg/operator/apis/aro.openshift.io/v1alpha1"
	aroclient "github.com/Azure/ARO-RP/pkg/operator/clientset/versioned"
	"github.com/Azure/ARO-RP/pkg/util/version"
)

// clock is used to set status condition timestamps.
// This variable makes it easier to test conditions.
var clock kubeclock.Clock = &kubeclock.RealClock{}

func SetCondition(ctx context.Context, arocli aroclient.Interface, cond *arov1alpha1.Condition, role string) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		cluster, err := arocli.AroV1alpha1().Clusters().Get(ctx, arov1alpha1.SingletonClusterName, metav1.GetOptions{})
		if err != nil {
			return err
		}

		changed := setCondition(cluster.Status.Conditions, *cond)

		if cleanStaleConditions(cluster, role) {
			changed = true
		}

		if !changed {
			return nil
		}

		_, err = arocli.AroV1alpha1().Clusters().UpdateStatus(ctx, cluster, metav1.UpdateOptions{})
		return err
	})
}

// cleanStaleConditions ensures that conditions no longer in use as defined by older operators
// are always removed from the cluster.status.conditions
func cleanStaleConditions(cluster *arov1alpha1.Cluster, role string) (changed bool) {
	conditions := make([]arov1alpha1.Condition, 0, len(cluster.Status.Conditions))

	// cleanup any old conditions
	current := map[arov1alpha1.ConditionType]bool{}
	for _, ct := range arov1alpha1.AllConditionTypes() {
		current[ct] = true
	}

	for _, cond := range cluster.Status.Conditions {
		if _, ok := current[cond.Type]; ok {
			conditions = append(conditions, cond)
		} else {
			changed = true
		}
	}

	cluster.Status.Conditions = conditions

	if role == operator.RoleMaster && cluster.Status.OperatorVersion != version.GitCommit {
		cluster.Status.OperatorVersion = version.GitCommit
		changed = true
	}

	return
}

// SetCondition adds (or updates) the set of conditions with the given
// condition. It returns a boolean value indicating whether the set condition
// is new or was a change to the existing condition with the same type.
func setCondition(conditions []arov1alpha1.Condition, newCond arov1alpha1.Condition) bool {
	newCond.LastTransitionTime = metav1.Time{Time: clock.Now()}

	for i, condition := range conditions {
		if condition.Type == newCond.Type {
			if condition.Status == newCond.Status {
				newCond.LastTransitionTime = condition.LastTransitionTime
			}
			changed := condition.Status != newCond.Status ||
				condition.Reason != newCond.Reason ||
				condition.Message != newCond.Message
			(conditions)[i] = newCond
			return changed
		}
	}
	conditions = append(conditions, newCond)
	return true
}
