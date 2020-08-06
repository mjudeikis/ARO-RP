package cluster

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"strconv"
)

const (
	masterRoleLabel = "node-role.kubernetes.io/master"
	workerRoleLabel = "node-role.kubernetes.io/worker"
)

// emitReportingMetric emits joined metric to be able to report better on
// all clusters state in single dashboard
func (mon *Monitor) emitReportingMetric(ctx context.Context) error {
	// daily only
	if mon.dailyRun {
		var masterCount, workerCount int
		for _, node := range mon.cache.ns.Items {
			if _, ok := node.Labels[masterRoleLabel]; ok {
				masterCount++
			}
			if _, ok := node.Labels[workerRoleLabel]; ok {
				workerCount++
			}
		}

		mon.emitGauge("cluster.summary", 1, map[string]string{
			"actualVersion":     actualVersion(mon.cache.cv),
			"desiredVersion":    desiredVersion(mon.cache.cv),
			"masterCount":       strconv.Itoa(masterCount),
			"workerCount":       strconv.Itoa(workerCount),
			"provisioningState": mon.oc.Properties.ProvisioningState.String(),
		})

	}

	return nil
}
