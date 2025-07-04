package main

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd/api"
)

var cmdSetContext = &cobra.Command{
	Use:   "set-context",
	Short: "Set context",
	RunE:  setContext,
}

func init() {
	cmdSetContext.Flags().StringVar(&memberName, "member-name", "", "member name")

	cmdFleet.AddCommand(cmdSetContext)
}

func setContext(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	err := validate(cmd.Flags())
	if err != nil {
		return err
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return err
	}

	return ensureKubeconfig(ctx, cred, func(kubeconfig *api.Config) {
		kubeconfig.CurrentContext = resourceID
	})
}
