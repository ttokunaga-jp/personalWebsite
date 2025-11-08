output "error_metric_name" {
  description = "Log-based metric tracking API errors"
  value       = google_logging_metric.api_error_rate.name
}

output "uptime_check_id" {
  description = "Resource ID of the API uptime check (empty if disabled)"
  value       = try(google_monitoring_uptime_check_config.api[0].id, "")
}

output "dashboard_id" {
  description = "Cloud Monitoring dashboard resource ID (empty if disabled)"
  value       = try(google_monitoring_dashboard.operations[0].id, "")
}

