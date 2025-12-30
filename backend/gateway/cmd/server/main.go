package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kakke18/platform-security-poc/backend/gateway/internal/config"
	"github.com/kakke18/platform-security-poc/backend/gateway/internal/server"
)

const (
	untilForcedTerminationTimeoutSec = 30
)

func main() {
	// 環境変数から設定を読み込む
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		panic(err)
	}

	// サーバーを作成
	srv, cleanup, err := server.New(cfg)
	if err != nil {
		slog.Error("Failed to initialize server", "error", err)
		panic(err)
	}
	defer cleanup()

	// サーバーを別のゴルーチンで起動
	go func() {
		slog.Info("Server is starting...", "port", cfg.Port)
		if err := srv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server failed to start", "error", err)
			panic(err)
		}
	}()

	// graceful shutdownのための処理
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	slog.Info("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), untilForcedTerminationTimeoutSec*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		return
	}

	slog.Info("Server exited properly")
}
