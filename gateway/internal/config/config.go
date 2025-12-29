package config

import (
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

	return &Config{
		IdentityAPIURL: identityAPIURL,
		Port:           port,
	}, nil
}
