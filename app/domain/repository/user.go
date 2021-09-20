package repository

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/entity"
	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/schema"
)

type UserRepository interface {
	Query(ctx context.Context, params schema.UserQueryParam) (entity.Users, *schema.PaginationResult, error)
	Get(ctx context.Context, id string) (*entity.User, error)
	Create(ctx context.Context, item entity.User) error
	Update(ctx context.Context, id string, item entity.User) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status int) error
	UpdatePassword(ctx context.Context, id, password string) error
}
