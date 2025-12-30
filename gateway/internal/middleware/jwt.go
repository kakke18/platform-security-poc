package middleware

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWKSKey はAuth0のJWKSから取得した単一の鍵を表す
type JWKSKey struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

// JWKS はAuth0のJSON Web Key Setを表す
type JWKS struct {
	Keys []JWKSKey `json:"keys"`
}

// JWTClaims はJWTのクレームを表す
type JWTClaims struct {
	jwt.RegisteredClaims
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

// JWTMiddleware はJWT検証ミドルウェアを提供する
type JWTMiddleware struct {
	domain    string
	audience  string
	keys      map[string]*rsa.PublicKey
	keysMu    sync.RWMutex
	lastFetch time.Time
}

// NewJWTMiddleware は新しいJWTミドルウェアを作成する
func NewJWTMiddleware(domain, audience string) (*JWTMiddleware, error) {
	m := &JWTMiddleware{
		domain:   domain,
		audience: audience,
		keys:     make(map[string]*rsa.PublicKey),
	}

	// 初期化時にJWKS鍵を取得
	if err := m.fetchJWKS(); err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}

	return m, nil
}

// fetchJWKS はAuth0からJWKSを取得する
func (m *JWTMiddleware) fetchJWKS() error {
	jwksURL := fmt.Sprintf("https://%s/.well-known/jwks.json", m.domain)

	resp, err := http.Get(jwksURL)
	if err != nil {
		return fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch JWKS: status=%d", resp.StatusCode)
	}

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return fmt.Errorf("failed to decode JWKS: %w", err)
	}

	m.keysMu.Lock()
	defer m.keysMu.Unlock()

	// JWKSをRSA公開鍵に変換
	for _, key := range jwks.Keys {
		if key.Kty != "RSA" || key.Use != "sig" {
			continue
		}

		// N（モジュラス）をデコード
		nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
		if err != nil {
			slog.Warn("failed to decode N", slog.String("kid", key.Kid), slog.String("error", err.Error()))
			continue
		}

		// E（指数）をデコード
		eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
		if err != nil {
			slog.Warn("failed to decode E", slog.String("kid", key.Kid), slog.String("error", err.Error()))
			continue
		}

		// RSA公開鍵を作成
		var e int
		for _, b := range eBytes {
			e = e<<8 + int(b)
		}

		pubKey := &rsa.PublicKey{
			N: new(big.Int).SetBytes(nBytes),
			E: e,
		}

		m.keys[key.Kid] = pubKey
	}

	m.lastFetch = time.Now()

	return nil
}

// getKey は指定されたkidの公開鍵を返す
func (m *JWTMiddleware) getKey(kid string) (*rsa.PublicKey, error) {
	m.keysMu.RLock()
	key, exists := m.keys[kid]
	m.keysMu.RUnlock()

	if exists {
		return key, nil
	}

	// 鍵が見つからない場合、JWKSを更新（最大1分に1回）
	m.keysMu.Lock()
	defer m.keysMu.Unlock()

	if time.Since(m.lastFetch) > time.Minute {
		if err := m.fetchJWKS(); err != nil {
			return nil, fmt.Errorf("failed to refresh JWKS: %w", err)
		}

		if key, exists := m.keys[kid]; exists {
			return key, nil
		}
	}

	return nil, fmt.Errorf("key with kid=%s not found", kid)
}

// VerifyToken はJWTトークンを検証する
func (m *JWTMiddleware) VerifyToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 署名方式を検証
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// ヘッダーからkidを取得
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid not found in token header")
		}

		// 公開鍵を取得
		return m.getKey(kid)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// オーディエンスを検証
	if m.audience != "" {
		found := false
		for _, aud := range claims.Audience {
			if aud == m.audience {
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("invalid audience")
		}
	}

	// 発行者を検証
	expectedIssuer := fmt.Sprintf("https://%s/", m.domain)
	if claims.Issuer != expectedIssuer {
		return nil, fmt.Errorf("invalid issuer: expected=%s, got=%s", expectedIssuer, claims.Issuer)
	}

	return claims, nil
}

// Middleware はJWTを検証し、subをX-User-IDヘッダーとして下流サービスに転送する
func (m *JWTMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Authorizationヘッダーからトークンを抽出
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// "Bearer <token>"をパース
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// トークンを検証
		claims, err := m.VerifyToken(tokenString)
		if err != nil {
			slog.Warn("JWT verification failed", slog.String("error", err.Error()))
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// ユーザーIDをヘッダーに追加（下流サービスで使用）
		r.Header.Set("X-User-ID", claims.Subject)

		next.ServeHTTP(w, r)
	})
}
