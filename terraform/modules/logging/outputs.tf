output "log_bucket_id" {
  description = "Centralized logging bucket identifier"
  value       = google_logging_project_bucket_config.central.bucket_id
}

output "log_bucket_name" {
  description = "Full name of the centralized logging bucket resource"
  value       = google_logging_project_bucket_config.central.name
}

output "bigquery_dataset_id" {
  description = "BigQuery dataset used for structured log exports"
  value       = google_bigquery_dataset.logs.dataset_id
}

output "archive_bucket" {
  description = "Archive bucket name for long-term log storage"
  value       = google_storage_bucket.archive.name
}

output "bucket_sink_name" {
  description = "Name of the sink feeding the centralized log bucket"
  value       = google_logging_project_sink.central_bucket.name
}

output "bigquery_sink_name" {
  description = "Name of the sink exporting logs to BigQuery"
  value       = google_logging_project_sink.bigquery.name
}

output "storage_sink_name" {
  description = "Name of the sink exporting logs to Cloud Storage"
  value       = google_logging_project_sink.storage.name
}

