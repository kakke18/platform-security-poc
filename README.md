# platform-security-poc

マルチテナントSaaSプラットフォームのセキュリティ機能を検証するためのPoCプロジェクトです。

## 概要

本プロジェクトは、BFF (Backend for Frontend) パターンによる認証・認可の一元管理を実装したセキュリティPoC環境です。

### 実装済み機能

- **Auth0認証**: OAuth 2.0 / OIDC ベースの認証（Authorization Code Flow + HttpOnly Cookie）
- **BFF (Gateway)**: JWT検証、認証・認可の一元管理、X-User-IDヘッダー付与
- **内部サービス**: Gatewayで検証済みのユーザーIDを信頼した処理
- **ユーザー情報管理**: Connect RPC (gRPC互換) によるAPI通信

### 将来実装予定

- mTLS (Gateway-Identity間の相互TLS認証)
- 特権ユーザー管理
- IPアドレス制限
- レートリミット
- 監査ログ

## アーキテクチャ

```
┌─────────────────────────────────────────────────┐
│              Frontend (Next.js 15)              │
│  - @auth0/nextjs-auth0 (HttpOnly Cookie)       │
│  - Connect RPC Client                           │
└──────────────────┬──────────────────────────────┘
                   │ HTTPS
                   │ Authorization: Bearer JWT
                   ▼
┌─────────────────────────────────────────────────┐
│          Gateway (BFF) - Port 8080              │
│  - JWT検証 (Auth0 JWKS)                         │
│  - 認証・認可の一元管理                          │
│  - リバースプロキシ                              │
└──────────────────┬──────────────────────────────┘
                   │ HTTP (内部ネットワーク)
                   │ X-User-ID ヘッダー付与
                   ▼
┌─────────────────────────────────────────────────┐
│         Identity API - Port 8081                │
│  - ユーザー情報管理                              │
│  - X-User-ID ヘッダーからユーザー取得            │
│  - ビジネスロジックのみ (JWT検証不要)            │
└──────────────────┬──────────────────────────────┘
                   │
        ┌──────────┴──────────┐
        ▼                     ▼
┌──────────────┐      ┌──────────────┐
│    Auth0     │      │  Database    │
│  (JWKS提供)  │      │ (Mock実装)   │
└──────────────┘      └──────────────┘
```

### セキュリティレイヤー

| レイヤー | 責務 | 実装内容 |
|----------|------|----------|
| **Frontend** | セッション管理 | HttpOnly Cookie、トークン自動リフレッシュ |
| **Gateway** | 認証・認可 | JWT検証、ユーザー情報取得、X-User-IDヘッダー付与 |
| **Identity** | ビジネスロジック | Gatewayからの信頼済みリクエスト処理 |

**注意**: 現在はGateway-Identity間はHTTP通信。本番環境ではmTLSまたはネットワーク分離を実装推奨。

## 技術スタック

| コンポーネント | 技術 | 備考 |
|----------------|------|------|
| Frontend | Next.js 15 (App Router) | @auth0/nextjs-auth0 |
| Gateway | Go 1.25 / Connect | JWT検証、リバースプロキシ |
| Identity API | Go 1.25 / Connect | ユーザー情報管理 |
| 認証 | Auth0 | OIDC Provider |
| プロトコル | Connect (Protocol Buffers) | gRPC互換 |
| コード生成 | Buf | Protobuf → Go/TypeScript |

## ディレクトリ構成

```
platform-security-poc/
├── proto/                      # Protocol Buffers 定義
│   └── identity/v1/
│       └── user.proto
├── frontend/                   # Next.js フロントエンド
│   ├── src/
│   │   ├── api/
│   │   │   ├── connect-client.ts    # Connect RPC client
│   │   │   └── generated/           # 自動生成 (buf)
│   │   ├── app/
│   │   │   ├── api/auth/[auth0]/   # Auth0 SDK endpoints
│   │   │   └── dashboard/
│   │   ├── features/
│   │   └── libs/auth/
│   └── config/
│       └── .env.local.json
├── gateway/                    # Gateway (BFF)
│   ├── cmd/server/
│   └── internal/
│       ├── config/
│       ├── middleware/
│       │   ├── jwt.go              # JWT検証
│       │   └── logging.go
│       └── server/
├── identity/                   # Identity API
│   ├── cmd/server/
│   └── internal/
│       ├── config/
│       ├── server/
│       ├── user/
│       │   ├── handler.go          # X-User-ID から取得
│       │   ├── repository.go
│       │   └── mock_repository.go
│       └── middleware/
│           └── logging.go
├── terraform/                  # Terraform設定（Auth0）
├── buf.gen.yaml                # Buf code generation config
├── buf.yaml                    # Buf config
├── Makefile
└── README.md
```

