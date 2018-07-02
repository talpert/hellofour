package client

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"time"

	"github.com/talpert/hellofour/util/apiclient"
)

const ID_HOST = "https://id.heroku.com"

type Client struct {
	clientSecret string
	*apiclient.APIClient
}

func NewClient(secret string) *Client {
	c := &Client{
		clientSecret: secret,
		APIClient:    apiclient.NewAPIClient(),
	}

	c.Headers = map[string]string{
		"Accept": "application/vnd.heroku+json; version=3",
	}

	return c
}

type Auth struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

//TODO: use an oauth lib for some of this shit
func (c *Client) GetAuth(ctx context.Context, grantToken string) (*Auth, error) {
	val := url.Values{}
	val.Set("grant_type", "authorization_code")
	val.Set("code", grantToken)
	val.Set("client_secret", c.clientSecret)

	auth := &Auth{}

	resp, err := c.MakeRequest(ctx, &apiclient.APIRequest{
		URL:     path.Join(ID_HOST, "/oauth/token"),
		Method:  "POST",
		Payload: val,
		Result:  auth,
	})
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf("failed to get auth: %v", resp.Err)
	}

	return auth, nil
}

//TODO: implement
func (c *Client) RefreshAuth(ctx context.Context) (*Auth, error) {
	return nil, nil
}

type AddonConfig struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (c *Client) UpdateConfig(ctx context.Context, auth *Auth, baseURL string, config []*AddonConfig) error {
	body := struct {
		Config []*AddonConfig `json:"config"`
	}{
		Config: config,
	}

	resp, err := c.MakeRequest(ctx, &apiclient.APIRequest{
		URL:     path.Join(baseURL, "/config"),
		Method:  "PATCH",
		Headers: map[string]string{"Authorization": fmt.Sprintf("Bearer %s", auth.AccessToken)},
		Payload: body,
	})
	if err != nil {
		return err
	}

	if !resp.Success {
		return fmt.Errorf("failed to update config: %v", resp.Err)
	}

	return nil
}

//{
//  "actions": {
//    "id": "01234567-89ab-cdef-0123-456789abcdef",
//    "label": "Example",
//    "action": "example",
//    "url": "http://example.com?resource_id=:resource_id",
//    "requires_owner": true
//  },
//  "addon_service": {
//    "id": "01234567-89ab-cdef-0123-456789abcdef",
//    "name": "heroku-postgresql"
//  },
//  "app": {
//    "id": "01234567-89ab-cdef-0123-456789abcdef",
//    "name": "example"
//  },
//  "billed_price": {
//    "cents": 0,
//    "contract": false,
//    "unit": "month"
//  },
//  "config_vars": [
//    "FOO",
//    "BAZ"
//  ],
//  "created_at": "2012-01-01T12:00:00Z",
//  "id": "01234567-89ab-cdef-0123-456789abcdef",
//  "name": "acme-inc-primary-database",
//  "plan": {
//    "id": "01234567-89ab-cdef-0123-456789abcdef",
//    "name": "heroku-postgresql:dev"
//  },
//  "provider_id": "abcd1234",
//  "state": "provisioned",
//  "updated_at": "2012-01-01T12:00:00Z",
//  "web_url": "https://postgres.heroku.com/databases/01234567-89ab-cdef-0123-456789abcdef"
//}

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

func (c *Client) CallDone(ctx context.Context, baseURL string, auth *Auth) (*ProvisionResponse, error) {
	pr := &ProvisionResponse{}

	resp, err := c.MakeRequest(ctx, &apiclient.APIRequest{
		URL:    path.Join(baseURL, "/actions/provision"),
		Method: "POST",
		Result: pr,
	})
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf("failed to get auth: %v", resp.Err)
	}

	return pr, nil
}
