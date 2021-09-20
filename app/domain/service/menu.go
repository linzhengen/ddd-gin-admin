package service

import (
	"context"
	"fmt"
	"os"

	"github.com/linzhengen/ddd-gin-admin/app/domain/factory"

	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/contextx"

	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/errors"

	"github.com/linzhengen/ddd-gin-admin/app/domain/repository"
	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/schema"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/uuid"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/yaml"
)

type Menu interface {
	InitData(ctx context.Context, dataFile string) error
	Query(ctx context.Context, params schema.MenuQueryParam) (*schema.MenuQueryResult, error)
	Get(ctx context.Context, id string) (*schema.Menu, error)
	QueryActions(ctx context.Context, id string) (schema.MenuActions, error)
	Create(ctx context.Context, item schema.Menu) (*schema.IDResult, error)
	Update(ctx context.Context, id string, item schema.Menu) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status int) error
}

func NewMenu(
	transRepo repository.TransRepository,
	menuRepo repository.MenuRepository,
	menuActionRepo repository.MenuActionRepository,
	menuActionResourceRepo repository.MenuActionResourceRepository,
	menuFactory factory.Menu,
	menuActionFactory factory.MenuAction,
	menuActionResourceFactory factory.MenuActionResource,
) Menu {
	return &menu{
		transRepo:                 transRepo,
		menuRepo:                  menuRepo,
		menuActionRepo:            menuActionRepo,
		menuActionResourceRepo:    menuActionResourceRepo,
		menuFactory:               menuFactory,
		menuActionFactory:         menuActionFactory,
		menuActionResourceFactory: menuActionResourceFactory,
	}
}

type menu struct {
	transRepo                 repository.TransRepository
	menuRepo                  repository.MenuRepository
	menuActionRepo            repository.MenuActionRepository
	menuActionResourceRepo    repository.MenuActionResourceRepository
	menuFactory               factory.Menu
	menuActionFactory         factory.MenuAction
	menuActionResourceFactory factory.MenuActionResource
}

func (a *menu) InitData(ctx context.Context, dataFile string) error {
	_, pr, err := a.menuRepo.Query(ctx, schema.MenuQueryParam{
		PaginationParam: schema.PaginationParam{OnlyCount: true},
	})
	if err != nil {
		return err
	}
	if pr.Total > 0 {
		return nil
	}

	data, err := a.readData(dataFile)
	if err != nil {
		return err
	}

	return a.createMenus(ctx, "", data)
}

