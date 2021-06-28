package monitoring

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/ugorji/go/codec"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	ctrl "sigs.k8s.io/controller-runtime"

	arov1alpha1 "github.com/Azure/ARO-RP/pkg/operator/apis/aro.openshift.io/v1alpha1"
	arofake "github.com/Azure/ARO-RP/pkg/operator/clientset/versioned/fake"
	"github.com/Azure/ARO-RP/pkg/util/cmp"
)

var (
	cmMetadata = metav1.ObjectMeta{Name: "cluster-monitoring-config", Namespace: "openshift-monitoring"}

	arocli = arofake.NewSimpleClientset(&arov1alpha1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name: arov1alpha1.SingletonClusterName,
		},
		Spec: arov1alpha1.ClusterSpec{
			Features: arov1alpha1.FeaturesSpec{
				ReconcileMonitoring: true,
			},
		},
	})
)

func TestReconcileMonitoringConfig(t *testing.T) {
	log := logrus.NewEntry(logrus.StandardLogger())
	type test struct {
		name       string
		reconciler func() *Reconciler
		wantConfig string
	}

	for _, tt := range []*test{
		{
			name: "ConfigMap does not exist - enable",
			reconciler: func() *Reconciler {
				return &Reconciler{
					arocli:        arocli,
					kubernetescli: fake.NewSimpleClientset(),
					log:           log,
					jsonHandle:    new(codec.JsonHandle),
				}
			},
			wantConfig: `{}`,
		},
		{
			name: "empty config.yaml",
			reconciler: func() *Reconciler {
				return &Reconciler{
					arocli: arocli,
					kubernetescli: fake.NewSimpleClientset(&corev1.ConfigMap{
						ObjectMeta: cmMetadata,
						Data: map[string]string{
							"config.yaml": ``,
						},
					}),
					log:        log,
					jsonHandle: new(codec.JsonHandle),
				}
			},
			wantConfig: `{}`,
		},
		{
			name: "settings restored to default and extra fields are preserved",
			reconciler: func() *Reconciler {
				return &Reconciler{
					arocli: arocli,
					kubernetescli: fake.NewSimpleClientset(&corev1.ConfigMap{
						ObjectMeta: cmMetadata,
						Data: map[string]string{
							"config.yaml": `
prometheusK8s:
  extraField: prometheus
  retention: 1d
  volumeClaimTemplate:
    metadata:
      name: meh
    spec:
      resources:
        requests:
          storage: 50Gi
      storageClassName: fast
      volumeMode: Filesystem
alertmanagerMain:
  extraField: yeet
  volumeClaimTemplate:
    metadata:
      name: slowest-storage
    spec:
      resources:
        requests:
          storage: 50Gi
        storageClassName: snail-mail
        volumeMode: Filesystem
`,
						},
					}),
					log:        log,
					jsonHandle: new(codec.JsonHandle),
				}
			},
			wantConfig: `
alertmanagerMain:
  extraField: yeet
prometheusK8s:
  extraField: prometheus
`,
		},
		{
			name: "empty volumeClaimTemplate struct is cleared out",
			reconciler: func() *Reconciler {
				return &Reconciler{
					arocli: arocli,
					kubernetescli: fake.NewSimpleClientset(&corev1.ConfigMap{
						ObjectMeta: cmMetadata,
						Data: map[string]string{
							"config.yaml": `
alertmanagerMain:
  volumeClaimTemplate: {}
  extraField: alertmanager
prometheusK8s:
  volumeClaimTemplate: {}
  bugs: not-here
`,
						},
					}),
					log:        log,
					jsonHandle: new(codec.JsonHandle),
				}
			},
			wantConfig: `
alertmanagerMain:
  extraField: alertmanager
prometheusK8s:
  bugs: not-here
`,
		},
		{
			name: "empty volumeClaimTemplate cleared out part 2",
			reconciler: func() *Reconciler {
				return &Reconciler{
					arocli: arocli,
					kubernetescli: fake.NewSimpleClientset(&corev1.ConfigMap{
						ObjectMeta: cmMetadata,
						Data: map[string]string{
							"config.yaml": `
alertmanagerMain:
  volumeClaimTemplate:
    spec:
      requests: 5Gi
  extraField: alertmanager
  somethingElse: {}
  hello: true
prometheusK8s:
  volumeClaimTemplate: {}
  bugs: not-here
`,
						},
					}),
					log:        log,
					jsonHandle: new(codec.JsonHandle),
				}
			},
			wantConfig: `
alertmanagerMain:
  extraField: alertmanager
  hello: true
  somethingElse: {}
prometheusK8s:
  bugs: not-here
`,
		},
		{
			name: "other monitoring components are configured",
			reconciler: func() *Reconciler {
				return &Reconciler{
					arocli: arocli,
					kubernetescli: fake.NewSimpleClientset(&corev1.ConfigMap{
						ObjectMeta: cmMetadata,
						Data: map[string]string{
							"config.yaml": `
alertmanagerMain:
  nodeSelector:
    foo: bar
somethingElse:
  configured: true
`,
						},
					}),
					log:        log,
					jsonHandle: new(codec.JsonHandle),
				}
			},
			wantConfig: `
alertmanagerMain:
  nodeSelector:
    foo: bar
somethingElse:
  configured: true
`,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			r := tt.reconciler()
			request := ctrl.Request{}
			request.Name = "cluster-monitoring-config"
			request.Namespace = "openshift-monitoring"

			_, err := r.Reconcile(ctx, request)
			if err != nil {
				t.Fatal(err)
			}

			cm, err := r.kubernetescli.CoreV1().ConfigMaps("openshift-monitoring").Get(ctx, "cluster-monitoring-config", metav1.GetOptions{})
			if err != nil {
				t.Fatal(err)
			}

			if strings.TrimSpace(cm.Data["config.yaml"]) != strings.TrimSpace(tt.wantConfig) {
				t.Error(cm.Data["config.yaml"])
			}
		})
	}
}

