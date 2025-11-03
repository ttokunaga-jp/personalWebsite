output "collections" {
  description = "Computed collection names with prefix applied"
  value = {
    research_blog_entries = local.collections.research_blog_entries
    tech_relationships    = local.collections.tech_relationships
    meeting_reservations  = local.collections.meeting_reservations
  }
}
