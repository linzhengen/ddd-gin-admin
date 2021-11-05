package user

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/pagination"
)

type Repository interface {
	Query(ctx context.Context, params QueryParams) (Users, *pagination.Pagination, error)
	Get(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, item *User) error
	Update(ctx context.Context, id string, item *User) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status int) error
	UpdatePassword(ctx context.Context, id, password string) error
}
