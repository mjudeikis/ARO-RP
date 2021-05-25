package conditions

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/clock"

	arov1alpha1 "github.com/Azure/ARO-RP/pkg/operator/apis/aro.openshift.io/v1alpha1"
	aroclient "github.com/Azure/ARO-RP/pkg/operator/clientset/versioned"
	arofake "github.com/Azure/ARO-RP/pkg/operator/clientset/versioned/fake"
	"github.com/Azure/ARO-RP/pkg/util/cmp"
)

func TestSetCondition(t *testing.T) {
	ctx := context.Background()
	role := "master"
	objectName := "cluster"

	kubeclock = &clock.FakeClock{}
	var transitionTime metav1.Time = metav1.Time{Time: kubeclock.Now()}

	for _, tt := range []struct {
		name      string
		aroclient aroclient.Interface
		input     corev1.PodCondition

		expected []corev1.PodCondition
		wantErr  error
	}{
		{
			name: "noop",
			aroclient: arofake.NewSimpleClientset(&arov1alpha1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: objectName,
				},
			}),
			expected: []corev1.PodCondition{},
		},
		{
			name: "noop with condition",
			aroclient: arofake.NewSimpleClientset(&arov1alpha1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: objectName,
				},
				Status: arov1alpha1.ClusterStatus{
					Conditions: []corev1.PodCondition{
						{
							Type:   arov1alpha1.InternetReachableFromMaster,
							Status: corev1.ConditionFalse,
						},
					},
				},
			}),
			expected: []corev1.PodCondition{
				{
					Type:   arov1alpha1.InternetReachableFromMaster,
					Status: corev1.ConditionFalse,
				},
			},
		},
		{
			name: "with condition change",
			aroclient: arofake.NewSimpleClientset(&arov1alpha1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: objectName,
				},
				Status: arov1alpha1.ClusterStatus{
					Conditions: []corev1.PodCondition{
						{
							Type:   arov1alpha1.InternetReachableFromMaster,
							Status: corev1.ConditionFalse,
						},
					},
				},
			}),
			input: corev1.PodCondition{
				Type:   arov1alpha1.InternetReachableFromMaster,
				Status: corev1.ConditionTrue,
			},
			expected: []corev1.PodCondition{
				{
					Type:               arov1alpha1.InternetReachableFromMaster,
					Status:             corev1.ConditionTrue,
					LastTransitionTime: transitionTime,
				},
			},
		},
		{
			name: "preserve with condition change",
			aroclient: arofake.NewSimpleClientset(&arov1alpha1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: objectName,
				},
				Status: arov1alpha1.ClusterStatus{
					Conditions: []corev1.PodCondition{
						{
							Type:   arov1alpha1.InternetReachableFromMaster,
							Status: corev1.ConditionFalse,
						},
						{
							Type:   arov1alpha1.MachineValid,
							Status: corev1.ConditionFalse,
						},
					},
				},
			}),
			input: corev1.PodCondition{
				Type:   arov1alpha1.InternetReachableFromMaster,
				Status: corev1.ConditionTrue,
			},
			expected: []corev1.PodCondition{
				{
					Type:               arov1alpha1.InternetReachableFromMaster,
					Status:             corev1.ConditionTrue,
					LastTransitionTime: transitionTime,
				},
				{
					Type:   arov1alpha1.MachineValid,
					Status: corev1.ConditionFalse,
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

func TestIsConditions(t *testing.T) {
	for _, tt := range []struct {
		name       string
		conditions []corev1.PodCondition
		t          corev1.PodConditionType
		f          func([]corev1.PodCondition, corev1.PodConditionType) bool
		expect     bool
	}{
		{
			name: "IsTrue - non-existing",
			conditions: []corev1.PodCondition{
				{
					Type:   arov1alpha1.InternetReachableFromWorker,
					Status: corev1.ConditionTrue,
				},
			},
			t:      arov1alpha1.InternetReachableFromMaster,
			f:      IsTrue,
			expect: false,
		},
		{
			name: "IsTrue - true",
			conditions: []corev1.PodCondition{
				{
					Type:   arov1alpha1.InternetReachableFromMaster,
					Status: corev1.ConditionTrue,
				},
			},
			t:      arov1alpha1.InternetReachableFromMaster,
			f:      IsTrue,
			expect: true,
		},
		{
			name: "IsTrue - false",
			conditions: []corev1.PodCondition{
				{
					Type:   arov1alpha1.InternetReachableFromMaster,
					Status: corev1.ConditionFalse,
				},
			},
			t:      arov1alpha1.InternetReachableFromMaster,
			f:      IsTrue,
			expect: false,
		},
		{
			name: "IsFalse - true",
			conditions: []corev1.PodCondition{
				{
					Type:   arov1alpha1.InternetReachableFromMaster,
					Status: corev1.ConditionFalse,
				},
			},
			t:      arov1alpha1.InternetReachableFromMaster,
			f:      IsFalse,
			expect: true,
		},
		{
			name: "IsFalse - false",
			conditions: []corev1.PodCondition{
				{
					Type:   arov1alpha1.InternetReachableFromMaster,
					Status: corev1.ConditionTrue,
				},
			},
			t:      arov1alpha1.InternetReachableFromMaster,
			f:      IsFalse,
			expect: false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.f(tt.conditions, tt.t)
			if result != tt.expect {
				t.Fatalf("expected %t, got %t", tt.expect, result)
			}
		})
	}
}
