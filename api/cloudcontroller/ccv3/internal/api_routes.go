package internal

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	APILinksRequest      = "API Links"
	TasksRequest         = "Tasks"
	RunTaskRequest       = "Run Task"
	AppTaskRequest       = "App's Task"
	TerminateTaskRequest = "Cancel a Task"
	ApplicationsRequest  = "List applications"
)

var routes map[string]Route = map[string]Route{
	APILinksRequest:      Route{Path: "", Method: "GET", Resource: RootResource},
	TasksRequest:         Route{Path: "", Method: "GET", Resource: TasksResource},
	AppTaskRequest:       Route{Path: "/v3/apps/:app_guid/tasks", Method: "GET", Resource: RootResource},
	RunTaskRequest:       Route{Path: "/v3/apps/:app_guid/tasks", Method: "POST", Resource: RootResource},
	TerminateTaskRequest: Route{Path: "/v3/tasks/:task_guid/cancel", Method: "PUT", Resource: RootResource},
	ApplicationsRequest:  Route{Path: "/v3/apps", Method: "GET", Resource: RootResource},
}

const (
	RootResource  = "root"
	TasksResource = "tasks"
)

//-----------------------------------------------------------------------------------------------------------------------

type Route struct {
	Resource string
	Path     string
	Method   string //ENUM?
}

func (route *Route) CreatePath(params map[string]string) (string, error) {
	components := strings.Split(route.Path, "/")
	for i, c := range components {
		if len(c) == 0 {
			continue
		}
		if c[0] == ':' {
			val, ok := params[c[1:]]
			if !ok {
				return "", fmt.Errorf("missing param %s", c)
			}
			components[i] = val
		}
	}

	u, err := url.Parse(strings.Join(components, "/"))
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

type Router struct {
	routes    map[string]Route
	resources map[string]string
}

func NewRouter(APIURL string) *Router {
	return &Router{
		resources: map[string]string{
			"root": APIURL,
		},
		routes: routes,
	}
}

func (router *Router) CreateRequest(requestName string, params map[string]string, body io.Reader) (*http.Request, error) {
	route, exist := router.routes[requestName]
	if !exist {
		return &http.Request{}, fmt.Errorf("No route exists with the name %s", requestName)
	}
	path, err := route.CreatePath(params)
	if err != nil {
		return &http.Request{}, err
	}

	url := router.resources[route.Resource] + "/" + strings.TrimLeft(path, "/")

	req, err := http.NewRequest(route.Method, url, body)
	if err != nil {
		return &http.Request{}, err
	}

	return req, nil
}

func (router *Router) UpdateResources(resources map[string]string) {
	for resourceName, resourceURL := range resources {
		router.resources[resourceName] = resourceURL
	}
}
