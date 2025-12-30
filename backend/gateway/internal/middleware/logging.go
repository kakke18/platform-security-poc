package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// responseWriter はステータスコードを取得するためにhttp.ResponseWriterをラップする
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    int64
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.written += int64(n)
	return n, err
}

// AccessLog はHTTPリクエストをメソッド、パス、ステータス、時間、クライアント情報と共にログ出力する
func AccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// ステータスコードを取得するためにレスポンスライターをラップ
		wrapped := newResponseWriter(w)

		// リクエストを処理
		next.ServeHTTP(wrapped, r)

		// 処理時間を計算
		duration := time.Since(start)

		// クライアントIPを取得（X-Forwarded-Forヘッダーを考慮）
		clientIP := r.RemoteAddr
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			clientIP = xff
		} else if xri := r.Header.Get("X-Real-IP"); xri != "" {
			clientIP = xri
		}

		// アクセスログを出力
		slog.Info("access",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", wrapped.statusCode),
			slog.Duration("duration", duration),
			slog.String("client_ip", clientIP),
			slog.String("user_agent", r.UserAgent()),
			slog.Int64("bytes", wrapped.written),
		)
	})
}
