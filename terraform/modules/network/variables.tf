variable "project_id" {
  description = "GCP project ID"
  type        = string
}

variable "region" {
  description = "GCP region"
  type        = string
}

variable "network_name" {
  description = "Name for the VPC network hosting serverless resources"
  type        = string
  default     = "serverless-network"
}

variable "app_subnet_cidr" {
  description = "CIDR block for workloads that require private IP access"
  type        = string
  default     = "10.8.0.0/24"
}

variable "connector_subnet_cidr" {
  description = "Dedicated /28 range for the Serverless VPC connector"
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
  description = "Minimum throughput for the Serverless VPC connector in Mbps"
  type        = number
  default     = 200
}

variable "vpc_connector_max_throughput" {
  description = "Maximum throughput for the Serverless VPC connector in Mbps"
  type        = number
  default     = 300
}

variable "private_service_connect_name" {
  description = "Name for the private service access allocation"
  type        = string
  default     = "serverless-psc"
}

variable "private_service_connect_prefix_length" {
  description = "Prefix length for the allocated IP range used by private service access"
  type        = number
  default     = 16
}
