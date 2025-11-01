resource "google_storage_bucket" "this" {
  project  = var.project_id
  name     = var.bucket_name
  location = var.location

  uniform_bucket_level_access = true
  force_destroy               = var.force_destroy
  labels                      = var.labels
  public_access_prevention    = var.public_access_prevention

  versioning {
    enabled = var.versioning
  }

  dynamic "cors" {
    for_each = var.cors
    content {
      origin          = cors.value.origin
      method          = cors.value.method
      response_header = cors.value.response_header
      max_age_seconds = cors.value.max_age_seconds
    }
  }

  dynamic "lifecycle_rule" {
    for_each = var.lifecycle_rules
    content {
      action {
        type = lifecycle_rule.value.action
      }
      condition {
        age                        = try(lifecycle_rule.value.condition.age, null)
        matches_prefix             = try(lifecycle_rule.value.condition.matches_prefix, null)
        matches_suffix             = try(lifecycle_rule.value.condition.matches_suffix, null)
        num_newer_versions         = try(lifecycle_rule.value.condition.num_newer_versions, null)
        with_state                 = try(lifecycle_rule.value.condition.with_state, null)
        created_before             = try(lifecycle_rule.value.condition.created_before, null)
        custom_time_before         = try(lifecycle_rule.value.condition.custom_time_before, null)
        days_since_custom_time     = try(lifecycle_rule.value.condition.days_since_custom_time, null)
        days_since_noncurrent_time = try(lifecycle_rule.value.condition.days_since_noncurrent_time, null)
      }
    }
  }

  dynamic "logging" {
    for_each = var.log_bucket != "" ? [1] : []
    content {
      log_bucket        = var.log_bucket
      log_object_prefix = var.log_object_prefix
    }
  }

  dynamic "encryption" {
    for_each = var.kms_key_name != "" ? [1] : []
    content {
      default_kms_key_name = var.kms_key_name
    }
  }
}

resource "google_storage_bucket_iam_binding" "public_read" {
  count   = var.enable_public_read ? 1 : 0
  bucket  = google_storage_bucket.this.name
  role    = "roles/storage.objectViewer"
  members = ["allUsers"]
}
