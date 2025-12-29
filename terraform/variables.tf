variable "auth0_domain" {
  description = "Auth0 tenant domain"
  type        = string
}

variable "auth0_management_client_id" {
  description = "Auth0 Management API client ID"
  type        = string
}

variable "auth0_management_client_secret" {
  description = "Auth0 Management API client secret"
  type        = string
  sensitive   = true
}

variable "api_identifier" {
  description = "API identifier (audience)"
  type        = string
  default     = "https://api.platform-security.local"
}

variable "callback_urls" {
  description = "Allowed callback URLs"
  type        = list(string)
  default = [
    "http://localhost:3000/auth/callback"
  ]
}

variable "logout_urls" {
  description = "Allowed logout URLs"
  type        = list(string)
  default = [
    "http://localhost:3000"
  ]
}

variable "web_origins" {
  description = "Allowed web origins"
  type        = list(string)
  default = [
    "http://localhost:3000"
  ]
}
