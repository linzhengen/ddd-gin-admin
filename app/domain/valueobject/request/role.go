package request

type RoleQuery struct {
	Pagination
	OrderFields
	IDs        []string `form:"-"`          // IDs
	Name       string   `form:"-"`          // Name
	QueryValue string   `form:"queryValue"` // Query Search Values
	UserID     string   `form:"-"`          // User ID
	Status     int      `form:"status"`     // Status(1:enable 2:disable)
}

type RoleMenuQuery struct {
	Pagination
	OrderFields
	RoleID  string
	RoleIDs []string
}
