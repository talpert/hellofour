package apiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/InVisionApp/go-logger"
	shim "github.com/InVisionApp/go-logger/shims/logrus"
)

type APIClient struct {
	Headers map[string]string

	client *http.Client
	log    log.Logger
}

// Create a new APIClient
// host string should be in the form of `scheme://host:port`
// if a client is not provided, a default one will be used
func NewAPIClient() *APIClient {
	ac := &APIClient{
		client: &http.Client{},
		log:    shim.New(nil),
	}

	ac.client.Timeout = time.Second * 10

	return ac
}

type APIRequest struct {
	URL string

	Method  string
	Headers map[string]string
	Payload interface{}
	Result  interface{}

	//internal
	ctx context.Context
}

type APIResponse struct {
	StatusCode int
	Success    bool // status 200 range
	Body       []byte
	Err        error
	Headers    http.Header
	Cookies    []*http.Cookie
	Response   *http.Response
}

func (a *APIClient) MakeRequest(ctx context.Context, r *APIRequest) (*APIResponse, error) {
	r.ctx = ctx

	resp, err := a.do(r)
	if err != nil {
		a.log.Error(err)
		return nil, err
	}

	apiResp := &APIResponse{
		StatusCode: resp.StatusCode,
		Success:    !(resp.StatusCode > 299 || resp.StatusCode < 200),
		Headers:    resp.Header,
		Cookies:    resp.Cookies(),
		Response:   resp,
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		a.log.Errorf("failed to read response from %s: %v", r.URL, err)
		return nil, err
	}

	apiResp.Body = b

	if !apiResp.Success {
		return a.respondUnsuccessful(r, apiResp)
	}

	if r.Result == nil {
		return apiResp, nil
	}

	return a.unmarshalJSON(r, apiResp)
}

// #######
// Helpers
// #######

// Error condition for when http response is > 299 and < 200
func (a *APIClient) respondUnsuccessful(r *APIRequest, apiResp *APIResponse) (*APIResponse, error) {

	errBody := struct {
		Message string `json:"message"`
	}{}
	if err := json.Unmarshal(apiResp.Body, &errBody); err != nil {
		a.log.Errorf("failed to unmarshal error response from %s: %v", r.URL, err)
		return nil, fmt.Errorf("failed to unmarshal error response: %s", apiResp.Body)
	}

	apiResp.Err = errors.New(errBody.Message)

	a.log.WithFields(log.Fields{
		"message":  errBody.Message,
		"full_url": r.URL,
		"method":   r.Method,
	}).Errorf("api request failure")

	return apiResp, nil
}

// Standard JSON unmarshal
func (a *APIClient) unmarshalJSON(r *APIRequest, apiResp *APIResponse) (*APIResponse, error) {

	if err := json.Unmarshal(apiResp.Body, r.Result); err != nil {
		a.log.Errorf("failed to unmarshal response from %s: %v", r.URL, err)
		return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
	}
	return apiResp, nil
}

func (a *APIClient) do(r *APIRequest) (*http.Response, error) {
	payloadBody, err := parsePayload(r.Payload)
	if err != nil {
		return nil, fmt.Errorf("error parsing payload: %v", err)
	}

	var req *http.Request
	var resp *http.Response

	req, err = http.NewRequest(r.Method, r.URL, payloadBody)
	if err != nil {
		return nil, fmt.Errorf("failed to construct request: %v", err)
	}

	// if content-type is provided, do not inject application/json content-type
	appendHeaderIfNotExist(r, "Content-Type", "application/json")

	for name, value := range r.Headers {
		req.Header.Set(name, value)
	}

	resp, err = a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed request (%s): %v", r.URL, err)
	}

	return resp, nil
}

func appendHeaderIfNotExist(r *APIRequest, name, val string) {
	h := strings.ToLower(name)
	for k, _ := range r.Headers {
		if strings.ToLower(k) == h {
			return
		}
	}

	if r.Headers == nil {
		r.Headers = map[string]string{}
	}

	r.Headers[name] = val

	return
}

func parsePayload(b interface{}) (io.Reader, error) {
	if b == nil {
		return nil, nil
	}

	switch b.(type) {
	case []byte:
		return bytes.NewReader(b.([]byte)), nil
	case string:
		return bytes.NewReader([]byte(b.(string))), nil
	case url.Values:
		return bytes.NewReader([]byte(b.(url.Values).Encode())), nil
	default:
		jb, err := json.Marshal(b)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal json body: %v", err)
		}

		return bytes.NewReader(jb), nil
	}
}
