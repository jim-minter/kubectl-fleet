package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	resourceID        string
	subscriptionID    string
	resourceGroupName string
	fleetName         string
	memberName        string
	help              bool
)

var cmdFleet = &cobra.Command{
	Use:                "kubectl-fleet",
	Short:              "fleet",
	RunE:               fleet,
	Args:               cobra.ArbitraryArgs,
	SilenceErrors:      true,
	SilenceUsage:       true,
	DisableFlagParsing: true,
}

func init() {
	cmdFleet.PersistentFlags().BoolVarP(&help, "help", "h", false, "help")
	cmdFleet.PersistentFlags().MarkHidden("help")

	cmdFleet.PersistentFlags().StringVar(&resourceID, "resource-id", "", "resource ID")
	cmdFleet.PersistentFlags().StringVar(&subscriptionID, "subscription", "", "subscription ID")
	cmdFleet.PersistentFlags().StringVar(&resourceGroupName, "resource-group", "", "resource group")
	cmdFleet.PersistentFlags().StringVar(&fleetName, "fleet-name", "", "fleet name")

	cmdFleet.Flags().StringVar(&memberName, "member-name", "", "member name")
}

func fleet(cmd *cobra.Command, args []string) error {
	args, kubectlArgs := splitArgs(cmd.Flags(), args)

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

	// ensure the cluster's co-ordinates are in the kubeconfig
	err = ensureKubeconfig(ctx, cred, nil)
	if err != nil {
		return err
	}

	// now launch kubectl with --context referring to the cluster
	command := exec.Command("kubectl", append([]string{"--context", resourceID}, kubectlArgs...)...)
	command.Stdin, command.Stdout, command.Stderr = os.Stdin, os.Stdout, os.Stderr

	return command.Run()
}

// split out the flags on the command that belong to us from those that don't
func splitArgs(flags *pflag.FlagSet, args []string) (ourArgs, theirArgs []string) {
	for i := 0; i < len(args); i++ {
		var f *pflag.Flag
		if strings.HasPrefix(args[i], "--") {
			f = flags.Lookup(strings.SplitN(args[i][2:], "=", 2)[0])
		} else if strings.HasPrefix(args[i], "-") {
			f = flags.ShorthandLookup(strings.SplitN(args[i][1:], "=", 2)[0])
		}

		if f == nil {
			theirArgs = append(theirArgs, args[i])
			continue
		}

		if strings.Contains(args[i], "=") || f.NoOptDefVal != "" || i == len(args)-1 {
			ourArgs = append(ourArgs, args[i])
		} else {
			ourArgs = append(ourArgs, args[i], args[i+1])
			i++
		}
	}

	return
}
