package application

import (
	"context"
	"os"

	"github.com/linzhengen/ddd-gin-admin/domain/repository"

	"github.com/google/wire"
	"github.com/linzhengen/ddd-gin-admin/domain/schema"
	"github.com/linzhengen/ddd-gin-admin/infrastructure/contextx"
	"github.com/linzhengen/ddd-gin-admin/pkg/errors"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/uuid"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/yaml"
)

var MenuSet = wire.NewSet(wire.Struct(new(Menu), "*"))

type Menu struct {
	TransModel              repository.TransRepository
	MenuModel               repository.MenuRepository
	MenuActionModel         repository.MenuActionRepository
	MenuActionResourceModel repository.MenuActionResourceRepository
}

func (a *Menu) InitData(ctx context.Context, dataFile string) error {
	result, err := a.MenuModel.Query(ctx, schema.MenuQueryParam{
		PaginationParam: schema.PaginationParam{OnlyCount: true},
	})
	if err != nil {
		return err
	}
	if result.PageResult.Total > 0 {
		return nil
	}

	data, err := a.readData(dataFile)
	if err != nil {
		return err
	}

	return a.createMenus(ctx, "", data)
}

func (a *Menu) readData(name string) (schema.MenuTrees, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data schema.MenuTrees
	d := yaml.NewDecoder(file)
	d.SetStrict(true)
	err = d.Decode(&data)
	return data, err
}

