// Package ccv3 is experimental, do not bounce!
package ccv3

import "code.cloudfoundry.org/cli/api/cloudcontroller"

type Warnings []string

type Client struct {
	cloudControllerURL string
	tasksEndpoint      string
	endpoints          map[string]string

	connection cloudcontroller.Connection
}

// NewClient returns a new CloudControllerClient.
func NewClient() *Client {
	return new(Client)
}
