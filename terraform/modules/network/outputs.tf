output "network" {
  description = "VPC network name"
  value       = google_compute_network.serverless.name
}

output "vpc_connector" {
  description = "Serverless VPC connector name"
  value       = google_vpc_access_connector.serverless.name
}
