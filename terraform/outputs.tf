output "api_identifier" {
  description = "Auth0 API identifier (audience)"
  value       = auth0_resource_server.platform_api.identifier
}

output "frontend_client_id" {
  description = "Frontend application client ID"
  value       = auth0_client.frontend_app.client_id
}

output "frontend_client_secret_note" {
  description = "How to retrieve frontend application client secret"
  value       = "Client Secret must be retrieved from Auth0 Dashboard: Applications → ${auth0_client.frontend_app.name} → Settings → Client Secret"
}

output "auth0_domain" {
  description = "Auth0 tenant domain"
  value       = var.auth0_domain
}
