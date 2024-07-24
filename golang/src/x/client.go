package x

import (
	"context"
	"log"
	"net/http"

	"github.com/ChimeraCoder/anaconda"
	"github.com/Higakinn/festival-crawler/util"
	"github.com/dghubble/oauth1"
	twitter "github.com/g8rswimmer/go-twitter/v2"
)

const (
	X_API_KEY                 = ""
	X_API_KEY_SECRET          = ""
	X_API_BEARER_TOKEN        = ""
	X_API_ACCESS_TOKEN        = ""
	X_API_ACCESS_TOKEN_SECRET = ""
)

type XClient struct {
	Client *twitter.Client
	//画像upload用
	Api *anaconda.TwitterApi
}

type authorizer struct{}

func (a *authorizer) Add(req *http.Request) {}

func NewXClient() *XClient {
	config := oauth1.NewConfig(X_API_KEY, X_API_KEY_SECRET)
	httpClient := config.Client(oauth1.NoContext, &oauth1.Token{
		Token:       X_API_ACCESS_TOKEN,
		TokenSecret: X_API_ACCESS_TOKEN_SECRET,
	})

	return &XClient{
		Client: &twitter.Client{
			Authorizer: &authorizer{},
			Client:     httpClient,
			Host:       "https://api.twitter.com",
		},
		Api: anaconda.NewTwitterApiWithCredentials(X_API_ACCESS_TOKEN, X_API_ACCESS_TOKEN_SECRET, X_API_KEY, X_API_KEY_SECRET),
	}
}

func (xc *XClient) Post(ctx context.Context, text string, imgUrl string) (string, error) {
	var req twitter.CreateTweetRequest
	if imgUrl == "" {
		req = twitter.CreateTweetRequest{
			Text: text,
		}
	} else {
		img, err := util.SendGetHTTPRequestForBase64Image(imgUrl)
		if err != nil {
			log.Fatal(err)
			return "", err
		}
		media, err := xc.Api.UploadMedia(img)
		if err != nil {
			log.Fatal(err)
			return "", err
		}
		req = twitter.CreateTweetRequest{
			Text: text,
			Media: &twitter.CreateTweetMedia{
				IDs: []string{
					media.MediaIDString,
				},
			},
		}
	}
	tweetResponse, err := xc.Client.CreateTweet(ctx, req)
	if err != nil {
		log.Panicf("create tweet error: %v", err)
		return "", err
	}
	return tweetResponse.Tweet.ID, nil
}
