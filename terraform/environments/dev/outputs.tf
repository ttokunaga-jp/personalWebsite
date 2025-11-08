output "api_service_url" {
  description = "Cloud Run URL for the API"
  value       = module.api.url
}

output "frontend_service_url" {
  description = "Cloud Run URL for the frontend"
  value       = module.frontend.url
}

output "api_service_account_email" {
  description = "Service account running the API Cloud Run service"
  value       = module.api.service_account
}

output "frontend_service_account_email" {
  description = "Service account running the frontend Cloud Run service"
  value       = module.frontend.service_account
}

output "cloud_sql_connection_name" {
  description = "Connection name for the Cloud SQL instance"
  value       = module.cloudsql.instance_connection_name
}

output "cloud_sql_database_name" {
  description = "Default database created for the application"
  value       = module.cloudsql.database_name
}

output "cloud_sql_database_user" {
  description = "Database user provisioned for the application"
  value       = module.cloudsql.database_user
}

output "db_password_secret" {
  description = "Secret Manager resource storing the database password"
  value       = module.cloudsql.db_password_secret
}

output "assets_bucket_name" {
  description = "GCS bucket used for static assets"
  value       = module.assets_bucket.bucket_name
}

output "dns_name_servers" {
  description = "Name servers for delegating the managed DNS zone"
  value       = module.dns.name_servers
}

output "monitoring_uptime_check_id" {
  description = "Uptime check ID for the API endpoint"
  value       = module.monitoring.uptime_check_id
}

output "monitoring_dashboard_id" {
  description = "Cloud Monitoring dashboard resource ID"
  value       = module.monitoring.dashboard_id
}

output "firestore_collections" {
  description = "Computed Firestore collection names (with prefix applied)"
  value       = module.firestore.collections
}

output "logging_log_bucket_id" {
  description = "Centralized logging bucket ID"
  value       = module.logging.log_bucket_id
}

output "logging_bigquery_dataset_id" {
  description = "BigQuery dataset receiving log exports"
  value       = module.logging.bigquery_dataset_id
}

output "backup_bucket_name" {
  description = "Backup bucket receiving nightly transfers"
  value       = module.backup.backup_bucket_name
}

output "backup_transfer_job_name" {
  description = "Storage Transfer Service job responsible for nightly backups"
  value       = module.backup.transfer_job_name
}
