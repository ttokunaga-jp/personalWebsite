variable "project_id" {
  description = "GCP project ID"
  type        = string
}

variable "error_metric_name" {
  description = "Name of the log-based metric used for error tracking"
  type        = string
  default     = "cloud_run_api_error_count"
}

variable "notification_channels" {
  description = "Monitoring notification channel IDs"
  type        = list(string)
  default     = []
}

variable "error_threshold_per_minute" {
  description = "Threshold for number of 5xx responses per minute before alerting"
  type        = number
  default     = 5
}

variable "api_service_name" {
  description = "Cloud Run service name for the API, used in monitoring filters"
  type        = string
}

variable "api_uptime_check" {
  description = "Configuration for the API uptime check"
  type = object({
    display_name = string
    host         = string
    path         = string
    port         = number
    use_ssl      = bool
    regions      = list(string)
  })
  default = {
    display_name = "API Uptime"
    host         = ""
    path         = "/healthz"
    port         = 443
    use_ssl      = true
    regions      = ["USA"]
  }
}

variable "error_log_filter" {
  description = "Custom filter for the log-based error metric (optional)"
  type        = string
  default     = null
}

variable "latency_threshold_ms" {
  description = "Request latency threshold (milliseconds) that triggers the latency alert"
  type        = number
  default     = 800
}

variable "dashboard_display_name" {
  description = "Display name for the Cloud Monitoring dashboard"
  type        = string
  default     = "Personal Website Operations Overview"
}

variable "dashboard_enabled" {
  description = "Create the operational dashboard when true"
  type        = bool
  default     = true
}

variable "sql_backup_alert_enabled" {
  description = "Enable alerting on Cloud SQL backup failures"
  type        = bool
  default     = true
}
