# コンテンツ管理リファクタリング設計

## 目的
- 公開 SPA・管理 SPA 間でプロフィール / プロジェクト / 研究・ブログ / お問い合わせのデータ構造を統一し、ハードコードを排除する。
- 技術セットを単一カタログで集中管理し、プロフィールのスキルチップ、研究・ブログのタグ、プロジェクト技術スタックで再利用できるようにする。
- 研究とブログを統合した「研究・ブログ」モデルを定義し、外部公開記事を参照する構造へ移行する。
- カレンダー連携と通知を備えたお問い合わせ予約フローを実装し、予約状況をユーザー自身でも確認できるようにする。
- ホーム画面をコンテンツブロック／チップ表示に再構成し、プロフィールと連動したナビゲーション体験を提供する。

## スコープと前提
- 対象フェーズは本番運用を見据えた実装フェーズ。MySQL / Firestore / インメモリの 3 実装を共通データモデルに合わせる。
- 既存のブログテーブルは廃止し、新しい研究・ブログテーブルへマイグレーションする。
- 旧スキルテーブル・プロジェクト技術テーブルは統合カタログを採用し、不要な列は削除する。
- DB には ISO8601 互換の `DATETIME(3)` を採用し、終了日 `NULL` で現在在籍を表現する。
- LocalizedText は `{ ja?: string, en?: string }` の JSON を API 層で扱い、MySQL では `*_ja`/`*_en` 列で保持する。

## データモデル

### 技術セット (`tech_catalog`)

| カラム | 型 | 制約 | 説明 |
| --- | --- | --- | --- |
| `id` | BIGINT UNSIGNED | PK, AUTO_INCREMENT | 技術エントリ ID |
| `slug` | VARCHAR(64) | UNIQUE, NOT NULL | 英小文字＋ハイフンで構成される識別子 |
| `display_name` | VARCHAR(128) | NOT NULL | 表示名（多言語不要） |
| `category` | VARCHAR(64) | NULL | 任意カテゴリ（例: backend, ml） |
| `level` | ENUM('beginner','intermediate','advanced') | NOT NULL | 習熟度 |
| `icon` | VARCHAR(128) | NULL | アイコン URL / Emoji |
| `sort_order` | INT | DEFAULT 0 | 表示順制御 |
| `is_active` | TINYINT(1) | DEFAULT 1 | 表示可否 |
| `created_at` / `updated_at` | DATETIME(3) | DEFAULT CURRENT_TIMESTAMP(3) | 監査用 |

関連テーブル: `tech_relationships`

| カラム | 型 | 制約 | 説明 |
| --- | --- | --- | --- |
| `id` | BIGINT UNSIGNED | PK, AUTO_INCREMENT | 関連 ID |
| `entity_type` | ENUM('profile_section','project','research_blog') | NOT NULL | 紐付け対象 |
| `entity_id` | BIGINT UNSIGNED | NOT NULL | 対象エンティティ ID |
| `tech_id` | BIGINT UNSIGNED | FK (`tech_catalog.id`) | 技術 ID |
| `context` | ENUM('primary','supporting') | DEFAULT 'primary' | 表示優先度 |
| `note` | VARCHAR(255) | NULL | 追加説明 |
| `sort_order` | INT | DEFAULT 0 | 並び順 |
| `created_at` | DATETIME(3) | DEFAULT CURRENT_TIMESTAMP(3) | 作成日時 |

`entity_type = 'profile_section'` の場合は後述の `profile_tech_sections.id` を参照する。

### プロフィール領域

#### `profiles`

