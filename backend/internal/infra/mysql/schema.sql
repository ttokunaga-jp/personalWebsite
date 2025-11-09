-- This file is embedded at build time. Keep it in sync with deploy/mysql/schema.sql.
-- 技術カタログ
CREATE TABLE IF NOT EXISTS tech_catalog (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  slug VARCHAR(64) NOT NULL UNIQUE,
  display_name VARCHAR(128) NOT NULL,
  category VARCHAR(64) NULL,
  level ENUM('beginner','intermediate','advanced') NOT NULL,
  icon VARCHAR(128) NULL,
  sort_order INT DEFAULT 0,
  is_active TINYINT(1) DEFAULT 1,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- プロフィール
CREATE TABLE IF NOT EXISTS profiles (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  display_name VARCHAR(255) NOT NULL,
  headline_ja VARCHAR(255) NULL,
  headline_en VARCHAR(255) NULL,
  summary_ja TEXT NULL,
  summary_en TEXT NULL,
  avatar_url VARCHAR(512) NULL,
  location_ja VARCHAR(255) NULL,
  location_en VARCHAR(255) NULL,
  theme_mode ENUM('light','dark','system') DEFAULT 'system',
  theme_accent_color VARCHAR(32) NULL,
  lab_name_ja VARCHAR(255) NULL,
  lab_name_en VARCHAR(255) NULL,
  lab_advisor_ja VARCHAR(255) NULL,
  lab_advisor_en VARCHAR(255) NULL,
  lab_room_ja VARCHAR(255) NULL,
  lab_room_en VARCHAR(255) NULL,
  lab_url VARCHAR(512) NULL,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS profile_affiliations (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  profile_id BIGINT UNSIGNED NOT NULL,
  kind ENUM('affiliation','community') NOT NULL,
  name VARCHAR(255) NOT NULL,
  url VARCHAR(512) NULL,
  started_at DATETIME(3) NOT NULL,
  description_ja VARCHAR(255) NULL,
  description_en VARCHAR(255) NULL,
  sort_order INT DEFAULT 0,
  INDEX idx_profile_affiliations_profile_kind (profile_id, kind, sort_order),
  CONSTRAINT fk_profile_affiliations_profile FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS profile_work_history (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  profile_id BIGINT UNSIGNED NOT NULL,
  organization_ja VARCHAR(255) NOT NULL,
  organization_en VARCHAR(255) NOT NULL,
  role_ja VARCHAR(255) NOT NULL,
  role_en VARCHAR(255) NOT NULL,
  summary_ja TEXT NULL,
  summary_en TEXT NULL,
  started_at DATETIME(3) NOT NULL,
  ended_at DATETIME(3) NULL,
  external_url VARCHAR(512) NULL,
  sort_order INT DEFAULT 0,
  INDEX idx_profile_work_history_profile (profile_id, sort_order),
  CONSTRAINT fk_profile_work_history_profile FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS profile_social_links (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  profile_id BIGINT UNSIGNED NOT NULL,
  provider ENUM('github','zenn','linkedin','x','email','website','other') NOT NULL,
  label_ja VARCHAR(255) NULL,
  label_en VARCHAR(255) NULL,
  url VARCHAR(512) NOT NULL,
  is_footer TINYINT(1) DEFAULT 0,
  sort_order INT DEFAULT 0,
  INDEX idx_profile_social_links_profile (profile_id, sort_order),
  CONSTRAINT fk_profile_social_links_profile FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS profile_tech_sections (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  profile_id BIGINT UNSIGNED NOT NULL,
  title_ja VARCHAR(255) NULL,
  title_en VARCHAR(255) NULL,
  layout ENUM('chips','list') DEFAULT 'chips',
  breakpoint VARCHAR(32) DEFAULT 'lg',
  sort_order INT DEFAULT 0,
  INDEX idx_profile_tech_sections_profile (profile_id, sort_order),
  CONSTRAINT fk_profile_tech_sections_profile FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS tech_relationships (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  entity_type ENUM('profile_section','project','research_blog') NOT NULL,
  entity_id BIGINT UNSIGNED NOT NULL,
  tech_id BIGINT UNSIGNED NOT NULL,
  context ENUM('primary','supporting') DEFAULT 'primary',
  note VARCHAR(255) NULL,
  sort_order INT DEFAULT 0,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  INDEX idx_tech_relationships_entity (entity_type, entity_id),
  INDEX idx_tech_relationships_tech (tech_id),
  CONSTRAINT fk_tech_relationships_tech FOREIGN KEY (tech_id) REFERENCES tech_catalog(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- プロジェクト
CREATE TABLE IF NOT EXISTS projects (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  slug VARCHAR(128) NOT NULL UNIQUE,
  title_ja VARCHAR(255) NOT NULL,
  title_en VARCHAR(255) NOT NULL,
  summary_ja TEXT NOT NULL,
  summary_en TEXT NOT NULL,
  description_ja LONGTEXT NULL,
  description_en LONGTEXT NULL,
  cover_image_url VARCHAR(512) NULL,
  primary_link_url VARCHAR(512) NULL,
  period_start DATE NULL,
  period_end DATE NULL,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  published TINYINT(1) DEFAULT 0,
  highlight TINYINT(1) DEFAULT 0,
  sort_order INT DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS project_links (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  project_id BIGINT UNSIGNED NOT NULL,
  link_type ENUM('repo','demo','article','slides','other') NOT NULL,
  label_ja VARCHAR(255) NULL,
  label_en VARCHAR(255) NULL,
  url VARCHAR(512) NOT NULL,
  sort_order INT DEFAULT 0,
  INDEX idx_project_links_project (project_id, sort_order),
  CONSTRAINT fk_project_links_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 研究・ブログ
CREATE TABLE IF NOT EXISTS research_blog_entries (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  slug VARCHAR(128) NOT NULL UNIQUE,
  kind ENUM('research','blog') NOT NULL,
  title_ja VARCHAR(255) NOT NULL,
  title_en VARCHAR(255) NOT NULL,
  overview_ja TEXT NULL,
  overview_en TEXT NULL,
  outcome_ja TEXT NULL,
  outcome_en TEXT NULL,
  outlook_ja TEXT NULL,
  outlook_en TEXT NULL,
  external_url VARCHAR(512) NOT NULL,
  published_at DATETIME(3) NOT NULL,
  highlight_image_url VARCHAR(512) NULL,
  image_alt_ja VARCHAR(255) NULL,
  image_alt_en VARCHAR(255) NULL,
  is_draft TINYINT(1) DEFAULT 0,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS research_blog_tags (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  entry_id BIGINT UNSIGNED NOT NULL,
  tag VARCHAR(64) NOT NULL,
  sort_order INT DEFAULT 0,
  INDEX idx_research_blog_tags_entry (entry_id, sort_order),
  CONSTRAINT fk_research_blog_tags_entry FOREIGN KEY (entry_id) REFERENCES research_blog_entries(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS research_blog_links (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  entry_id BIGINT UNSIGNED NOT NULL,
  link_type ENUM('paper','slides','video','code','external') NOT NULL,
  label_ja VARCHAR(255) NULL,
  label_en VARCHAR(255) NULL,
  url VARCHAR(512) NOT NULL,
  sort_order INT DEFAULT 0,
  INDEX idx_research_blog_links_entry (entry_id, sort_order),
  CONSTRAINT fk_research_blog_links_entry FOREIGN KEY (entry_id) REFERENCES research_blog_entries(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS research_blog_assets (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  entry_id BIGINT UNSIGNED NOT NULL,
  asset_url VARCHAR(512) NOT NULL,
  caption_ja VARCHAR(255) NULL,
  caption_en VARCHAR(255) NULL,
  sort_order INT DEFAULT 0,
  INDEX idx_research_blog_assets_entry (entry_id, sort_order),
  CONSTRAINT fk_research_blog_assets_entry FOREIGN KEY (entry_id) REFERENCES research_blog_entries(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ホーム画面設定
CREATE TABLE IF NOT EXISTS home_page_config (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  profile_id BIGINT UNSIGNED NOT NULL,
  hero_subtitle_ja VARCHAR(255) NULL,
  hero_subtitle_en VARCHAR(255) NULL,
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  CONSTRAINT fk_home_page_config_profile FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS home_quick_links (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  config_id BIGINT UNSIGNED NOT NULL,
  section ENUM('profile','research_blog','projects','contact') NOT NULL,
  label_ja VARCHAR(255) NOT NULL,
  label_en VARCHAR(255) NOT NULL,
  description_ja TEXT NULL,
  description_en TEXT NULL,
  cta_ja VARCHAR(128) NOT NULL,
  cta_en VARCHAR(128) NOT NULL,
  target_url VARCHAR(512) NOT NULL,
  sort_order INT DEFAULT 0,
  INDEX idx_home_quick_links_config (config_id, sort_order),
  CONSTRAINT fk_home_quick_links_config FOREIGN KEY (config_id) REFERENCES home_page_config(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS home_chip_sources (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  config_id BIGINT UNSIGNED NOT NULL,
  source_type ENUM('tech','affiliation','community') NOT NULL,
  limit_count INT DEFAULT 0,
  label_ja VARCHAR(255) NULL,
  label_en VARCHAR(255) NULL,
  sort_order INT DEFAULT 0,
  INDEX idx_home_chip_sources_config (config_id, sort_order),
  CONSTRAINT fk_home_chip_sources_config FOREIGN KEY (config_id) REFERENCES home_page_config(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS contact_form_settings (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  hero_title_ja VARCHAR(255) NULL,
  hero_title_en VARCHAR(255) NULL,
  hero_description_ja TEXT NULL,
  hero_description_en TEXT NULL,
  topics JSON NOT NULL,
  consent_text_ja TEXT NOT NULL,
  consent_text_en TEXT NOT NULL,
  minimum_lead_hours INT DEFAULT 24,
  recaptcha_public_key VARCHAR(128) NULL,
  support_email VARCHAR(255) NOT NULL,
  calendar_timezone VARCHAR(64) NOT NULL,
  google_calendar_id VARCHAR(255) NULL,
  booking_window_days INT DEFAULT 30,
  meeting_url_template TEXT NULL,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE contact_form_settings
  ADD COLUMN meeting_url_template TEXT NULL AFTER booking_window_days;

CREATE TABLE IF NOT EXISTS meeting_reservations (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  topic VARCHAR(255) NULL,
  message TEXT NULL,
  start_at DATETIME(3) NOT NULL,
  end_at DATETIME(3) NOT NULL,
  duration_minutes INT NOT NULL,
  google_event_id VARCHAR(255) NULL,
  google_calendar_status ENUM('pending','confirmed','declined','cancelled') DEFAULT 'pending',
  status ENUM('pending','confirmed','cancelled') DEFAULT 'pending',
  confirmation_sent_at DATETIME(3) NULL,
  last_notification_sent_at DATETIME(3) NULL,
  lookup_hash CHAR(64) NOT NULL,
  cancellation_reason TEXT NULL,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  INDEX idx_meeting_reservations_start_at (start_at),
  INDEX idx_meeting_reservations_email (email),
  INDEX idx_meeting_reservations_lookup (lookup_hash)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS meeting_notifications (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  reservation_id BIGINT UNSIGNED NOT NULL,
  notification_type ENUM('confirmation_email','reminder_email','calendar_invite','cancellation_email') NOT NULL,
  status ENUM('pending','sent','failed') DEFAULT 'pending',
  error_message TEXT NULL,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  CONSTRAINT fk_meeting_notifications_reservation FOREIGN KEY (reservation_id) REFERENCES meeting_reservations(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE profile_social_links
  MODIFY COLUMN provider ENUM('github','zenn','linkedin','x','email','website','other') NOT NULL;

-- ブラックリスト / 休業設定（既存資産を継続利用）
CREATE TABLE IF NOT EXISTS blacklist (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  email VARCHAR(255) NOT NULL UNIQUE,
  reason VARCHAR(255) NULL,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS schedule_blackouts (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  start_time DATETIME(3) NOT NULL,
  end_time DATETIME(3) NOT NULL,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 管理者セッション
CREATE TABLE IF NOT EXISTS admin_sessions (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  session_id_hash CHAR(64) NOT NULL UNIQUE,
  subject VARCHAR(255) NOT NULL,
  email VARCHAR(320) NOT NULL,
  roles JSON NOT NULL,
  user_agent VARCHAR(512) NULL,
  ip_address VARCHAR(45) NULL,
  expires_at DATETIME(3) NOT NULL,
  last_accessed_at DATETIME(3) NOT NULL,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  revoked_at DATETIME(3) NULL,
  INDEX idx_admin_sessions_expires (expires_at),
  INDEX idx_admin_sessions_last_accessed (last_accessed_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO profiles (
  display_name,
  headline_ja,
  headline_en,
  summary_ja,
  summary_en,
  avatar_url,
  location_ja,
  location_en,
  theme_mode,
  theme_accent_color,
  lab_name_ja,
  lab_name_en,
  lab_advisor_ja,
  lab_advisor_en,
  lab_room_ja,
  lab_room_en,
  lab_url
)
SELECT
  'Takumi Tokunaga',
  '',
  '',
  '',
  '',
  '',
  '',
  '',
  'system',
  NULL,
  '',
  '',
  '',
  '',
  '',
  '',
  ''
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM profiles);

INSERT INTO contact_form_settings (
  hero_title_ja,
  hero_title_en,
  hero_description_ja,
  hero_description_en,
  topics,
  consent_text_ja,
  consent_text_en,
  minimum_lead_hours,
  recaptcha_public_key,
  support_email,
  calendar_timezone,
  google_calendar_id,
  booking_window_days,
  meeting_url_template
)
SELECT
  '',
  '',
  '',
  '',
  '[]',
  '',
  '',
  24,
  '',
  'support@example.com',
  'Asia/Tokyo',
  '',
  30,
  ''
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM contact_form_settings);

INSERT INTO home_page_config (
  profile_id,
  hero_subtitle_ja,
  hero_subtitle_en
)
SELECT
  p.id,
  '',
  ''
FROM profiles p
WHERE NOT EXISTS (SELECT 1 FROM home_page_config)
ORDER BY p.id
LIMIT 1;
