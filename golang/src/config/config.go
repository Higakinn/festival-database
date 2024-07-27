package config

import (
	"github.com/caarlos0/env/v11"
)

// NOTION_API_TOKEN = os.getenv("NOTION_API_TOKEN")
// NOTION_DATABASE_ID = os.getenv("NOTION_DATABASE_ID")
// ## X関連の環境変数
// X_API_KEY = os.getenv("X_API_KEY")
// X_API_KEY_SECRET = os.getenv("X_API_KEY_SECRET")
// X_API_BEARER_TOKEN = os.getenv("X_API_BEARER_TOKEN")
// X_API_ACCESS_TOKEN = os.getenv("X_API_ACCESS_TOKEN")
// X_API_ACCESS_TOKEN_SECRET = os.getenv("X_API_ACCESS_TOKEN_SECRET")
type Config struct {
	NotionApiToken        string `env:"NOTION_API_TOKEN"`
	NotionDBId            string `env:"NOTION_DATABASE_ID" envDefault:"c8a42984e70b4f4b86930f6824f97450"`
	XApiKey               string `env:"X_API_KEY"`
	XApiKeySecret         string `env:"X_API_KEY_SECRET"`
	XApiBearerToken       string `env:"X_API_BEARER_TOKEN"`
	XApiAccessToken       string `env:"X_API_ACCESS_TOKEN"`
	XApiAccessTokenSecret string `env:"X_API_ACCESS_TOKEN_SECRET"`
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
