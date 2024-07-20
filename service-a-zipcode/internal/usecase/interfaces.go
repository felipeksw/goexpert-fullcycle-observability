package usecase

import (
	"context"

	"github.com/felipeksw/goexpert-fullcycle-observability/service-a-zipcode/internal/dto"
)

type Finder interface {
	Execute(ctx context.Context, zipcode string) (*dto.LocaleOutput, error)
}
