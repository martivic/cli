package v3

import (
	"fmt"

	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv2"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	"code.cloudfoundry.org/cli/api/cloudcontroller/wrapper"
	"code.cloudfoundry.org/cli/api/uaa"
	"code.cloudfoundry.org/cli/commands"
	"code.cloudfoundry.org/cli/commands/flags"
	"code.cloudfoundry.org/cli/commands/v2/common"
)

type RunTaskCommand struct {
	RequiredArgs flags.RunTaskArgs `positional-args:"yes"`
	usage        interface{}       `usage:"CF_NAME run-task APP_NAME COMMAND"`

	UI     commands.UI
	Config commands.Config
}

func (cmd *RunTaskCommand) Setup(config commands.Config, ui commands.UI) error {
	cmd.UI = ui
	cmd.Config = config
	return nil
}

func (cmd *RunTaskCommand) Execute(args []string) error {
	user, err := cmd.Config.CurrentUser()
	if err != nil {
		return err
	}
	cmd.UI.DisplayText("Creating task for app {{.App}} in org {{.Org}} / space {{.Space}} as {{.User}}...", map[string]interface{}{
		"App":   cmd.RequiredArgs.AppName,
		"Org":   cmd.Config.TargetedOrganization().Name,
		"Space": cmd.Config.TargetedSpace().Name,
		"User":  user.Name,
	})

	client := ccv3.NewClient()
	_, err = client.TargetCF(cmd.Config.Target(), true)
	if err != nil {
		return err
	}

	v2client, err := common.NewCloudControllerClient(cmd.Config)
	if err != nil {
		return err
	}
	uaaClient := uaa.NewClient(v2client.AuthorizationEndpoint(), cmd.Config)
	client.WrapConnection(wrapper.NewUAAAuthentication(uaaClient))

	apps, _, err := v2client.GetApplications([]ccv2.Query{{
		Filter:   ccv2.NameFilter,
		Operator: ccv2.EqualOperator,
		Value:    cmd.RequiredArgs.AppName,
	}})
	if err != nil {
		return err
	}

	if len(apps) == 0 {
		return nil
	}

	task, err := client.RunTaskByApplication(apps[0].GUID, fmt.Sprintf("{\"command\":\"%s\"}", cmd.RequiredArgs.Command))
	if err != nil {
		return err
	}

	cmd.UI.DisplayText("Created task: {{.Name}}\t\tSequence: {{.SequenceID}}\t\tState: {{.State}}", map[string]interface{}{
		"Name":       task.Name,
		"SequenceID": task.SequenceID,
		"State":      task.State,
	})

	return nil
}
