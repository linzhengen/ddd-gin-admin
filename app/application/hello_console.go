package application

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/repository"
)

type HelloConsole interface {
	GetUserName(ctx context.Context, userId string) (string, error)
}

func NewHelloConsole(userRepo repository.UserRepository) HelloConsole {
	return &helloConsole{
		userRepo: userRepo,
	}
}

type helloConsole struct {
	userRepo repository.UserRepository
}

func (a *helloConsole) GetUserName(ctx context.Context, userId string) (string, error) {
	user, err := a.userRepo.Get(ctx, userId)
	if err != nil {
		return "", err
	}
	return user.UserName, nil
}
