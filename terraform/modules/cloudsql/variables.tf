variable "project_id" {
  description = "GCP project ID"
  type        = string
}

variable "region" {
  description = "Region for the Cloud SQL instance"
  type        = string
}

variable "instance_name" {
  description = "Name of the Cloud SQL instance"
  type        = string
  default     = "personal-db"
}

variable "database_version" {
  description = "Cloud SQL database engine version"
  type        = string
  default     = "MYSQL_8_0"
}

variable "tier" {
  description = "Machine tier for the Cloud SQL instance"
  type        = string
  default     = "db-custom-2-4096"
}

variable "disk_type" {
  description = "Disk type for Cloud SQL"
  type        = string
  default     = "PD_SSD"
}

variable "disk_size_gb" {
  description = "Disk size in GB"
  type        = number
  default     = 20
}

variable "availability_type" {
  description = "Availability configuration (ZONAL or REGIONAL)"
  type        = string
  default     = "ZONAL"
}

variable "maintenance_day" {
  description = "Day of week for maintenance (1=Monday)"
  type        = number
  default     = 7
}

variable "maintenance_hour" {
  description = "Hour (0-23) for maintenance window"
  type        = number
  default     = 3
}

variable "backup_enabled" {
  description = "Whether automated backups are enabled"
  type        = bool
  default     = true
}

variable "point_in_time_recovery" {
  description = "Enable point in time recovery (requires binary logs for MySQL)"
  type        = bool
  default     = true
}

variable "db_name" {
  description = "Default database name"
  type        = string
  default     = "portfolio"
}

variable "user_name" {
  description = "Database user for the application"
  type        = string
  default     = "app_user"
}

variable "password_length" {
  description = "Length of the generated user password"
  type        = number
  default     = 32
}

variable "vpc_network" {
  description = "VPC network self link for private IP"
  type        = string
}

variable "labels" {
  description = "Labels applied to created resources"
  type        = map(string)
  default     = {}
}

variable "deletion_protection" {
  description = "Enable deletion protection for the Cloud SQL instance"
  type        = bool
  default     = true
}

variable "private_service_connection" {
  description = "Service Networking connection resource ID to ensure private service access is established"
  type        = string
}
