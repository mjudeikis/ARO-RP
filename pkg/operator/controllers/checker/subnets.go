package checker

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"fmt"

	maoclient "github.com/openshift/machine-api-operator/pkg/generated/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	azureproviderv1beta1 "sigs.k8s.io/cluster-api-provider-azure/pkg/apis/azureprovider/v1beta1"
)

func getSubnetIDs(ctx context.Context, vnetID string, clustercli maoclient.Interface) ([]string, error) {
	// TODO - there exists a possibility that the subnets discovered here are not part of the same vnet or even in the same subscription
	// it is an unlikely circumstance but still possible so this logic should account for it.
	subnetNames := make(map[string]struct{})
	subnetIDs := []string{}

	machines, err := clustercli.MachineV1beta1().Machines(machineSetsNamespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return []string{}, err
	}

	for _, machine := range machines.Items {
		if machine.Spec.ProviderSpec.Value == nil {
			return []string{}, fmt.Errorf("machine %s: provider spec missing", machine.Name)
		}

		o, _, err := scheme.Codecs.UniversalDeserializer().Decode(machine.Spec.ProviderSpec.Value.Raw, nil, nil)
		if err != nil {
			return []string{}, err
		}

		machineProviderSpec, ok := o.(*azureproviderv1beta1.AzureMachineProviderSpec)
		if !ok {
			// This should never happen: codecs uses scheme that has only one registered type
			// and if something is wrong with the provider spec - decoding should fail
			return []string{}, fmt.Errorf("machine %s: failed to read provider spec: %T", machine.Name, o)
		}

		subnetNames[machineProviderSpec.Subnet] = struct{}{}
	}

	for k := range subnetNames {
		subnetIDs = append(subnetIDs, vnetID+"/subnets/"+k)
	}

	return subnetIDs, nil
}
