variable "project_id" {
  description = "GCP project ID"
  type        = string
}

variable "region" {
  description = "GCP region"
  type        = string
}

variable "service_name" {
  description = "Cloud Run service name"
  type        = string
}

variable "image" {
  description = "Container image"
  type        = string
}

variable "vpc_connector" {
  description = "Serverless VPC connector"
  type        = string
}

variable "port" {
  description = "Container port exposed by the service"
  type        = number
  default     = 8100
}

variable "min_instance_count" {
  description = "Minimum number of instances for the Cloud Run service"
  type        = number
  default     = 0
}

variable "max_instance_count" {
  description = "Maximum number of instances for the Cloud Run service"
  type        = number
  default     = 5
}

variable "concurrency" {
  description = "Maximum requests per instance"
  type        = number
  default     = 80
}

variable "timeout_seconds" {
  description = "Request timeout for the service"
  type        = number
  default     = 30
}

variable "cpu" {
  description = "CPU limit for the container"
  type        = string
  default     = "1"
}

variable "memory" {
  description = "Memory limit for the container"
  type        = string
  default     = "512Mi"
}

variable "env_vars" {
  description = "Plain environment variables"
  type        = map(string)
  default     = {}
}

variable "secret_env_vars" {
  description = "Environment variables sourced from Secret Manager"
  type = map(object({
    secret  = string
    version = string
  }))
  default = {}
}

variable "cloud_sql_instances" {
  description = "List of Cloud SQL instance connection names to mount"
  type        = list(string)
  default     = []
}

variable "service_account_roles" {
  description = "Additional project roles to bind to the service account"
  type        = list(string)
  default     = []
}

variable "vpc_egress" {
  description = "Egress settings for the VPC connector"
  type        = string
  default     = "ALL_TRAFFIC"
}

variable "execution_environment" {
  description = "Execution environment for Cloud Run (Gen1 or Gen2)"
  type        = string
  default     = "EXECUTION_ENVIRONMENT_GEN2"
}

variable "ingress" {
  description = "Ingress policy for the Cloud Run service"
  type        = string
  default     = "INGRESS_TRAFFIC_ALL"
}

variable "labels" {
  description = "Labels applied to the Cloud Run service"
  type        = map(string)
  default     = {}
}
