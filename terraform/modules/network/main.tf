resource "google_compute_network" "serverless" {
  project                 = var.project_id
  name                    = var.network_name
  auto_create_subnetworks = false
  routing_mode            = "REGIONAL"
}

resource "google_compute_subnetwork" "app" {
  project                  = var.project_id
  name                     = "${var.network_name}-app-${var.region}"
  ip_cidr_range            = var.app_subnet_cidr
  region                   = var.region
  network                  = google_compute_network.serverless.id
  private_ip_google_access = true
  stack_type               = "IPV4_ONLY"
}

resource "google_compute_subnetwork" "connector" {
  project       = var.project_id
  name          = "${var.network_name}-connector-${var.region}"
  ip_cidr_range = var.connector_subnet_cidr
  region        = var.region
  network       = google_compute_network.serverless.id
}

resource "google_vpc_access_connector" "serverless" {
  project = var.project_id
  name    = var.vpc_connector_name
  region  = var.region

  subnet {
    name = google_compute_subnetwork.connector.name
  }

  machine_type   = var.vpc_connector_machine_type
  min_throughput = var.vpc_connector_min_throughput
  max_throughput = var.vpc_connector_max_throughput
}

resource "google_compute_global_address" "private_service_connect" {
  project       = var.project_id
  name          = var.private_service_connect_name
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = var.private_service_connect_prefix_length
  network       = google_compute_network.serverless.id
}

resource "google_service_networking_connection" "private_vpc_connection" {
  network                 = google_compute_network.serverless.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_service_connect.name]
}
