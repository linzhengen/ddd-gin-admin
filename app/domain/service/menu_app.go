package service

import (
	"context"
	"os"

	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/contextx"
	errors2 "github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/errors"

	"github.com/linzhengen/ddd-gin-admin/app/domain/entity"

	"github.com/linzhengen/ddd-gin-admin/pkg/util/structure"

	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/response"

	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/request"

	"github.com/linzhengen/ddd-gin-admin/app/domain/repository"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/uuid"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/yaml"
)

type Menu interface {
	InitData(ctx context.Context, dataFile string) error
	Query(ctx context.Context, req request.MenuQuery) (*response.MenuQuery, error)
	Get(ctx context.Context, id string) (*response.Menu, error)
	QueryActions(ctx context.Context, id string) (response.MenuActions, error)
	Create(ctx context.Context, req request.Menu) (*response.IDResult, error)
	Update(ctx context.Context, id string, req request.Menu) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status int) error
}

func NewMenu(
	transRepo repository.TransRepository,
	menuRepo repository.MenuRepository,
	menuActionRepo repository.MenuActionRepository,
	menuActionResourceRepo repository.MenuActionResourceRepository,
) Menu {
	return &menu{
		transRepo:              transRepo,
		menuRepo:               menuRepo,
		menuActionRepo:         menuActionRepo,
		menuActionResourceRepo: menuActionResourceRepo,
	}
}

type menu struct {
	transRepo              repository.TransRepository
	menuRepo               repository.MenuRepository
	menuActionRepo         repository.MenuActionRepository
	menuActionResourceRepo repository.MenuActionResourceRepository
}

func (a *menu) InitData(ctx context.Context, dataFile string) error {
	_, pager, err := a.menuRepo.Query(ctx, request.MenuQuery{
		Pagination: request.Pagination{OnlyCount: true},
	})
	if err != nil {
		return err
	}
	if pager.Total > 0 {
		return nil
	}

	data, err := a.readData(dataFile)
	if err != nil {
		return err
	}

	return a.createMenus(ctx, "", data)
}

func (a *menu) readData(name string) (response.MenuTrees, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data response.MenuTrees
	d := yaml.NewDecoder(file)
	d.SetStrict(true)
	err = d.Decode(&data)
	return data, err
}