| カラム | 型 | 制約 | 説明 |
| --- | --- | --- | --- |
| `id` | BIGINT UNSIGNED | PK, AUTO_INCREMENT | プロフィール ID |
| `display_name` | VARCHAR(255) | NOT NULL | 表示名（ホーム 1 段目） |
| `headline_ja` / `headline_en` | VARCHAR(255) | NULL | キャッチコピー |
| `summary_ja` / `summary_en` | TEXT | NULL | 自己紹介文 |
| `avatar_url` | VARCHAR(512) | NULL | プロフィール画像 |
| `location_ja` / `location_en` | VARCHAR(255) | NULL | 所在地 |
| `theme_mode` | ENUM('light','dark','system') | DEFAULT 'system' | テーマモード |
| `theme_accent_color` | VARCHAR(32) | NULL | アクセントカラー |
| `lab_name_ja` / `lab_name_en` | VARCHAR(255) | NULL | 研究室名 |
| `lab_advisor_ja` / `lab_advisor_en` | VARCHAR(255) | NULL | 指導教員 |
| `lab_room_ja` / `lab_room_en` | VARCHAR(255) | NULL | 居室番号等 |
| `lab_url` | VARCHAR(512) | NULL | 研究室 URL |
| `created_at` / `updated_at` | DATETIME(3) | DEFAULT CURRENT_TIMESTAMP(3) | 監査用 |

#### `profile_affiliations`

| カラム | 型 | 制約 | 説明 |
| --- | --- | --- | --- |
| `id` | BIGINT UNSIGNED | PK, AUTO_INCREMENT | 行 ID |
| `profile_id` | BIGINT UNSIGNED | FK (`profiles.id`), NOT NULL | プロフィール参照 |
| `kind` | ENUM('affiliation','community') | NOT NULL | 所属 or コミュニティ |
| `name` | VARCHAR(255) | NOT NULL | 名称 |
| `url` | VARCHAR(512) | NULL | 公式リンク |
| `started_at` | DATETIME(3) | NOT NULL | 開始日時 |
| `description_ja` / `description_en` | VARCHAR(255) | NULL | 補足説明 |
| `sort_order` | INT | DEFAULT 0 | 表示順 |

終了日は保持せず、現在所属のみを管理する（要求 4）。

#### `profile_work_history`

| カラム | 型 | 制約 | 説明 |
| --- | --- | --- | --- |
| `id` | BIGINT UNSIGNED | PK, AUTO_INCREMENT | 行 ID |
| `profile_id` | BIGINT UNSIGNED | FK (`profiles.id`), NOT NULL | プロフィール参照 |
| `organization_ja` / `organization_en` | VARCHAR(255) | NOT NULL | 組織名 |
| `role_ja` / `role_en` | VARCHAR(255) | NOT NULL | 役割 |
| `summary_ja` / `summary_en` | TEXT | NULL | 詳細説明 |
| `started_at` | DATETIME(3) | NOT NULL | 開始日時 |
| `ended_at` | DATETIME(3) | NULL | 終了日時（NULL = 現在） |
| `external_url` | VARCHAR(512) | NULL | 関連 URL |
| `sort_order` | INT | DEFAULT 0 | 表示順 |

#### `profile_social_links`

| カラム | 型 | 制約 | 説明 |
| --- | --- | --- | --- |
| `id` | BIGINT UNSIGNED | PK, AUTO_INCREMENT | 行 ID |
| `profile_id` | BIGINT UNSIGNED | FK (`profiles.id`), NOT NULL | プロフィール参照 |
| `provider` | ENUM('github','zenn','linkedin','x','email','other') | NOT NULL | サービス種別 |
| `label_ja` / `label_en` | VARCHAR(255) | NULL | 表示ラベル |
| `url` | VARCHAR(512) | NOT NULL | リンク |
| `is_footer` | TINYINT(1) | DEFAULT 0 | フッター表示フラグ |
| `sort_order` | INT | DEFAULT 0 | 表示順 |

GitHub / Zenn / LinkedIn の 3 件は必須で `is_footer=1` とする（要求 7）。

#### `profile_tech_sections`

