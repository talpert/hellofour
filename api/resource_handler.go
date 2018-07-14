package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/InVisionApp/rye"
	"github.com/gorilla/mux"
	hc "github.com/talpert/hellofour/dal/heroku/client"
	ht "github.com/talpert/hellofour/dal/heroku/types"
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
	OAuthGrant  *ht.OAuthGrant         `json:"oauth_grant"`
	Options     map[string]interface{} `json:"options"`
	Plan        string                 `json:"plan"`
	Region      string                 `json:"region"`
	UUID        string                 `json:"uuid"`
	// optional
	LogInputURL   string `json:"log_input_url"`
	LogDrainToken string `json:"log_drain_token"`
}

type ProvisionResponse struct {
	ID string `json:"id"`
	// optional
	Message     string `json:"message"`
	LogDrainURL string `json:"log_drain_url"`
}

func (a *API) createHandler(rw http.ResponseWriter, r *http.Request) *rye.Response {
	req := &ProvisionRequest{}

	if err := decodeJSONInput(r.Body, req, log); err != nil {
		return err
	}

	if err := validateProvision(req); err != nil {
		return &rye.Response{
			Err:        err,
			StatusCode: http.StatusBadRequest,
		}
	}

	msg := fmt.Sprintf("Accepted new addon provision request {%s} with options: %v", req.UUID, req.Options)

	resp := &ProvisionResponse{
		ID:      req.UUID,
		Message: msg,
	}

	// respond as accepted for processing
	respondAsJSON(rw, http.StatusAccepted, resp, log)

	// now we are async
	log.Info(msg)

	// provision!
	go a.Provision(r.Context(), req)

	return nil
}

func validateProvision(req *ProvisionRequest) error {
	//TODO: flesh out

	if req.OAuthGrant == nil {
		return fmt.Errorf("must have an oauth grant")
	}

	return nil
}

func (a *API) Provision(ctx context.Context, request *ProvisionRequest) {
	// do the auth first so not to provision junk
	auth, err := a.Deps.HerokuClient.GetAuth(ctx, request.OAuthGrant)
	if err != nil {
		log.Errorf("failed to authenticate: %v", err)
		return
	}

	//TODO: put the real provisioning here

	log.Infof("Created a new addon {%s}", request.UUID)

	//set a config val
	a.Deps.HerokuClient.UpdateConfig(ctx, auth, request.CallbackURL, []*hc.AddonConfig{
		{Name: "SPECIAL_VAR", Value: "super value"},
	})

	// if success...
	a.Finished(ctx, request.CallbackURL, auth)
}

func (a *API) Finished(ctx context.Context, url string, auth *ht.Auth) {
	//call api to report done
	resp, err := a.Deps.HerokuClient.CallDone(ctx, url, auth)
	if err != nil {
		log.Errorf("failed to authenticate: %v", err)
		return
	}

	log.Info("finished provisioning for %s: %+v", resp.ID, resp)
}

func (a *API) updateHandler(rw http.ResponseWriter, r *http.Request) *rye.Response {
	accountID := mux.Vars(r)["accountID"]

	reqBody := &ProvisionRequest{}

	if err := decodeJSONInput(r.Body, reqBody, log); err != nil {
		return err
	}

	msg := fmt.Sprintf("Updating an addon {%s} with options: %v", accountID, reqBody.Options)
	log.Info(msg)

	resp := &ProvisionResponse{
		Message: msg,
	}

	respondAsJSON(rw, http.StatusOK, resp, log)
	return nil
}

func (a *API) deleteHandler(rw http.ResponseWriter, r *http.Request) *rye.Response {
	accountID := mux.Vars(r)["accountID"]

	msg := fmt.Sprintf("Deleting an addon {%s}", accountID)
	log.Info(msg)

	resp := &ProvisionResponse{
		Message: msg,
	}

	respondAsJSON(rw, http.StatusOK, resp, log)
	return nil
}
