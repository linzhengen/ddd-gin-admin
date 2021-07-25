package application

import (
	"context"

	"github.com/google/wire"
	repo "github.com/linzhengen/ddd-gin-admin/domain/repository"
	"github.com/linzhengen/ddd-gin-admin/domain/schema"
	"github.com/linzhengen/ddd-gin-admin/pkg/errors"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/uuid"
)

var DemoSet = wire.NewSet(wire.Struct(new(Demo), "*"))

type Demo struct {
	DemoModel *repo.Demo
}

func (a *Demo) Query(ctx context.Context, params schema.DemoQueryParam, opts ...schema.DemoQueryOptions) (*schema.DemoQueryResult, error) {
	return a.DemoModel.Query(ctx, params, opts...)
}

func (a *Demo) Get(ctx context.Context, id string, opts ...schema.DemoQueryOptions) (*schema.Demo, error) {
	item, err := a.DemoModel.Get(ctx, id, opts...)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}

	return item, nil
}

func (a *Demo) checkCode(ctx context.Context, code string) error {
	result, err := a.DemoModel.Query(ctx, schema.DemoQueryParam{
		PaginationParam: schema.PaginationParam{
			OnlyCount: true,
		},
		Code: code,
	})
	if err != nil {
		return err
	} else if result.PageResult.Total > 0 {
		return errors.New400Response("编号已经存在")
	}

	return nil
}

func (a *Demo) Create(ctx context.Context, item schema.Demo) (*schema.IDResult, error) {
	err := a.checkCode(ctx, item.Code)
	if err != nil {
		return nil, err
	}

	item.ID = uuid.MustString()
	err = a.DemoModel.Create(ctx, item)
	if err != nil {
		return nil, err
	}

	return schema.NewIDResult(item.ID), nil
}

func (a *Demo) Update(ctx context.Context, id string, item schema.Demo) error {
	oldItem, err := a.DemoModel.Get(ctx, id)
	if err != nil {
		return err
	}
	if oldItem == nil {
		return errors.ErrNotFound
	}
	if oldItem.Code != item.Code {
		if err := a.checkCode(ctx, item.Code); err != nil {
			return err
		}
	}
	item.ID = oldItem.ID
	item.Creator = oldItem.Creator
	item.CreatedAt = oldItem.CreatedAt

	return a.DemoModel.Update(ctx, id, item)
}

func (a *Demo) Delete(ctx context.Context, id string) error {
	oldItem, err := a.DemoModel.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	return a.DemoModel.Delete(ctx, id)
}

func (a *Demo) UpdateStatus(ctx context.Context, id string, status int) error {
	oldItem, err := a.DemoModel.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	return a.DemoModel.UpdateStatus(ctx, id, status)
}