func (a *menu) createMenus(ctx context.Context, parentID string, list response.MenuTrees) error {
	return a.transRepo.Exec(ctx, func(ctx context.Context) error {
		for _, item := range list {
			menuRes := response.Menu{
				Name:       item.Name,
				Sequence:   item.Sequence,
				Icon:       item.Icon,
				Router:     item.Router,
				ParentID:   parentID,
				Status:     1,
				ShowStatus: 1,
				Actions:    item.Actions,
			}
			if v := item.ShowStatus; v > 0 {
				menuRes.ShowStatus = v
			}
			menuReq := request.Menu{}
			structure.Copy(menuRes, menuReq)
			nsitem, err := a.Create(ctx, menuReq)
			if err != nil {
				return err
			}

			if item.Children != nil && len(*item.Children) > 0 {
				err := a.createMenus(ctx, nsitem.ID, *item.Children)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (a *menu) Query(ctx context.Context, req request.MenuQuery) (*response.MenuQuery, error) {
	menuActions, _, err := a.menuActionRepo.Query(ctx, request.MenuActionQuery{})
	if err != nil {
		return nil, err
	}

	menus, page, err := a.menuRepo.Query(ctx, req)
	if err != nil {
		return nil, err
	}
	menuActionsRes := response.MenuActions{}
	structure.Copy(menuActions, menuActionsRes)
	menusRes := response.Menus{}
	structure.Copy(menus, menusRes)
	menusRes.FillMenuAction(menuActionsRes.ToMenuIDMap())
	return &response.MenuQuery{
		Data:       menusRes,
		PageResult: page,
	}, nil
}

func (a *menu) Get(ctx context.Context, id string) (*response.Menu, error) {
	menu, err := a.menuRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if menu == nil {
		return nil, errors2.ErrNotFound
	}

	actions, err := a.QueryActions(ctx, id)
	if err != nil {
		return nil, err
	}
	menuRes := new(response.Menu)
	structure.Copy(menu, menuRes)
	menuRes.Actions = actions

	return menuRes, nil
}

func (a *menu) QueryActions(ctx context.Context, id string) (response.MenuActions, error) {
	menuActions, _, err := a.menuActionRepo.Query(ctx, request.MenuActionQuery{
		MenuID: id,
	})
	if err != nil {
		return nil, err
	}
	if len(menuActions) == 0 {
		return nil, nil
	}

	resources, _, err := a.menuActionResourceRepo.Query(ctx, request.MenuActionResourceQuery{
		MenuID: id,
	})
	if err != nil {
		return nil, err
	}
	menuActionsRes := response.MenuActions{}
	structure.Copy(menuActions, menuActionsRes)
	resourcesRes := response.MenuActionResources{}
	structure.Copy(resources, resourcesRes)
	menuActionsRes.FillResources(resourcesRes.ToActionIDMap())

	return menuActionsRes, nil
}

func (a *menu) checkName(ctx context.Context, item request.Menu) error {
	_, page, err := a.menuRepo.Query(ctx, request.MenuQuery{
		Pagination: request.Pagination{
			OnlyCount: true,
		},
		ParentID: &item.ParentID,
		Name:     item.Name,
	})
	if err != nil {
		return err
	}
	if page.Total > 0 {
		return errors2.New400Response("The menu name already exists")
	}
	return nil
}

func (a *menu) Create(ctx context.Context, req request.Menu) (*response.IDResult, error) {
	if err := a.checkName(ctx, req); err != nil {
		return nil, err
	}

	parentPath, err := a.getParentPath(ctx, req.ParentID)
	if err != nil {
		return nil, err
	}
	req.ParentPath = parentPath
	req.ID = uuid.MustString()

	err = a.transRepo.Exec(ctx, func(ctx context.Context) error {
		err := a.createActions(ctx, req.ID, req.Actions)
		if err != nil {
			return err
		}
		itemEntity := entity.Menu{}
		structure.Copy(req, itemEntity)
		return a.menuRepo.Create(ctx, itemEntity)
	})
	if err != nil {
		return nil, err
	}

	return response.NewIDResult(req.ID), nil
}

func (a *menu) createActions(ctx context.Context, menuID string, menuActions request.MenuActions) error {
	for _, item := range menuActions {
		item.ID = uuid.MustString()
		item.MenuID = menuID
		menuActionEntity := entity.MenuAction{}
		structure.Copy(item, menuActionEntity)
		err := a.menuActionRepo.Create(ctx, menuActionEntity)
		if err != nil {
			return err
		}

		for _, ritem := range item.Resources {
			ritem.ID = uuid.MustString()
			ritem.ActionID = item.ID
			resourceEntity := entity.MenuActionResource{}
			structure.Copy(ritem, resourceEntity)
			err := a.menuActionResourceRepo.Create(ctx, resourceEntity)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *menu) getParentPath(ctx context.Context, parentID string) (string, error) {
	if parentID == "" {
		return "", nil
	}

	pitem, err := a.menuRepo.Get(ctx, parentID)
	if err != nil {
		return "", err
	}
	if pitem == nil {
		return "", errors2.ErrInvalidParent
	}

	return a.joinParentPath(*pitem.ParentPath, pitem.ID), nil
}

func (a *menu) joinParentPath(parent, id string) string {
	if parent != "" {
		return parent + "/" + id
	}
	return id
}

func (a *menu) Update(ctx context.Context, id string, item request.Menu) error {
	if id == item.ParentID {
		return errors2.ErrInvalidParent
	}

	oldItem, err := a.Get(ctx, id)
	if err != nil {
		return err
	}
	if oldItem == nil {
		return errors2.ErrNotFound
	}
	if oldItem.Name != item.Name {
		if err := a.checkName(ctx, item); err != nil {
			return err
		}
	}

	item.ID = oldItem.ID
	item.Creator = oldItem.Creator
	item.CreatedAt = oldItem.CreatedAt

	if oldItem.ParentID != item.ParentID {
		parentPath, err := a.getParentPath(ctx, item.ParentID)
		if err != nil {
			return err
		}
		item.ParentPath = parentPath
	} else {
		item.ParentPath = oldItem.ParentPath
	}

	return a.transRepo.Exec(ctx, func(ctx context.Context) error {
		err := a.updateActions(ctx, id, oldItem.Actions, item.Actions)
		if err != nil {
			return err
		}

		err = a.updateChildParentPath(ctx, *oldItem, item)
		if err != nil {
			return err
		}

		return a.menuRepo.Update(ctx, id, item)
	})
}

func (a *menu) updateActions(ctx context.Context, menuID string, oldItems response.MenuActions, newItems request.MenuActions) error {
	addActions, delActions, updateActions := a.compareActions(oldItems, newItems)

	err := a.createActions(ctx, menuID, addActions)
	if err != nil {
		return err
	}

	for _, item := range delActions {
		err := a.menuActionRepo.Delete(ctx, item.ID)
		if err != nil {
			return err
		}

		err = a.menuActionResourceRepo.DeleteByActionID(ctx, item.ID)
		if err != nil {
			return err
		}
	}

	mOldItems := oldItems.ToMap()
	for _, item := range updateActions {
		oitem := mOldItems[item.Code]
		// only update action name
		if item.Name != oitem.Name {
			oitem.Name = item.Name
			err := a.menuActionRepo.Update(ctx, item.ID, *oitem)
			if err != nil {
				return err
			}
		}

		// update new and delete, not update
		addResources, delResources := a.compareResources(oitem.Resources, item.Resources)
		for _, aritem := range addResources {
			aritem.ID = uuid.MustString()
			aritem.ActionID = oitem.ID
			err := a.menuActionResourceRepo.Create(ctx, *aritem)
			if err != nil {
				return err
			}
		}

		for _, ditem := range delResources {
			err := a.menuActionResourceRepo.Delete(ctx, ditem.ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (a *menu) compareActions(oldActions, newActions response.MenuActions) (addList, delList, updateList response.MenuActions) {
	mOldActions := oldActions.ToMap()
	mNewActions := newActions.ToMap()

	for k, item := range mNewActions {
		if _, ok := mOldActions[k]; ok {
			updateList = append(updateList, item)
			delete(mOldActions, k)
			continue
		}
		addList = append(addList, item)
	}

	for _, item := range mOldActions {
		delList = append(delList, item)
	}
	return
}

func (a *menu) compareResources(oldResources, newResources response.MenuActionResources) (addList, delList response.MenuActionResources) {
	mOldResources := oldResources.ToMap()
	mNewResources := newResources.ToMap()

	for k, item := range mNewResources {
		if _, ok := mOldResources[k]; ok {
			delete(mOldResources, k)
			continue
		}
		addList = append(addList, item)
	}

	for _, item := range mOldResources {
		delList = append(delList, item)
	}
	return
}

func (a *menu) updateChildParentPath(ctx context.Context, oldItem response.Menu, newItem request.Menu) error {
	if oldItem.ParentID == newItem.ParentID {
		return nil
	}

	opath := a.joinParentPath(oldItem.ParentPath, oldItem.ID)
	menus, _, err := a.menuRepo.Query(contextx.NewNoTrans(ctx), request.MenuQuery{
		PrefixParentPath: opath,
	})
	if err != nil {
		return err
	}
	menusRes := response.Menus{}
	structure.Copy(menus, menusRes)
	npath := a.joinParentPath(newItem.ParentPath, newItem.ID)
	for _, menu := range menusRes {
		err = a.menuRepo.UpdateParentPath(ctx, menu.ID, npath+menu.ParentPath[len(opath):])
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *menu) Delete(ctx context.Context, id string) error {
	oldItem, err := a.menuRepo.Get(ctx, id)
	if err != nil {
		return err
	}
	if oldItem == nil {
		return errors2.ErrNotFound
	}

	_, page, err := a.menuRepo.Query(ctx, request.MenuQuery{
		Pagination: request.Pagination{OnlyCount: true},
		ParentID:   &id,
	})
	if err != nil {
		return err
	}
	if page.Total > 0 {
		return errors2.ErrNotAllowDeleteWithChild
	}

	return a.transRepo.Exec(ctx, func(ctx context.Context) error {
		err = a.menuActionResourceRepo.DeleteByMenuID(ctx, id)
		if err != nil {
			return err
		}

		err := a.menuActionRepo.DeleteByMenuID(ctx, id)
		if err != nil {
			return err
		}

		return a.menuRepo.Delete(ctx, id)
	})
}

func (a *menu) UpdateStatus(ctx context.Context, id string, status int) error {
	oldItem, err := a.menuRepo.Get(ctx, id)
	if err != nil {
		return err
	}
	if oldItem == nil {
		return errors2.ErrNotFound
	}

	return a.menuRepo.UpdateStatus(ctx, id, status)
}
