package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/InVisionApp/rye"
	"github.com/gorilla/mux"
	hc "github.com/talpert/hellofour/dal/heroku/client"
	ht "github.com/talpert/hellofour/dal/heroku/types"
	rt "github.com/talpert/hellofour/manager/resource/types"
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

type ProvisionResponse struct {
	ID string `json:"id"`
	// optional
	Message     string `json:"message"`
	LogDrainURL string `json:"log_drain_url"`
}

func (a *API) createHandler(rw http.ResponseWriter, r *http.Request) *rye.Response {
	req := &rt.ProvisionRequest{}

	if err := decodeJSONInput(r.Body, req, log); err != nil {
		return err
	}

	if err := req.Validate(); err != nil {
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

	log.Info(msg)
	// printing this for now until we have a queryable DB
	log.Debugf("Auth: %+v", req.OAuthGrant)

	// trying synchronously
	if err := a.Provision(r.Context(), req); err != nil {
		emsg := fmt.Errorf("failed to provision: %v", err)
		log.Error(emsg)
		return &rye.Response{
			Err:        emsg,
			StatusCode: http.StatusInternalServerError,
		}
	}

	// created
	respondAsJSON(rw, http.StatusCreated, resp, log)

	//// respond as accepted for processing
	//respondAsJSON(rw, http.StatusAccepted, resp, log)
	//// now we are async
	//
	//// provision!
	//go func() {
	//	if err := a.Provision(r.Context(), req); err != nil {
	//		log.Errorf("failed to provision: %v", err)
	//	}
	//}()

	return nil
}

func (a *API) Provision(ctx context.Context, request *rt.ProvisionRequest) error {
	// do the auth first so not to provision junk
	auth, err := a.Deps.HerokuClient.GetAuth(ctx, request.OAuthGrant)
	if err != nil {
		return fmt.Errorf("failed to authenticate: %v", err)
	}

	//TODO: put the real provisioning here

	log.Infof("Created a new addon {%s}", request.UUID)

	//set a config val
	a.Deps.HerokuClient.UpdateConfig(ctx, auth, request.CallbackURL, []*hc.AddonConfig{
		{Name: "SPECIAL_VAR", Value: "super value"},
	})

	// if success...
	return a.Finished(ctx, request.CallbackURL, auth)
}

func (a *API) Finished(ctx context.Context, url string, auth *ht.Auth) error {
	//call api to report done
	resp, err := a.Deps.HerokuClient.CallDone(ctx, url, auth)
	if err != nil {
		return fmt.Errorf("failed to authenticate: %v", err)
	}

	log.Info("finished provisioning for %s: %+v", resp.ID, resp)

	return nil
}

func (a *API) updateHandler(rw http.ResponseWriter, r *http.Request) *rye.Response {
	accountID := mux.Vars(r)["accountID"]

	reqBody := &rt.ProvisionRequest{}

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
