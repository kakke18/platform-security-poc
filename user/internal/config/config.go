package config

import (
	"fmt"
	"os"
)

// Config はアプリケーション設定
type Config struct {
	Port string
}

// Load は環境変数から設定を読み込む
func Load() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	cfg := &Config{
		Port: port,
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

// validate は設定の妥当性を検証する
func (c *Config) validate() error {
	if c.Port == "" {
		return fmt.Errorf("PORT is required")
	}
	return nil
}
