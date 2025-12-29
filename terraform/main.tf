terraform {
  required_version = ">= 1.0"

  required_providers {
    auth0 = {
      source  = "auth0/auth0"
      version = "~> 1.0"
    }
  }
}

provider "auth0" {
  domain        = var.auth0_domain
  client_id     = var.auth0_management_client_id
  client_secret = var.auth0_management_client_secret
}

# Auth0 API（バックエンドAPI保護用）
resource "auth0_resource_server" "platform_api" {
  name       = "Platform Security API"
  identifier = var.api_identifier

  token_lifetime                                  = 86400 # 24時間
  skip_consent_for_verifiable_first_party_clients = true
}

# Auth0 API Scopes
resource "auth0_resource_server_scope" "read_profile" {
  resource_server_identifier = auth0_resource_server.platform_api.identifier
  scope                      = "read:profile"
  description                = "Read user profile"
}

resource "auth0_resource_server_scope" "write_profile" {
  resource_server_identifier = auth0_resource_server.platform_api.identifier
  scope                      = "write:profile"
  description                = "Update user profile"
}

# Auth0 Application（Regular Web App）
resource "auth0_client" "frontend_app" {
  name        = "Platform Security Frontend"
  description = "Frontend application for platform security PoC"
  app_type    = "regular_web"

  callbacks           = var.callback_urls
  allowed_logout_urls = var.logout_urls
  web_origins         = var.web_origins

  jwt_configuration {
    alg = "RS256"
  }

  # トークン設定
  oidc_conformant = true

  # Authorization Code Flow with Client Secret
  grant_types = [
    "authorization_code",
    "refresh_token"
  ]

  # Refresh Tokenの設定
  refresh_token {
    rotation_type       = "rotating"
    expiration_type     = "expiring"
    leeway              = 0
    token_lifetime      = 2592000 # 30日
    idle_token_lifetime = 1296000 # 15日
  }
}

# ApplicationにAPIへのアクセス権限を付与
resource "auth0_client_grant" "frontend_api_grant" {
  client_id = auth0_client.frontend_app.id
  audience  = auth0_resource_server.platform_api.identifier

  scopes = [
    "read:profile",
    "write:profile"
  ]
}
