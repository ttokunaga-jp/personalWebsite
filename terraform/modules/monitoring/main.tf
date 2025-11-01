locals {
  enable_uptime_check = var.api_uptime_check.host != ""
  error_filter        = coalesce(var.error_log_filter, "resource.type=\"cloud_run_revision\" AND severity>=ERROR")
}

resource "google_logging_project_bucket_config" "this" {
  project        = var.project_id
  location       = var.log_location
  bucket_id      = var.log_bucket_id
  retention_days = var.log_retention_days
  description    = "Centralized log bucket for the personal website services."
}

resource "google_logging_project_sink" "cloud_run" {
  project                = var.project_id
  name                   = var.log_sink_name
  destination            = "logging.googleapis.com/projects/${var.project_id}/locations/${var.log_location}/buckets/${google_logging_project_bucket_config.this.bucket_id}"
  filter                 = "resource.type=\"cloud_run_revision\""
  unique_writer_identity = true
}

resource "google_project_iam_member" "sink_writer" {
  project = var.project_id
  role    = "roles/logging.bucketWriter"
  member  = google_logging_project_sink.cloud_run.writer_identity
}

resource "google_logging_metric" "api_error_rate" {
  project     = var.project_id
  name        = var.error_metric_name
  description = "Count of Cloud Run API error logs for SLO tracking."
  filter      = local.error_filter
  metric_descriptor {
    value_type  = "INT64"
    metric_kind = "DELTA"
    unit        = "1"
  }
  label_extractors = {}
}

resource "google_monitoring_alert_policy" "api_errors" {
  project               = var.project_id
  display_name          = "Cloud Run API 5xx spike"
  combiner              = "OR"
  notification_channels = var.notification_channels

  documentation {
    content   = "Cloud Run API is returning 5xx responses above the configured threshold. Investigate service logs and Cloud SQL availability."
    mime_type = "text/markdown"
  }

  conditions {
    display_name = "High Cloud Run API 5xx rate"

    condition_threshold {
      comparison      = "COMPARISON_GT"
      duration        = "300s"
      threshold_value = var.error_threshold_per_minute
      trigger {
        count = 1
      }
      filter = "resource.type=\"cloud_run_revision\" AND metric.type=\"run.googleapis.com/request_count\" AND metric.label.\"response_code_class\"=\"5xx\" AND resource.label.\"service_name\"=\"${var.api_service_name}\""

      aggregations {
        alignment_period     = "60s"
        per_series_aligner   = "ALIGN_DELTA"
        cross_series_reducer = "REDUCE_SUM"
        group_by_fields      = ["resource.label.\"service_name\""]
      }
    }
  }
}

resource "google_monitoring_uptime_check_config" "api" {
  count = local.enable_uptime_check ? 1 : 0

  project      = var.project_id
  display_name = var.api_uptime_check.display_name
  timeout      = "10s"
  period       = "60s"

  selected_regions = var.api_uptime_check.regions

  http_check {
    path         = var.api_uptime_check.path
    port         = var.api_uptime_check.port
    use_ssl      = var.api_uptime_check.use_ssl
    validate_ssl = true
  }

  monitored_resource {
    type = "uptime_url"
    labels = {
      host = var.api_uptime_check.host
    }
  }
}