| カラム | 型 | 制約 | 説明 |
| --- | --- | --- | --- |
| `id` | BIGINT UNSIGNED | PK, AUTO_INCREMENT | セクション ID |
| `profile_id` | BIGINT UNSIGNED | FK (`profiles.id`), NOT NULL | プロフィール参照 |
| `title_ja` / `title_en` | VARCHAR(255) | NULL | セクション名 |
| `layout` | ENUM('chips','list') | DEFAULT 'chips' | 表示方法（チップ表示対応） |
| `breakpoint` | VARCHAR(32) | DEFAULT 'lg' | 横並び切替のブレークポイント |
| `sort_order` | INT | DEFAULT 0 | 表示順 |

`tech_relationships` と組み合わせてスキルチップを制御する。画面幅が一定以上の場合は `profile_work_history` と並列表示する（要求 6）。

### 研究・ブログ領域

#### `research_blog_entries`

| カラム | 型 | 制約 | 説明 |
| --- | --- | --- | --- |
| `id` | BIGINT UNSIGNED | PK, AUTO_INCREMENT | エントリ ID |
| `slug` | VARCHAR(128) | UNIQUE, NOT NULL | URL スラッグ |
| `kind` | ENUM('research','blog') | NOT NULL | コンテンツ種別（要求 2） |
| `title_ja` / `title_en` | VARCHAR(255) | NOT NULL | タイトル |
| `overview_ja` / `overview_en` | TEXT | NULL | 概要 |
| `outcome_ja` / `outcome_en` | TEXT | NULL | 成果 |
| `outlook_ja` / `outlook_en` | TEXT | NULL | 展望 |
| `external_url` | VARCHAR(512) | NOT NULL | 本文公開先 URL |
| `published_at` | DATETIME(3) | NOT NULL | 公開日時 |
| `highlight_image_url` | VARCHAR(512) | NULL | カバー画像 |
| `image_alt_ja` / `image_alt_en` | VARCHAR(255) | NULL | 画像代替テキスト |
| `created_at` / `updated_at` | DATETIME(3) | DEFAULT CURRENT_TIMESTAMP(3) | 監査用 |
| `is_draft` | TINYINT(1) | DEFAULT 0 | 下書きフラグ |

補助テーブル:

- `research_blog_tags`: 任意タグ（`tag` VARCHAR(64)）
- `research_blog_links`: 追加リンク（`label` Localized, `url`, `type` ENUM('paper','slides','video','code','external')）
- `research_blog_assets`: 画像・キャプション管理
- `tech_relationships` (`entity_type='research_blog'`)

技術タグは `tech_catalog` を選択して紐付ける（要求 8）。

### プロジェクト領域

#### `projects`

| カラム | 型 | 制約 | 説明 |
| --- | --- | --- | --- |
| `id` | BIGINT UNSIGNED | PK, AUTO_INCREMENT | プロジェクト ID |
| `slug` | VARCHAR(128) | UNIQUE, NOT NULL | URL スラッグ |
| `title_ja` / `title_en` | VARCHAR(255) | NOT NULL | プロジェクト名 |
| `summary_ja` / `summary_en` | TEXT | NOT NULL | 概要 |
| `description_ja` / `description_en` | LONGTEXT | NULL | 詳細 |
| `cover_image_url` | VARCHAR(512) | NULL | カバー画像 |
| `primary_link_url` | VARCHAR(512) | NULL | メインリンク |
| `period_start` | DATE | NULL | 開始日 |
| `period_end` | DATE | NULL | 終了日 |
| `created_at` / `updated_at` | DATETIME(3) | DEFAULT CURRENT_TIMESTAMP(3) | 監査用 |
| `published` | TINYINT(1) | DEFAULT 0 | 公開フラグ |
| `highlight` | TINYINT(1) | DEFAULT 0 | ホーム強調表示 |
| `sort_order` | INT | DEFAULT 0 | 表示順 |

補助テーブル:

