package main

import (
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd/api"
)

var cmdSetContext = &cobra.Command{
	Use:                "set-context",
	Short:              "set context",
	RunE:               setContext,
	SilenceErrors:      true,
	SilenceUsage:       true,
	DisableFlagParsing: true,
}

func init() {
	cmdSetContext.Flags().StringVar(&memberName, "member-name", "", "member name")

	cmdFleet.AddCommand(cmdSetContext)
}

func setContext(cmd *cobra.Command, args []string) error {
	err := cmd.Flags().Parse(args)
	if err != nil {
		return err
	}

	if help {
		return cmd.Help()
	}

	ctx := cmd.Context()

	err = validate(cmd.Flags())
	if err != nil {
		return err
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return err
	}

	err = ensureKubeconfig(ctx, cred, func(kubeconfig *api.Config) {
		kubeconfig.CurrentContext = resourceID
	})
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "context set to %s\n", resourceID)

	return nil
}
