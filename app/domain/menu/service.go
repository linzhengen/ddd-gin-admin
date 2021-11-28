package menu

import (
	"context"
	"fmt"

	"github.com/linzhengen/ddd-gin-admin/app/domain/contextx"
	"github.com/linzhengen/ddd-gin-admin/app/domain/menu/menuaction"
	"github.com/linzhengen/ddd-gin-admin/app/domain/menu/menuactionresource"
	"github.com/linzhengen/ddd-gin-admin/app/domain/pagination"
	"github.com/linzhengen/ddd-gin-admin/app/domain/trans"

	"github.com/linzhengen/ddd-gin-admin/app/domain/errors"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/uuid"
)

type Service interface {
	Query(ctx context.Context, params QueryParam) (Menus, *pagination.Pagination, error)
	Get(ctx context.Context, id string) (*Menu, error)
	QueryActions(ctx context.Context, id string) (menuaction.MenuActions, error)
	Create(ctx context.Context, item *Menu) (string, error)
	Update(ctx context.Context, id string, item *Menu) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status int) error
}

func NewService(
	transRepo trans.Repository,
	menuRepo Repository,
	menuActionRepo menuaction.Repository,
	menuActionResourceRepo menuactionresource.Repository,
) Service {
	return &service{
		transRepo:              transRepo,
		menuRepo:               menuRepo,
		menuActionRepo:         menuActionRepo,
		menuActionResourceRepo: menuActionResourceRepo,
	}
}

type service struct {
	transRepo              trans.Repository
	menuRepo               Repository
	menuActionRepo         menuaction.Repository
	menuActionResourceRepo menuactionresource.Repository
}

func (a *service) Query(ctx context.Context, params QueryParam) (Menus, *pagination.Pagination, error) {
	menuActionResult, _, err := a.menuActionRepo.Query(ctx, menuaction.QueryParam{})
	if err != nil {
		return nil, nil, err
	}

	menuResult, pr, err := a.menuRepo.Query(ctx, params)
	if err != nil {
		return nil, nil, err
	}
	menuResult.FillMenuAction(menuActionResult.ToMenuIDMap())
	return menuResult, pr, nil
}

func (a *service) Get(ctx context.Context, id string) (*Menu, error) {
	item, err := a.menuRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errors.ErrNotFound
	}

	actions, err := a.QueryActions(ctx, id)
	if err != nil {
		return nil, err
	}
	item.Actions = actions

	return item, nil
}

func (a *service) QueryActions(ctx context.Context, id string) (menuaction.MenuActions, error) {
	result, _, err := a.menuActionRepo.Query(ctx, menuaction.QueryParam{
		MenuID: id,
	})
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}

	resourceResult, _, err := a.menuActionResourceRepo.Query(ctx, menuactionresource.QueryParam{
		MenuID: id,
	})
	if err != nil {
		return nil, err
	}
	result.FillResources(resourceResult.ToMenuActionIDMap())
	return result, nil
}

func (a *service) Create(ctx context.Context, item *Menu) (string, error) {
	if err := a.checkName(ctx, item); err != nil {
		return "", err
	}

	parentPath, err := a.getParentPath(ctx, item.ParentID)
	if err != nil {
		return "", err
	}
	item.ParentPath = parentPath
	item.ID = uuid.MustString()

	err = a.transRepo.Exec(ctx, func(ctx context.Context) error {
		err := a.createActions(ctx, item.ID, item.Actions)
		if err != nil {
			return err
		}

		return a.menuRepo.Create(ctx, item)
	})
	if err != nil {
		return "", err
	}

	return item.ID, nil
}

func (a *service) Update(ctx context.Context, id string, item *Menu) error {
	if id == item.ParentID {
		return errors.ErrInvalidParent
	}

	oldItem, err := a.Get(ctx, id)
	if err != nil {
		return err
	}
	if oldItem == nil {
		return errors.ErrNotFound
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

		err = a.updateChildParentPath(ctx, oldItem, item)
		if err != nil {
			return err
		}

		return a.menuRepo.Update(ctx, id, item)
	})
}

