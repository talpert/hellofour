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
