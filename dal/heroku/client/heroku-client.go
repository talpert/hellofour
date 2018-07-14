package client

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/talpert/hellofour/dal/heroku/types"
	"github.com/talpert/hellofour/util/apiclient"
)

const ID_HOST = "https://id.heroku.com"

type IClient interface {
	GetAuth(ctx context.Context, grant *types.OAuthGrant) (*types.Auth, error)
}

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

//TODO: use an oauth lib for some of this shit
func (c *Client) GetAuth(ctx context.Context, grant *types.OAuthGrant) (*types.Auth, error) {
	if grant == nil {
		return nil, errors.New("OAuth grant may not be nil")
	}

	val := url.Values{}
	val.Set("grant_type", "authorization_code")
	val.Set("code", grant.Code)
	val.Set("client_secret", c.clientSecret)

	auth := &types.Auth{}

	resp, err := c.MakeRequest(ctx, &apiclient.APIRequest{
		URL:         ID_HOST + "/oauth/token",
		Method:      "POST",
		ContentType: "application/x-www-form-urlencoded",
		Payload:     val,
		Result:      auth,
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
func (c *Client) RefreshAuth(ctx context.Context) (*types.Auth, error) {
	return nil, nil
}

type AddonConfig struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (c *Client) UpdateConfig(ctx context.Context, auth *types.Auth, baseURL string, config []*AddonConfig) error {
	body := struct {
		Config []*AddonConfig `json:"config"`
	}{
		Config: config,
	}

	resp, err := c.MakeRequest(ctx, &apiclient.APIRequest{
		URL:     baseURL + "/config",
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

func (c *Client) CallDone(ctx context.Context, baseURL string, auth *types.Auth) (*types.ProvisionResponse, error) {
	pr := &types.ProvisionResponse{}

	resp, err := c.MakeRequest(ctx, &apiclient.APIRequest{
		URL:    baseURL + "/actions/provision",
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
