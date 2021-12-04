package response

import (
	"strings"
	"time"

	"github.com/linzhengen/ddd-gin-admin/pkg/util/structure"

	"github.com/linzhengen/ddd-gin-admin/app/domain/menu"

	"github.com/linzhengen/ddd-gin-admin/app/domain/pagination"

	"github.com/linzhengen/ddd-gin-admin/pkg/util/json"
)

type Menu struct {
	ID         string      `json:"id"`                                         // ID
	Name       string      `json:"name" binding:"required"`                    // Name
	Sequence   int         `json:"sequence"`                                   // Sequence
	Icon       string      `json:"icon"`                                       // Icon
	Router     string      `json:"router"`                                     // Router
	ParentID   string      `json:"parent_id"`                                  // Parent ID
	ParentPath string      `json:"parent_path"`                                // Parent Path
	ShowStatus int         `json:"show_status" binding:"required,max=2,min=1"` // Show Status(1:show 2:hide)
	Status     int         `json:"status" binding:"required,max=2,min=1"`      // Menu Status(1:enable 2:disable)
	Memo       string      `json:"memo"`                                       // Memo
	Creator    string      `json:"creator"`                                    // Creator
	CreatedAt  time.Time   `json:"created_at"`                                 // CreatedAt
	UpdatedAt  time.Time   `json:"updated_at"`                                 // UpdatedAt
	Actions    MenuActions `json:"actions"`                                    // Actions
}

func (a *Menu) String() string {
	return json.MarshalToString(a)
}

func (a *Menu) ToDomain() *menu.Menu {
	item := new(menu.Menu)
	structure.Copy(a, item)
	return item
}

func MenuFromDomain(menu *menu.Menu) *Menu {
	item := new(Menu)
	structure.Copy(menu, item)
	return item
}

type MenuQueryResult struct {
	Data       Menus
	PageResult *pagination.Pagination
}

type Menus []*Menu

func MenusFromDomain(menus menu.Menus) Menus {
	list := make([]*Menu, len(menus))
	for i, item := range menus {
		ts := new(Menu)
		structure.Copy(item, ts)
		list[i] = ts
	}
	return list
}

func (a Menus) Len() int {
	return len(a)
}

func (a Menus) Less(i, j int) bool {
	return a[i].Sequence > a[j].Sequence
}

func (a Menus) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a Menus) ToMap() map[string]*Menu {
	m := make(map[string]*Menu)
	for _, item := range a {
		m[item.ID] = item
	}
	return m
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

func (a Menus) ToTree() MenuTrees {
	list := make(MenuTrees, len(a))
	for i, item := range a {
		list[i] = &MenuTree{
			ID:         item.ID,
			Name:       item.Name,
			Icon:       item.Icon,
			Router:     item.Router,
			ParentID:   item.ParentID,
			ParentPath: item.ParentPath,
			Sequence:   item.Sequence,
			ShowStatus: item.ShowStatus,
			Status:     item.Status,
			Actions:    item.Actions,
		}
	}
	return list.ToTree()
}

func (a Menus) FillMenuAction(mActions map[string]MenuActions) Menus {
	for _, item := range a {
		if v, ok := mActions[item.ID]; ok {
			item.Actions = v
		}
	}
	return a
}

// ----------------------------------------MenuTree--------------------------------------

type MenuTree struct {
	ID         string      `yaml:"-" json:"id"`
	Name       string      `yaml:"name" json:"name"`
	Icon       string      `yaml:"icon" json:"icon"`
	Router     string      `yaml:"router,omitempty" json:"router"`
	ParentID   string      `yaml:"-" json:"parent_id"`
	ParentPath string      `yaml:"-" json:"parent_path"`
	Sequence   int         `yaml:"sequence" json:"sequence"`
	ShowStatus int         `yaml:"-" json:"show_status"` // 1:show 2:hide
	Status     int         `yaml:"-" json:"status"`      // 1:enable 2:disable
	Actions    MenuActions `yaml:"actions,omitempty" json:"actions"`
	Children   *MenuTrees  `yaml:"children,omitempty" json:"children,omitempty"`
}

type MenuTrees []*MenuTree

func (a MenuTrees) ToTree() MenuTrees {
	mi := make(map[string]*MenuTree)
	for _, item := range a {
		mi[item.ID] = item
	}

	var list MenuTrees
	for _, item := range a {
		if item.ParentID == "" {
			list = append(list, item)
			continue
		}
		if pitem, ok := mi[item.ParentID]; ok {
			if pitem.Children == nil {
				children := MenuTrees{item}
				pitem.Children = &children
				continue
			}
			*pitem.Children = append(*pitem.Children, item)
		}
	}
	return list
}

// ----------------------------------------MenuAction--------------------------------------

type MenuAction struct {
	ID        string              `yaml:"-" json:"id"`
	MenuID    string              `yaml:"-" binding:"required" json:"menu_id"`
	Code      string              `yaml:"code" binding:"required" json:"code"`
	Name      string              `yaml:"name" binding:"required" json:"name"`
	Resources MenuActionResources `yaml:"resources,omitempty" json:"resources"`
}

type MenuActionQueryResult struct {
	Data       MenuActions
	PageResult *pagination.Pagination
}

type MenuActions []*MenuAction

func (a MenuActions) ToMap() map[string]*MenuAction {
	m := make(map[string]*MenuAction)
	for _, item := range a {
		m[item.Code] = item
	}
	return m
}

func (a MenuActions) FillResources(mResources map[string]MenuActionResources) {
	for i, item := range a {
		a[i].Resources = mResources[item.ID]
	}
}

func (a MenuActions) ToMenuIDMap() map[string]MenuActions {
	m := make(map[string]MenuActions)
	for _, item := range a {
		m[item.MenuID] = append(m[item.MenuID], item)
	}
	return m
}

// ----------------------------------------MenuActionResource--------------------------------------

type MenuActionResource struct {
	ID       string `yaml:"-" json:"id"`
	ActionID string `yaml:"-" json:"action_id"`
	Method   string `yaml:"method" binding:"required" json:"method"`
	Path     string `yaml:"path" binding:"required" json:"path"`
}

type MenuActionResourceQueryResult struct {
	Data       MenuActionResources
	PageResult *pagination.Pagination
}

type MenuActionResources []*MenuActionResource

func (a MenuActionResources) ToMap() map[string]*MenuActionResource {
	m := make(map[string]*MenuActionResource)
	for _, item := range a {
		m[item.Method+item.Path] = item
	}
	return m
}

func (a MenuActionResources) ToActionIDMap() map[string]MenuActionResources {
	m := make(map[string]MenuActionResources)
	for _, item := range a {
		m[item.ActionID] = append(m[item.ActionID], item)
	}
	return m
}
