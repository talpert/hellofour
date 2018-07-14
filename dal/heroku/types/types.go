package types

import "time"

type Auth struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type OAuthGrant struct {
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expires_at"`
	Type      string    `json:"type"`
}

//TODO: split this out
type ProvisionResponse struct {
	Actions struct {
		ID            string `json:"id"`
		Label         string `json:"label"`
		Action        string `json:"action"`
		URL           string `json:"url"`
		RequiresOwner bool   `json:"requires_owner"`
	} `json:"actions"`
	AddonService struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"addon_service"`
	App struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"app"`
	BilledPrice struct {
		Cents    int    `json:"cents"`
		Contract bool   `json:"contract"`
		Unit     string `json:"unit"`
	} `json:"billed_price"`
	ConfigVars []string  `json:"config_vars"`
	CreatedAt  time.Time `json:"created_at"`
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Plan       struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"plan"`
	ProviderID string    `json:"provider_id"`
	State      string    `json:"state"`
	UpdatedAt  time.Time `json:"updated_at"`
	WebURL     string    `json:"web_url"`
}
