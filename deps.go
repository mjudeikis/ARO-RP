// +build tools

// tools is a dummy package that will be ignored for builds, but included for dependencies
package main

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	_ "github.com/alvaroloes/enumer"                      // used to in enum type generation
	_ "github.com/go-bindata/go-bindata/go-bindata"       // used to in static content generation
	_ "github.com/golang/mock/mockgen"                    // used to in tests
	_ "github.com/jim-minter/go-cosmosdb/cmd/gencosmosdb" // used to in database client generation
	_ "github.com/onsi/ginkgo"                            // used to in tests
	_ "github.com/onsi/gomega"                            // used to in tests
	_ "golang.org/x/tools/cmd/goimports"                  // used to in verify tests
	_ "k8s.io/code-generator/cmd/client-gen"              // used to in operator code generation
	_ "sigs.k8s.io/controller-tools/cmd/controller-gen"   // used to in operator  code generation
)
