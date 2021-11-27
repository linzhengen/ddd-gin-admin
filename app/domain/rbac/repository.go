package rbac

import "context"

type Repository interface {
	ListRolesPolices(ctx context.Context) ([]string, error)
	ListUsersPolices(ctx context.Context) ([]string, error)
}
