locals {
  record_map = {
    for record in var.record_sets :
    "${record.name}_${record.type}" => record
  }
}

resource "google_dns_managed_zone" "this" {
  project     = var.project_id
  name        = var.name
  dns_name    = var.dns_name
  description = var.description
  visibility  = var.visibility
  labels      = var.labels
}

resource "google_dns_record_set" "records" {
  for_each = local.record_map

  project      = var.project_id
  managed_zone = google_dns_managed_zone.this.name
  name         = each.value.name
  type         = each.value.type
  ttl          = each.value.ttl
  rrdatas      = each.value.rrdatas
}
