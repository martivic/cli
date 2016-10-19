package ccv3

import (
	"bytes"
	"net/url"
	"strings"
	"time"

	"code.cloudfoundry.org/cli/api/cloudcontroller"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3/internal"
)

/*
{
  "pagination": {
    "total_results": 3,
    "total_pages": 2,
    "first": {
      "href": "https://api.example.org/v3/tasks?page=1&per_page=2"
    },
    "last": {
      "href": "https://api.example.org/v3/tasks?page=2&per_page=2"
    },
    "next": {
      "href": "https://api.example.org/v3/tasks?page=2&per_page=2"
    },
    "previous": null
  },
  "resources": [
  ]
}
*/

type Task struct {
	CreatedAt   time.Time `json:"created_at"`
	DropletGUID string    `json:"droplet_guid"`
	GUID        string    `json:"guid"`
	Links       struct {
		App     Link `json:"app"`
		Droplet Link `json:"app"`
	} `json:"links"`
	Memory int    `json:"memory_in_mb"`
	Name   string `json:"name"`
	Result struct {
		Reason string `json:"failure_reason"`
	} `json:"result"`
	SequenceID int        `json:"sequence_id"`
	State      string     `json:"state"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

func FormatQueryParameters(queries url.Values, queryName string, value ...string) url.Values {
	if queries == nil {
		queries = url.Values{}
	}
	queries.Add(queryName, strings.Join(value, ","))

	return queries
}

func (client *Client) GetTasks(queries map[string]string) ([]Task, error) {
	parameters := url.Values{}
	for queryName, queryValue := range queries {
		parameters.Add(queryName, queryValue)
	}

	request := cloudcontroller.Request{
		RequestName: internal.TasksRequest,
		Query:       parameters,
	}

	fullTasksList := []Task{}

	for {
		var tasks []Task
		wrapper := PaginatedWrapper{
			Resources: &tasks,
		}
		response := cloudcontroller.Response{
			Result: &wrapper,
		}

		err := client.connection.Make(request, &response)
		if err != nil {
			return nil, err
		}
		fullTasksList = append(fullTasksList, tasks...)

		if wrapper.Pagination.Next.URL == "" {
			break
		}
		request = cloudcontroller.Request{
			URL:    wrapper.Pagination.Next.URL,
			Method: "GET",
		}
	}

	return fullTasksList, nil
}

func (client *Client) GetApplicationTasks(appGUID string, queries map[string]string) ([]Task, error) {
	parameters := url.Values{}
	for queryName, queryValue := range queries {
		parameters.Add(queryName, queryValue)
	}

	request := cloudcontroller.Request{
		RequestName: internal.AppTaskRequest,
		Query:       parameters,
		Params: map[string]string{
			"app_guid": appGUID,
		},
	}

	fullTasksList := []Task{}

	for {
		var tasks []Task
		wrapper := PaginatedWrapper{
			Resources: &tasks,
		}
		response := cloudcontroller.Response{
			Result: &wrapper,
		}

		err := client.connection.Make(request, &response)
		if err != nil {
			return nil, err
		}
		fullTasksList = append(fullTasksList, tasks...)

		if wrapper.Pagination.Next.URL == "" {
			break
		}
		request = cloudcontroller.Request{
			URL:    wrapper.Pagination.Next.URL,
			Method: "GET",
		}
	}

	return fullTasksList, nil
}

func (client *Client) RunTaskByApplication(appGUID string, command string) (Task, error) {
	request := cloudcontroller.Request{
		RequestName: internal.RunTaskRequest,
		Params: map[string]string{
			"app_guid": appGUID,
		},
		Body: bytes.NewBufferString(command),
	}

	var newTask Task
	response := cloudcontroller.Response{
		Result: &newTask,
	}

	err := client.connection.Make(request, &response)
	if err != nil {
		return Task{}, err
	}

	return newTask, nil
}

func (client *Client) TerminateTask(taskGUID string) (Task, error) {
	request := cloudcontroller.Request{
		RequestName: internal.TerminateTaskRequest,
		Params: map[string]string{
			"task_guid": taskGUID,
		},
	}

	var task Task
	response := cloudcontroller.Response{
		Result: &task,
	}

	err := client.connection.Make(request, &response)
	if err != nil {
		return Task{}, err
	}

	return task, nil
}
