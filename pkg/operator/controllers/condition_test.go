package controllers

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"reflect"
	"testing"

	arov1alpha1 "github.com/Azure/ARO-RP/pkg/operator/apis/aro.openshift.io/v1alpha1"
	aroclient "github.com/Azure/ARO-RP/pkg/operator/clientset/versioned"
	"github.com/Azure/ARO-RP/pkg/operator/clientset/versioned/fake"
	"github.com/Azure/ARO-RP/test/util/cmp"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeclock "k8s.io/apimachinery/pkg/util/clock"
)

func TestSetCondition(t *testing.T) {
	ctx := context.Background()
	role := "master"
	objectName := "cluster"

	clock = &kubeclock.FakeClock{}
	var transitionTime metav1.Time = metav1.Time{Time: clock.Now()}

	for _, tt := range []struct {
		name      string
		aroclient aroclient.Interface
		input     arov1alpha1.Condition

		expected []arov1alpha1.Condition
		wantErr  error
	}{
		{
			name: "noop",
			aroclient: fake.NewSimpleClientset(&arov1alpha1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: objectName,
				},
			}),
			expected: []arov1alpha1.Condition{},
		},
		{
			name: "noop with condition",
			aroclient: fake.NewSimpleClientset(&arov1alpha1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: objectName,
				},
				Status: arov1alpha1.ClusterStatus{
					Conditions: []arov1alpha1.Condition{
						{
							Type:   arov1alpha1.InternetReachableFromMaster,
							Status: v1.ConditionFalse,
						},
					},
				},
			}),
			expected: []arov1alpha1.Condition{
				{
					Type:   arov1alpha1.InternetReachableFromMaster,
					Status: v1.ConditionFalse,
				},
			},
		},
		{
			name: "with condition change",
			aroclient: fake.NewSimpleClientset(&arov1alpha1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: objectName,
				},
				Status: arov1alpha1.ClusterStatus{
					Conditions: []arov1alpha1.Condition{
						{
							Type:   arov1alpha1.InternetReachableFromMaster,
							Status: v1.ConditionFalse,
						},
					},
				},
			}),
			input: arov1alpha1.Condition{
				Type:   arov1alpha1.InternetReachableFromMaster,
				Status: v1.ConditionTrue,
			},
			expected: []arov1alpha1.Condition{
				{
					Type:               arov1alpha1.InternetReachableFromMaster,
					Status:             v1.ConditionTrue,
					LastTransitionTime: transitionTime,
				},
			},
		},
		{
			name: "preserve with condition change",
			aroclient: fake.NewSimpleClientset(&arov1alpha1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: objectName,
				},
				Status: arov1alpha1.ClusterStatus{
					Conditions: []arov1alpha1.Condition{
						{
							Type:   arov1alpha1.InternetReachableFromMaster,
							Status: v1.ConditionFalse,
						},
						{
							Type:   arov1alpha1.MachineValid,
							Status: v1.ConditionFalse,
						},
					},
				},
			}),
			input: arov1alpha1.Condition{
				Type:   arov1alpha1.InternetReachableFromMaster,
				Status: v1.ConditionTrue,
			},
			expected: []arov1alpha1.Condition{
				{
					Type:               arov1alpha1.InternetReachableFromMaster,
					Status:             v1.ConditionTrue,
					LastTransitionTime: transitionTime,
				},
				{
					Type:   arov1alpha1.MachineValid,
					Status: v1.ConditionFalse,
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {

			err := SetCondition(ctx, tt.aroclient, &tt.input, role)
			if err != nil && tt.wantErr != nil {
				t.Fatal(err.Error())
			}

			result, err := tt.aroclient.AroV1alpha1().Clusters().Get(ctx, objectName, metav1.GetOptions{})
			if err != nil {
				t.Fatal(err.Error())
			}

			if !reflect.DeepEqual(result.Status.Conditions, tt.expected) {
				t.Fatal(cmp.Diff(result.Status.Conditions, tt.expected))
			}
		})
	}
}
