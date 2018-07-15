package resource

import (
	"context"

	"github.com/talpert/hellofour/dal/resource/types"
)

type IManager interface {
	Provision(ctx context.Context, request *types.ProvisionRequest) error
}
