variable "project_id" {
  description = "GCP project ID"
  type        = string
}

variable "source_bucket" {
  description = "Name of the bucket to back up"
  type        = string
}

variable "backup_bucket_name" {
  description = "Destination bucket for replicated objects"
  type        = string
}

variable "backup_bucket_location" {
  description = "Location/region for the backup bucket"
  type        = string
  default     = "asia-northeast1"
}

variable "backup_bucket_storage_class" {
  description = "Storage class for the backup bucket"
  type        = string
  default     = "NEARLINE"
}

variable "schedule_hour" {
  description = "Hour (0-23) when the transfer job should run"
  type        = number
  default     = 3
}

variable "schedule_minute" {
  description = "Minute (0-59) when the transfer job should run"
  type        = number
  default     = 0
}

variable "labels" {
  description = "Labels applied to created resources"
  type        = map(string)
  default     = {}
}
