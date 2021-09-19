package repository

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/response"

	"github.com/linzhengen/ddd-gin-admin/app/domain/entity"
	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/request"
)

type RoleRepository interface {
	Query(ctx context.Context, req request.RoleQuery) (entity.Roles, *response.Pagination, error)
	Get(ctx context.Context, id string) (*entity.Role, error)
	Create(ctx context.Context, item entity.Role) error
	Update(ctx context.Context, id string, item entity.Role) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status int) error
}
