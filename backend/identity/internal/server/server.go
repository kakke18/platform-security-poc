package server

import (
	"context"
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/kakke18/platform-security-poc/backend/identity/gen/identity/v1/identityv1connect"
	"github.com/kakke18/platform-security-poc/backend/identity/internal/config"
	"github.com/kakke18/platform-security-poc/backend/identity/internal/middleware"
	"github.com/kakke18/platform-security-poc/backend/identity/internal/user"
	"github.com/kakke18/platform-security-poc/backend/identity/internal/workspace"
	"github.com/kakke18/platform-security-poc/backend/identity/internal/workspaceuser"
)

// Server はHTTPサーバーを表す
type Server struct {
	config     *config.Config
	httpServer *http.Server
}

// New は新しいサーバーを作成する
func New(cfg *config.Config) (*Server, error) {
	// ユーザー機能を初期化
	userRepo := user.NewMockRepository()
	userHandler := user.NewHandler(userRepo)

	// Workspace機能を初期化
	workspaceRepo := workspace.NewMockRepository()

	// WorkspaceUser機能を初期化
	workspaceUserRepo := workspaceuser.NewMockRepository()
	workspaceUserHandler := workspaceuser.NewHandler(workspaceUserRepo, workspaceRepo)

	// マルチプレクサを作成
	mux := http.NewServeMux()

	// UserServiceを登録（Gatewayで認証済み）
	userPath, userConnectHandler := identityv1connect.NewUserServiceHandler(userHandler)
	mux.Handle(userPath, userConnectHandler)

	// WorkspaceUserServiceを登録（Gatewayで認証済み）
	workspaceUserPath, workspaceUserConnectHandler := identityv1connect.NewWorkspaceUserServiceHandler(workspaceUserHandler)
	mux.Handle(workspaceUserPath, workspaceUserConnectHandler)

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
