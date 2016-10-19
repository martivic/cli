package ccv3

import (
	"net/url"
	"time"

	"code.cloudfoundry.org/cli/api/cloudcontroller"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3/internal"
)

type Application struct {
	CreatedAt            time.Time `json:"created_at"`
	DesiredState         string    `json:"desired_state"`
	EnvironmentVariables struct {
		RedactedMessage string `json:"redacted_message"`
	} `json:"environment_variables"`
	GUID      string `json:"guid"`
	Lifecycle struct {
		Type string `json:"type"`
		Data struct {
			Buildpack string `json:"buildpack"`
			Stack     string `json:"stack"`
		} `json:"data"`
	} `json:"lifecycle"`
	Links struct {
		Space         Link `json:"space"`
		Processes     Link `json:"processes"`
		RouteMappings Link `json:"route_mappings"`
		Packages      Link `json:"packages"`
		Droplet       Link `json:"droplet"`
		Droplets      Link `json:"droplets"`
		Tasks         Link `json:"tasks"`
		Start         Link `json:"start"`
		Stop          Link `json:"stop"`
	} `json:"links"`
	Name                  string     `json:"name"`
	TotalDesiredInstances int        `json:"total_desired_instances"`
	UpdatedAt             *time.Time `json:"updated_at"`
}

func (client *Client) GetApplications(queries url.Values) ([]Application, error) {
	request := cloudcontroller.Request{
		RequestName: internal.ApplicationsRequest,
		Query:       queries,
	}

	fullApplicationList := []Application{}

	for {
		var applications []Application
		wrapper := PaginatedWrapper{
			Resources: &applications,
		}
		response := cloudcontroller.Response{
			Result: &wrapper,
		}

		err := client.connection.Make(request, &response)
		if err != nil {
			return nil, err
		}
		fullApplicationList = append(fullApplicationList, applications...)

		if wrapper.Pagination.Next.URL == "" {
			break
		}
		request = cloudcontroller.Request{
			URL:    wrapper.Pagination.Next.URL,
			Method: "GET",
		}
	}

	return fullApplicationList, nil
}
