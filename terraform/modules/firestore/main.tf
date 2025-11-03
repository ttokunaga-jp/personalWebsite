locals {
  normalized_prefix = trimspace(var.collection_prefix) == "" ? "" : "${trimspace(var.collection_prefix)}_"

  collections = {
    research_blog_entries = "${local.normalized_prefix}research_blog_entries"
    tech_relationships    = "${local.normalized_prefix}tech_relationships"
    meeting_reservations  = "${local.normalized_prefix}meeting_reservations"
  }
}

# Index: research/blog entries by kind + draft flag + published date
resource "google_firestore_index" "research_kind_published" {
  project    = var.project_id
  database   = var.database_id
  collection = local.collections.research_blog_entries

  fields {
    field_path = "kind"
    order      = "ASCENDING"
  }
  fields {
    field_path = "isDraft"
    order      = "ASCENDING"
  }
  fields {
    field_path = "publishedAt"
    order      = "DESCENDING"
  }
}

# Index: research/blog entries by tags (array contains) and published date
resource "google_firestore_index" "research_tags" {
  project    = var.project_id
  database   = var.database_id
  collection = local.collections.research_blog_entries

  fields {
    field_path  = "kind"
    order       = "ASCENDING"
  }
  fields {
    field_path   = "tags"
    array_config = "CONTAINS"
  }
  fields {
    field_path = "publishedAt"
    order      = "DESCENDING"
  }
}

# Index: tech relationships by entity and sort order
resource "google_firestore_index" "tech_relationships_entity" {
  project    = var.project_id
  database   = var.database_id
  collection = local.collections.tech_relationships

  fields {
    field_path = "entityType"
    order      = "ASCENDING"
  }
  fields {
    field_path = "entityId"
    order      = "ASCENDING"
  }
  fields {
    field_path = "sortOrder"
    order      = "ASCENDING"
  }
}

# Index: meeting reservation lookup (hash + createdAt)
resource "google_firestore_index" "meeting_lookup" {
  project    = var.project_id
  database   = var.database_id
  collection = local.collections.meeting_reservations

  fields {
    field_path = "lookupHash"
    order      = "ASCENDING"
  }
  fields {
    field_path = "createdAt"
    order      = "DESCENDING"
  }
}
