package handler

import "github.com/google/wire"

var APISet = wire.NewSet(
	HealthCheckSet,
	LoginSet,
	MenuSet,
	RoleSet,
	UserSet,
)
