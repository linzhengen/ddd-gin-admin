package request

type OrderDirection int

const (
	OrderByASC  OrderDirection = 1
	OrderByDESC OrderDirection = 2
)

type OrderField struct {
	Key       string
	Direction OrderDirection
}

type OrderFields []*OrderField

func (a OrderFields) AddIdSortKey() OrderFields {
	return append(a, NewOrderField("id", OrderByASC))
}

type Pagination struct {
	Pagination bool
	OnlyCount  bool
	Current    uint
	PageSize   uint
}

func (a Pagination) GetCurrent() uint {
	return a.Current
}

func (a Pagination) GetPageSize() uint {
	pageSize := a.PageSize
	if a.PageSize == 0 {
		pageSize = 100
	}
	return pageSize
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
