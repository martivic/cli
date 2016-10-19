package v3

import (
	"fmt"
	"net/url"

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

	v3client := ccv3.NewClient()
	_, err = v3client.TargetCF(cmd.Config.Target(), true)
	if err != nil {
		return err
	}

	v2client, err := common.NewCloudControllerClient(cmd.Config)
	if err != nil {
		return err
	}
	uaaClient := uaa.NewClient(v2client.AuthorizationEndpoint(), cmd.Config)
	v3client.WrapConnection(wrapper.NewUAAAuthentication(uaaClient))

	queries := url.Values{
		"space_guids": []string{cmd.Config.TargetedSpace().GUID},
		"names":       []string{cmd.RequiredArgs.AppName},
	}
	apps, err := v3client.GetApplications(queries)
	if err != nil {
		return err
	}

	if len(apps) == 0 {
		fmt.Println("no apps found")
		return nil
	}

	task, err := v3client.RunTaskByApplication(apps[0].GUID, fmt.Sprintf("{\"command\":\"%s\"}", cmd.RequiredArgs.Command))
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
