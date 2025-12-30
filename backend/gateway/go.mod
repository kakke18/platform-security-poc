module github.com/kakke18/platform-security-poc/backend/gateway

go 1.25.5

require (
	connectrpc.com/connect v1.19.1
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/kakke18/platform-security-poc/backend/identity v0.0.0-00010101000000-000000000000
	github.com/kakke18/platform-security-poc/backend/user v0.0.0-00010101000000-000000000000
	github.com/rs/cors v1.11.1
	google.golang.org/protobuf v1.36.11
)

replace (
	github.com/kakke18/platform-security-poc/backend/identity => ../identity
	github.com/kakke18/platform-security-poc/backend/user => ../user
)
