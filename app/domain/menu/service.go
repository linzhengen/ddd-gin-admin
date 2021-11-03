package menu

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/menu/menuactionresource"

	"github.com/linzhengen/ddd-gin-admin/app/domain/menu/menuaction"

	"github.com/linzhengen/ddd-gin-admin/app/domain/trans"

	"github.com/linzhengen/ddd-gin-admin/app/domain/pagination"

	"github.com/linzhengen/ddd-gin-admin/app/domain/errors"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/uuid"
)

type Service interface {
	Create(ctx context.Context, item *Menu) (string, error)
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

		err = a.updateChildParentPath(ctx, *oldItem, item)
		if err != nil {
			return err
		}

		return a.menuRepo.Update(ctx, id, *a.menuFactory.ToEntity(&item))
	})
}

func (a *service) checkName(ctx context.Context, item *Menu) error {
	_, pr, err := a.menuRepo.Query(ctx, QueryParam{
		PaginationParam: pagination.Param{
			OnlyCount: true,
		},
		ParentID: item.ParentID,
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
