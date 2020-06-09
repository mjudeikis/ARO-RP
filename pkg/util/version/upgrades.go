package version

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

type Upgrade struct {
	Version    *Version
	PullSpec   string
	MustGather string
	Latest     bool
}

var (
	Upgrades = []Upgrade{
		{
			Version:    NewVersion(4, 3, 18),
			PullSpec:   "quay.io/openshift-release-dev/ocp-release@sha256:1f0fd38ac0640646ab8e7fec6821c8928341ad93ac5ca3a48c513ab1fb63bc4b",
			MustGather: "quay.io/openshift-release-dev/ocp-v4.0-art-dev@sha256:2e10ad0fc17f39c7a83aac32a725c78d7dd39cd9bbe3ec5ca0b76dcaa98416fa",
		},
		{
			Version:    NewVersion(4, 4, 8),
			PullSpec:   "quay.io/openshift-release-dev/ocp-release@sha256:26760772b0128f11b3c3c13eeebf1238cf9b86ea5377808febdd14d4d00a2c31",
			MustGather: "quay.io/openshift-release-dev/ocp-v4.0-art-dev@sha256:4206d834810d5e8e5fd28445c4ee81d27a1265fe01022b00f0e5193d95fb5bc2",
			Latest:     true,
		},
	}
)
