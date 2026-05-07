package rbac

import "context"

type Repository interface {
	ListRolesPolicies(ctx context.Context) ([]string, error)
	ListUsersPolicies(ctx context.Context) ([]string, error)
}
