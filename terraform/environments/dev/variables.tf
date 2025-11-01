variable "project_id" {
  description = "GCP project ID"
  type        = string
}

variable "region" {
  description = "Primary region for regional resources"
  type        = string
  default     = "asia-northeast1"
}

variable "environment" {
  description = "Deployment environment identifier"
  type        = string
  default     = "dev"
}

variable "additional_labels" {
  description = "Additional labels applied to resources"
  type        = map(string)
  default     = {}
}

variable "enabled_apis" {
  description = "List of APIs that must be enabled for the project"
  type        = list(string)
  default = [
    "run.googleapis.com",
    "artifactregistry.googleapis.com",
    "cloudbuild.googleapis.com",
    "compute.googleapis.com",
    "vpcaccess.googleapis.com",
    "servicenetworking.googleapis.com",
    "sqladmin.googleapis.com",
    "secretmanager.googleapis.com",
    "dns.googleapis.com",
    "logging.googleapis.com",
    "monitoring.googleapis.com",
    "iam.googleapis.com"
  ]
}

# Network
variable "network_name" {
  description = "Name of the VPC network used by serverless workloads"
  type        = string
  default     = "serverless-network"
}

variable "app_subnet_cidr" {
  description = "CIDR range for application subnet"
  type        = string
  default     = "10.8.0.0/24"
}

variable "connector_subnet_cidr" {
  description = "CIDR range for VPC connector subnet"
  type        = string
  default     = "10.8.8.0/28"
}

variable "vpc_connector_name" {
  description = "Name of the Serverless VPC connector"
  type        = string
  default     = "serverless-connector"
}

variable "vpc_connector_machine_type" {
  description = "Machine type for the Serverless VPC connector"
  type        = string
  default     = "e2-micro"
}

variable "vpc_connector_min_throughput" {
  description = "Minimum throughput for the VPC connector"
  type        = number
  default     = 200
}

variable "vpc_connector_max_throughput" {
  description = "Maximum throughput for the VPC connector"
  type        = number
  default     = 300
}

variable "private_service_connect_name" {
  description = "Name for the private service access allocation"
  type        = string
  default     = "serverless-psc"
}

variable "private_service_connect_prefix_length" {
  description = "Prefix length for private service access ranges"
  type        = number
  default     = 16
}

# Cloud SQL
variable "db_instance_name" {
  description = "Cloud SQL instance name"
  type        = string
  default     = "personal-db"
}

variable "db_version" {
  description = "Cloud SQL database version"
  type        = string
  default     = "MYSQL_8_0"
}

variable "db_tier" {
  description = "Machine tier for Cloud SQL"
  type        = string
  default     = "db-custom-2-4096"
}

variable "db_disk_type" {
  description = "Cloud SQL disk type"
  type        = string
  default     = "PD_SSD"
}

variable "db_disk_size_gb" {
  description = "Cloud SQL disk size in GB"
  type        = number
  default     = 20
}

variable "db_availability_type" {
  description = "Cloud SQL availability type (ZONAL or REGIONAL)"
  type        = string
  default     = "ZONAL"
}

variable "db_maintenance_day" {
  description = "Maintenance window day (1=Monday)"
  type        = number
  default     = 7
}

variable "db_maintenance_hour" {
  description = "Maintenance window hour (0-23)"
  type        = number
  default     = 3
}

variable "db_backup_enabled" {
  description = "Enable automated backups"
  type        = bool
  default     = true
}

variable "db_point_in_time_recovery" {
  description = "Enable point-in-time recovery"
  type        = bool
  default     = true
}

variable "db_name" {
  description = "Default database name"
  type        = string
  default     = "portfolio"
}

variable "db_user" {
  description = "Application database user"
  type        = string
  default     = "app_user"
}

variable "db_password_length" {
  description = "Generated database password length"
  type        = number
  default     = 32
}

variable "db_deletion_protection" {
  description = "Enable deletion protection on the Cloud SQL instance"
  type        = bool
  default     = true
}

# Storage
variable "assets_bucket_name" {
  description = "Override for the assets bucket name (leave null to auto-generate)"
  type        = string
  default     = null
}

variable "storage_location" {
  description = "Location for the assets bucket"
  type        = string
  default     = "asia-northeast1"
}

variable "storage_force_destroy" {
  description = "Allow Terraform to destroy non-empty buckets"
  type        = bool
  default     = false
}

variable "storage_enable_versioning" {
  description = "Enable object versioning for the assets bucket"
  type        = bool
  default     = true
}

variable "storage_public_access_prevention" {
  description = "Public access prevention mode for the bucket"
  type        = string
  default     = "enforced"
}

variable "storage_cors" {
  description = "CORS rules for the assets bucket"
  type = list(object({
    origin          = list(string)
    method          = list(string)
    response_header = list(string)
    max_age_seconds = number
  }))
  default = []
}

