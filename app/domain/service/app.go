package service

import (
	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/response"

	"github.com/linzhengen/ddd-gin-admin/configs"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/hash"
)

func getRootUser() *response.User {
	user := configs.C.Root
	return &response.User{
		ID:       user.UserName,
		UserName: user.UserName,
		RealName: user.RealName,
		Password: hash.MD5String(user.Password),
	}
}

func checkIsRootUser(userID string) bool {
	return getRootUser().ID == userID
}
