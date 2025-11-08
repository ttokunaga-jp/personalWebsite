output "backup_bucket_name" {
  description = "Name of the backup bucket receiving replicated objects"
  value       = google_storage_bucket.backup.name
}

output "transfer_job_name" {
  description = "Identifier for the Storage Transfer Service job"
  value       = google_storage_transfer_job.asset_backup.name
}

