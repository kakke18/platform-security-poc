module github.com/kakke18/platform-security-poc/backend/identity

go 1.25.5

require (
	connectrpc.com/connect v1.19.1
	github.com/kakke18/platform-security-poc/backend/gen v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.47.0
)

require (
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251029180050-ab9386a59fda // indirect
	google.golang.org/grpc v1.78.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/kakke18/platform-security-poc/backend/gen => ../gen
