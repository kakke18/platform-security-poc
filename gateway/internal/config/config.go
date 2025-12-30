package config

import (
	"fmt"
	"os"
)

const (
	defaultIdentityAPIURL = "http://localhost:8081"
	defaultPort           = "8080"
)

// Config はアプリケーション設定を保持する
type Config struct {
	// IdentityAPIURL はIdentity APIのベースURL
	IdentityAPIURL string

	// Port はサーバーのポート番号
	Port string

	// Auth0Domain はAuth0のドメイン
	Auth0Domain string

	// Auth0Audience はAuth0のオーディエンス
	Auth0Audience string
}

// Load は環境変数から設定を読み込む
func Load() (*Config, error) {
	identityAPIURL := os.Getenv("IDENTITY_API_URL")
	if identityAPIURL == "" {
		identityAPIURL = defaultIdentityAPIURL
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	auth0Domain := os.Getenv("AUTH0_DOMAIN")
	if auth0Domain == "" {
		return nil, fmt.Errorf("AUTH0_DOMAIN must be set")
	}

	auth0Audience := os.Getenv("AUTH0_AUDIENCE")
	if auth0Audience == "" {
		return nil, fmt.Errorf("AUTH0_AUDIENCE must be set")
	}

	return &Config{
		IdentityAPIURL: identityAPIURL,
		Port:           port,
		Auth0Domain:    auth0Domain,
		Auth0Audience:  auth0Audience,
	}, nil
}
