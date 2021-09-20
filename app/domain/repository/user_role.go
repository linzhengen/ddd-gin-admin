package repository

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/entity"

	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/schema"
)

type UserRoleRepository interface {
	Query(ctx context.Context, params schema.UserRoleQueryParam) (entity.UserRoles, *schema.PaginationResult, error)
	Get(ctx context.Context, id string) (*entity.UserRole, error)
	Create(ctx context.Context, item entity.UserRole) error
	Update(ctx context.Context, id string, item entity.UserRole) error
	Delete(ctx context.Context, id string) error
	DeleteByUserID(ctx context.Context, userID string) error
}
