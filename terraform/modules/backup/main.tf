data "google_project" "current" {
  project_id = var.project_id
}

locals {
  start_year  = tonumber(formatdate("YYYY", timestamp()))
  start_month = tonumber(formatdate("MM", timestamp()))
  start_day   = tonumber(formatdate("DD", timestamp()))
  transfer_sa = "service-${data.google_project.current.number}@gcp-sa-storagetransfer.iam.gserviceaccount.com"
}

resource "google_storage_bucket" "backup" {
  name          = var.backup_bucket_name
  project       = var.project_id
  location      = var.backup_bucket_location
  storage_class = var.backup_bucket_storage_class

  labels = var.labels

  uniform_bucket_level_access = true
  force_destroy               = false

  versioning {
    enabled = true
  }
}

resource "google_storage_bucket_iam_member" "source_viewer" {
  bucket = var.source_bucket
  role   = "roles/storage.objectViewer"
  member = "serviceAccount:${local.transfer_sa}"
}

resource "google_storage_bucket_iam_member" "backup_admin" {
  bucket = google_storage_bucket.backup.name
  role   = "roles/storage.objectAdmin"
  member = "serviceAccount:${local.transfer_sa}"
}

resource "google_storage_transfer_job" "asset_backup" {
  project     = var.project_id
  description = "Nightly backup for ${var.source_bucket}"
  status      = "ENABLED"

  schedule {
    schedule_start_date {
      year  = local.start_year
      month = local.start_month
      day   = local.start_day
    }
    start_time_of_day {
      hours   = var.schedule_hour
      minutes = var.schedule_minute
      seconds = 0
      nanos   = 0
    }
  }

  transfer_spec {
    gcs_data_source {
      bucket_name = var.source_bucket
    }
    gcs_data_sink {
      bucket_name = google_storage_bucket.backup.name
    }
    object_conditions {}
    transfer_options {
      overwrite_objects_already_existing_in_sink = true
      delete_objects_unique_in_sink              = false
    }
  }
}
