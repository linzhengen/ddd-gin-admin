package mysql

import (
	"context"

	"github.com/jinzhu/gorm"

	"github.com/linzhengen/ddd-gin-admin/infrastructure/s"

	"github.com/linzhengen/ddd-gin-admin/infrastructure/hash"

	"github.com/linzhengen/ddd-gin-admin/infrastructure/persistence/mysql/model"

	"github.com/linzhengen/ddd-gin-admin/domain/repository"
	"github.com/linzhengen/ddd-gin-admin/domain/schema"
	"github.com/linzhengen/ddd-gin-admin/infrastructure/errors"
)

// UserRepository is struct.
type UserRepository struct {
	UserModel *model.User
}

// NewUserRepository constructor.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		UserModel: model.NewUser(db),
	}
}

var _ repository.UserRepository = &UserRepository{}

// Query ...
func (a *UserRepository) Query(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserQueryResult, error) {
	return a.UserModel.Query(ctx, params, opts...)
}

// QueryShow ...
func (a *UserRepository) QueryShow(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserShowQueryResult, error) {
	userResult, err := a.UserModel.Query(ctx, params, opts...)
	if err != nil {
		return nil, err
	} else if userResult == nil {
		return nil, nil
	}

	result := &schema.UserShowQueryResult{
		PageResult: userResult.PageResult,
	}
	if len(userResult.Data) == 0 {
		return result, nil
	}

	result.Data = userResult.Data.ToUserShows()
	return result, nil
}

// Get ...
func (a *UserRepository) Get(ctx context.Context, recordID string, opts ...schema.UserQueryOptions) (*schema.User, error) {
	item, err := a.UserModel.Get(ctx, recordID, opts...)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}
	return item, nil
}

func (a *UserRepository) checkUserName(ctx context.Context, userName string) error {
	result, err := a.UserModel.Query(ctx, schema.UserQueryParam{
		UserName: userName,
	}, schema.UserQueryOptions{
		PageParam: &schema.PaginationParam{PageSize: -1},
	})
	if err != nil {
		return err
	} else if result.PageResult.Total > 0 {
		return errors.New400Response("user name already exists")
	}
	return nil
}

func (a *UserRepository) getUpdate(ctx context.Context, recordID string) (*schema.User, error) {
	nitem, err := a.Get(ctx, recordID)
	if err != nil {
		return nil, err
	}
	return nitem, nil
}

// Create ...
func (a *UserRepository) Create(ctx context.Context, item schema.User) (*schema.User, error) {
	if item.Password == "" {
		return nil, errors.New400Response("密码不允许为空")
	}

	err := a.checkUserName(ctx, item.UserName)
	if err != nil {
		return nil, err
	}

	item.Password = hash.SHA1HashString(item.Password)
	item.RecordID = s.MustUUID()
	err = a.UserModel.Create(ctx, item)
	if err != nil {
		return nil, err
	}

	return a.getUpdate(ctx, item.RecordID)
}

// Update update user.
func (a *UserRepository) Update(ctx context.Context, recordID string, item schema.User) (*schema.User, error) {
	oldItem, err := a.UserModel.Get(ctx, recordID)
	if err != nil {
		return nil, err
	} else if oldItem == nil {
		return nil, errors.ErrNotFound
	} else if oldItem.UserName != item.UserName {
		err := a.checkUserName(ctx, item.UserName)
		if err != nil {
			return nil, err
		}
	}

	if item.Password != "" {
		item.Password = hash.SHA1HashString(item.Password)
	}

	err = a.UserModel.Update(ctx, recordID, item)
	if err != nil {
		return nil, err
	}

	return a.getUpdate(ctx, recordID)
}

// Delete delete user.
func (a *UserRepository) Delete(ctx context.Context, recordID string) error {
	oldItem, err := a.UserModel.Get(ctx, recordID)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	err = a.UserModel.Delete(ctx, recordID)
	if err != nil {
		return err
	}
	return nil
}

// UpdateStatus update user status.
func (a *UserRepository) UpdateStatus(ctx context.Context, recordID string, status int) error {
	oldItem, err := a.UserModel.Get(ctx, recordID)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	err = a.UserModel.UpdateStatus(ctx, recordID, status)
	if err != nil {
		return err
	}
	return nil
}
