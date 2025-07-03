package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/spf13/pflag"
	"k8s.io/client-go/util/homedir"
)

func validate(flags *pflag.FlagSet) error {
	// either the user should have specified --resource-id...
	if resourceID != "" {
		rid, err := arm.ParseResourceID(resourceID)
		if err != nil {
			return fmt.Errorf("invalid --resource-id: %w", err)
		}
		if subscriptionID != "" {
			return fmt.Errorf("--resource-id and --subscription are mutually exclusive")
		}
		if resourceGroupName != "" {
			return fmt.Errorf("--resource-id and --resource-group are mutually exclusive")
		}
		if fleetName != "" {
			return fmt.Errorf("--resource-id and --fleet-name are mutually exclusive")
		}

		switch strings.ToLower(rid.ResourceType.String()) {
		case "microsoft.containerservice/fleets":
			if memberName != "" {
				resourceID += "/members/" + memberName
				memberName = ""
			}

			_, err = arm.ParseResourceID(resourceID)
			if err != nil {
				return fmt.Errorf("invalid --member-name: %w", err)
			}

		case "microsoft.containerservice/fleets/members":
			if memberName != "" {
				return fmt.Errorf("can't specify member via --resource-id and --member-name")
			}

		default:
			return fmt.Errorf("invalid --resource-id type %q", rid.ResourceType)
		}

		resourceID = strings.ToLower(resourceID)

		return nil
	}

	// ...or they should have specified some combination of --subscription, --resource-group, --fleet-name and --member-name
	defaults(flags)

	if subscriptionID == "" {
		return fmt.Errorf("--subscription could not be determined and neither --resource-id or --subscription was set")
	}
	if resourceGroupName == "" {
		return fmt.Errorf("either --resource-id or --resource-group must be set")
	}
	if fleetName == "" {
		return fmt.Errorf("either --resource-id or --fleet-name must be set")
	}

	resourceID = "/subscriptions/" + subscriptionID + "/resourceGroups/" + resourceGroupName + "/providers/Microsoft.ContainerService/fleets/" + fleetName

	_, err := arm.ParseResourceID(resourceID)
	if err != nil {
		return fmt.Errorf("invalid --subscription, --resource-group and/or --fleet-name: %w", err)
	}

	if memberName != "" {
		resourceID += "/members/" + memberName
		memberName = ""
	}

	_, err = arm.ParseResourceID(resourceID)
	if err != nil {
		return fmt.Errorf("invalid --member-name: %w", err)
	}

	subscriptionID = ""
	resourceGroupName = ""
	fleetName = ""
	memberName = ""

	resourceID = strings.ToLower(resourceID)

	return nil
}

func defaults(flags *pflag.FlagSet) {
	// if none of the subscription, resource-group or fleet-name have been specified, try to get them from the current context
	if !flags.Lookup("subscription").Changed && !flags.Lookup("resource-group").Changed && !flags.Lookup("fleet-name").Changed {
		maybeResourceID, err := getKubeconfigCurrentContext()
		if err == nil {
			rid, err := arm.ParseResourceID(maybeResourceID)
			if err == nil {
				switch rid.ResourceType.String() {
				case "microsoft.containerservice/fleets":
					subscriptionID = rid.SubscriptionID
					resourceGroupName = rid.ResourceGroupName
					fleetName = rid.Name

				case "microsoft.containerservice/fleets/members":
					subscriptionID = rid.SubscriptionID
					resourceGroupName = rid.ResourceGroupName
					fleetName = rid.Parent.Name

					// if in addition the member-name has not been specified, get that from the current context too
					if flags.Lookup("member-name") != nil && !flags.Lookup("member-name").Changed {
						memberName = rid.Name
					}
				}

				return
			}
		}
	}

	// alternatively, if we have a resource-group and fleet-name but no subscription, try to get that from az
	if !flags.Lookup("subscription").Changed && flags.Lookup("resource-group").Changed && flags.Lookup("fleet-name").Changed {
		b, err := os.ReadFile(filepath.Join(homedir.HomeDir(), ".azure/azureProfile.json"))
		if err != nil {
			return
		}

		b = bytes.TrimPrefix(b, []byte{0xef, 0xbb, 0xbf}) // UTF-8 BOM

		type subscription struct {
			ID        string `json:"id"`
			IsDefault bool   `json:"isDefault"`
		}

		var azureProfile *struct {
			Subscriptions []*subscription `json:"subscriptions"`
		}

		err = json.Unmarshal(b, &azureProfile)
		if err != nil {
			return
		}

		for _, sub := range azureProfile.Subscriptions {
			if sub.IsDefault {
				subscriptionID = sub.ID
				return
			}
		}
	}
}
