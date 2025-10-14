resource "google_compute_network" "serverless" {
  name                    = "serverless-network"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "serverless" {
  name          = "serverless-subnet"
  ip_cidr_range = "10.8.0.0/28"
  region        = var.region
  network       = google_compute_network.serverless.id
}

resource "google_vpc_access_connector" "serverless" {
  name          = "serverless-connector"
  network       = google_compute_network.serverless.name
  region        = var.region
  ip_cidr_range = "10.8.0.0/28"
}
