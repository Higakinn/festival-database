package main

// 実装参考：　https://hatappi.blog/entry/2017/12/24/134451

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Higakinn/festival-crawler/cmd"
	"github.com/Higakinn/festival-crawler/cmd/di"
	"github.com/Higakinn/festival-crawler/config"
	"github.com/Higakinn/festival-crawler/pkg/x"
	"github.com/dstotijn/go-notion"
	"github.com/urfave/cli"
)

func main() {
	// 環境変数の読み込み
	ctx := context.Background()
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config error")
	}
	var (
		suffix string
	)

	xClient := cmd.XClient{
		Client: *x.NewXClient(cfg.XApiKey, cfg.XApiKeySecret, cfg.XApiAccessToken, cfg.XApiAccessTokenSecret),
	}

	notionClient := notion.NewClient(cfg.NotionApiToken)
	// di.InitCLIUseCases()
	cli_usecase := di.InitCLIUseCases(notionClient, &xClient, cfg.NotionDBId)

	{
		app := cli.NewApp()
		app.Name = "festival cli"
		app.Usage = "祭禮情報の操作を提供します"
		app.Version = "0.1.0"

		app.Flags = []cli.Flag{
			cli.StringFlag{
				Name:        "suffix, s",
				Value:       "!!!",
				Usage:       "text after speaking something",
				Destination: &suffix,
				EnvVar:      "SUFFIX",
			},
		}

		app.Commands = []cli.Command{
			{
				Name:  "get",
				Usage: "festival controll",
				Subcommands: []cli.Command{
					{
						Name:  "unposted",
						Usage: "get unposted festival data",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:  "dry-run",
								Usage: "if dry-run==false なら ダミーデータ取得",
							},
						},
						Action: func(c *cli.Context) error {
							err := cli_usecase.FestivalUseCase.UnposetedList(ctx, c.Bool("dry-run"))
							if err != nil {
								return err
							}
							return nil
						},
					},
					{
						Name:  "today",
						Usage: "今日から始めるの祭りを取得します",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:  "dry-run",
								Usage: "if dry-run==false なら ダミーデータ取得",
							},
						},
						Action: func(c *cli.Context) error {
							fmt.Println(c.Bool("dry-run"))
							err := cli_usecase.FestivalUseCase.HoldTodayList(ctx, c.Bool("dry-run"))
							if err != nil {
								return err
							}
							return nil
						},
					},
				},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "dry-run",
						Usage: "if dry-run==false なら ダミーデータ取得",
					},
				},
				Action: func(c *cli.Context) error {
					return nil
				},
			},
			{
				Name:  "notify",
				Usage: "festival controll",
				Subcommands: []cli.Command{
					{
						Name:  "unposted",
						Usage: "notify unposted festival data",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:  "dry-run",
								Usage: "if dry-run==false なら ダミーデータ取得",
							},
						},
						Action: func(c *cli.Context) error {
							fmt.Println(c.Bool("dry-run"))
							err := cli_usecase.FestivalUseCase.NofityUnposetedList(ctx, c.Bool("dry-run"))
							if err != nil {
								return err
							}
							return nil
						},
					},
					{
						Name:  "today",
						Usage: "今日から始めるの祭りを通知します",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:  "dry-run",
								Usage: "if dry-run==false なら ダミーデータ取得",
							},
						},
						Action: func(c *cli.Context) error {
							fmt.Println(c.Bool("dry-run"))
							err := cli_usecase.FestivalUseCase.NofityHoldTodayList(ctx, c.Bool("dry-run"))
							if err != nil {
								return err
							}
							return nil
						},
					},
				},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "dry-run",
						Usage: "if dry-run==false なら ダミーデータ取得",
					},
				},
				Action: func(c *cli.Context) error {
					return nil
				},
			},
		}
		app.Run(os.Args)
	}
}
