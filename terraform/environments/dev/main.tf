terraform {
  required_version = ">= 1.5.0"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.17"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

locals {
  common_labels = merge(
    {
      environment = var.environment
      workload    = "personal-website"
    },
    var.additional_labels
  )

  assets_bucket_name = (
    var.assets_bucket_name != null && var.assets_bucket_name != ""
    ) ? var.assets_bucket_name : replace(
    lower(format("%s-%s-assets", var.project_id, var.environment)),
    "_",
    "-"
  )
  api_base_raw = trimspace(var.public_api_base_url)
  api_base_normalized = local.api_base_raw != "" ? regexreplace(local.api_base_raw, "/+$", "") : ""
  api_base_canonical = local.api_base_normalized == "" ? "" : (
    can(regex(".*/api(/.*)?$", local.api_base_normalized))
    ? local.api_base_normalized
    : "${local.api_base_normalized}/api"
  )
  api_proxy_pass  = local.api_base_canonical == "" ? "" : "${local.api_base_canonical}/"
  admin_login_url = local.api_base_canonical == "" ? "" : "${local.api_base_canonical}/admin/auth/login"
  google_redirect_url = local.api_base_canonical == "" ? "" : "${local.api_base_canonical}/admin/auth/callback"
}

module "project_services" {
  source     = "../../modules/project_services"
  project_id = var.project_id
  services   = var.enabled_apis
}

module "network" {
  source     = "../../modules/network"
  project_id = var.project_id
  region     = var.region

  network_name                          = var.network_name
  app_subnet_cidr                       = var.app_subnet_cidr
  connector_subnet_cidr                 = var.connector_subnet_cidr
  vpc_connector_name                    = var.vpc_connector_name
  vpc_connector_machine_type            = var.vpc_connector_machine_type
  vpc_connector_min_throughput          = var.vpc_connector_min_throughput
  vpc_connector_max_throughput          = var.vpc_connector_max_throughput
  private_service_connect_name          = var.private_service_connect_name
  private_service_connect_prefix_length = var.private_service_connect_prefix_length

  depends_on = [module.project_services]
}

module "cloudsql" {
  source = "../../modules/cloudsql"

  project_id                 = var.project_id
  region                     = var.region
  instance_name              = var.db_instance_name
  database_version           = var.db_version
  tier                       = var.db_tier
  disk_type                  = var.db_disk_type
  disk_size_gb               = var.db_disk_size_gb
  availability_type          = var.db_availability_type
  maintenance_day            = var.db_maintenance_day
  maintenance_hour           = var.db_maintenance_hour
  backup_enabled             = var.db_backup_enabled
  point_in_time_recovery     = var.db_point_in_time_recovery
  db_name                    = var.db_name
  user_name                  = var.db_user
  password_length            = var.db_password_length
  vpc_network                = module.network.network_self_link
  labels                     = local.common_labels
  deletion_protection        = var.db_deletion_protection
  private_service_connection = module.network.service_networking_connection

  depends_on = [module.network]
}

module "assets_bucket" {
  source = "../../modules/storage"

  project_id               = var.project_id
  bucket_name              = local.assets_bucket_name
  location                 = var.storage_location
  labels                   = local.common_labels
  force_destroy            = var.storage_force_destroy
  versioning               = var.storage_enable_versioning
  public_access_prevention = var.storage_public_access_prevention
  cors                     = var.storage_cors
  lifecycle_rules          = var.storage_lifecycle_rules
  log_bucket               = var.storage_log_bucket
  log_object_prefix        = var.storage_log_object_prefix
  kms_key_name             = var.storage_kms_key
  enable_public_read       = var.storage_enable_public_read

  depends_on = [module.project_services]
}

module "dns" {
  source = "../../modules/dns"

  project_id  = var.project_id
  name        = var.dns_zone_name
  dns_name    = var.dns_domain
  description = "Managed zone for personal website"
  visibility  = var.dns_visibility
  labels      = local.common_labels
  record_sets = var.dns_records

  depends_on = [module.project_services]
}

module "monitoring" {
  source = "../../modules/monitoring"

  project_id                 = var.project_id
  log_location               = var.log_location
  log_bucket_id              = var.log_bucket_id
  log_retention_days         = var.log_retention_days
  log_sink_name              = var.log_sink_name
  error_metric_name          = var.log_error_metric_name
  notification_channels      = var.notification_channels
  error_threshold_per_minute = var.monitoring_error_threshold_per_minute
  api_service_name           = var.api_service_name
  api_uptime_check           = var.api_uptime_check
  error_log_filter           = var.monitoring_error_log_filter

  depends_on = [module.project_services]
}

module "api" {
  source     = "../../modules/cloudrun/api"
  project_id = var.project_id
  region     = var.region

  service_name          = var.api_service_name
  image                 = var.api_image
  port                  = var.api_port
  min_instance_count    = var.api_min_instances
  max_instance_count    = var.api_max_instances
  concurrency           = var.api_concurrency
  timeout_seconds       = var.api_timeout_seconds
  cpu                   = var.api_cpu
  memory                = var.api_memory
  vpc_connector         = module.network.vpc_connector
  vpc_egress            = var.api_vpc_egress
  execution_environment = var.api_execution_environment
  ingress               = var.api_ingress
  labels                = local.common_labels
  cloud_sql_instances   = [module.cloudsql.instance_connection_name]
  env_vars = merge(
    {
      ENVIRONMENT                 = var.environment
      DB_NAME                     = module.cloudsql.database_name
      DB_USER                     = module.cloudsql.database_user
      DB_INSTANCE_CONNECTION_NAME = module.cloudsql.instance_connection_name
      STORAGE_BUCKET              = module.assets_bucket.bucket_name
    },
    trimspace(var.admin_redirect_uri) != "" ? {
      APP_ADMIN_REDIRECT_URI             = var.admin_redirect_uri
      APP_AUTH_ADMIN_DEFAULT_REDIRECT_URI = var.admin_redirect_uri
    } : {},
    local.google_redirect_url != "" ? {
      APP_GOOGLE_REDIRECT_URL = local.google_redirect_url
    } : {},
    var.api_additional_env
  )
  secret_env_vars = merge(
    {
      DB_PASSWORD = {
        secret  = module.cloudsql.db_password_secret
        version = "latest"
      }
    },
    trimspace(var.admin_allowed_emails_secret) != "" ? {
      APP_ADMIN_ALLOWED_EMAILS = {
        secret  = var.admin_allowed_emails_secret
        version = "latest"
      }
      APP_AUTH_ADMIN_ALLOWED_EMAILS = {
        secret  = var.admin_allowed_emails_secret
        version = "latest"
      }
    } : {},
    var.api_secret_env
  )
  service_account_roles = concat(
    [
      "roles/cloudsql.client",
      "roles/logging.logWriter",
      "roles/monitoring.metricWriter"
    ],
    var.api_additional_roles
  )

  depends_on = [
    module.project_services,
    module.network,
    module.cloudsql
  ]
}

module "frontend" {
  source     = "../../modules/cloudrun/frontend"
  project_id = var.project_id
  region     = var.region

  service_name          = var.frontend_service_name
  image                 = var.frontend_image
  port                  = var.frontend_port
  min_instance_count    = var.frontend_min_instances
  max_instance_count    = var.frontend_max_instances
  concurrency           = var.frontend_concurrency
  timeout_seconds       = var.frontend_timeout_seconds
  cpu                   = var.frontend_cpu
  memory                = var.frontend_memory
  ingress               = var.frontend_ingress
  execution_environment = var.frontend_execution_environment
  labels                = local.common_labels
  env_vars = merge(
    {
      ENVIRONMENT = var.environment
    },
    local.api_base_canonical != "" ? {
      VITE_API_BASE_URL    = local.api_base_canonical
      API_PROXY_PASS       = local.api_proxy_pass
      VITE_ADMIN_LOGIN_URL = local.admin_login_url
    } : {},
    var.frontend_additional_env
  )
  secret_env_vars = var.frontend_secret_env
  service_account_roles = concat(
    [
      "roles/logging.logWriter",
      "roles/monitoring.metricWriter"
    ],
    var.frontend_additional_roles
  )

  depends_on = [module.project_services]
}

resource "google_secret_manager_secret_iam_member" "api_db_password" {
  project   = var.project_id
  secret_id = module.cloudsql.db_password_secret_name
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${module.api.service_account}"
}

resource "google_storage_bucket_iam_member" "frontend_asset_reader" {
  bucket = module.assets_bucket.bucket_name
  role   = "roles/storage.objectViewer"
  member = "serviceAccount:${module.frontend.service_account}"
}
