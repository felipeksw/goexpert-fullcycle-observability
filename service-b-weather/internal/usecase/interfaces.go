package usecase

import "context"

type Finder interface {
	Execute(ctx context.Context, query string) (interface{}, error)
}
