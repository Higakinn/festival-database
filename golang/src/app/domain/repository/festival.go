package repository

import (
	"context"

	"github.com/Higakinn/festival-crawler/app/domain/models"
)

type FestivalRepository interface {
	FindUnPosted(ctx context.Context) ([]*models.Festival, error)
	FindUnQuoted(ctx context.Context) ([]*models.Festival, error)
	Save(ctx context.Context, festival *models.Festival) error
}
