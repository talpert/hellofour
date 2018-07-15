package types

import (
	"fmt"

	ht "github.com/talpert/hellofour/dal/heroku/types"
)

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

func (p *ProvisionRequest) Validate() error {
	//TODO: flesh out

	if p.OAuthGrant == nil {
		return fmt.Errorf("must have an oauth grant")
	}

	return nil
}
