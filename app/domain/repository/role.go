package repository

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/entity"

	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/schema"
)

type RoleRepository interface {
	Query(ctx context.Context, params schema.RoleQueryParam) (entity.Roles, *schema.PaginationResult, error)
	Get(ctx context.Context, id string) (*entity.Role, error)
	Create(ctx context.Context, item entity.Role) error
	Update(ctx context.Context, id string, item entity.Role) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status int) error
}
