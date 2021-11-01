package trans

import "context"

type Repository interface {
	Exec(ctx context.Context, fn func(context.Context) error) error
}
