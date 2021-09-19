package repository

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/response"

	"github.com/linzhengen/ddd-gin-admin/app/domain/entity"

	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/request"
)

type UserRepository interface {
	Query(ctx context.Context, req request.UserQuery) (entity.Users, *response.Pagination, error)
	Get(ctx context.Context, id string) (*entity.User, error)
	Create(ctx context.Context, item entity.User) error
	Update(ctx context.Context, id string, item entity.User) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status int) error
	UpdatePassword(ctx context.Context, id, password string) error
}
