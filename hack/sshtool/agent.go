package main

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"crypto/x509"
	"net"

	"golang.org/x/crypto/ssh/agent"
)

const agentPath = "/tmp/agent"

func (s *sshTool) agent() (func() error, error) {
	key, err := x509.ParsePKCS1PrivateKey(s.oc.Properties.SSHKey)
	if err != nil {
		return nil, err
	}

	keyring := agent.NewKeyring()

	err = keyring.Add(agent.AddedKey{
		PrivateKey: key,
	})
	if err != nil {
		return nil, err
	}

	l, err := net.Listen("unix", agentPath)
	if err != nil {
		return nil, err
	}

	go func() error {
		for {
			c, err := l.Accept()
			if err != nil {
				return err
			}

			go agent.ServeAgent(keyring, c)
		}
	}()

	return l.Close, nil
}
