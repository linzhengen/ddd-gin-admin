package repository

import "github.com/google/wire"

var RepoSet = wire.NewSet(
	DemoSet,
	MenuActionResourceSet,
	MenuActionSet,
	MenuSet,
	RoleMenuSet,
	RoleSet,
	TransSet,
	UserRoleSet,
	UserSet,
)
