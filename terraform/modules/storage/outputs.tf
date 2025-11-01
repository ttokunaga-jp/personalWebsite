output "bucket_name" {
  description = "Name of the created bucket"
  value       = google_storage_bucket.this.name
}

output "bucket_self_link" {
  description = "Self link for the bucket"
  value       = google_storage_bucket.this.self_link
}

output "bucket_url" {
  description = "gs:// URL for the bucket"
  value       = "gs://${google_storage_bucket.this.name}"
}
