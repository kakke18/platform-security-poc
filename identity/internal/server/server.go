package server

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/kakke18/platform-security-poc/identity/gen/identity/v1/identityv1connect"
	"github.com/kakke18/platform-security-poc/identity/internal/config"
	"github.com/kakke18/platform-security-poc/identity/internal/middleware"
	"github.com/kakke18/platform-security-poc/identity/internal/user"
)

// Server はHTTPサーバーを表す
type Server struct {
	config     *config.Config
	httpServer *http.Server
}

// New は新しいサーバーを作成する
func New(cfg *config.Config) (*Server, error) {
	// JWTミドルウェアを初期化
	jwtMiddleware, err := middleware.NewJWTMiddleware(cfg.Auth0Domain, cfg.Auth0Audience)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize JWT middleware: %w", err)
	}

	// ユーザー機能を初期化
	userRepo := user.NewMockRepository()
	userHandler := user.NewHandler(userRepo)

	// マルチプレクサを作成
	mux := http.NewServeMux()

	// UserServiceをJWTミドルウェアと共に登録
	userPath, userConnectHandler := identityv1connect.NewUserServiceHandler(userHandler)
	mux.Handle(userPath, jwtMiddleware.Middleware(userConnectHandler))

	// ヘルスチェックエンドポイント
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// ハンドラーチェーンを構築: AccessLog -> h2c -> mux
	finalHandler := middleware.AccessLog(h2c.NewHandler(mux, &http2.Server{}))

	httpServer := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: finalHandler,
	}

	return &Server{
		config:     cfg,
		httpServer: httpServer,
	}, nil
}

// Run はサーバーを起動する
func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown はサーバーをグレースフルシャットダウンする
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
