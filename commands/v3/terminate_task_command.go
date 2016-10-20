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

type TerminateTaskCommand struct {
	RequiredArgs flags.TerminateTaskArgs `positional-args:"yes"`
	usage        interface{}             `usage:"CF_NAME tasks"`

	UI     commands.UI
	Config commands.Config
}

func (cmd *TerminateTaskCommand) Setup(config commands.Config, ui commands.UI) error {
	cmd.UI = ui
	cmd.Config = config
	return nil
}

func (cmd *TerminateTaskCommand) Execute(args []string) error {
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
	queries.Add("sequence_ids", cmd.RequiredArgs.TaskID)
	tasks, err := v3client.GetApplicationTasks(apps[0].GUID, queries)
	if err != nil {
		return err
	}

	user, err := cmd.Config.CurrentUser()
	if err != nil {
		return err
	}

	cmd.UI.DisplayText("Terminating task {{.SequenceID}} of app {{.AppName}} in org {{.Org}} / space {{.Space}} as {{.User}}", map[string]interface{}{
		"SequenceID": tasks[0].SequenceID,
		"AppName":    apps[0].Name,
		"Org":        cmd.Config.TargetedOrganization().Name,
		"Space":      cmd.Config.TargetedSpace().Name,
		"User":       user.Name,
	})

	task, err := v3client.TerminateTask(tasks[0].GUID)
	if err != nil {
		return err
	}

	fmt.Println("Deleted: ", task)

	return nil
}