func (a *service) checkName(ctx context.Context, item *Menu) error {
	_, pr, err := a.menuRepo.Query(ctx, QueryParam{
		PaginationParam: pagination.Param{
			OnlyCount: true,
		},
		ParentID: &item.ParentID,
		Name:     item.Name,
	})
	if err != nil {
		return err
	}
	if pr.Total > 0 {
		return errors.New400Response("The menu name already exists")
	}
	return nil
}

func (a *service) getParentPath(ctx context.Context, parentID string) (string, error) {
	if parentID == "" {
		return "", nil
	}

	pitem, err := a.menuRepo.Get(ctx, parentID)
	if err != nil {
		return "", err
	}
	if pitem == nil {
		return "", errors.ErrInvalidParent
	}

	return a.joinParentPath(pitem.ParentPath, pitem.ID), nil
}

func (a *service) joinParentPath(parent, id string) string {
	if parent != "" {
		return parent + "/" + id
	}
	return id
}

func (a *service) createActions(ctx context.Context, menuID string, items menuaction.MenuActions) error {
	for _, item := range items {
		item.ID = uuid.MustString()
		item.MenuID = menuID
		err := a.menuActionRepo.Create(ctx, item)
		if err != nil {
			return err
		}

		for _, ritem := range item.Resources {
			ritem.ID = uuid.MustString()
			ritem.ActionID = item.ID
			err := a.menuActionResourceRepo.Create(ctx, ritem)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *service) updateActions(ctx context.Context, menuID string, oldItems, newItems menuaction.MenuActions) error {
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
			err := a.menuActionRepo.Update(ctx, item.ID, oitem)
			if err != nil {
				return err
			}
		}

		// update new and delete, not update
		addResources, delResources := a.compareResources(oitem.Resources, item.Resources)
		for _, aritem := range addResources {
			aritem.ID = uuid.MustString()
			aritem.ActionID = oitem.ID
			err := a.menuActionResourceRepo.Create(ctx, aritem)
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

func (a *service) compareActions(oldActions, newActions menuaction.MenuActions) (addList, delList, updateList menuaction.MenuActions) {
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

func (a *service) compareResources(oldResources, newResources menuactionresource.MenuActionResources) (addList, delList menuactionresource.MenuActionResources) {
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

func (a *service) updateChildParentPath(ctx context.Context, oldItem, newItem *Menu) error {
	if oldItem.ParentID == newItem.ParentID {
		return nil
	}

	opath := a.joinParentPath(oldItem.ParentPath, oldItem.ID)
	result, _, err := a.menuRepo.Query(contextx.NewNoTrans(ctx), QueryParam{
		PrefixParentPath: opath,
	})
	if err != nil {
		return err
	}

	npath := a.joinParentPath(newItem.ParentPath, newItem.ID)
	for _, menu := range result {
		parentPath := menu.ParentPath
		err = a.menuRepo.UpdateParentPath(ctx, menu.ID, fmt.Sprintf("%s%s", npath, parentPath[len(opath):]))
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *service) Delete(ctx context.Context, id string) error {
	oldItem, err := a.menuRepo.Get(ctx, id)
	if err != nil {
		return err
	}
	if oldItem == nil {
		return errors.ErrNotFound
	}

	_, pr, err := a.menuRepo.Query(ctx, QueryParam{
		PaginationParam: pagination.Param{OnlyCount: true},
		ParentID:        &id,
	})
	if err != nil {
		return err
	}
	if pr.Total > 0 {
		return errors.ErrNotAllowDeleteWithChild
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

func (a *service) UpdateStatus(ctx context.Context, id string, status int) error {
	oldItem, err := a.menuRepo.Get(ctx, id)
	if err != nil {
		return err
	}
	if oldItem == nil {
		return errors.ErrNotFound
	}

	return a.menuRepo.UpdateStatus(ctx, id, status)
}
