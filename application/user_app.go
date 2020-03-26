package application

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/domain/repository"
	"github.com/linzhengen/ddd-gin-admin/domain/schema"
)

type userApp struct {
	us repository.UserRepository
}

var _ UserRepository = &userApp{}

// UserRepository ...
type UserRepository interface {
	Query(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserQueryResult, error)
	QueryShow(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserShowQueryResult, error)
	Get(ctx context.Context, recordID string, opts ...schema.UserQueryOptions) (*schema.User, error)
	Create(ctx context.Context, item schema.User) (*schema.User, error)
	Update(ctx context.Context, recordID string, item schema.User) (*schema.User, error)
	Delete(ctx context.Context, recordID string) error
	UpdateStatus(ctx context.Context, recordID string, status int) error
}

func (u *userApp) Query(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserQueryResult, error) {
	return u.us.Query(ctx, params, opts...)
}

func (u *userApp) QueryShow(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserShowQueryResult, error) {
	return u.us.QueryShow(ctx, params, opts...)
}

func (u *userApp) Get(ctx context.Context, recordID string, opts ...schema.UserQueryOptions) (*schema.User, error) {
	return u.us.Get(ctx, recordID, opts...)
}

func (u *userApp) Create(ctx context.Context, item schema.User) (*schema.User, error) {
	return u.us.Create(ctx, item)
}

func (u *userApp) Update(ctx context.Context, recordID string, item schema.User) (*schema.User, error) {
	return u.us.Update(ctx, recordID, item)
}

func (u *userApp) Delete(ctx context.Context, recordID string) error {
	return u.us.Delete(ctx, recordID)
}

func (u *userApp) UpdateStatus(ctx context.Context, recordID string, status int) error {
	return u.us.UpdateStatus(ctx, recordID, status)
}
