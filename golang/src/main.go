package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Higakinn/festival-crawler/config"
	festival "github.com/Higakinn/festival-crawler/domain/festival"
	"github.com/Higakinn/festival-crawler/infra/notion"
	"github.com/Higakinn/festival-crawler/notification"
)

func main() {
	// 環境変数の読み込み
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config error")
	}
	// contextの定義
	ctx := context.Background()

	var repo festival.FestivalRepository
	repo = notion.NewFestivalRepository(ctx, cfg.NotionApiToken, cfg.NotionDBId)

	notificationService := &notification.NotificationService{}
	// XPluginを追加
	xPlugin := notification.NewXPlugin(cfg.XApiKey, cfg.XApiKeySecret, cfg.XApiAccessToken, cfg.XApiAccessTokenSecret)
	notificationService.AddNotificationPlugin(xPlugin)

	// TODO: LinePluginを追加
	// linePlugin := &n.LinePlugin{}
	// notificationService.AddNotificationPlugin(linePlugin)
	fmt.Println("find unposeted festival data")
	festivals, err := repo.FindUnPosted()
	if err != nil {
		log.Fatalf("")
	}
	for _, festival := range festivals {
		fmt.Println("festival data post")
		notificationService.Notify(ctx, &festival)
		repo.Save(festival)
	}

	fmt.Println("find unquoted festival data")

	afestivals, err := repo.FindUnQuoted()
	if err != nil {
		log.Fatalf("")
	}
	for _, festival := range afestivals {
		fmt.Println("festival data quote post")
		notificationService.Notify(ctx, &festival)
		repo.Save(festival)
	}
}
