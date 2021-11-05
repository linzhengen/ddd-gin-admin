package response

import "github.com/linzhengen/ddd-gin-admin/app/domain/pagination"

type StatusText string

func (t StatusText) String() string {
	return string(t)
}

const (
	OKStatus    StatusText = "OK"
	ErrorStatus StatusText = "ERROR"
	FailStatus  StatusText = "FAIL"
)

type StatusResult struct {
	Status StatusText `json:"status"` // Result status
}

type ErrorResult struct {
	Error ErrorItem `json:"error"` // Error
}

type ErrorItem struct {
	Code    int    `json:"code"`    // Error Code
	Message string `json:"message"` // Error Message
}

type ListResult struct {
	List       interface{}            `json:"list"`                 // List
	Pagination *pagination.Pagination `json:"pagination,omitempty"` // Pagination
}

func NewIDResult(id string) *IDResult {
	return &IDResult{
		ID: id,
	}
}

type IDResult struct {
	ID string `json:"id"`
}
