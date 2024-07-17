package notification

import "fmt"

// NotificationPlugin インターフェース
type NotificationPlugin interface {
	Notify(message string)
}

// NotificationService struct
type NotificationService struct {
	plugins []NotificationPlugin
}

// AddNotificationPlugin プラグインを追加する関数
func (n *NotificationService) AddNotificationPlugin(plugin NotificationPlugin) {
	n.plugins = append(n.plugins, plugin)
}

// Notify 全てのプラグインに通知を送信する関数
func (n *NotificationService) Notify(message string) {
	for _, plugin := range n.plugins {
		plugin.Notify(message)
	}
}

// XPlugin struct
type XPlugin struct {
	// Xに対する認証情報などの設定はここに含まれる場合もあります
}

// Notify Slackに通知を送る関数
func (s *XPlugin) Notify(message string) {
	// Slackへの通知処理を行う
	fmt.Println("X: ", message)
}

// LinePlugin struct
type LinePlugin struct {
	// Lineに対する認証情報などの設定はここに含まれる場合もあります
}

// Notify Lineに通知を送る関数
func (l *LinePlugin) Notify(message string) {
	// Lineへの通知処理を行う
	fmt.Println("Line: ", message)
}
