package application

import "github.com/google/wire"

var ServiceSet = wire.NewSet(
	LoginSet,
	MenuSet,
	RoleSet,
	UserSet,
)
