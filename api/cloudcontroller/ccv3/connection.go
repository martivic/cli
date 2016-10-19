package ccv3

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"code.cloudfoundry.org/cli/api/cloudcontroller"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3/internal"
)

type Connection struct {
	HTTPClient *http.Client
	router     *internal.Router
}

func NewConnection(router *internal.Router, skipSSLValidation bool) *Connection {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: skipSSLValidation,
		},
	}

	return &Connection{
		HTTPClient: &http.Client{Transport: tr},
		router:     router,
	}
}

func (connection *Connection) Make(passedRequest cloudcontroller.Request, passedResponse *cloudcontroller.Response) error {
	req, err := connection.createHTTPRequest(passedRequest)
	if err != nil {
		return err
	}

	response, err := connection.HTTPClient.Do(req)
	if err != nil {
		fmt.Println("error from CC", err)
		return connection.processRequestErrors(err)
	}

	defer response.Body.Close()

	return connection.populateResponse(response, passedResponse)
}

func (connection *Connection) createHTTPRequest(passedRequest cloudcontroller.Request) (*http.Request, error) {
	var request *http.Request
	var err error

	body := &bytes.Buffer{}
	if passedRequest.Body != nil {
		body = passedRequest.Body
	}

	if passedRequest.URL != "" {
		request, err = http.NewRequest(
			passedRequest.Method,
			passedRequest.URL,
			body,
		)
	} else {
		request, err = connection.router.CreateRequest(
			passedRequest.RequestName,
			passedRequest.Params,
			body,
		)
		if err == nil {
			request.URL.RawQuery = passedRequest.Query.Encode()
		}
	}
	if err != nil {
		return nil, err
	}

	if passedRequest.Header != nil {
		request.Header = passedRequest.Header
	}

	request.Header.Set("accept", "application/json")
	request.Header.Set("content-type", "application/json")

	// request.Header.Set("Connection", "close")
	// request.Header.Set("User-Agent", "go-cli "+cf.Version+" / "+runtime.GOOS)

	return request, nil
}

func (connection *Connection) processRequestErrors(err error) error {
	switch e := err.(type) {
	case *url.Error:
		if _, ok := e.Err.(x509.UnknownAuthorityError); ok {
			return cloudcontroller.UnverifiedServerError{
			// URL: connection.URL, //TODO: Figure this one out
			}
		}
		return cloudcontroller.RequestError{Err: e}
	default:
		return err
	}
}

func (connection *Connection) populateResponse(response *http.Response, passedResponse *cloudcontroller.Response) error {
	if rawWarnings := response.Header.Get("X-Cf-Warnings"); rawWarnings != "" {
		passedResponse.Warnings = []string{}
		for _, warning := range strings.Split(rawWarnings, ",") {
			warningTrimmed := strings.Trim(warning, " ")
			passedResponse.Warnings = append(passedResponse.Warnings, warningTrimmed)
		}
	}

	err := connection.handleStatusCodes(response)
	if err != nil {
		return err
	}

	if passedResponse.Result != nil {
		rawBytes, _ := ioutil.ReadAll(response.Body)
		passedResponse.RawResponse = rawBytes

		decoder := json.NewDecoder(bytes.NewBuffer(rawBytes))
		decoder.UseNumber()
		err = decoder.Decode(passedResponse.Result)
		if err != nil {
			return err
		}
	}

	return nil
}

func (*Connection) handleStatusCodes(response *http.Response) error {
	if response.StatusCode >= 400 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}

		return cloudcontroller.RawCCError{
			StatusCode:  response.StatusCode,
			RawResponse: body,
		}
	}

	return nil
}
