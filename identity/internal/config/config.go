package config

import (
	"os"
)

const (
	defaultPort = "8081"
)

// Config はアプリケーション設定を保持する
type Config struct {
	// Port はサーバーのポート番号
	Port string
}

// Load は環境変数から設定を読み込む
func Load() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	return &Config{
		Port: port,
	}, nil
}
