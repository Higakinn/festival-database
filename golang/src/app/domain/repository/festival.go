package repository

import (
	"context"
	"time"

	"github.com/Higakinn/festival-crawler/app/domain/models"
)

type FestivalRepository interface {
	FindByIsPost(ctx context.Context, isPost bool) ([]*models.Festival, error)
	FindByDate(ctx context.Context, date time.Time) ([]*models.Festival, error)
	Save(ctx context.Context, festival *models.Festival) error
}
