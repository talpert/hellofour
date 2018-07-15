package types

import ht "github.com/talpert/hellofour/dal/heroku/types"

type Resource struct {
	UUID        string
	CallbackURL string
	Name        string
	Auth        *ht.Auth
	Options     map[string]interface{}
	Plan        string
	Region      string
}