func TestReconcilePVC(t *testing.T) {
	log := logrus.NewEntry(logrus.StandardLogger())
	tests := []struct {
		name       string
		reconciler func() *Reconciler
		want       []corev1.PersistentVolumeClaim
		wantErr    error
	}{
		{
			name: "Should delete the prometheus PVCs",
			reconciler: func() *Reconciler {
				return &Reconciler{
					arocli: arocli,
					kubernetescli: fake.NewSimpleClientset(&corev1.PersistentVolumeClaim{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "prometheus-k8s-db-prometheus-k8s-0",
							Namespace: "openshift-monitoring",
							Labels: map[string]string{
								"app":        "prometheus",
								"prometheus": "k8s",
							},
						},
					},
						&corev1.PersistentVolumeClaim{
							ObjectMeta: metav1.ObjectMeta{
								Name:      "prometheus-k8s-db-prometheus-k8s-1",
								Namespace: "openshift-monitoring",
								Labels: map[string]string{
									"app":        "prometheus",
									"prometheus": "k8s",
								},
							},
						}),
					log:        log,
					jsonHandle: new(codec.JsonHandle),
				}
			},
			want:    nil,
			wantErr: nil,
		},
		{
			name: "Should preserve 1 pvc",
			reconciler: func() *Reconciler {
				return &Reconciler{
					arocli: arocli,
					kubernetescli: fake.NewSimpleClientset(&corev1.PersistentVolumeClaim{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "prometheus-k8s-db-prometheus-k8s-0",
							Namespace: "openshift-monitoring",
							Labels: map[string]string{
								"app":        "prometheus",
								"prometheus": "k8s",
							},
						},
					},
						&corev1.PersistentVolumeClaim{
							ObjectMeta: metav1.ObjectMeta{
								Name:      "random-pvc",
								Namespace: "openshift-monitoring",
								Labels: map[string]string{
									"app": "random",
								},
							},
						}),
					log:        log,
					jsonHandle: new(codec.JsonHandle),
				}
			},
			want: []corev1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "random-pvc",
						Namespace: "openshift-monitoring",
						Labels: map[string]string{
							"app": "random",
						},
					},
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			r := tt.reconciler()
			request := ctrl.Request{}
			request.Name = "cluster-monitoring-config"
			request.Namespace = "openshift-monitoring"

			_, err := r.Reconcile(ctx, request)
			if err != nil {
				t.Fatal(err)
			}

			pvcList, err := r.kubernetescli.CoreV1().PersistentVolumeClaims(monitoringName.Namespace).List(context.Background(), metav1.ListOptions{})
			if err != nil {
				t.Fatalf("Unexpected error during list of PVCs: %v", err)
			}

			if !reflect.DeepEqual(pvcList.Items, tt.want) {
				t.Error(cmp.Diff(pvcList.Items, tt.want))
			}
		})
	}
}
