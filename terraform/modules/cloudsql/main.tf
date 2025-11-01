resource "random_password" "app_user" {
  length  = var.password_length
  special = false
}

resource "terraform_data" "service_networking_dependency" {
  triggers_replace = [var.private_service_connection]
}

resource "google_sql_database_instance" "this" {
  project             = var.project_id
  name                = var.instance_name
  region              = var.region
  database_version    = var.database_version
  deletion_protection = var.deletion_protection
  depends_on          = [terraform_data.service_networking_dependency]

  settings {
    tier              = var.tier
    availability_type = var.availability_type
    disk_type         = var.disk_type
    disk_size         = var.disk_size_gb
    user_labels       = var.labels

    ip_configuration {
      ipv4_enabled                                  = false
      require_ssl                                   = true
      private_network                               = var.vpc_network
      enable_private_path_for_google_cloud_services = true
    }

    backup_configuration {
      enabled                        = var.backup_enabled
      start_time                     = format("%02d:00", var.maintenance_hour)
      point_in_time_recovery_enabled = var.point_in_time_recovery
    }

    maintenance_window {
      day          = var.maintenance_day
      hour         = var.maintenance_hour
      update_track = "stable"
    }
  }
}

resource "google_sql_database" "default" {
  project  = var.project_id
  name     = var.db_name
  instance = google_sql_database_instance.this.name
}

resource "google_sql_user" "app" {
  project  = var.project_id
  instance = google_sql_database_instance.this.name
  name     = var.user_name
  password = random_password.app_user.result
  host     = "%"
}

resource "google_secret_manager_secret" "db_password" {
  project   = var.project_id
  secret_id = "${var.instance_name}-db-password"

  replication {
    auto {}
  }

  labels = var.labels
}

resource "google_secret_manager_secret_version" "db_password" {
  secret      = google_secret_manager_secret.db_password.id
  secret_data = random_password.app_user.result
}
