package schema

import "time"

type Demo struct {
	ID        string    `json:"id"`                                    // ID
	Code      string    `json:"code" binding:"required"`               // Code
	Name      string    `json:"name" binding:"required"`               // Nmae
	Memo      string    `json:"memo"`                                  // Memo
	Status    int       `json:"status" binding:"required,max=2,min=1"` // Status(1:enable 2:disable)
	Creator   string    `json:"creator"`                               // Creator
	CreatedAt time.Time `json:"created_at"`                            // CreatedAt
	UpdatedAt time.Time `json:"updated_at"`                            // UpdatedAt
}

type DemoQueryParam struct {
	PaginationParam
	Code       string `form:"-"`
	QueryValue string `form:"queryValue"`
}

type DemoQueryOptions struct {
	OrderFields []*OrderField
}

type DemoQueryResult struct {
	Data       []*Demo
	PageResult *PaginationResult
}
