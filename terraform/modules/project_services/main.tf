resource "google_project_service" "enabled" {
  for_each = toset(var.services)

  project = var.project_id
  service = each.value

  disable_on_destroy         = false
  disable_dependent_services = false
}
