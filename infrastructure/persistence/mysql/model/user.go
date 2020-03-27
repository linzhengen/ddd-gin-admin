package model

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/linzhengen/ddd-gin-admin/domain/entity"
	"github.com/linzhengen/ddd-gin-admin/domain/schema"
	"github.com/linzhengen/ddd-gin-admin/infrastructure/errors"
)

// NewUser constructor.
func NewUser(db *gorm.DB) *User {
	return &User{db}
}

// User ...
type User struct {
	db *gorm.DB
}

func (a *User) getQueryOption(opts ...schema.UserQueryOptions) schema.UserQueryOptions {
	var opt schema.UserQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query ...
func (a *User) Query(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserQueryResult, error) {
	db := entity.GetUserDB(ctx, a.db)
	if v := params.UserName; v != "" {
		db = db.Where("user_name=?", v)
	}
	if v := params.LikeUserName; v != "" {
		db = db.Where("user_name LIKE ?", "%"+v+"%")
	}
	if v := params.LikeRealName; v != "" {
		db = db.Where("real_name LIKE ?", "%"+v+"%")
	}
	if v := params.Status; v > 0 {
		db = db.Where("status=?", v)
	}
	db = db.Order("id DESC")

	opt := a.getQueryOption(opts...)
	var list entity.Users
	pr, err := WrapPageQuery(ctx, db, opt.PageParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	qr := &schema.UserQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaUsers(),
	}

	return qr, nil
}

// Get ...
func (a *User) Get(ctx context.Context, recordID string, opts ...schema.UserQueryOptions) (*schema.User, error) {
	var item entity.User
	ok, err := FindOne(ctx, entity.GetUserDB(ctx, a.db).Where("record_id=?", recordID), &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	sitem := item.ToSchemaUser()

	return sitem, nil
}

// Create ...
func (a *User) Create(ctx context.Context, item schema.User) error {
	return ExecTrans(ctx, a.db, func(ctx context.Context) error {
		sitem := entity.SchemaUser(item)
		result := entity.GetUserDB(ctx, a.db).Create(sitem.ToUser())
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}
		return nil
	})
}

// Update ...
func (a *User) Update(ctx context.Context, recordID string, item schema.User) error {
	return ExecTrans(ctx, a.db, func(ctx context.Context) error {
		sitem := entity.SchemaUser(item)
		omits := []string{"record_id", "creator"}
		if sitem.Password == "" {
			omits = append(omits, "password")
		}

		result := entity.GetUserDB(ctx, a.db).Where("record_id=?", recordID).Omit(omits...).Updates(sitem.ToUser())
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}
		return nil
	})
}

// Delete ...
func (a *User) Delete(ctx context.Context, recordID string) error {
	return ExecTrans(ctx, a.db, func(ctx context.Context) error {
		result := entity.GetUserDB(ctx, a.db).Where("record_id=?", recordID).Delete(entity.User{})
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}
		return nil
	})
}

// UpdateStatus ...
func (a *User) UpdateStatus(ctx context.Context, recordID string, status int) error {
	result := entity.GetUserDB(ctx, a.db).Where("record_id=?", recordID).Update("status", status)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// UpdatePassword ...
func (a *User) UpdatePassword(ctx context.Context, recordID, password string) error {
	result := entity.GetUserDB(ctx, a.db).Where("record_id=?", recordID).Update("password", password)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}
