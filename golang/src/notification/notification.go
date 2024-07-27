// package notification

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"time"

// 	festival "github.com/Higakinn/festival-crawler/domain/festival"
// 	"github.com/Higakinn/festival-crawler/pkg/x"
// )

// // NotificationPlugin インターフェース
// type NotificationPlugin interface {
// 	Notify(ctx context.Context, festival *festival.Festival)
// }

// // NotificationService struct
// type NotificationService struct {
// 	plugins []NotificationPlugin
// }

// // AddNotificationPlugin プラグインを追加する関数
// func (n *NotificationService) AddNotificationPlugin(plugin NotificationPlugin) {
// 	n.plugins = append(n.plugins, plugin)
// }

// // Notify 全てのプラグインに通知を送信する関数
// func (n *NotificationService) Notify(ctx context.Context, content *festival.Festival) {
// 	for _, plugin := range n.plugins {
// 		plugin.Notify(ctx, content)
// 	}
// }

// // XPlugin struct
// type XPlugin struct {
// 	// Xに対する認証情報などの設定はここに含まれる場合もあります
// 	client x.XClient
// }

// // Notify implements NotificationPlugin.
// func (s XPlugin) Notify(ctx context.Context, festival *festival.Festival) {
// 	// TODO: 通知処理にいれる内容ではないので、外だしする。
// 	// 通知処理は内容をそのまま該当のプラットフォームに通知する処理だけが書かれているべき。

// 	// 投稿済み(postidが空ではない) かつ 引用投稿されていない(repostidが空)　の場合に引用リポストを行う
// 	if festival.PostId != "" && festival.RepostId == "" {
// 		fmt.Println("x quoted post")
// 		c := genQuoteRepostContent(*festival)
// 		repostId := ""
// 		repostId, err := s.client.Post(ctx, c, "")
// 		if err != nil {
// 			log.Fatal(err)
// 			return
// 		}
// 		festival.RepostId = repostId
// 		return
// 	}

// 	// 　通常のポストを行う
// 	fmt.Println(festival.PosterUrl)
// 	fmt.Println("x post")
// 	content := genPostContent(*festival)
// 	fmt.Println(content)
// 	postId, err := s.client.Post(ctx, content, festival.PosterUrl)
// 	if err != nil {
// 		log.Fatal(err)
// 		return
// 	}
// 	festival.PostId = postId
// }

// func NewXPlugin(XApiKey string, XApiKeySecret string, XApiAcessToken string, XApiAcessTokenSecret string) XPlugin {
// 	return XPlugin{client: *x.NewXClient(XApiKey, XApiKeySecret, XApiAcessToken, XApiAcessTokenSecret)}
// }

// // Notify Slackに通知を送る関数
// // func (s *XPlugin) Notify(ctx context.Context, content *festival.Festival) {
// // 	// Slackへの通知処理を行う
// // 	// content.Id = "1815109059567292703"

// // 	// s.client.PostWithImg()
// // 	c := content.PostContent()
// // 	// fmt.Println(c)
// // 	s.client.Post(ctx, c)
// // 	content.PostId = "huga"
// // 	fmt.Println("X: ", content)
// // }

// // LinePlugin struct
// type LinePlugin struct {
// 	// Lineに対する認証情報などの設定はここに含まれる場合もあります
// }

// // Notify Lineに通知を送る関数
// func (l *LinePlugin) Notify(ctx context.Context, content *festival.Festival) {
// 	// Lineへの通知処理を行う
// 	fmt.Println("Line: ", content)
// }

// func genPostContent(festival festival.Festival) string {
// 	date := festival.StartDate.Format(time.DateOnly) + " ~ " + festival.EndDate.Format(time.DateOnly)
// 	if festival.StartDate == festival.EndDate {
// 		date = festival.StartDate.Format(time.DateOnly)
// 	}
// 	return fmt.Sprintf(`【🏮祭り情報🏮】
// #%s

// ■ 開催期間
// ・%s

// ■ 開催場所
// ・%s

// ■ アクセス
// ・%s
// ■ 参考
// %s
// `,
// 		festival.Name, date, festival.Region, festival.Access, festival.Url)
// }
// func genQuoteRepostContent(festival festival.Festival) string {
// 	date := festival.StartDate.Format(time.DateOnly) + " ~ " + festival.EndDate.Format(time.DateOnly)
// 	if festival.StartDate == festival.EndDate {
// 		date = festival.StartDate.Format(time.DateOnly)
// 	}
// 	return fmt.Sprintf(`【#%s】
// #%s 始まります！

// ■ 開催期間
// ・%s

// ■ アクセス
// ・%s

// %s
// `,
// 		festival.Region, festival.Name, date, festival.Access, festival.XUrl)
// }
