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

type TasksCommand struct {
	RequiredArgs flags.AppName `positional-args:"yes"`
	usage        interface{}   `usage:"CF_NAME tasks"`

	UI     commands.UI
	Config commands.Config
}

func (cmd *TasksCommand) Setup(config commands.Config, ui commands.UI) error {
	cmd.UI = ui
	cmd.Config = config
	return nil
}

func (cmd *TasksCommand) Execute(args []string) error {
	cmd.UI.DisplayText("Getting tasks for application {{.AppName}}", map[string]interface{}{
		"AppName": cmd.RequiredArgs.AppName,
	})

	v3client := ccv3.NewClient()
	_, err := v3client.TargetCF(cmd.Config.Target(), true)
	if err != nil {
		return err
	}

	v2client, err := common.NewCloudControllerClient(cmd.Config)
	if err != nil {
		return err
	}
	uaaClient := uaa.NewClient(v2client.AuthorizationEndpoint(), cmd.Config)
	v3client.WrapConnection(wrapper.NewUAAAuthentication(uaaClient))

	queries := url.Values{}
	queries.Add("space_guids", cmd.Config.TargetedSpace().GUID)
	queries.Add("names", cmd.RequiredArgs.AppName)
	apps, err := v3client.GetApplications(queries)
	if err != nil {
		return err
	}

	if len(apps) == 0 {
		fmt.Println("application not found")
		return nil
	}

	queries = url.Values{}
	queries.Add("app_guids", apps[0].GUID)
	tasks, err := v3client.GetTasks(queries)
	if err != nil {
		return err
	}

	for _, task := range tasks {
		cmd.UI.DisplayText("Got task: {{.Name}}\t\tSequence: {{.SequenceID}}\t\tState: {{.State}}", map[string]interface{}{
			"Name":       task.Name,
			"SequenceID": task.SequenceID,
			"State":      task.State,
		})
	}

	return nil
}
