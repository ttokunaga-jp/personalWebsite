locals {
  env_vars              = var.env_vars
  secret_env_vars       = var.secret_env_vars
  cloud_sql_instances   = var.cloud_sql_instances
  has_cloud_sql         = length(local.cloud_sql_instances) > 0
  service_account_roles = toset(var.service_account_roles)
}

resource "google_service_account" "api" {
  project      = var.project_id
  account_id   = "${var.service_name}-sa"
  display_name = "Cloud Run API Service Account"
  description  = "Service account for the personal API running on Cloud Run."
}

resource "google_project_iam_member" "api_service_account_roles" {
  for_each = local.service_account_roles

  project = var.project_id
  role    = each.value
  member  = "serviceAccount:${google_service_account.api.email}"
}

resource "google_cloud_run_v2_service" "api" {
  project  = var.project_id
  name     = var.service_name
  location = var.region
  ingress  = var.ingress
  labels   = var.labels

  template {
    labels                           = var.labels
    service_account                  = google_service_account.api.email
    execution_environment            = var.execution_environment
    max_instance_request_concurrency = var.concurrency
    timeout                          = "${var.timeout_seconds}s"

    scaling {
      min_instance_count = var.min_instance_count
      max_instance_count = var.max_instance_count
    }

    dynamic "volumes" {
      for_each = local.has_cloud_sql ? [1] : []
      content {
        name = "cloudsql"
        cloud_sql_instance {
          instances = local.cloud_sql_instances
        }
      }
    }

    vpc_access {
      connector = var.vpc_connector
      egress    = var.vpc_egress
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

      dynamic "volume_mounts" {
        for_each = local.has_cloud_sql ? [1] : []
        content {
          name       = "cloudsql"
          mount_path = "/cloudsql"
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

resource "google_cloud_run_service_iam_binding" "api_invoker" {
  location = google_cloud_run_v2_service.api.location
  service  = google_cloud_run_v2_service.api.name
  role     = "roles/run.invoker"
  members  = ["allUsers"]
}
