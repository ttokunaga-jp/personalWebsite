CREATE TABLE IF NOT EXISTS profile (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name_ja VARCHAR(255) NOT NULL,
    name_en VARCHAR(255) NOT NULL,
    title_ja VARCHAR(255) NOT NULL,
    title_en VARCHAR(255) NOT NULL,
    affiliation_ja VARCHAR(255) NOT NULL,
    affiliation_en VARCHAR(255) NOT NULL,
    lab_ja VARCHAR(255) NOT NULL,
    lab_en VARCHAR(255) NOT NULL,
    summary_ja TEXT NOT NULL,
    summary_en TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS profile_skills (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    profile_id BIGINT UNSIGNED NOT NULL,
    skill_ja VARCHAR(255) NOT NULL,
    skill_en VARCHAR(255) NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_profile_skills_profile
        FOREIGN KEY (profile_id) REFERENCES profile(id)
        ON DELETE CASCADE
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

INSERT INTO profile (
    id,
    name_ja,
    name_en,
    title_ja,
    title_en,
    affiliation_ja,
    affiliation_en,
    lab_ja,
    lab_en,
    summary_ja,
    summary_en
) VALUES (
    1,
    '高橋 拓海',
    'Takumi Takahashi',
    'フルスタックエンジニア',
    'Full-stack Engineer',
    '株式会社サンプル / プロダクト開発部',
    'Sample Inc. / Product Engineering',
    'クラウド・AI応用研究室',
    'Cloud & Applied AI Lab',
    'クラウドネイティブな分散システムと AI 応用を中心に、研究開発から実運用まで横断的に取り組んでいます。Go / TypeScript を軸に、高品質で保守性の高いサービス構築をリードしています。',
    'Leads research and delivery of cloud-native distributed systems and applied AI solutions, focusing on Go and TypeScript to build maintainable, high-quality services.'
)
ON DUPLICATE KEY UPDATE
    name_ja = VALUES(name_ja),
    name_en = VALUES(name_en),
    title_ja = VALUES(title_ja),
    title_en = VALUES(title_en),
    affiliation_ja = VALUES(affiliation_ja),
    affiliation_en = VALUES(affiliation_en),
    lab_ja = VALUES(lab_ja),
    lab_en = VALUES(lab_en),
    summary_ja = VALUES(summary_ja),
    summary_en = VALUES(summary_en);

INSERT INTO profile_skills (
    profile_id,
    skill_ja,
    skill_en,
    sort_order
) VALUES
    (1, 'Go / Gin / Clean Architecture', 'Go / Gin / Clean Architecture', 1),
    (1, 'React / TypeScript / Next.js', 'React / TypeScript / Next.js', 2),
    (1, 'GCP / Cloud Run / Cloud SQL', 'GCP / Cloud Run / Cloud SQL', 3),
    (1, 'CI/CD / Terraform / GitHub Actions', 'CI/CD / Terraform / GitHub Actions', 4)
ON DUPLICATE KEY UPDATE
    skill_ja = VALUES(skill_ja),
    skill_en = VALUES(skill_en),
    sort_order = VALUES(sort_order);
