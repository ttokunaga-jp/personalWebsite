resource "google_service_account" "api" {
  account_id   = "${var.service_name}-sa"
  display_name = "Cloud Run API Service Account"
}

resource "google_cloud_run_v2_service" "api" {
  name     = var.service_name
  location = var.region

  template {
    service_account = google_service_account.api.email

    scaling {
      min_instance_count = 0
      max_instance_count = 3
    }

    containers {
      image = var.image
      ports {
        container_port = 8100
      }
      env {
        name  = "PORT"
        value = "8100"
      }
    }

    vpc_access {
      connector = var.vpc_connector
    }
  }
}

resource "google_cloud_run_service_iam_binding" "api_invoker" {
  location = google_cloud_run_v2_service.api.location
  service  = google_cloud_run_v2_service.api.name
  role     = "roles/run.invoker"
  members  = ["allUsers"]
}
