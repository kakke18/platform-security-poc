# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# プロジェクトファイルをコピー
COPY project-api/go.mod project-api/go.sum ./
RUN go mod download

COPY project-api/ ./

# ビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/server ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /app

COPY --from=builder /app/bin/server /app/server

EXPOSE 8082

CMD ["/app/server"]
