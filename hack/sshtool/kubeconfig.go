package main

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"k8s.io/client-go/tools/clientcmd"
)

const kubeconfigPath = "/tmp/admin.kubeconfig"

func (s *sshTool) kubeconfig() error {
	kubeconfig := s.oc.Properties.AROServiceKubeconfig
	if kubeconfig == nil {
		kubeconfig = s.oc.Properties.AdminKubeconfig
	}

	config, err := clientcmd.Load(kubeconfig)
	if err != nil {
		return err
	}

	for _, cluster := range config.Clusters {
		cluster.Server = "https://" + s.oc.Properties.NetworkProfile.APIServerPrivateEndpointIP + ":6443"
		cluster.CertificateAuthorityData = nil
		cluster.InsecureSkipTLSVerify = true
	}

	return clientcmd.WriteToFile(*config, kubeconfigPath)
}