- `project_links`: 複数リンク（`label` Localized, `url`, `type` ENUM('repo','demo','article','slides','other')）
- `tech_relationships` (`entity_type='project'`)

一覧では `title` / `summary` / 技術名 / `created_at` / `primary_link_url` を利用する（要求 9）。

### ホーム画面設定

#### `home_page_config`

| カラム | 型 | 制約 | 説明 |
| --- | --- | --- | --- |
| `id` | BIGINT UNSIGNED | PK, AUTO_INCREMENT | 設定 ID |
| `profile_id` | BIGINT UNSIGNED | FK (`profiles.id`), NOT NULL | 参照プロフィール |
| `hero_subtitle_ja` / `hero_subtitle_en` | VARCHAR(255) | NULL | 1 段目補助テキスト |
| `updated_at` | DATETIME(3) | DEFAULT CURRENT_TIMESTAMP(3) | 更新時刻 |

#### `home_quick_links`

| カラム | 型 | 制約 | 説明 |
| --- | --- | --- | --- |
| `id` | BIGINT UNSIGNED | PK, AUTO_INCREMENT | 行 ID |
| `config_id` | BIGINT UNSIGNED | FK (`home_page_config.id`), NOT NULL | 親設定 |
| `section` | ENUM('profile','research_blog','projects','contact') | NOT NULL | 対象セクション（要求 5） |
| `label_ja` / `label_en` | VARCHAR(255) | NOT NULL | 見出し |
| `description_ja` / `description_en` | TEXT | NULL | 説明 |
| `cta_ja` / `cta_en` | VARCHAR(128) | NOT NULL | ボタン文言 |
| `target_url` | VARCHAR(512) | NOT NULL | 遷移先 |
| `sort_order` | INT | DEFAULT 0 | 表示順 |

#### `home_chip_sources`

| カラム | 型 | 制約 | 説明 |
| --- | --- | --- | --- |
| `id` | BIGINT UNSIGNED | PK, AUTO_INCREMENT | 行 ID |
| `config_id` | BIGINT UNSIGNED | FK (`home_page_config.id`), NOT NULL | 親設定 |
| `source_type` | ENUM('affiliation','community','tech') | NOT NULL | チップ元 |
| `limit_count` | INT | DEFAULT 6 | 表示数 |
| `label_ja` / `label_en` | VARCHAR(255) | NULL | セクション見出し |
| `sort_order` | INT | DEFAULT 0 | 表示順 |

これにより 3 段目以降で所属 / コミュニティ / スキルチップを任意構成できる（要求 5）。

### お問い合わせ・予約領域

#### `contact_form_settings`

| カラム | 型 | 制約 | 説明 |
| --- | --- | --- | --- |
| `id` | BIGINT UNSIGNED | PK, AUTO_INCREMENT | 設定 ID |
| `hero_title_ja` / `hero_title_en` | VARCHAR(255) | NULL | フォーム見出し |
| `hero_description_ja` / `hero_description_en` | TEXT | NULL | 概要 |
| `topics` | JSON | NOT NULL | 相談トピック配列 |
| `consent_text_ja` / `consent_text_en` | TEXT | NOT NULL | 同意文言 |
| `minimum_lead_hours` | INT | DEFAULT 24 | 予約受付の最低猶予 |
| `recaptcha_public_key` | VARCHAR(128) | NULL | reCAPTCHA |
| `support_email` | VARCHAR(255) | NOT NULL | サポート窓口 |
| `calendar_timezone` | VARCHAR(64) | NOT NULL | タイムゾーン |
| `google_calendar_id` | VARCHAR(255) | NULL | 連携カレンダー ID |
| `booking_window_days` | INT | DEFAULT 30 | 予約可能期間 |
| `created_at` / `updated_at` | DATETIME(3) | DEFAULT CURRENT_TIMESTAMP(3) | 監査用 |

#### `meeting_reservations`