variable "storage_lifecycle_rules" {
  description = "Lifecycle rules for the assets bucket"
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

variable "storage_log_bucket" {
  description = "Bucket for access logs (empty string disables logging)"
  type        = string
  default     = ""
}

variable "storage_log_object_prefix" {
  description = "Prefix for access log objects"
  type        = string
  default     = "storage-logs"
}

variable "storage_kms_key" {
  description = "KMS key for bucket encryption (empty for Google-managed)"
  type        = string
  default     = ""
}

variable "storage_enable_public_read" {
  description = "Allow unauthenticated read access to the bucket"
  type        = bool
  default     = false
}

# DNS
variable "dns_zone_name" {
  description = "Cloud DNS managed zone name"
  type        = string
  default     = "personal-site"
}

variable "dns_domain" {
  description = "DNS name for the zone (must end with a dot)"
  type        = string
  default     = "example.com."
}

variable "dns_visibility" {
  description = "DNS zone visibility (public or private)"
  type        = string
  default     = "public"
}

variable "dns_records" {
  description = "Record sets to create in the DNS zone"
  type = list(object({
    name    = string
    type    = string
    ttl     = number
    rrdatas = list(string)
  }))
  default = []
}

# Monitoring & Logging
variable "log_location" {
  description = "Location of the centralized log bucket"
  type        = string
  default     = "global"
}

variable "log_bucket_id" {
  description = "Log bucket ID for centralized logs"
  type        = string
  default     = "personal-logs"
}

variable "log_retention_days" {
  description = "Retention in days for centralized logs"
  type        = number
  default     = 30
}

variable "log_sink_name" {
  description = "Log sink name aggregating Cloud Run logs"
  type        = string
  default     = "cloud-run-logs"
}

variable "log_error_metric_name" {
  description = "Log-based metric name for API errors"
  type        = string
  default     = "cloud_run_api_error_count"
}

variable "notification_channels" {
  description = "Monitoring notification channel IDs"
  type        = list(string)
  default     = []
}

variable "monitoring_error_threshold_per_minute" {
  description = "Threshold of 5xx responses per minute triggering the alert"
  type        = number
  default     = 5
}

variable "monitoring_error_log_filter" {
  description = "Optional custom filter for the error log metric"
  type        = string
  default     = null
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

# Cloud Run API
variable "api_service_name" {
  description = "Cloud Run service name for the API"
  type        = string
  default     = "personal-api"
}

variable "api_image" {
  description = "Container image for the API service"
  type        = string
}

variable "api_port" {
  description = "Container port exposed by the API"
  type        = number
  default     = 8100
}

variable "api_min_instances" {
  description = "Minimum number of API instances"
  type        = number
  default     = 0
}

variable "api_max_instances" {
  description = "Maximum number of API instances"
  type        = number
  default     = 5
}

variable "api_concurrency" {
  description = "Maximum concurrent requests per API instance"
  type        = number
  default     = 80
}

variable "api_timeout_seconds" {
  description = "API request timeout in seconds"
  type        = number
  default     = 30
}

variable "api_cpu" {
  description = "CPU limit for the API container"
  type        = string
  default     = "1"
}

variable "api_memory" {
  description = "Memory limit for the API container"
  type        = string
  default     = "512Mi"
}

variable "api_vpc_egress" {
  description = "VPC egress setting for the API"
  type        = string
  default     = "ALL_TRAFFIC"
}

variable "api_execution_environment" {
  description = "Execution environment for the API Cloud Run service"
  type        = string
  default     = "EXECUTION_ENVIRONMENT_GEN2"
}

variable "api_ingress" {
  description = "Ingress policy for the API Cloud Run service"
  type        = string
  default     = "INGRESS_TRAFFIC_ALL"
}

variable "api_additional_env" {
  description = "Extra plain environment variables for the API"
  type        = map(string)
  default     = {}
}

variable "api_secret_env" {
  description = "Additional secret-sourced environment variables for the API"
  type = map(object({
    secret  = string
    version = string
  }))
  default = {}
}

variable "api_additional_roles" {
  description = "Additional IAM roles for the API service account"
  type        = list(string)
  default     = []
}

# Cloud Run Frontend
variable "frontend_service_name" {
  description = "Cloud Run service name for the frontend"
  type        = string
  default     = "personal-frontend"
}

variable "frontend_image" {
  description = "Container image for the frontend service"
  type        = string
}

variable "frontend_port" {
  description = "Container port exposed by the frontend"
  type        = number
  default     = 8080
}

variable "frontend_min_instances" {
  description = "Minimum number of frontend instances"
  type        = number
  default     = 0
}

variable "frontend_max_instances" {
  description = "Maximum number of frontend instances"
  type        = number
  default     = 3
}

variable "frontend_concurrency" {
  description = "Maximum concurrent requests per frontend instance"
  type        = number
  default     = 250
}

variable "frontend_timeout_seconds" {
  description = "Frontend request timeout in seconds"
  type        = number
  default     = 30
}

variable "frontend_cpu" {
  description = "CPU limit for the frontend container"
  type        = string
  default     = "1"
}

variable "frontend_memory" {
  description = "Memory limit for the frontend container"
  type        = string
  default     = "512Mi"
}

variable "frontend_ingress" {
  description = "Ingress policy for the frontend Cloud Run service"
  type        = string
  default     = "INGRESS_TRAFFIC_ALL"
}

variable "frontend_execution_environment" {
  description = "Execution environment for the frontend Cloud Run service"
  type        = string
  default     = "EXECUTION_ENVIRONMENT_GEN2"
}

variable "frontend_additional_env" {
  description = "Plain environment variables for the frontend"
  type        = map(string)
  default     = {}
}

variable "frontend_secret_env" {
  description = "Secret-sourced environment variables for the frontend"
  type = map(object({
    secret  = string
    version = string
  }))
  default = {}
}

variable "frontend_additional_roles" {
  description = "Additional IAM roles for the frontend service account"
  type        = list(string)
  default     = []
}

variable "public_api_base_url" {
  description = "Public base URL for the API consumed by the frontend"
  type        = string
  default     = ""
}