func (a *Menu) createMenus(ctx context.Context, parentID string, list schema.MenuTrees) error {
	return a.TransModel.Exec(ctx, func(ctx context.Context) error {
		for _, item := range list {
			sitem := schema.Menu{
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
				sitem.ShowStatus = v
			}

			nsitem, err := a.Create(ctx, sitem)
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

func (a *Menu) Query(ctx context.Context, params schema.MenuQueryParam, opts ...schema.MenuQueryOptions) (*schema.MenuQueryResult, error) {
	menuActionResult, err := a.MenuActionModel.Query(ctx, schema.MenuActionQueryParam{})
	if err != nil {
		return nil, err
	}

	result, err := a.MenuModel.Query(ctx, params, opts...)
	if err != nil {
		return nil, err
	}
	result.Data.FillMenuAction(menuActionResult.Data.ToMenuIDMap())
	return result, nil
}

func (a *Menu) Get(ctx context.Context, id string, opts ...schema.MenuQueryOptions) (*schema.Menu, error) {
	item, err := a.MenuModel.Get(ctx, id, opts...)
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

func (a *Menu) QueryActions(ctx context.Context, id string) (schema.MenuActions, error) {
	result, err := a.MenuActionModel.Query(ctx, schema.MenuActionQueryParam{
		MenuID: id,
	})
	if err != nil {
		return nil, err
	}
	if len(result.Data) == 0 {
		return nil, nil
	}

	resourceResult, err := a.MenuActionResourceModel.Query(ctx, schema.MenuActionResourceQueryParam{
		MenuID: id,
	})
	if err != nil {
		return nil, err
	}

	result.Data.FillResources(resourceResult.Data.ToActionIDMap())

	return result.Data, nil
}

func (a *Menu) checkName(ctx context.Context, item schema.Menu) error {
	result, err := a.MenuModel.Query(ctx, schema.MenuQueryParam{
		PaginationParam: schema.PaginationParam{
			OnlyCount: true,
		},
		ParentID: &item.ParentID,
		Name:     item.Name,
	})
	if err != nil {
		return err
	}
	if result.PageResult.Total > 0 {
		return errors.New400Response("The menu name already exists")
	}
	return nil
}

func (a *Menu) Create(ctx context.Context, item schema.Menu) (*schema.IDResult, error) {
	if err := a.checkName(ctx, item); err != nil {
		return nil, err
	}

	parentPath, err := a.getParentPath(ctx, item.ParentID)
	if err != nil {
		return nil, err
	}
	item.ParentPath = parentPath
	item.ID = uuid.MustString()

	err = a.TransModel.Exec(ctx, func(ctx context.Context) error {
		err := a.createActions(ctx, item.ID, item.Actions)
		if err != nil {
			return err
		}

		return a.MenuModel.Create(ctx, item)
	})
	if err != nil {
		return nil, err
	}

	return schema.NewIDResult(item.ID), nil
}

func (a *Menu) createActions(ctx context.Context, menuID string, items schema.MenuActions) error {
	for _, item := range items {
		item.ID = uuid.MustString()
		item.MenuID = menuID
		err := a.MenuActionModel.Create(ctx, *item)
		if err != nil {
			return err
		}

		for _, ritem := range item.Resources {
			ritem.ID = uuid.MustString()
			ritem.ActionID = item.ID
			err := a.MenuActionResourceModel.Create(ctx, *ritem)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *Menu) getParentPath(ctx context.Context, parentID string) (string, error) {
	if parentID == "" {
		return "", nil
	}

	pitem, err := a.MenuModel.Get(ctx, parentID)
	if err != nil {
		return "", err
	}
	if pitem == nil {
		return "", errors.ErrInvalidParent
	}

	return a.joinParentPath(pitem.ParentPath, pitem.ID), nil
}

func (a *Menu) joinParentPath(parent, id string) string {
	if parent != "" {
		return parent + "/" + id
	}
	return id
}

func (a *Menu) Update(ctx context.Context, id string, item schema.Menu) error {
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

	return a.TransModel.Exec(ctx, func(ctx context.Context) error {
		err := a.updateActions(ctx, id, oldItem.Actions, item.Actions)
		if err != nil {
			return err
		}

		err = a.updateChildParentPath(ctx, *oldItem, item)
		if err != nil {
			return err
		}

		return a.MenuModel.Update(ctx, id, item)
	})
}

func (a *Menu) updateActions(ctx context.Context, menuID string, oldItems, newItems schema.MenuActions) error {
	addActions, delActions, updateActions := a.compareActions(oldItems, newItems)

	err := a.createActions(ctx, menuID, addActions)
	if err != nil {
		return err
	}

	for _, item := range delActions {
		err := a.MenuActionModel.Delete(ctx, item.ID)
		if err != nil {
			return err
		}

		err = a.MenuActionResourceModel.DeleteByActionID(ctx, item.ID)
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
			err := a.MenuActionModel.Update(ctx, item.ID, *oitem)
			if err != nil {
				return err
			}
		}

		// update new and delete, not update
		addResources, delResources := a.compareResources(oitem.Resources, item.Resources)
		for _, aritem := range addResources {
			aritem.ID = uuid.MustString()
			aritem.ActionID = oitem.ID
			err := a.MenuActionResourceModel.Create(ctx, *aritem)
			if err != nil {
				return err
			}
		}

		for _, ditem := range delResources {
			err := a.MenuActionResourceModel.Delete(ctx, ditem.ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (a *Menu) compareActions(oldActions, newActions schema.MenuActions) (addList, delList, updateList schema.MenuActions) {
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

func (a *Menu) compareResources(oldResources, newResources schema.MenuActionResources) (addList, delList schema.MenuActionResources) {
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

func (a *Menu) updateChildParentPath(ctx context.Context, oldItem, newItem schema.Menu) error {
	if oldItem.ParentID == newItem.ParentID {
		return nil
	}

	opath := a.joinParentPath(oldItem.ParentPath, oldItem.ID)
	result, err := a.MenuModel.Query(contextx.NewNoTrans(ctx), schema.MenuQueryParam{
		PrefixParentPath: opath,
	})
	if err != nil {
		return err
	}

	npath := a.joinParentPath(newItem.ParentPath, newItem.ID)
	for _, menu := range result.Data {
		err = a.MenuModel.UpdateParentPath(ctx, menu.ID, npath+menu.ParentPath[len(opath):])
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Menu) Delete(ctx context.Context, id string) error {
	oldItem, err := a.MenuModel.Get(ctx, id)
	if err != nil {
		return err
	}
	if oldItem == nil {
		return errors.ErrNotFound
	}

	result, err := a.MenuModel.Query(ctx, schema.MenuQueryParam{
		PaginationParam: schema.PaginationParam{OnlyCount: true},
		ParentID:        &id,
	})
	if err != nil {
		return err
	}
	if result.PageResult.Total > 0 {
		return errors.ErrNotAllowDeleteWithChild
	}

	return a.TransModel.Exec(ctx, func(ctx context.Context) error {
		err = a.MenuActionResourceModel.DeleteByMenuID(ctx, id)
		if err != nil {
			return err
		}

		err := a.MenuActionModel.DeleteByMenuID(ctx, id)
		if err != nil {
			return err
		}

		return a.MenuModel.Delete(ctx, id)
	})
}

func (a *Menu) UpdateStatus(ctx context.Context, id string, status int) error {
	oldItem, err := a.MenuModel.Get(ctx, id)
	if err != nil {
		return err
	}
	if oldItem == nil {
		return errors.ErrNotFound
	}

	return a.MenuModel.UpdateStatus(ctx, id, status)
}
