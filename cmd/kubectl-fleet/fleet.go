package main

import (
	"github.com/spf13/cobra"
)

var (
	resourceID        string
	subscriptionID    string
	resourceGroupName string
	fleetName         string
	memberName        string
)

var cmdFleet = &cobra.Command{
	Use:           "kubectl-fleet",
	Short:         "fleet",
	SilenceErrors: true,
	SilenceUsage:  true,
}

func init() {
	cmdFleet.PersistentFlags().StringVar(&resourceID, "resource-id", "", "resource ID")
	cmdFleet.PersistentFlags().StringVar(&subscriptionID, "subscription", "", "subscription ID")
	cmdFleet.PersistentFlags().StringVar(&resourceGroupName, "resource-group", "", "resource group")
	cmdFleet.PersistentFlags().StringVar(&fleetName, "fleet-name", "", "fleet name")
}
