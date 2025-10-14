output "url" {
  description = "Cloud Run URL"
  value       = google_cloud_run_v2_service.frontend.uri
}

output "service_account" {
  description = "Service account email"
  value       = google_service_account.frontend.email
}
