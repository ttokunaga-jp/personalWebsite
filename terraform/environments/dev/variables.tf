variable "project_id" {
  description = "GCP project ID"
  type        = string
}

variable "region" {
  description = "GCP region for Cloud Run"
  type        = string
  default     = "asia-northeast1"
}

variable "api_image" {
  description = "Container image for the API service"
  type        = string
}

variable "frontend_image" {
  description = "Container image for the frontend service"
  type        = string
}
