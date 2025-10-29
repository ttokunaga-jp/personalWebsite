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
    '徳永 拓未',
    'Takumi Tokunaga',
    '立命館大学 情報理工学部 実世界情報コース / フルスタックエンジニア',
    'Ritsumeikan Univ. Real-world Information Program / Full-stack Engineer',
    '立命館大学 情報理工学部 実世界情報コース',
    'College of Information Science and Engineering, Ritsumeikan University',
    '木村研究室（RM²C モバイルコンピューティング／リアリティメディア研究室）',
    'Kimura Laboratory (RM²C Mobile Computing / Reality Media Lab)',
    'ロボティクスと XR を核に実世界とサイバー空間をつなぐ体験設計を探究しつつ、RAG や検索基盤、教育向けプロダクトを Go / TypeScript / GCP で開発しています。起業経験と大規模サービス開発の現場経験を掛け合わせ、企画から運用まで一気通貫で価値提供することを目指しています。',
    'Exploring ways to connect the physical and digital worlds via robotics and XR while building RAG-powered learning tools, search infrastructure, and education products with Go, TypeScript, and GCP. Combines entrepreneurial experience with large-scale service development to deliver value end-to-end.'
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
    (1, 'RAG / 検索基盤開発', 'RAG / Retrieval Systems', 1),
    (1, 'Go / Gin / gRPC', 'Go / Gin / gRPC', 2),
    (1, 'React / TypeScript / Tailwind', 'React / TypeScript / Tailwind', 3),
    (1, 'GCP / Cloud Run / Terraform', 'GCP / Cloud Run / Terraform', 4),
    (1, 'ROS / ロボットシステム', 'ROS / Robotics', 5)
ON DUPLICATE KEY UPDATE
    skill_ja = VALUES(skill_ja),
    skill_en = VALUES(skill_en),
    sort_order = VALUES(sort_order);
