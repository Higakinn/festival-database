package repository

import (
	"fmt"
)

// データモデルの定義
type Festival struct {
	ID     string
	Name   string
	Url    string
	Region string
	Date   string
}

// リポジトリインターフェース
type FestivalRepository interface {
	FindNonNotifiedData() ([]Festival, error)
}

// NotionDB リポジトリの実装
type NotionDB struct {
	// NotionDBへの接続情報やクライアントなどのフィールドがここに含まれます
}

// GetDataByConditions インターフェースの実装
func (n *NotionDB) FindNonNotifiedData() ([]Festival, error) {
	// NotionDBからのデータ取得処理を実装します
	// condition1とcondition2の条件に合致するデータを取得し、Dataのスライスとして返します
	// エラーハンドリングもここで行います
	fmt.Println("notion dbからデータ取得")
	return []Festival{{ID: "ID", Name: "hoge", Url: "https://example.com", Region: "愛媛県", Date: "2024-07-17 ~ 2024-07-18"}}, nil
}

// MySQL リポジトリの実装
// type MySQL struct {
// 	// MySQLへの接続情報やクライアントなどのフィールドがここに含まれます
// }

// // GetDataByConditions インターフェースの実装
// func (m *MySQL) GetDataByConditions(condition1 string, condition2 int) ([]Festival, error) {
// 	// MySQLからのデータ取得処理を実装します
// 	// condition1とcondition2の条件に合致するデータを取得し、Dataのスライスとして返します
// 	// エラーハンドリングもここで行います
// 	return []Festival{}, nil
// }
