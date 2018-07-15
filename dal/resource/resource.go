package resource

import (
	"context"

	"github.com/talpert/hellofour/dal/resource/types"
)

type IDAL interface {
	Insert(ctx context.Context, resource *types.Resource) error
	Get(ctx context.Context, uuid string) (*types.Resource, error)
}