func (a *menu) readData(name string) (schema.MenuTrees, error) {
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

func (a *menu) createMenus(ctx context.Context, parentID string, list schema.MenuTrees) error {
	return a.transRepo.Exec(ctx, func(ctx context.Context) error {
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

func (a *menu) Query(ctx context.Context, params schema.MenuQueryParam) (*schema.MenuQueryResult, error) {
	menuActionResult, pr, err := a.menuActionRepo.Query(ctx, schema.MenuActionQueryParam{})
	if err != nil {
		return nil, err
	}

	result, _, err := a.menuRepo.Query(ctx, params)
	if err != nil {
		return nil, err
	}
	return &schema.MenuQueryResult{
		Data:       a.menuFactory.ToSchemaList(result).FillMenuAction(a.menuActionFactory.ToSchemaList(menuActionResult).ToMenuIDMap()),
		PageResult: pr,
	}, nil
}

func (a *menu) Get(ctx context.Context, id string) (*schema.Menu, error) {
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
	menuSchema := a.menuFactory.ToSchema(item)
	menuSchema.Actions = actions

	return menuSchema, nil
}

func (a *menu) QueryActions(ctx context.Context, id string) (schema.MenuActions, error) {
	result, _, err := a.menuActionRepo.Query(ctx, schema.MenuActionQueryParam{
		MenuID: id,
	})
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}

	resourceResult, _, err := a.menuActionResourceRepo.Query(ctx, schema.MenuActionResourceQueryParam{
		MenuID: id,
	})
	if err != nil {
		return nil, err
	}

	menuActions := a.menuActionFactory.ToSchemaList(result)
	menuActions.FillResources(a.menuActionResourceFactory.ToSchemaList(resourceResult).ToActionIDMap())

	return menuActions, nil
}

func (a *menu) checkName(ctx context.Context, item schema.Menu) error {
	_, pr, err := a.menuRepo.Query(ctx, schema.MenuQueryParam{
		PaginationParam: schema.PaginationParam{
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

func (a *menu) Create(ctx context.Context, item schema.Menu) (*schema.IDResult, error) {
	if err := a.checkName(ctx, item); err != nil {
		return nil, err
	}

	parentPath, err := a.getParentPath(ctx, item.ParentID)
	if err != nil {
		return nil, err
	}
	item.ParentPath = parentPath
	item.ID = uuid.MustString()

	err = a.transRepo.Exec(ctx, func(ctx context.Context) error {
		err := a.createActions(ctx, item.ID, item.Actions)
		if err != nil {
			return err
		}

		return a.menuRepo.Create(ctx, *a.menuFactory.ToEntity(&item))
	})
	if err != nil {
		return nil, err
	}

	return schema.NewIDResult(item.ID), nil
}

func (a *menu) createActions(ctx context.Context, menuID string, items schema.MenuActions) error {
	for _, item := range items {
		item.ID = uuid.MustString()
		item.MenuID = menuID
		err := a.menuActionRepo.Create(ctx, *a.menuActionFactory.ToEntity(item))
		if err != nil {
			return err
		}

		for _, ritem := range item.Resources {
			ritem.ID = uuid.MustString()
			ritem.ActionID = item.ID
			err := a.menuActionResourceRepo.Create(ctx, *a.menuActionResourceFactory.ToEntity(ritem))
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
		return "", errors.ErrInvalidParent
	}

	return a.joinParentPath(a.menuFactory.ToSchema(pitem).ParentPath, pitem.ID), nil
}

func (a *menu) joinParentPath(parent, id string) string {
	if parent != "" {
		return parent + "/" + id
	}
	return id
}

func (a *menu) Update(ctx context.Context, id string, item schema.Menu) error {
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

		err = a.updateChildParentPath(ctx, *oldItem, item)
		if err != nil {
			return err
		}

		return a.menuRepo.Update(ctx, id, *a.menuFactory.ToEntity(&item))
	})
}

func (a *menu) updateActions(ctx context.Context, menuID string, oldItems, newItems schema.MenuActions) error {
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
			err := a.menuActionRepo.Update(ctx, item.ID, *a.menuActionFactory.ToEntity(oitem))
			if err != nil {
				return err
			}
		}

		// update new and delete, not update
		addResources, delResources := a.compareResources(oitem.Resources, item.Resources)
		for _, aritem := range addResources {
			aritem.ID = uuid.MustString()
			aritem.ActionID = oitem.ID
			err := a.menuActionResourceRepo.Create(ctx, *a.menuActionResourceFactory.ToEntity(aritem))
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

func (a *menu) compareActions(oldActions, newActions schema.MenuActions) (addList, delList, updateList schema.MenuActions) {
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

func (a *menu) compareResources(oldResources, newResources schema.MenuActionResources) (addList, delList schema.MenuActionResources) {
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

func (a *menu) updateChildParentPath(ctx context.Context, oldItem, newItem schema.Menu) error {
	if oldItem.ParentID == newItem.ParentID {
		return nil
	}

	opath := a.joinParentPath(oldItem.ParentPath, oldItem.ID)
	result, _, err := a.menuRepo.Query(contextx.NewNoTrans(ctx), schema.MenuQueryParam{
		PrefixParentPath: opath,
	})
	if err != nil {
		return err
	}

	npath := a.joinParentPath(newItem.ParentPath, newItem.ID)
	for _, menu := range result {
		parentPath := *(menu.ParentPath)
		err = a.menuRepo.UpdateParentPath(ctx, menu.ID, fmt.Sprintf("%s%s", npath, parentPath[len(opath):]))
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
		return errors.ErrNotFound
	}

	_, pr, err := a.menuRepo.Query(ctx, schema.MenuQueryParam{
		PaginationParam: schema.PaginationParam{OnlyCount: true},
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

func (a *menu) UpdateStatus(ctx context.Context, id string, status int) error {
	oldItem, err := a.menuRepo.Get(ctx, id)
	if err != nil {
		return err
	}
	if oldItem == nil {
		return errors.ErrNotFound
	}

	return a.menuRepo.UpdateStatus(ctx, id, status)
}
