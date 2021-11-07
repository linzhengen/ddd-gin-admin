package userrole

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/pagination"
)

type Repository interface {
	Query(ctx context.Context, params QueryParam) (UserRoles, *pagination.Pagination, error)
	Get(ctx context.Context, id string) (*UserRole, error)
	Create(ctx context.Context, item *UserRole) error
	Update(ctx context.Context, id string, item *UserRole) error
	Delete(ctx context.Context, id string) error
	DeleteByUserID(ctx context.Context, userID string) error
}
