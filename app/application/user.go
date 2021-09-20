package application

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/service"
	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/schema"
)

type User interface {
	Query(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserQueryResult, error)
	QueryShow(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserShowQueryResult, error)
	Get(ctx context.Context, id string, opts ...schema.UserQueryOptions) (*schema.User, error)
	Create(ctx context.Context, item schema.User) (*schema.IDResult, error)
	Update(ctx context.Context, id string, item schema.User) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status int) error
}

func NewUser(
	userSvc service.User,
) User {
	return &user{
		userSvc: userSvc,
	}
}

type user struct {
	userSvc service.User
}

func (u user) Query(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserQueryResult, error) {
	return u.userSvc.Query(ctx, params, opts...)
}

func (u user) QueryShow(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserShowQueryResult, error) {
	return u.userSvc.QueryShow(ctx, params, opts...)
}

func (u user) Get(ctx context.Context, id string, opts ...schema.UserQueryOptions) (*schema.User, error) {
	return u.userSvc.Get(ctx, id, opts...)
}

func (u user) Create(ctx context.Context, item schema.User) (*schema.IDResult, error) {
	return u.userSvc.Create(ctx, item)
}

func (u user) Update(ctx context.Context, id string, item schema.User) error {
	return u.userSvc.Update(ctx, id, item)
}

func (u user) Delete(ctx context.Context, id string) error {
	return u.userSvc.Delete(ctx, id)
}

func (u user) UpdateStatus(ctx context.Context, id string, status int) error {
	return u.userSvc.UpdateStatus(ctx, id, status)
}
