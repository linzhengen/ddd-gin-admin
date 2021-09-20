package repository

import (
	"context"
)

type TransRepository interface {
	Exec(ctx context.Context, fn func(context.Context) error) error
}
