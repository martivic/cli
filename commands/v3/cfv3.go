package v3

import "code.cloudfoundry.org/cli/commands/v2"

var Commands CommandList

type CommandList struct {
	v2.CommandList
	TerminateTask TerminateTaskCommand `command:"terminate-task" description:"Cancel a running task"`
	Tasks         TasksCommand         `command:"tasks" description:"Display a list of tasks"`
	RunTask       RunTaskCommand       `command:"run-task" description:"Run a one-off task on an app"`
}
