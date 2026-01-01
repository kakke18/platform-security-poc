module github.com/kakke18/platform-security-poc/backend/user

go 1.25.5

require (
	connectrpc.com/connect v1.19.1
	github.com/kakke18/platform-security-poc/backend/gen v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.33.0
)

require (
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/kakke18/platform-security-poc/backend/gen => ../gen
