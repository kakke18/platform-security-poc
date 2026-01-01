package server

import (
	"context"
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/kakke18/platform-security-poc/backend/gen/user/v1/userv1connect"
	"github.com/kakke18/platform-security-poc/backend/user/internal/config"
	"github.com/kakke18/platform-security-poc/backend/user/internal/middleware"
	"github.com/kakke18/platform-security-poc/backend/user/internal/tenantuser"
)

// Server はHTTPサーバーを表す
type Server struct {
	config     *config.Config
	httpServer *http.Server
}

// New は新しいサーバーを作成する
func New(cfg *config.Config) (*Server, error) {
	// TenantUser機能を初期化
	tenantUserRepo := tenantuser.NewMockRepository()
	tenantUserHandler := tenantuser.NewHandler(tenantUserRepo)

	// マルチプレクサを作成
	mux := http.NewServeMux()

	// TenantUserServiceを登録（Gatewayで認証済み）
	tenantUserPath, tenantUserConnectHandler := userv1connect.NewTenantUserServiceHandler(tenantUserHandler)
	mux.Handle(tenantUserPath, tenantUserConnectHandler)

	// ヘルスチェックエンドポイント
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// ハンドラーチェーンを構築: h2c -> AccessLog -> mux
	// h2cハンドラーを最外層に配置することで、HTTP/2接続の確立を優先
	finalHandler := h2c.NewHandler(middleware.AccessLog(mux), &http2.Server{})

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
