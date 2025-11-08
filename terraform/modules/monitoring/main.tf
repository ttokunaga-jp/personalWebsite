locals {
  enable_uptime_check = var.api_uptime_check.host != ""
  error_filter        = coalesce(var.error_log_filter, "resource.type=\"cloud_run_revision\" AND severity>=ERROR")
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

resource "google_monitoring_alert_policy" "api_latency" {
  project               = var.project_id
  display_name          = "Cloud Run API latency regression"
  combiner              = "OR"
  notification_channels = var.notification_channels

  documentation {
    content   = "Request latency for the public API exceeded ${var.latency_threshold_ms} ms. Review recent deployments, inspect Cloud Run revisions, and consult the incident runbook."
    mime_type = "text/markdown"
  }

  conditions {
    display_name = "High API latency"

    condition_threshold {
      comparison      = "COMPARISON_GT"
      duration        = "600s"
      threshold_value = var.latency_threshold_ms
      trigger {
        count = 1
      }
      filter = "resource.type=\"cloud_run_revision\" AND metric.type=\"run.googleapis.com/request_latencies\" AND resource.label.\"service_name\"=\"${var.api_service_name}\""

      aggregations {
        alignment_period     = "60s"
        per_series_aligner   = "ALIGN_PERCENTILE_95"
        cross_series_reducer = "REDUCE_MEAN"
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

resource "google_logging_metric" "sql_backup_failure" {
  count       = var.sql_backup_alert_enabled ? 1 : 0
  project     = var.project_id
  name        = "cloudsql_backup_failures"
  description = "Count of Cloud SQL automated backup failures."
  filter      = "resource.type=\"cloudsql_database\" AND protoPayload.serviceName=\"cloudsql.googleapis.com\" AND protoPayload.status.code!=0"

  metric_descriptor {
    value_type  = "INT64"
    metric_kind = "DELTA"
    unit        = "1"
  }
}

resource "google_monitoring_alert_policy" "sql_backup_failure" {
  count                 = var.sql_backup_alert_enabled ? 1 : 0
  project               = var.project_id
  display_name          = "Cloud SQL automated backup failed"
  combiner              = "OR"
  notification_channels = var.notification_channels

  documentation {
    content   = "Automated Cloud SQL backup reported failures. Follow the database backup runbook to restore redundancy."
    mime_type = "text/markdown"
  }

  conditions {
    display_name = "Backup failure detected"

    condition_threshold {
      comparison      = "COMPARISON_GT"
      duration        = "0s"
      threshold_value = 0
      trigger {
        count = 1
      }
      filter = "resource.type=\"logging_log\" AND metric.type=\"logging.googleapis.com/user/${google_logging_metric.sql_backup_failure[0].name}\""

      aggregations {
        alignment_period   = "300s"
        per_series_aligner = "ALIGN_RATE"
      }
    }
  }
}

resource "google_monitoring_dashboard" "operations" {
  count   = var.dashboard_enabled ? 1 : 0
  project = var.project_id
  dashboard_json = templatefile("${path.module}/dashboards/service_overview.json.tmpl", {
    DISPLAY_NAME      = var.dashboard_display_name
    PROJECT_ID        = var.project_id
    API_SERVICE_NAME  = var.api_service_name
    LATENCY_THRESHOLD = var.latency_threshold_ms
    ERROR_POLICY_NAME = google_monitoring_alert_policy.api_errors.name
    UPTIME_CHECK_ID   = try(google_monitoring_uptime_check_config.api[0].id, "")
  })
}
