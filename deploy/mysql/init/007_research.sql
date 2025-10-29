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
    'RAG による学習支援システムの信頼性向上に関する研究',
    'Improving Reliability of RAG-based Learning Assistants',
    '大学 LMS の資料を対象に、Docling を活用した抽出とベクトル検索を組み合わせ、引用リンク提示でハルシネーションを抑制する学習支援基盤を設計。',
    'Designed a learning assistant that blends Docling-based parsing with hybrid retrieval to surface source-backed answers and reduce hallucinations for LMS course materials.',
    '# 研究概要\n\nLMS から自動取得した講義資料を Docling で構造化し、Elasticsearch と Qdrant を統合したハイブリッド検索でベクトル類似度とキーワードスコアを重み付けします。回答には参照元のリンクを必ず添付し、ユーザーが検証しやすい UI を React で実装しました。RAG パイプライン全体を可観測化し、再現性のある評価ワークフローを整備しています。',
    '# Overview\n\nCourse materials fetched from the LMS are parsed with Docling, indexed into Elasticsearch and Qdrant, and served via a hybrid scoring pipeline. Responses cite their sources by design, with a React UI that keeps verification one click away. Observability and repeatable evaluation workflows were introduced to quantify reliability improvements.',
    2025,
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
