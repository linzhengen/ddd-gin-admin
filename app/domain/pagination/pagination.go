package pagination

type Pagination struct {
	Total    int
	Current  uint
	PageSize uint
}

type Param struct {
	Pagination bool
	OnlyCount  bool
	Current    uint
	PageSize   uint
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
