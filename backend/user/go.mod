module github.com/kakke18/platform-security-poc/backend/user

go 1.25

require (
	connectrpc.com/connect v1.19.1
	github.com/kakke18/platform-security-poc/backend/gateway v0.0.0-00010101000000-000000000000
	github.com/kakke18/platform-security-poc/backend/identity v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.33.0
	google.golang.org/protobuf v1.36.11
)

require golang.org/x/text v0.21.0 // indirect

replace (
	github.com/kakke18/platform-security-poc/backend/gateway => ../gateway
	github.com/kakke18/platform-security-poc/backend/identity => ../identity
)
