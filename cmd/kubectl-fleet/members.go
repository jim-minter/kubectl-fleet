package main

import (
	"fmt"
	"sort"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservicefleet/armcontainerservicefleet"
	"github.com/spf13/cobra"
)

var cmdMembers = &cobra.Command{
	Use:   "members",
	Short: "members",
	RunE:  members,
}

func init() {
	cmdFleet.AddCommand(cmdMembers)
}

func members(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	err := validate(cmd.Flags())
	if err != nil {
		return err
	}

	rid, err := arm.ParseResourceID(resourceID)
	if err != nil {
		return err
	}

	if rid.ResourceType.String() == "microsoft.containerservice/fleets/members" {
		rid = rid.Parent
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return err
	}

	membersCli, err := armcontainerservicefleet.NewFleetMembersClient(rid.SubscriptionID, cred, nil)
	if err != nil {
		return err
	}

	var memberNames []string

	pager := membersCli.NewListByFleetPager(rid.ResourceGroupName, rid.Name, nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return err
		}

		for _, member := range page.Value {
			memberNames = append(memberNames, *member.Name)
		}
	}

	sort.Strings(memberNames)

	fmt.Println("NAME")
	for _, memberName := range memberNames {
		fmt.Println(memberName)
	}

	return nil
}
