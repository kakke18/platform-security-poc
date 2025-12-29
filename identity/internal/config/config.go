package config

import (
	"fmt"
	"os"
)

const (
	defaultPort = "8081"
)

// Config はアプリケーション設定を保持する
type Config struct {
	// Auth0Domain はAuth0のドメイン (例: your-tenant.auth0.com)
	Auth0Domain string

	// Auth0Audience はAuth0のAPI識別子
	Auth0Audience string

	// Port はサーバーのポート番号
	Port string
}

// Load は環境変数から設定を読み込む
func Load() (*Config, error) {
	auth0Domain := os.Getenv("AUTH0_DOMAIN")
	if auth0Domain == "" {
		return nil, fmt.Errorf("AUTH0_DOMAIN must be set")
	}

	auth0Audience := os.Getenv("AUTH0_AUDIENCE")
	if auth0Audience == "" {
		return nil, fmt.Errorf("AUTH0_AUDIENCE must be set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	return &Config{
		Auth0Domain:   auth0Domain,
		Auth0Audience: auth0Audience,
		Port:          port,
	}, nil
}
