package main

import (
	"fmt"
	"log"

	n "github.com/Higakinn/festival-crawler/notification"
	r "github.com/Higakinn/festival-crawler/repository"
)

func main() {
	// NotificationServiceのインスタンスを作成
	notificationService := &n.NotificationService{}

	// XPluginを追加
	xPlugin := &n.XPlugin{}
	notificationService.AddNotificationPlugin(xPlugin)

	// LinePluginを追加
	linePlugin := &n.LinePlugin{}
	notificationService.AddNotificationPlugin(linePlugin)

	// データの取得に使用するリポジトリを選択
	var repo r.FestivalRepository
	repo = &r.NotionDB{} // NotionDBを使用する場合
	// repo = &MySQL{} // MySQLを使用する場合

	// 情報を使ってデータを取得
	data, err := repo.FindNonNotifiedData()
	if err != nil {
		// エラーハンドリング
		log.Fatal(err)
	}

	// 取得したデータの処理
	for _, d := range data {
		fmt.Println(d)
		// 通知の実行
		notificationService.Notify("Hello, World!")
	}
}
