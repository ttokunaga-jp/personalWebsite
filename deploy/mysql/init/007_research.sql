CREATE TABLE IF NOT EXISTS research (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title_ja VARCHAR(255) NOT NULL,
    title_en VARCHAR(255) NOT NULL,
    summary_ja TEXT NOT NULL,
    summary_en TEXT NOT NULL,
    content_md_ja MEDIUMTEXT NOT NULL,
    content_md_en MEDIUMTEXT NOT NULL,
    year INT NOT NULL,
    published TINYINT(1) NOT NULL DEFAULT 0,
    sort_order INT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

INSERT INTO research (
    id,
    title_ja,
    title_en,
    summary_ja,
    summary_en,
    content_md_ja,
    content_md_en,
    year,
    published,
    sort_order
) VALUES (
    1,
    '大規模時系列データのオンライン異常検知',
    'Online Anomaly Detection for Large-scale Time Series',
    'エッジデバイスから収集した数百万系列のデータに対して、リアルタイムな異常検知を行うためのストリーミング処理基盤を構築。',
    'Designed a streaming platform capable of running real-time anomaly detection across millions of edge-device time series.',
    '# 研究概要\n\nGo と Rust を用いたストリーミング処理によって、従来数時間を要していたバッチ解析を数秒レイテンシまで短縮しました。Cloud Run と Pub/Sub を活用し、スケールするイベントドリブンなアーキテクチャを採用しています。',
    '# Overview\n\nA streaming architecture written in Go and Rust reduced batch analysis latency from hours to seconds. Leveraged Cloud Run and Pub/Sub for an event-driven design that scales with demand.',
    2024,
    TRUE,
    10
)
ON DUPLICATE KEY UPDATE
    title_ja = VALUES(title_ja),
    title_en = VALUES(title_en),
    summary_ja = VALUES(summary_ja),
    summary_en = VALUES(summary_en),
    content_md_ja = VALUES(content_md_ja),
    content_md_en = VALUES(content_md_en),
    year = VALUES(year),
    published = VALUES(published),
    sort_order = VALUES(sort_order);
