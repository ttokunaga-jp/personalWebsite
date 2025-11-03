variable "project_id" {
  description = "GCP project ID"
  type        = string
}

variable "database_id" {
  description = "Firestore database ID (usually (default))"
  type        = string
  default     = "(default)"
}

variable "collection_prefix" {
  description = "Optional prefix (environment) added to collection names"
  type        = string
  default     = ""
}