| カラム | 型 | 制約 | 説明 |
| --- | --- | --- | --- |
| `id` | BIGINT UNSIGNED | PK, AUTO_INCREMENT | 予約 ID |
| `name` | VARCHAR(255) | NOT NULL | 予約者名 |
| `email` | VARCHAR(255) | NOT NULL | 予約者メール |
| `topic` | VARCHAR(255) | NULL | 選択トピック |
| `message` | TEXT | NULL | 補足メッセージ |
| `start_at` | DATETIME(3) | NOT NULL | 予約開始 |
| `end_at` | DATETIME(3) | NOT NULL | 予約終了 |
| `duration_minutes` | INT | NOT NULL | 所要時間 |
| `google_event_id` | VARCHAR(255) | NULL | Google Calendar イベント |
| `google_calendar_status` | ENUM('pending','confirmed','declined','cancelled') | DEFAULT 'pending' | Google 側ステータス |
| `status` | ENUM('pending','confirmed','cancelled') | DEFAULT 'pending' | 社内フロー状態 |
| `confirmation_sent_at` | DATETIME(3) | NULL | 招待送信日時 |
| `last_notification_sent_at` | DATETIME(3) | NULL | 直近通知日時 |
| `lookup_hash` | CHAR(64) | NOT NULL | 氏名＋メールのハッシュ（要求 10） |
| `cancellation_reason` | TEXT | NULL | キャンセル理由 |
| `created_at` / `updated_at` | DATETIME(3) | DEFAULT CURRENT_TIMESTAMP(3) | 監査用 |

`lookup_hash` は `SHA2(CONCAT(LOWER(email), ':', LOWER(name)), 256)` を想定し、メールアドレスと氏名での本人確認検索を高速化する。

既存テーブル `blacklist` / `schedule_blackouts` は継続利用し、予約時に重複を排除する。

#### `meeting_notifications`

予約通知ログを保持し、再送制御を行う。

| カラム | 型 | 制約 | 説明 |
| --- | --- | --- | --- |
| `id` | BIGINT UNSIGNED | PK, AUTO_INCREMENT | ログ ID |
| `reservation_id` | BIGINT UNSIGNED | FK (`meeting_reservations.id`), NOT NULL | 予約参照 |
| `notification_type` | ENUM('confirmation_email','reminder_email','calendar_invite','cancellation_email') | NOT NULL | 通知種別 |
| `status` | ENUM('pending','sent','failed') | DEFAULT 'pending' | 送信状態 |
| `error_message` | TEXT | NULL | 失敗理由 |
| `created_at` | DATETIME(3) | DEFAULT CURRENT_TIMESTAMP(3) | 記録時刻 |

### インデックス設計
- `tech_catalog.slug`, `projects.slug`, `research_blog_entries.slug` にユニークインデックス。
- `profile_affiliations` は `(profile_id, kind, sort_order)` の複合インデックスで取得効率化。
- `tech_relationships` は `(entity_type, entity_id)`、`tech_id` にそれぞれインデックス。
- `meeting_reservations` は `(start_at)`, `(lookup_hash)`, `(email)` インデックスを追加。

## API / サービス仕様
- `GET /api/v1/public/home`: `profiles`, `home_page_config`, 技術チップ情報を結合しホーム表示データを返却。
- `GET /api/v1/public/profile`: プロフィール JSON を返却。所属／コミュニティは `profile_affiliations` を `kind` でフィルタし、経歴は終了日 NULL で現在在籍を表現。
- `GET /api/v1/public/research-blog`: `kind` フィルタ、タグ／技術での絞り込み、公開日の並び替えに対応。
- `GET /api/v1/public/projects`: 公開フラグのみ返却し、技術タグは `tech_catalog` を解決した配列を含める。
- `GET /api/v1/public/contact/config`: `contact_form_settings` を公開し、カレンダー設定・トピック・同意文言を取得。
- `GET /api/v1/public/contact/availability`: `meeting_reservations`, `schedule_blackouts`, Google Calendar のビジー情報を統合して空き枠を返却。
- `POST /api/v1/public/contact/reservations`: 予約登録後に Google Calendar イベント作成とメール通知（確認・カレンダー招待）を発火。予約確定時に `lookup_hash` を返し、ユーザーは `GET /api/v1/public/contact/reservations/{hash}` で確認できる。
- 管理 API（`/api/admin/*`）はプロフィール・研究・ブログ・プロジェクト・お問い合わせ設定の CRUD とファイルアップロードに対応。

