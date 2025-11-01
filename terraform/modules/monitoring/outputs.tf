output "log_bucket_id" {
  description = "Log bucket ID storing centralized service logs"
  value       = google_logging_project_bucket_config.this.bucket_id
}

output "log_sink_writer_identity" {
  description = "Service account used by the Cloud Run log sink"
  value       = google_logging_project_sink.cloud_run.writer_identity
}

output "uptime_check_id" {
  description = "Resource ID of the API uptime check (empty if disabled)"
  value       = try(google_monitoring_uptime_check_config.api[0].id, "")
}
