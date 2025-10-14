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

module "network" {
  source     = "../../modules/network"
  project_id = var.project_id
  region     = var.region
}

module "api" {
  source        = "../../modules/cloudrun/api"
  project_id    = var.project_id
  region        = var.region
  service_name  = "personal-api"
  image         = var.api_image
  vpc_connector = module.network.vpc_connector
}

module "frontend" {
  source        = "../../modules/cloudrun/frontend"
  project_id    = var.project_id
  region        = var.region
  service_name  = "personal-frontend"
  image         = var.frontend_image
}
