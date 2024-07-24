package main

import (
	"context"
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

	festivals, err := repo.FindUnPosted()
	if err != nil {
		log.Fatalf("")
	}
	for _, festival := range festivals {
		notificationService.Notify(ctx, &festival)
		repo.Save(festival)
	}

	afestivals, err := repo.FindUnQuoted()
	if err != nil {
		log.Fatalf("")
	}
	for _, festival := range afestivals {
		notificationService.Notify(ctx, &festival)
		repo.Save(festival)
	}
}
