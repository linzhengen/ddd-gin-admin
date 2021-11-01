package schema

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
	List       interface{}       `json:"list"`                 // List
	Pagination *PaginationResult `json:"pagination,omitempty"` // Pagination
}

type PaginationResult struct {
	Total    int  `json:"total"`    // Total count
	Current  uint `json:"current"`  // Current Page
	PageSize uint `json:"pageSize"` // Page Size
}

type PaginationParam struct {
	Pagination bool `form:"-"`                                     // Pagination
	OnlyCount  bool `form:"-"`                                     // Only count
	Current    uint `form:"current,default=1"`                     // Current page
	PageSize   uint `form:"pageSize,default=10" binding:"max=100"` // Page size
}

func (a PaginationParam) GetCurrent() uint {
	return a.Current
}

func (a PaginationParam) GetPageSize() uint {
	pageSize := a.PageSize
	if a.PageSize == 0 {
		pageSize = 100
	}
	return pageSize
}

type OrderDirection int

const (
	OrderByASC  OrderDirection = 1
	OrderByDESC OrderDirection = 2
)

// NewOrderFieldWithKeys create order fields with keys(default asc)
// keys: sort keys
// directions: sort order
func NewOrderFieldWithKeys(keys []string, directions ...map[int]OrderDirection) []*OrderField {
	m := make(map[int]OrderDirection)
	if len(directions) > 0 {
		m = directions[0]
	}

	fields := make([]*OrderField, len(keys))
	for i, key := range keys {
		d := OrderByASC
		if v, ok := m[i]; ok {
			d = v
		}

		fields[i] = NewOrderField(key, d)
	}

	return fields
}

func NewOrderFields(orderFields ...*OrderField) []*OrderField {
	return orderFields
}

func NewOrderField(key string, d OrderDirection) *OrderField {
	return &OrderField{
		Key:       key,
		Direction: d,
	}
}

type OrderField struct {
	Key       string
	Direction OrderDirection
}

type OrderFields []*OrderField

func (a OrderFields) AddIdSortField() OrderFields {
	return append(a, NewOrderField("id", OrderByDESC))
}

func NewIDResult(id string) *IDResult {
	return &IDResult{
		ID: id,
	}
}

type IDResult struct {
	ID string `json:"id"`
}