output "instance_connection_name" {
  description = "Connection name used by Cloud Run to connect to the instance"
  value       = google_sql_database_instance.this.connection_name
}

output "instance_name" {
  description = "Cloud SQL instance name"
  value       = google_sql_database_instance.this.name
}

output "database_name" {
  description = "Default database created for the application"
  value       = google_sql_database.default.name
}

output "database_user" {
  description = "Database user for the application"
  value       = google_sql_user.app.name
}

output "db_password_secret" {
  description = "Secret Manager resource ID storing the database password"
  value       = google_secret_manager_secret.db_password.id
}

output "db_password_secret_name" {
  description = "Secret Manager secret name for referencing in other modules"
  value       = google_secret_manager_secret.db_password.secret_id
}
