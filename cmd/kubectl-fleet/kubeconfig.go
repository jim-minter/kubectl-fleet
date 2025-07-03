package main

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservicefleet/armcontainerservicefleet"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func getKubeconfigCurrentContext() (string, error) {
	kubeconfig, err := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	if err != nil {
		return "", err
	}

	return kubeconfig.CurrentContext, nil
}

func ensureKubeconfig(ctx context.Context, cred azcore.TokenCredential, mutate func(*api.Config)) error {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()

	kubeconfig, err := rules.Load()
	if err != nil {
		return err
	}

	if _, ok := kubeconfig.Clusters[resourceID]; !ok {
		k, err := getKubeconfig(ctx, cred, resourceID)
		if err != nil {
			return err
		}

		kubeconfig.Clusters[resourceID] = k.Clusters[k.Contexts[k.CurrentContext].Cluster]
		kubeconfig.AuthInfos[resourceID] = k.AuthInfos[k.Contexts[k.CurrentContext].AuthInfo]
		kubeconfig.Contexts[resourceID] = &api.Context{
			Cluster:  resourceID,
			AuthInfo: resourceID,
		}

		err = clientcmd.ModifyConfig(rules, *kubeconfig, false)
		if err != nil {
			return err
		}
	}

	if mutate != nil {
		mutate(kubeconfig)

		err = clientcmd.ModifyConfig(rules, *kubeconfig, false)
		if err != nil {
			return err
		}
	}

	return nil
}

func getKubeconfig(ctx context.Context, cred azcore.TokenCredential, resourceID string) (*api.Config, error) {
	rid, err := arm.ParseResourceID(resourceID)
	if err != nil {
		return nil, err
	}

	switch rid.ResourceType.String() {
	case "microsoft.containerservice/fleets":
		cli, err := armcontainerservicefleet.NewFleetsClient(rid.SubscriptionID, cred, nil)
		if err != nil {
			return nil, err
		}

		creds, err := cli.ListCredentials(ctx, rid.ResourceGroupName, rid.Name, nil)
		if err != nil {
			return nil, err
		}

		return clientcmd.Load(creds.Kubeconfigs[0].Value)

	case "microsoft.containerservice/fleets/members":
		membersCli, err := armcontainerservicefleet.NewFleetMembersClient(rid.SubscriptionID, cred, nil)
		if err != nil {
			return nil, err
		}

		member, err := membersCli.Get(ctx, rid.ResourceGroupName, rid.Parent.Name, rid.Name, nil)
		if err != nil {
			return nil, err
		}

		rid, err = arm.ParseResourceID(*member.Properties.ClusterResourceID)
		if err != nil {
			return nil, err
		}

		cli, err := armcontainerservice.NewManagedClustersClient(rid.SubscriptionID, cred, nil)
		if err != nil {
			return nil, err
		}

		creds, err := cli.ListClusterUserCredentials(ctx, rid.ResourceGroupName, rid.Name, nil)
		if err != nil {
			return nil, err
		}

		return clientcmd.Load(creds.Kubeconfigs[0].Value)
	}

	return nil, fmt.Errorf("unimplemented resource type %s", rid.ResourceType)
}
