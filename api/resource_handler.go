package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/InVisionApp/rye"
	"github.com/talpert/hellofour/util"
)

//{
//	"callback_url": "https://api.heroku.com/addons/01234567-89ab-cdef-0123-456789abcdef",
//	"name": "acme-inc-primary-database",
//	"oauth_grant": {
//		"code": "01234567-89ab-cdef-0123-456789abcdef",
//		"expires_at": "2016-03-03T18:01:31-0800",
//		"type": "authorization_code"
//},
//	"options": { "foo" : "bar", "baz" : "true" },
//	"plan": "basic",
//	"region": "amazon-web-services::us-east-1",
//	"uuid": "01234567-89ab-cdef-0123-456789abcdef",
//	"log_input_url": "https://token:t.01234567-89ab-cdef-0123-456789abcdef@1.us.logplex.io/logs",
//	"log_drain_token": "d.01234567-89ab-cdef-0123-456789abcdef"
//}

type ProvisionRequest struct {
	CallbackURL string                 `json:"callback_url"`
	Name        string                 `json:"name"`
	OAuthGrant  *OAuthGrant            `json:"oauth_grant"`
	Options     map[string]interface{} `json:"options"`
	Plan        string                 `json:"plan"`
	Region      string                 `json:"region"`
	UUID        string                 `json:"uuid"`
	// optional
	LogInputURL   string `json:"log_input_url"`
	LogDrainToken string `json:"log_drain_token"`
}

type OAuthGrant struct {
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expires_at"`
	Type      string    `json:"type"`
}

type ProvisionResponse struct {
	ID string `json:"id"`
	// optional
	Message     string `json:"message"`
	LogDrainURL string `json:"log_drain_url"`
}

func (a *API) resourceHandler(rw http.ResponseWriter, r *http.Request) *rye.Response {
	reqBody := &ProvisionRequest{}

	if err := decodeJSONInput(r.Body, reqBody, log); err != nil {
		return err
	}

	id, err := Provision(r.Context(), reqBody)
	if err != nil {
		return &rye.Response{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	msg := fmt.Sprintf("Created a new addon {%s} with options: %v", id, reqBody.Options)
	log.Info(msg)

	resp := &ProvisionResponse{
		ID:      id,
		Message: msg,
	}

	respondAsJSON(rw, http.StatusAccepted, resp, log)
	return nil
}

func Provision(ctx context.Context, request *ProvisionRequest) (string, error) {
	//TODO: put the real provisioning here

	return util.GenerateUUID().String(), nil
}
