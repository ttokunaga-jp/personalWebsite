-- Migration: backfill project tech memberships into tech_relationships
-- Phase: Tech catalog integration (Phase 3)

SET @migration_start = NOW(3);

-- Determine legacy source table name if it exists
SET @src_table = (
  SELECT CASE
    WHEN EXISTS (
      SELECT 1 FROM information_schema.tables
      WHERE table_schema = DATABASE() AND table_name = 'project_tech_stack'
    ) THEN 'project_tech_stack'
    WHEN EXISTS (
      SELECT 1 FROM information_schema.tables
      WHERE table_schema = DATABASE() AND table_name = 'legacy_project_tech_stack'
    ) THEN 'legacy_project_tech_stack'
    ELSE ''
  END
);

-- Conditionally backfill only when a legacy table is present
SET @sql = IF(
  @src_table <> '',
  CONCAT(
    'INSERT INTO tech_relationships (',
    '  entity_type, entity_id, tech_id, context, note, sort_order, created_at',
    ') ',
    'SELECT ',
    "  'project' AS entity_type, ",
    '  pts.project_id AS entity_id, ',
    '  tc.id AS tech_id, ',
    "  'primary' AS context, ",
    '  NULL AS note, ',
    '  COALESCE(pts.sort_order, 0) AS sort_order, ',
    '  ''', @migration_start, ''' ',
    'FROM ', @src_table, ' pts ',
    'JOIN tech_catalog tc ON LOWER(tc.display_name) = LOWER(pts.label) ',
    'WHERE NOT EXISTS (',
    '  SELECT 1 FROM tech_relationships existing ',
    "  WHERE existing.entity_type = 'project' ",
    '    AND existing.entity_id = pts.project_id ',
    '    AND existing.tech_id = tc.id',
    ')'
  ),
  'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- NOTE: Review legacy rows without matching tech_catalog entries manually when applicable.
