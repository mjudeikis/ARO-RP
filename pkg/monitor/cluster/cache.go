package cluster

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (mon *Monitor) initCache(ctx context.Context) error {
	var err error
	if mon.configcli != nil {
		mon.cache.cv, err = mon.configcli.ConfigV1().ClusterVersions().Get("version", metav1.GetOptions{})
		if err != nil {
			return err
		}
		mon.cache.cos, err = mon.configcli.ConfigV1().ClusterOperators().List(metav1.ListOptions{})
		if err != nil {
			return err
		}
	}
	if mon.cli != nil {
		mon.cache.ns, err = mon.cli.CoreV1().Nodes().List(metav1.ListOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}
