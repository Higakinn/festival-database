package di

import "github.com/Higakinn/festival-crawler/app/cli_usecase"

type CLIUseCases struct {
	FestivalUseCase *cli_usecase.FestivalUseCase
}

func NewCLIUseCases(FestivalUseCase *cli_usecase.FestivalUseCase) *CLIUseCases {
	return &CLIUseCases{
		FestivalUseCase: FestivalUseCase,
	}
}
