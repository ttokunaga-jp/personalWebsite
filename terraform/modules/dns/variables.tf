variable "project_id" {
  description = "GCP project ID"
  type        = string
}

variable "name" {
  description = "Unique name for the DNS managed zone"
  type        = string
}

variable "dns_name" {
  description = "DNS name suffix for the managed zone (must end with a dot)"
  type        = string
}

variable "description" {
  description = "Description for the DNS zone"
  type        = string
  default     = "Managed zone for personal website"
}

variable "visibility" {
  description = "Zone visibility (public or private)"
  type        = string
  default     = "public"
}

variable "labels" {
  description = "Labels applied to the DNS zone"
  type        = map(string)
  default     = {}
}

variable "record_sets" {
  description = "Record sets to create within the managed zone"
  type = list(object({
    name    = string
    type    = string
    ttl     = number
    rrdatas = list(string)
  }))
  default = []
}
