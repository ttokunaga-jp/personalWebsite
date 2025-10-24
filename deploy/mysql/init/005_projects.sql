CREATE TABLE IF NOT EXISTS projects (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title_ja VARCHAR(255) NOT NULL,
    title_en VARCHAR(255) NOT NULL,
    description_ja TEXT NOT NULL,
    description_en TEXT NOT NULL,
    link_url VARCHAR(512) NOT NULL DEFAULT '',
    year INT NOT NULL,
    published TINYINT(1) NOT NULL DEFAULT 0,
    sort_order INT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS project_tech_stack (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    project_id BIGINT UNSIGNED NOT NULL,
    label VARCHAR(255) NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_project_stack_project
        FOREIGN KEY (project_id) REFERENCES projects(id)
        ON DELETE CASCADE
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

INSERT INTO projects (
    id,
    title_ja,
    title_en,
    description_ja,
    description_en,
    link_url,
    year,
    published,
    sort_order
) VALUES (
    1,
    'リアルタイム行動解析ダッシュボード',
    'Real-time Activity Analytics Dashboard',
    'IoT センサーから収集したデータをリアルタイムで可視化し、異常検知を行うための社内向けダッシュボード。イベントドリブンなアーキテクチャとストリーミング処理を採用。',
    'An internal dashboard that ingests IoT sensor data, performs real-time anomaly detection, and visualises metrics in a streaming-first architecture.',
    'https://example.com/projects/realtime-analytics',
    2024,
    TRUE,
    10
)
ON DUPLICATE KEY UPDATE
    title_ja = VALUES(title_ja),
    title_en = VALUES(title_en),
    description_ja = VALUES(description_ja),
    description_en = VALUES(description_en),
    link_url = VALUES(link_url),
    year = VALUES(year),
    published = VALUES(published),
    sort_order = VALUES(sort_order);

DELETE FROM project_tech_stack WHERE project_id = 1;

INSERT INTO project_tech_stack (project_id, label, sort_order) VALUES
    (1, 'Go / Gin', 1),
    (1, 'React / TypeScript', 2),
    (1, 'GCP Pub/Sub / Cloud Run', 3),
    (1, 'BigQuery / Looker Studio', 4);
