package util

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
)

func SendGetHTTPRequestForBase64Image(sendURL string) (string, error) {
	// 画像URLにHTTPリクエストを投げる
	client := &http.Client{}
	req, err := http.NewRequest("GET", sendURL, nil)
	if err != nil {
		return "", err
	}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// レスポンスからbyteデータを取得
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	// byteデータからBase64を取得
	base64Data := base64.StdEncoding.EncodeToString(body)
	return base64Data, nil
}
