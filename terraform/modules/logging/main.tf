resource "google_logging_project_bucket_config" "central" {
  project        = var.project_id
  location       = var.log_location
  bucket_id      = var.log_bucket_id
  retention_days = var.log_bucket_retention_days
  description    = "Centralized log bucket for the personal website workloads."
}

resource "google_bigquery_dataset" "logs" {
  dataset_id  = var.bigquery_dataset_id
  project     = var.project_id
  location    = var.bigquery_dataset_location
  description = "Structured Cloud Logging export for long-term analysis."
  labels      = var.labels
}

resource "google_storage_bucket" "archive" {
  name          = var.archive_bucket_name
  project       = var.project_id
  location      = var.archive_bucket_location
  storage_class = var.archive_storage_class

  force_destroy               = false
  uniform_bucket_level_access = true

  lifecycle_rule {
    action {
      type = "Delete"
    }
    condition {
      age = var.archive_retention_policy_days
    }
  }

  versioning {
    enabled = true
  }

  labels = var.labels
}

resource "google_logging_project_sink" "central_bucket" {
  project                = var.project_id
  name                   = var.bucket_sink_name
  destination            = "logging.googleapis.com/projects/${var.project_id}/locations/${var.log_location}/buckets/${google_logging_project_bucket_config.central.bucket_id}"
  filter                 = var.log_sink_filter
  unique_writer_identity = true
}

resource "google_project_iam_member" "central_bucket_writer" {
  project = var.project_id
  role    = "roles/logging.bucketWriter"
  member  = google_logging_project_sink.central_bucket.writer_identity
}

resource "google_logging_project_sink" "bigquery" {
  project                = var.project_id
  name                   = var.bigquery_sink_name
  destination            = "bigquery.googleapis.com/projects/${var.project_id}/datasets/${google_bigquery_dataset.logs.dataset_id}"
  filter                 = var.log_sink_filter
  unique_writer_identity = true
}

resource "google_bigquery_dataset_iam_member" "sink_writer" {
  dataset_id = google_bigquery_dataset.logs.dataset_id
  project    = var.project_id
  role       = "roles/bigquery.dataEditor"
  member     = google_logging_project_sink.bigquery.writer_identity
}

resource "google_logging_project_sink" "storage" {
  project                = var.project_id
  name                   = var.storage_sink_name
  destination            = "storage.googleapis.com/${google_storage_bucket.archive.name}"
  filter                 = var.log_sink_filter
  unique_writer_identity = true
}

resource "google_storage_bucket_iam_member" "archive_writer" {
  bucket = google_storage_bucket.archive.name
  role   = "roles/storage.objectCreator"
  member = google_logging_project_sink.storage.writer_identity
}
