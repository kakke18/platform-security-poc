package server

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/kakke18/platform-security-poc/backend/gateway/gen/gateway/v1/gatewayv1connect"
	"github.com/kakke18/platform-security-poc/backend/gateway/internal/config"
	"github.com/kakke18/platform-security-poc/backend/gateway/internal/me"
	"github.com/kakke18/platform-security-poc/backend/gateway/internal/middleware"
	"github.com/rs/cors"
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
		return nil, err
	}

	// Identity APIのURLをパース
	identityURL, err := url.Parse(cfg.IdentityAPIURL)
	if err != nil {
		return nil, err
	}

	// リバースプロキシを作成
	identityProxy := httputil.NewSingleHostReverseProxy(identityURL)

	// ヘッダーを保持するようにプロキシをカスタマイズ
	identityProxy.Director = func(req *http.Request) {
		req.URL.Scheme = identityURL.Scheme
		req.URL.Host = identityURL.Host
		req.Host = identityURL.Host
		// 元のパスを保持
		// req.URL.Path はすでに設定されている
	}

	// Me APIハンドラーを初期化
	meHandler := me.NewHandler(cfg.IdentityAPIURL, cfg.UserAPIURL)

	// マルチプレクサを作成
	mux := http.NewServeMux()

	// MeServiceを登録（JWT検証付き）
	mePath, meConnectHandler := gatewayv1connect.NewMeServiceHandler(meHandler)
	mux.Handle(mePath, jwtMiddleware.Middleware(meConnectHandler))

	// /identity.v1.UserService/* をIdentity APIにルーティング（JWT検証付き）
	mux.Handle("/identity.v1.UserService/", jwtMiddleware.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		identityProxy.ServeHTTP(w, r)
	})))

	// ヘルスチェックエンドポイント
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// CORS設定 - 開発用にすべてのヘッダーを許可
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           86400, // 24時間
	})

	// ハンドラーチェーンを構築: AccessLog -> CORS -> mux
	handler := middleware.AccessLog(c.Handler(mux))

	httpServer := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: handler,
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
