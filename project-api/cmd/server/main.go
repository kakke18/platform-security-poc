package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/kakke18/platform-security-poc2/project-api/pkg/gen/project/v1/projectv1connect"
	"github.com/kakke18/platform-security-poc2/project-api/internal/mock"
	"github.com/kakke18/platform-security-poc2/project-api/internal/service"
)

func main() {
	// モックストアを初期化
	store := mock.NewStore()

	// サービスを初期化
	projectService := service.NewProjectService(store)

	// HTTP ハンドラを設定
	mux := http.NewServeMux()
	path, handler := projectv1connect.NewProjectServiceHandler(projectService)
	mux.Handle(path, handler)

	// ヘルスチェックエンドポイント
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// ポート設定
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Project API server listening on %s", addr)

	// HTTP/2 サーバを起動 (h2c: HTTP/2 Cleartext)
	if err := http.ListenAndServe(
		addr,
		h2c.NewHandler(mux, &http2.Server{}),
	); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
