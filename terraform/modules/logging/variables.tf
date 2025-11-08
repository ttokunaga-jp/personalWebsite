variable "project_id" {
  description = "GCP project ID"
  type        = string
}

variable "log_location" {
  description = "Region for log bucket resources (e.g. global or region name)"
  type        = string
  default     = "global"
}

variable "log_bucket_id" {
  description = "Identifier for the centralized logging bucket"
  type        = string
}

variable "log_bucket_retention_days" {
  description = "Retention period (in days) for the centralized logging bucket"
  type        = number
  default     = 30
}

variable "log_sink_filter" {
  description = "Advanced log filter applied to all log router sinks"
  type        = string
  default     = "resource.type=\"cloud_run_revision\""
}

variable "bucket_sink_name" {
  description = "Name for the sink that routes logs into the centralized logging bucket"
  type        = string
  default     = "centralized-log-bucket"
}

variable "bigquery_sink_name" {
  description = "Name for the sink that routes logs into BigQuery"
  type        = string
  default     = "log-bq-export"
}

variable "storage_sink_name" {
  description = "Name for the sink that routes logs into Cloud Storage archive"
  type        = string
  default     = "log-storage-archive"
}

variable "bigquery_dataset_id" {
  description = "Dataset ID for structured log exports"
  type        = string
  default     = "app_logs"
}

variable "bigquery_dataset_location" {
  description = "BigQuery dataset location"
  type        = string
  default     = "asia-northeast1"
}

variable "archive_bucket_name" {
  description = "Cloud Storage bucket name for long-term log archive"
  type        = string
}

variable "archive_bucket_location" {
  description = "Location for the archive bucket"
  type        = string
  default     = "asia-northeast1"
}

variable "archive_storage_class" {
  description = "Storage class used for archive bucket"
  type        = string
  default     = "COLDLINE"
}

variable "archive_retention_policy_days" {
  description = "Object retention period for the archive bucket"
  type        = number
  default     = 90
}

variable "labels" {
  description = "Labels applied to logging resources"
  type        = map(string)
  default     = {}
}

