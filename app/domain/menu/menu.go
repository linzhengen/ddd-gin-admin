package menu

import (
	"strings"
	"time"

	"github.com/linzhengen/ddd-gin-admin/app/domain/menu/menuaction"

	"github.com/linzhengen/ddd-gin-admin/app/domain/pagination"
)

type Menu struct {
	ID         string
	Name       string
	Sequence   int
	Icon       string
	Router     string
	ParentID   string
	ParentPath string
	ShowStatus int
	Status     int
	Memo       string
	Creator    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
	Actions    menuaction.MenuActions
}

type Menus []*Menu

type QueryParam struct {
	PaginationParam  pagination.Param
	OrderFields      pagination.OrderFields
	IDs              []string
	Name             string
	PrefixParentPath string
	QueryValue       string
	ParentID         *string
	ShowStatus       int
	Status           int
}

func (a Menus) FillMenuAction(mActions map[string]menuaction.MenuActions) Menus {
	for _, item := range a {
		if v, ok := mActions[item.ID]; ok {
			item.Actions = v
		}
	}
	return a
}

func (a Menus) SplitParentIDs() []string {
	idList := make([]string, 0, len(a))
	mIDList := make(map[string]struct{})

	for _, item := range a {
		if _, ok := mIDList[item.ID]; ok || item.ParentPath == "" {
			continue
		}

		for _, pp := range strings.Split(item.ParentPath, "/") {
			if _, ok := mIDList[pp]; ok {
				continue
			}
			idList = append(idList, pp)
			mIDList[pp] = struct{}{}
		}
	}

	return idList
}
