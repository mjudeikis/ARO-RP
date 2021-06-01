package main

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/sirupsen/logrus"

	"github.com/Azure/ARO-RP/pkg/api"
	"github.com/Azure/ARO-RP/pkg/database"
	"github.com/Azure/ARO-RP/pkg/env"
	"github.com/Azure/ARO-RP/pkg/metrics/noop"
	"github.com/Azure/ARO-RP/pkg/util/azureclient/mgmt/network"
	"github.com/Azure/ARO-RP/pkg/util/encryption"
	"github.com/Azure/ARO-RP/pkg/util/keyvault"
	utillog "github.com/Azure/ARO-RP/pkg/util/log"
	"github.com/Azure/ARO-RP/pkg/util/stringutils"
)

type sshTool struct {
	log *logrus.Entry
	oc  *api.OpenShiftCluster
	sub *api.Subscription

	interfaces    network.InterfacesClient
	loadBalancers network.LoadBalancersClient

	clusterResourceGroup string
	infraID              string
}

func newSSHTool(log *logrus.Entry, env env.Interface, oc *api.OpenShiftCluster, sub *api.Subscription) (*sshTool, error) {
	r, err := azure.ParseResourceID(oc.ID)
	if err != nil {
		return nil, err
	}

	fpAuthorizer, err := env.FPAuthorizer(sub.Properties.TenantID, azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		return nil, err
	}

	infraID := oc.Properties.InfraID
	if infraID == "" {
		infraID = "aro"
	}

	return &sshTool{
		log: log,
		oc:  oc,
		sub: sub,

		interfaces:    network.NewInterfacesClient(env.Environment(), r.SubscriptionID, fpAuthorizer),
		loadBalancers: network.NewLoadBalancersClient(env.Environment(), r.SubscriptionID, fpAuthorizer),

		clusterResourceGroup: stringutils.LastTokenByte(oc.Properties.ClusterProfile.ResourceGroupID, '/'),
		infraID:              infraID,
	}, nil
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage:\n")
	fmt.Fprintf(os.Stderr, "  %s enable resourceid\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s disable resourceid\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s shell resourceid\n", os.Args[0])
}

func getCluster(ctx context.Context, log *logrus.Entry, _env env.Core, resourceID string) (*api.OpenShiftCluster, *api.Subscription, error) {
	msiAuthorizer, err := _env.NewMSIAuthorizer(env.MSIContextRP, _env.Environment().ResourceManagerEndpoint)
	if err != nil {
		return nil, nil, err
	}

	msiKVAuthorizer, err := _env.NewMSIAuthorizer(env.MSIContextRP, _env.Environment().ResourceIdentifiers.KeyVault)
	if err != nil {
		return nil, nil, err
	}

	serviceKeyvaultURI, err := keyvault.URI(_env, env.ServiceKeyvaultSuffix)
	if err != nil {
		return nil, nil, err
	}

	serviceKeyvault := keyvault.NewManager(msiKVAuthorizer, serviceKeyvaultURI)

	key, err := serviceKeyvault.GetBase64Secret(ctx, env.EncryptionSecretName)
	if err != nil {
		return nil, nil, err
	}

	aead, err := encryption.NewXChaCha20Poly1305(ctx, key)
	if err != nil {
		return nil, nil, err
	}

	dbAuthorizer, err := database.NewMasterKeyAuthorizer(ctx, _env, msiAuthorizer)
	if err != nil {
		return nil, nil, err
	}

	db, err := database.NewDatabaseClient(log.WithField("component", "database"), _env, dbAuthorizer, &noop.Noop{}, aead)
	if err != nil {
		return nil, nil, err
	}

	dbOpenShiftClusters, err := database.NewOpenShiftClusters(ctx, _env.IsLocalDevelopmentMode(), db)
	if err != nil {
		return nil, nil, err
	}

	dbSubscriptions, err := database.NewSubscriptions(ctx, _env.IsLocalDevelopmentMode(), db)
	if err != nil {
		return nil, nil, err
	}

	doc, err := dbOpenShiftClusters.Get(ctx, strings.ToLower(resourceID))
	if err != nil {
		return nil, nil, err
	}
	if doc == nil {
		return nil, nil, fmt.Errorf("resource %q not found", resourceID)
	}

	subDoc, err := dbSubscriptions.Get(ctx, strings.ToLower(resourceID))
	if err != nil {
		return nil, nil, err
	}
	if subDoc == nil {
		return nil, nil, fmt.Errorf("resource %q not found", resourceID)
	}

	return doc.OpenShiftCluster, subDoc.Subscription, nil
}

func run(ctx context.Context, log *logrus.Entry) error {
	if len(os.Args) != 3 {
		usage()
		os.Exit(2)
	}

	env, err := env.NewEnv(ctx, log)
	if err != nil {
		return err
	}

	oc, sub, err := getCluster(ctx, log, env, os.Args[2])
	if err != nil {
		return err
	}

	s, err := newSSHTool(log, env, oc, sub)
	if err != nil {
		return err
	}

	switch strings.ToLower(os.Args[1]) {
	case "disable":
		return s.disable(ctx)
	case "enable":
		return s.enable(ctx)
	case "shell":
		return s.shell(ctx)
	default:
		usage()
		os.Exit(2)
	}

	return nil
}

func main() {
	ctx := context.Background()
	log := utillog.GetLogger()

	if err := run(ctx, log); err != nil {
		log.Fatal(err)
	}
}
