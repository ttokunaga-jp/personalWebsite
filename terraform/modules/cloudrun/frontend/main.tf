locals {
  env_vars              = var.env_vars
  secret_env_vars       = var.secret_env_vars
  service_account_roles = toset(var.service_account_roles)
}

resource "google_service_account" "frontend" {
  project      = var.project_id
  account_id   = "${var.service_name}-sa"
  display_name = "Cloud Run Frontend Service Account"
  description  = "Service account for the personal frontend running on Cloud Run."
}

resource "google_project_iam_member" "frontend_service_account_roles" {
  for_each = local.service_account_roles

  project = var.project_id
  role    = each.value
  member  = "serviceAccount:${google_service_account.frontend.email}"
}

resource "google_cloud_run_v2_service" "frontend" {
  project  = var.project_id
  name     = var.service_name
  location = var.region
  ingress  = var.ingress
  labels   = var.labels

  template {
    labels                           = var.labels
    service_account                  = google_service_account.frontend.email
    execution_environment            = var.execution_environment
    max_instance_request_concurrency = var.concurrency
    timeout                          = "${var.timeout_seconds}s"

    scaling {
      min_instance_count = var.min_instance_count
      max_instance_count = var.max_instance_count
    }

    containers {
      image = var.image

      ports {
        container_port = var.port
      }

      env {
        name  = "PORT"
        value = tostring(var.port)
      }

      dynamic "env" {
        for_each = local.env_vars
        content {
          name  = env.key
          value = env.value
        }
      }

      dynamic "env" {
        for_each = local.secret_env_vars
        iterator = secret
        content {
          name = secret.key
          value_source {
            secret_key_ref {
              secret  = secret.value.secret
              version = secret.value.version
            }
          }
        }
      }

      resources {
        limits = {
          cpu    = var.cpu
          memory = var.memory
        }
      }
    }
  }

  traffic {
    percent = 100
    type    = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
  }
}

resource "google_cloud_run_service_iam_binding" "frontend_invoker" {
  location = google_cloud_run_v2_service.frontend.location
  service  = google_cloud_run_v2_service.frontend.name
  role     = "roles/run.invoker"
  members  = ["allUsers"]
}