## サービス概要

### Gateway (BFF)

フロントエンドとバックエンド間の仲介を行うBFF (Backend for Frontend)。

| 機能 | 説明 |
|------|------|
| JWT検証 | Auth0のJWKSを使用したトークン検証 |
| 認証・認可 | ユーザー認証とアクセス制御の一元管理 |
| リバースプロキシ | 内部APIへのリクエスト転送 |
| ヘッダー付与 | 検証済みユーザーIDを`X-User-ID`ヘッダーで転送 |

**セキュリティ実装**:
- Auth0のJWKSから公開鍵を取得してJWT署名検証
- トークン有効期限・発行者・オーディエンスの検証
- 検証済みユーザーID (`sub`) をヘッダーで下流に転送

### Identity API

ユーザー情報管理を担当する内部サービス。

| 機能 | 説明 |
|------|------|
| ユーザー情報取得 | `X-User-ID`ヘッダーからユーザー情報を返却 |

**セキュリティ実装**:
- Gatewayからの信頼済みリクエストのみ処理
- `X-User-ID`ヘッダーの存在チェックのみ（JWT検証不要）
- ビジネスロジックに専念

**注意**: JWT検証やAuth0連携は全てGatewayで実施。Identity APIは内部サービスとしてGatewayからの信頼済みリクエストのみを処理します。

## セットアップ

### 前提条件

- Go 1.25+
- Node.js 20+
- pnpm
- Auth0アカウント

### 1. Auth0の設定

1. [Auth0ダッシュボード](https://manage.auth0.com/)で **Regular Web Application** を作成
2. 以下のURLを設定：
   - **Allowed Callback URLs**: `http://localhost:3000/auth/callback`
   - **Allowed Logout URLs**: `http://localhost:3000`
   - **Allowed Web Origins**: `http://localhost:3000`
3. **Domain**、**Client ID**、**Client Secret** をコピー

### 2. 環境変数の設定

**フロントエンド用** (`frontend/config/.env.local.json`)：

```json
{
  "NEXT_PUBLIC_API_URL": "http://localhost:8080",
  "APP_BASE_URL": "http://localhost:3000",
  "AUTH0_DOMAIN": "your-tenant.auth0.com",
  "AUTH0_CLIENT_ID": "your_client_id",
  "AUTH0_CLIENT_SECRET": "your_client_secret",
  "AUTH0_SECRET": "your_generated_secret_here",
  "AUTH0_AUDIENCE": "your_api_identifier"
}
```

`AUTH0_SECRET` の生成:
```bash
openssl rand -hex 32
```

### 3. 起動

**Identity API と Gateway の環境変数設定** (direnv推奨):

```bash
# Identity API
cd identity
cp .env.example .env
# .env を編集してAuth0の認証情報を設定

# Gateway
cd ../gateway
cp .env.example .env
# 必要に応じて .env を編集
```

**ターミナル1: Identity API**

```bash
cd identity
go run ./cmd/server
```

**ターミナル2: Gateway**

```bash
cd gateway
go run ./cmd/server
```

**ターミナル3: Frontend**

```bash
cd frontend
pnpm install
pnpm dev
```

### 4. アクセス

- Frontend: http://localhost:3000
- Gateway: http://localhost:8080
- Identity API: http://localhost:8081

ブラウザで http://localhost:3000 を開くと、Auth0のログイン画面にリダイレクトされます。

## 認証フロー

本実装では `@auth0/nextjs-auth0` SDK を使用した OAuth 2.0 Authorization Code Flow を採用：

1. ユーザーが `/` にアクセス
2. Next.js Proxy（旧Middleware）でセッションをチェック
3. セッションがなければ Auth0 にリダイレクト
4. Auth0 で認証後、`/auth/callback` にコールバック
5. セッションを HttpOnly Cookie に保存
6. ダッシュボードにリダイレクト

**セキュリティ上の利点:**
- トークンは HttpOnly Cookie に保存（XSS攻撃から保護）
- CSRF保護が組み込み
- トークンの自動リフレッシュ

## 開発

### プロトコル定義の更新

```bash
# Protocol Buffersの再生成
buf generate
```

生成されたファイルは各サービスの `gen/` ディレクトリに配置されます（gitignore対象）。

### Terraform (Auth0管理)

```bash
cd terraform
terraform init
terraform plan
terraform apply
```
