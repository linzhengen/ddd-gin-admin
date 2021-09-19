package request

type UserQuery struct {
	Pagination
	OrderFields
	UserName   string   `form:"userName"`   // User Name
	QueryValue string   `form:"queryValue"` // Query search values
	Status     int      `form:"status"`     // Status(1:enable 2:disable)
	RoleIDs    []string `form:"-"`          // Role IDs
}

type UserRoleQuery struct {
	Pagination
	OrderFields
	UserID  string
	UserIDs []string
}
