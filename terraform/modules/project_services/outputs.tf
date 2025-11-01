output "enabled_services" {
  description = "APIs enabled by this module"
  value       = keys(google_project_service.enabled)
}
