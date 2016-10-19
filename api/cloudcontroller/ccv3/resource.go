package ccv3

import (
	"code.cloudfoundry.org/cli/api/cloudcontroller"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3/internal"
)

type Link struct {
	URL string `json:"href"`
}

type APIs struct {
	Links struct {
		CCV3Endpoint Link `json:"cloud_controller_v3"`
	} `json:"links"`
}

type CCV3Endpoints struct {
	Links map[string]Link `json:"links"`
}

// Resources returns back endpoint and API information from /v3.
func (client *Client) Resources() (CCV3Endpoints, Warnings, error) {
	request := cloudcontroller.Request{
		RequestName: internal.APILinksRequest,
	}

	var apis APIs
	response := cloudcontroller.Response{
		Result: &apis,
	}

	err := client.connection.Make(request, &response)
	if err != nil {
		return CCV3Endpoints{}, response.Warnings, err
	}

	request = cloudcontroller.Request{
		URL:    apis.Links.CCV3Endpoint.URL,
		Method: "GET",
	}

	var ccv3Endpoints CCV3Endpoints
	response = cloudcontroller.Response{
		Result: &ccv3Endpoints,
	}

	err = client.connection.Make(request, &response)
	if err != nil {
		return CCV3Endpoints{}, response.Warnings, err
	}

	return ccv3Endpoints, response.Warnings, nil
}
