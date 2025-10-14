output "api_service_url" {
  description = "Cloud Run URL for the API"
  value       = module.api.url
}

output "frontend_service_url" {
  description = "Cloud Run URL for the frontend"
  value       = module.frontend.url
}
