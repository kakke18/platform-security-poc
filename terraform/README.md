# Terraform - Auth0 Configuration

このディレクトリには、Auth0のアプリケーションとAPIをTerraformで管理するための設定が含まれています。

## 前提条件

- Terraform >= 1.0
- Auth0アカウント
- Auth0 Management APIの認証情報

## セットアップ

### 1. Auth0 Management APIの認証情報取得

1. [Auth0 Dashboard](https://manage.auth0.com/) にログイン
2. Applications → Applications → "Auth0 Management API" を選択
3. "Machine to Machine Applications" タブを選択
4. Terraformで使用するアプリケーションを作成または既存のものを選択
5. 以下の権限を付与:
   - `read:clients`
   - `create:clients`
   - `update:clients`
   - `delete:clients`
   - `read:resource_servers`
   - `create:resource_servers`
   - `update:resource_servers`
   - `delete:resource_servers`
   - `read:client_grants`
   - `create:client_grants`
   - `update:client_grants`
   - `delete:client_grants`

### 2. 環境変数ファイルの作成

```bash
cp terraform.tfvars.example terraform.tfvars
```

`terraform.tfvars` を編集して、実際の値を設定してください。

### 3. Terraform初期化

```bash
terraform init
```

### 4. 設定の確認

```bash
terraform plan
```

### 5. リソースの適用

```bash
terraform apply
```

## 作成されるリソース

### Auth0 API
- **名前**: Platform Security API
- **Identifier**: `https://api.platform-security.local` (変更可能)
- **Scopes**:
  - `read:profile`: ユーザープロファイルの読み取り
  - `write:profile`: ユーザープロファイルの更新

### Auth0 Application
- **名前**: Platform Security Frontend
- **タイプ**: Regular Web Application
- **Grant Types**: Authorization Code, Refresh Token
- **Refresh Token**: Rotating, 30日間有効

## 出力値の確認

```bash
# すべての出力値を表示
terraform output

# 特定の値を表示
terraform output frontend_client_id

# センシティブな値を表示
terraform output -raw frontend_client_secret
```

## フロントエンド環境変数への反映

Terraform適用後、以下のコマンドで環境変数を取得できます:

```bash
# クライアントID
terraform output -raw frontend_client_id

# クライアントシークレット
terraform output -raw frontend_client_secret

# API Identifier (Audience)
terraform output -raw api_identifier

# ドメイン
terraform output -raw auth0_domain
```

これらの値を `frontend/config/.env.local.json` に設定してください:

```json
{
  "AUTH0_SECRET": "use_openssl_rand_hex_32_to_generate_a_secret",
  "AUTH0_BASE_URL": "http://localhost:3000",
  "AUTH0_DOMAIN": "<terraform output -raw auth0_domain>",
  "AUTH0_CLIENT_ID": "<terraform output -raw frontend_client_id>",
  "AUTH0_CLIENT_SECRET": "<terraform output -raw frontend_client_secret>",
  "AUTH0_AUDIENCE": "<terraform output -raw api_identifier>"
}
```

## リソースの削除

```bash
terraform destroy
```

## ディレクトリ構成

```
terraform/
├── main.tf                    # メインのTerraform設定
├── variables.tf               # 変数定義
├── outputs.tf                 # 出力値定義
├── terraform.tfvars.example   # 環境変数の例
├── .gitignore                 # Git除外設定
└── README.md                  # このファイル
```
