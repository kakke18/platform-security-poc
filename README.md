# platform-security-poc

マルチテナントSaaSプラットフォームのセキュリティ機能を検証するためのPoCプロジェクトです。

## 概要

本プロジェクトは、以下のセキュリティ機能を実装・検証します：

- **Auth0認証**: OAuth 2.0 / OIDC ベースの認証（Authorization Code Flow + HttpOnly Cookie）
- **特権ユーザー**: 運営者向けの特別な認証ルール
- **IPアドレス制限**: ワークスペース単位のホワイトリスト制限
- **レートリミット**: グローバル + ワークスペース単位の制限
- **監査ログ**: 認証・ユーザー管理・セキュリティ設定の操作ログ

## アーキテクチャ

```
┌─────────────┐
│  Frontend   │
│  (Next.js)  │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Gateway   │
│ (Go/Connect)│
└──────┬──────┘
       │
       ├─────────────────┬─────────────────┐
       ▼                 ▼                 ▼
┌─────────────┐   ┌─────────────┐   ┌─────────────┐
│ Identity API│   │ Project API │   │  Task API   │
│ (Go/Connect)│   │ (Go/Connect)│   │ (Go/Connect)│
└──────┬──────┘   └──────┬──────┘   └──────┬──────┘
       │                 │                 │
       │                 └────────┬────────┘
       │                          │
       ├──────────────────────────┤
       ▼                          ▼
┌─────────────┐            ┌─────────────┐
│   Auth0     │            │ PostgreSQL  │
└─────────────┘            └─────────────┘
                                 │
                  ┌──────────────┼──────────────┐
                  ▼              ▼              ▼
              identity_db   project_db     task_db
```

## 技術スタック

| コンポーネント | 技術 |
|----------------|------|
| Client | Next.js 15 (App Router) |
| Gateway | Go 1.25 / Connect |
| Identity API | Go 1.25 / Connect |
| Project API | Go 1.25 / Connect |
| Task API | Go 1.25 / Connect |
| Database | PostgreSQL 17 |
| 認証 | Auth0 |
| プロトコル | Connect (Protocol Buffers) |

## ディレクトリ構成

```
platform-security-poc/
├── proto/                      # Protocol Buffers 定義
│   ├── identity/v1/
│   ├── project/v1/
│   └── task/v1/
├── frontend/                   # Next.js フロントエンド
├── identity/                   # Identity API
│   ├── cmd/server/
│   └── internal/
├── gateway/                    # Gateway
│   └── cmd/server/
├── terraform/                  # Terraform設定（Auth0）
├── Makefile
└── README.md
```

## サービス概要

### Identity API

認証・認可・ユーザー管理を担当するコアサービス。

| 機能 | 説明 |
|------|------|
| 認証 | トークン検証、ユーザー情報取得 |
| ユーザー管理 | ユーザーのCRUD操作 |
| ワークスペース管理 | 認証ポリシー、IP制限、レートリミット設定 |
| 監査ログ | 操作履歴の記録・提供 |

**注意**: 現在のPoC実装では、Auth0との認証フローはフロントエンド（Next.js）で完結しており、`@auth0/nextjs-auth0` SDKを使用してセキュアなHttpOnly Cookieベースのセッション管理を行っています。

### Project API

プロジェクト管理を担当するサンプルサービス。

| 機能 | 説明 |
|------|------|
| プロジェクト管理 | プロジェクトのCRUD操作 |
| メンバー管理 | プロジェクトメンバーの管理 |

### Task API

タスク管理を担当するサンプルサービス。

| 機能 | 説明 |
|------|------|
| タスク管理 | タスクのCRUD操作 |
| ステータス管理 | タスクの進捗管理 |

## セキュリティ機能

### 特権ユーザー

運営者向けの特別な認証ルール：

- 認証方式は `password_only` 固定
- IP制限から除外
- ワークスペースレートリミットから除外

### IPアドレス制限

ワークスペース単位でのホワイトリスト制限：

- CIDR表記対応（例: `192.168.1.0/24`）
- 特権ユーザーは除外
- 許可/拒否のログ出力

### レートリミット

2段階のレートリミット：

| レベル | 対象 | デフォルト |
|--------|------|------------|
| グローバル | 全リクエスト | 50,000 req/min |
| ワークスペース | 認証後リクエスト | 1,000 req/min |

### 監査ログ

以下の操作を記録：

- 認証: ログイン成功/失敗、ログアウト
- ユーザー管理: 作成、更新、削除
- セキュリティ管理: 認証ポリシー変更、IP制限変更、レートリミット変更

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
