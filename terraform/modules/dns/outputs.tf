output "zone_name" {
  description = "Managed zone name"
  value       = google_dns_managed_zone.this.name
}

output "dns_name" {
  description = "DNS name managed by the zone"
  value       = google_dns_managed_zone.this.dns_name
}

output "name_servers" {
  description = "Authoritative name servers for delegation"
  value       = google_dns_managed_zone.this.name_servers
}