## 管理 SPA 要件
- プロフィール編集画面は所属・コミュニティを FieldArray で管理し、開始日時は時刻含む入力（ISO8601）とする。
- 経歴フォームは終了日を空欄の場合に `NULL` を送信し、現在在籍を「Present」と表示。
- スキルセット編集は技術カタログを検索／選択し、セクション単位で並べ替え可能。技術カタログ自体の CRUD を別画面で提供。
- 研究・ブログフォームでは `kind` 選択、技術タグ（`tech_relationships`）、任意タグ追加、成果／展望のリッチテキスト入力、外部リンクのバリデーションを行う。
- プロジェクト編集は期間・公開状態・強調表示を設定し、リンク種別を選択式で入力。
- お問い合わせ設定画面はカレンダー ID、招待メール送信設定、リードタイム、通知テンプレートを管理。カレンダー連携テストの即時実行ボタンを提供。

## 公開 SPA 表示要件
- ホーム: 1 段目に `profiles.display_name` と `headline`。2 段目に `home_quick_links` の大型カードを表示し、3 段目以降で所属／コミュニティ／スキルのチップ行を出し分ける。
- プロフィール: 「所属・コミュニティ」「研究室」を上部で縦配置。画面幅が `profile_tech_sections.breakpoint` 以上では「経歴」「スキルセット」を横並びグリッド、それ以外は縦並び。GitHub / Zenn / LinkedIn はフッターに常時表示。
- 研究・ブログ: `kind` をタグとして明示し、技術タグを `tech_catalog` からチップ表示。概要→成果→展望→リンク→公開日時→画像の順で表示。
- プロジェクト: 一覧カードに概要、技術チップ、作成日、主要リンクを表示。ハイライトはホーム用レイアウトに流用。
- お問い合わせ: カレンダー UI（外部コンポーネント）で空き枠を表示し、予約確定で Google カレンダー招待・確認メール送信。`lookup_hash` で予約照会を提供。

## テスト指針
- MySQL / Firestore 向けリポジトリの CRUD と整合性をインテグレーションテストでカバー。
- 管理フォーム: FieldArray の追加・削除、技術カタログ検索、終了日 `NULL` シナリオをユニットテスト。
- 公開 SPA: Playwright でホームレイアウトのレスポンシブ切替、研究・ブログタグ表示、フッターリンク表示を検証。
- お問い合わせ: API モックで Google カレンダー失敗時のリトライ、ブラックリスト検知、予約重複防止をテスト。

## マイグレーションと運用
- 旧 `profile_*` テーブルから `profiles` ほか新テーブルへデータ移行する SQL / スクリプトを提供。終了日が存在しない行は `NULL` に変換。
- 旧 `blog_posts` / `research` は `research_blog_entries` へ統合し、`kind` を `blog` or `research` にマッピング。本文は外部公開済み URL を `external_url` に移し、Markdown 本文列は廃止。
- 既存プロジェクトの技術ラベルは `tech_catalog` に登録し、`tech_relationships` へ再紐付け。レベルは初期値として `intermediate` を設定し、後続で編集可能にする。
- 予約テーブルは `lookup_hash` をバッチで生成し、既存レコードを更新。Google イベント ID が欠落している場合は `status='pending'` として再同期ジョブを実行。
- 移行後は旧テーブルを削除し、Terraform / Firestore スキーマも同内容へ更新する。
