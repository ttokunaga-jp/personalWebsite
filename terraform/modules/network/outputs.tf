output "network_name" {
  description = "VPC network name"
  value       = google_compute_network.serverless.name
}

output "network_self_link" {
  description = "Self link for the VPC network"
  value       = google_compute_network.serverless.self_link
}

output "app_subnetwork" {
  description = "Subnetwork hosting private workloads"
  value       = google_compute_subnetwork.app.self_link
}

output "connector_subnetwork" {
  description = "Subnetwork used by the Serverless VPC connector"
  value       = google_compute_subnetwork.connector.self_link
}

output "vpc_connector" {
  description = "Serverless VPC connector resource"
  value       = google_vpc_access_connector.serverless.name
}

output "vpc_connector_id" {
  description = "Serverless VPC connector ID"
  value       = google_vpc_access_connector.serverless.id
}

output "private_service_connect_range" {
  description = "Allocated range name for private service access"
  value       = google_compute_global_address.private_service_connect.name
}

output "service_networking_connection" {
  description = "ID of the service networking connection enabling private services"
  value       = google_service_networking_connection.private_vpc_connection.id
}
