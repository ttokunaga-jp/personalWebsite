variable "project_id" {
  description = "GCP project ID"
  type        = string
}

variable "bucket_name" {
  description = "Storage bucket name"
  type        = string
}

variable "location" {
  description = "Bucket location"
  type        = string
}

variable "labels" {
  description = "Labels applied to the bucket"
  type        = map(string)
  default     = {}
}

variable "force_destroy" {
  description = "Whether to allow Terraform to destroy non-empty buckets"
  type        = bool
  default     = false
}

variable "versioning" {
  description = "Enable bucket object versioning"
  type        = bool
  default     = true
}

variable "public_access_prevention" {
  description = "Public access prevention mode (enforced or inherited)"
  type        = string
  default     = "enforced"
}

variable "cors" {
  description = "CORS rules for the bucket"
  type = list(object({
    origin          = list(string)
    method          = list(string)
    response_header = list(string)
    max_age_seconds = number
  }))
  default = []
}

variable "lifecycle_rules" {
  description = "Lifecycle rules to apply to the bucket"
  type = list(object({
    action = string
    condition = object({
      age                        = optional(number)
      matches_prefix             = optional(list(string))
      matches_suffix             = optional(list(string))
      num_newer_versions         = optional(number)
      with_state                 = optional(string)
      created_before             = optional(string)
      custom_time_before         = optional(string)
      days_since_custom_time     = optional(number)
      days_since_noncurrent_time = optional(number)
    })
  }))
  default = []
}

variable "log_bucket" {
  description = "Bucket used to store access logs (empty to disable logging)"
  type        = string
  default     = ""
}

variable "log_object_prefix" {
  description = "Prefix for log objects if logging is enabled"
  type        = string
  default     = "storage-logs"
}

variable "kms_key_name" {
  description = "KMS key used for default encryption (empty string to use Google-managed keys)"
  type        = string
  default     = ""
}

variable "enable_public_read" {
  description = "Whether to allow unauthenticated read access to objects"
  type        = bool
  default     = false
}
