resource "google_service_account" "frontend" {
  account_id   = "${var.service_name}-sa"
  display_name = "Cloud Run Frontend Service Account"
}

resource "google_cloud_run_v2_service" "frontend" {
  name     = var.service_name
  location = var.region

  template {
    service_account = google_service_account.frontend.email

    scaling {
      min_instance_count = 0
      max_instance_count = 3
    }

    containers {
      image = var.image
      ports {
        container_port = 8100
      }
    }
  }
}

resource "google_cloud_run_service_iam_binding" "frontend_invoker" {
  location = google_cloud_run_v2_service.frontend.location
  service  = google_cloud_run_v2_service.frontend.name
  role     = "roles/run.invoker"
  members  = ["allUsers"]
}
