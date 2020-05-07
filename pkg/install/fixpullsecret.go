package install

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"

	"github.com/Azure/ARO-RP/pkg/util/pullsecret"
)

func (i *Installer) fixPullSecret(ctx context.Context) error {
	// TODO: this function does not currently reapply a pull secret in
	// development mode.

	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		var ps *v1.Secret
		var err error
		var isCreate bool
		ps, err = i.kubernetescli.CoreV1().Secrets("openshift-config").Get("pull-secret", metav1.GetOptions{})
		if err != nil {
			if !errors.IsNotFound(err) {
				return err
			}
			ps = &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pull-secret",
					Namespace: "openshift-config",
				},
				Type: v1.SecretTypeDockerConfigJson,
				Data: map[string][]byte{
					v1.DockerConfigJsonKey: []byte{},
				},
			}
			isCreate = true
		}
		pullSecret, changed, err := pullsecret.SetRegistryProfiles(string(ps.Data[v1.DockerConfigJsonKey]), i.doc.OpenShiftCluster.Properties.RegistryProfiles...)
		if err != nil {
			return err
		}

		if !changed {
			return nil
		}

		ps.Data[v1.DockerConfigJsonKey] = []byte(pullSecret)

		if isCreate {
			_, err = i.kubernetescli.CoreV1().Secrets("openshift-config").Create(ps)
			return err
		}
		_, err = i.kubernetescli.CoreV1().Secrets("openshift-config").Update(ps)
		return err
	})
}
