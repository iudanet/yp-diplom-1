package repo

import "context"

type Repositories interface {
	Migrate(ctx context.Context) error
}
