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

INSERT INTO projects (id, title_ja, title_en, description_ja, description_en, link_url, year, published, sort_order) VALUES
    (
        1,
        'RAG 学習支援システム「ClassNav」',
        'ClassNav: RAG-powered Learning Assistant',
        '大学の LMS と連携し、講義資料を自動で取り込み要約・検索できる学習支援システム。Docling によるマルチフォーマット解析と RAG を組み合わせ、NotebookLM との差別化となる資料リンク提示や自動アップロード機能を実装。',
        'Learning assistant that syncs with the university LMS, parses lecture materials via Docling, and provides RAG-based answers with direct source linking and automatic ingestion.',
        'https://github.com/ttokunaga-jp',
        2024,
        TRUE,
        30
    ),
    (
        2,
        'searchService: ハイブリッド検索マイクロサービス',
        'searchService: Hybrid Retrieval Microservice',
        'Elasticsearch と Qdrant を統合し、キーワード・ベクトルを加重スコアリングする検索基盤。Kafka を用いた非同期インデックス更新、OpenTelemetry / Prometheus による可観測性、gRPC + HTTP API を備え、RAG サービスの共通モジュールとして運用。',
        'Hybrid retrieval service that blends Elasticsearch and Qdrant scoring, supports Kafka-driven asynchronous indexing, and exposes gRPC/HTTP APIs with full observability for RAG workloads.',
        'https://github.com/ttokunaga-jp',
        2024,
        TRUE,
        20
    ),
    (
        3,
        '個人ポートフォリオサイト（Go + React）',
        'Personal Portfolio Site (Go + React)',
        'Go (Gin) と React SPA で構築した個人ポートフォリオ。公開サイトと管理 GUI を分離し、予約フォーム、Google OAuth + JWT 認証、Cloud Build → Cloud Run の CI/CD、Terraform による IaC を備える。',
        'Full-stack personal site built with Go (Gin) and React SPA. Provides a public site and admin GUI with booking flows, Google OAuth + JWT auth, CI/CD on Cloud Build → Cloud Run, and Terraform-managed infrastructure.',
        'https://github.com/ttokunaga-jp',
        2025,
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

DELETE FROM project_tech_stack WHERE project_id IN (1, 2, 3);

INSERT INTO project_tech_stack (project_id, label, sort_order) VALUES
    (1, 'RAG / LangChain', 1),
    (1, 'Go / Fiber API', 2),
    (1, 'React / TypeScript / Tailwind', 3),
    (1, 'PostgreSQL / Redis', 4),
    (1, 'GCP Cloud Run / Secret Manager', 5),

    (2, 'Go / gRPC', 1),
    (2, 'Elasticsearch', 2),
    (2, 'Qdrant', 3),
    (2, 'Kafka', 4),
    (2, 'OpenTelemetry / Prometheus / Jaeger', 5),

    (3, 'Go / Gin / Fx', 1),
    (3, 'React / pnpm Workspace', 2),
    (3, 'MySQL', 3),
    (3, 'Cloud Build / Cloud Run', 4),
    (3, 'Terraform', 5);
