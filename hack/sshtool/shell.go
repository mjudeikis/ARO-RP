package main

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func (s *sshTool) shell(ctx context.Context, log *logrus.Entry) error {
	err := s.kubeconfig()
	if err != nil {
		return err
	}

	err = s.az(ctx, log)
	if err != nil {
		s.log.Warn(err)
	}

	done, err := s.agent()
	if err != nil {
		return err
	}
	defer done()

	// This is where We can add adhoc code todo things.

	fmt.Printf("ssh -A -p 2200 core@%s\n", s.oc.Properties.NetworkProfile.APIServerPrivateEndpointIP)

	c := &exec.Cmd{
		Path: "/bin/bash",
		Env: append(os.Environ(),
			fmt.Sprintf("KUBECONFIG=%s", kubeconfigPath),
			fmt.Sprintf("SSH_AUTH_SOCK=%s", agentPath),
		),
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	return c.Run()
}
