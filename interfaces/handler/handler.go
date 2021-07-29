package handler

import "github.com/google/wire"

var APISet = wire.NewSet(
	HealthCheckSet,
	DemoSet,
	LoginSet,
	MenuSet,
	RoleSet,
	UserSet,
)
