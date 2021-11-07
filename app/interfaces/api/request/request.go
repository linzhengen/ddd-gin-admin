package request

type PaginationParam struct {
	Pagination bool `form:"-"`                                     // Pagination
	OnlyCount  bool `form:"-"`                                     // Only count
	Current    uint `form:"current,default=1"`                     // Current page
	PageSize   uint `form:"pageSize,default=10" binding:"max=100"` // Page size
}
