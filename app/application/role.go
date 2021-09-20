package application

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/service"

	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/schema"
)

type Role interface {
	Query(ctx context.Context, params schema.RoleQueryParam, opts ...schema.RoleQueryOptions) (*schema.RoleQueryResult, error)
	Get(ctx context.Context, id string, opts ...schema.RoleQueryOptions) (*schema.Role, error)
	QueryRoleMenus(ctx context.Context, roleID string) (schema.RoleMenus, error)
	Create(ctx context.Context, item schema.Role) (*schema.IDResult, error)
	Update(ctx context.Context, id string, item schema.Role) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status int) error
}

func NewRole(
	roleSvc service.Role,
) Role {
	return &role{
		roleSvc: roleSvc,
	}
}

type role struct {
	roleSvc service.Role
}

func (r role) Query(ctx context.Context, params schema.RoleQueryParam, opts ...schema.RoleQueryOptions) (*schema.RoleQueryResult, error) {
	return r.roleSvc.Query(ctx, params, opts...)
}

func (r role) Get(ctx context.Context, id string, opts ...schema.RoleQueryOptions) (*schema.Role, error) {
	return r.roleSvc.Get(ctx, id, opts...)
}

func (r role) QueryRoleMenus(ctx context.Context, roleID string) (schema.RoleMenus, error) {
	return r.roleSvc.QueryRoleMenus(ctx, roleID)
}

func (r role) Create(ctx context.Context, item schema.Role) (*schema.IDResult, error) {
	return r.roleSvc.Create(ctx, item)
}

func (r role) Update(ctx context.Context, id string, item schema.Role) error {
	return r.roleSvc.Update(ctx, id, item)
}

func (r role) Delete(ctx context.Context, id string) error {
	return r.roleSvc.Delete(ctx, id)
}

func (r role) UpdateStatus(ctx context.Context, id string, status int) error {
	return r.roleSvc.UpdateStatus(ctx, id, status)
}
