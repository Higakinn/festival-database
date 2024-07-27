//go:build wireinject

// +wireinject

package di

import (
	"github.com/Higakinn/festival-crawler/app/cli_usecase"
	"github.com/dstotijn/go-notion"

	"github.com/Higakinn/festival-crawler/app/infrastructure/repository"
	"github.com/Higakinn/festival-crawler/cmd"

	"github.com/google/wire"
)

func InitCLIUseCases(
	client *notion.Client,
	// config *config.Config,
	xClient *cmd.XClient,
	dbName string,
) *CLIUseCases {
	wire.Build(
		NewCLIUseCases,
		cli_usecase.NewFestivalUseCase,
		repository.NewFestivalRepository,
	)
	return nil
}
