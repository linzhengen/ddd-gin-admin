package pagination

type Pagination struct {
	Total    int  `json:"total"`    // Total count
	Current  uint `json:"current"`  // Current Page
	PageSize uint `json:"pageSize"` // Page Size
}

type Param struct {
	Pagination bool `form:"-"`                                     // Pagination
	OnlyCount  bool `form:"-"`                                     // Only count
	Current    uint `form:"current,default=1"`                     // Current page
	PageSize   uint `form:"pageSize,default=10" binding:"max=100"` // Page size
}

func (a Param) GetCurrent() uint {
	return a.Current
}

func (a Param) GetPageSize() uint {
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

type OrderField struct {
	Key       string
	Direction OrderDirection
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

type OrderFields []*OrderField

func (a OrderFields) AddIdSortField() OrderFields {
	return append(a, NewOrderField("id", OrderByDESC))
}
