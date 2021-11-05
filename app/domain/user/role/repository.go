package role

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/pagination"
)

type Repository interface {
	Query(ctx context.Context, params QueryParam) (Roles, *pagination.Pagination, error)
	Get(ctx context.Context, id string) (*Role, error)
	Create(ctx context.Context, item *Role) error
	Update(ctx context.Context, id string, item *Role) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status int) error
}
