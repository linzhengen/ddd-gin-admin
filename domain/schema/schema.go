package schema

// HTTPStatusText 定义HTTP状态文本
type HTTPStatusText string

func (t HTTPStatusText) String() string {
	return string(t)
}

// Http status codes.
const (
	OKStatusText HTTPStatusText = "OK"
)

// HTTPError ...
type HTTPError struct {
	Error HTTPErrorItem `json:"error"`
}

// HTTPErrorItem ...
type HTTPErrorItem struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// HTTPStatus ...
type HTTPStatus struct {
	Status string `json:"status"`
}

// HTTPList ...
type HTTPList struct {
	List       interface{}     `json:"list"`
	Pagination *HTTPPagination `json:"pagination,omitempty"`
}

// HTTPPagination ...
type HTTPPagination struct {
	Total    int `json:"total"`
	Current  int `json:"current"`
	PageSize int `json:"pageSize"`
}

// PaginationParam ...
type PaginationParam struct {
	PageIndex int
	PageSize  int
}

// PaginationResult ...
type PaginationResult struct {
	Total int
}
