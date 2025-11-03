# コンテンツモデル再構築 移行計画（Phase 1）

## 目的
- 旧テーブル（`profile`, `projects`, `research`, `blog_posts`, `meetings`, `contact_messages` など）から新スキーマへ安全に切り替え、データ損失なく移行する。
- Firestore / MySQL / Terraform で一貫したコレクション・テーブル命名と設定を適用する。
- 本番切替後、旧データを参照しながら段階的にクリーンアップできる状態を作る。

## 前提条件
- `deploy/mysql/schema.sql` および `deploy/mysql/migrations/20240315_content_model_refactor.sql` がリポジトリに適用済みであること。
- バックアップ（Cloud SQL 自動バックアップ + 手動スナップショット）を直前に取得済みであること。
- Firestore は同一プロジェクト内でコレクションプレフィックスにより環境分離されていること。
- 予約／通知のリアルタイム連携を一時停止（メンテナンスモード）しておくこと。

## 移行フロー概要
1. **事前準備**
   - `gcloud sql backups create` で手動バックアップ取得。
   - `deploy/mysql/migrations/20240315_content_model_refactor.sql` をステージング DB へ適用し、アプリの smoke test を実施。
   - Firestore では `deploy/firestore/collections.yaml` と `deploy/firestore/indexes.yaml` をレビューし、必要なアクセストークン/ロールを確認。

2. **スキーマ適用**
   - 本番 Cloud SQL に対してメンテナンス時間を確保し、アプリを read-only またはメンテナンスモードに切り替え。
   - `mysql -u <user> -p -h <host> <database> < deploy/mysql/migrations/20240315_content_model_refactor.sql` を実行。旧テーブルは `legacy_*` テーブルとして退避される。
   - CI 用スクリプト `scripts/db/apply_migrations.sh` と同等の処理をローカルで実行し、スキーマ適用が成功することを確認。

3. **データ移送**
   - `backend/cmd/tools/contentmodelrefactor` を用いて以下を実施（`--dry-run` でプレビュー可能）:
     1. 旧 `legacy_profile` から `profiles` へコピー（`display_name` は `COALESCE(name_ja, name_en)`）。
     2. `legacy_profile_skills` から一意なスキルを抽出し `tech_catalog` へ登録。レベルは暫定的に `intermediate`。
     3. 新規 `profile_tech_sections` を生成し、`tech_relationships` に紐付け。
     4. `legacy_projects` / `legacy_project_tech_stack` を `projects` + `tech_relationships`へ移送、`slug` 未設定の場合は slugify する。
     5. `legacy_research` / `legacy_blog_posts` を `research_blog_entries` へ統合。`research` → kind=`research`、`blog_posts` → kind=`blog`。本文は `external_url` へ移す（GitHub/Zenn など外部 URL を設定）。
     6. `legacy_meetings` を `meeting_reservations` へコピーし、`lookup_hash` を `SHA2(LOWER(email) || ':' || LOWER(name), 256)` で生成。`calendar_event_id` は `google_event_id` へ移行。
     7. `legacy_contact_messages` を `meeting_notifications` へはコピーせず、必要に応じてバックアップとして保持。
   - データ移送後、各テーブルで件数一致と Spot Check（`SELECT *`）を行う。

4. **Firestore 構造反映**
   - `gcloud firestore indexes composite update deploy/firestore/indexes.yaml` を実行し、複合インデックスを登録。
   - 管理 SPA から API を通じてデータをプッシュするため、初期データは MySQL を真実ソースとし、バックエンド同期機能で Firestore へコピーする（後続フェーズ）。

5. **アプリケーション更新**
   - フェーズ 2 以降で実装するバックエンド/API 更新をデプロイし、新テーブルを参照することを確認。
   - メンテナンスモードを解除し、Smoke Test → Regression Test → 予約フォーム実動確認を順に実施。

6. **フォローアップ**
   - `legacy_*` テーブルは 2 週間程度保持し、RCA 完了後に削除する（Terraform / SQL スクリプトで追って対応）。
   - Firestore 旧コレクションが存在する場合は、移行完了後にサンプルデータをエクスポートしてから削除。

## スクリプト雛形の使い方
- `backend/cmd/tools/contentmodelrefactor`:
  ```bash
  cd backend
  go run ./cmd/tools/contentmodelrefactor --dsn "$APP_DATABASE_DSN" --dry-run
  go run ./cmd/tools/contentmodelrefactor --dsn "$APP_DATABASE_DSN" --apply
  ```
  `--dry-run` は予定されている INSERT/UPDATE 件数のみを計算し、実際の書き込みは行わない。
- `scripts/db/apply_migrations.sh`:
  CI で利用している MySQL 8.0 コンテナによるスキーマ検証を、本番適用前にローカルでも再現できる。

## 検証ポイント
- `SELECT COUNT(*)` が旧テーブルと新テーブルで一致しているか。
- `tech_catalog.slug` が一意であり、非 ASCII 文字の slug は手動確認済みか。
- `meeting_reservations.lookup_hash` で事前予約が引けるか（API 経由）。
- Firestore インデックスが `READY` 状態になっているか (`gcloud firestore indexes composite list`)。
- 管理 SPA からスキル追加 → 技術カタログ参照 → 公開 SPA 反映まで一連の流れが成功するか。

## ロールバック戦略
- MySQL: `legacy_*` テーブルへリネーム済みのため、再リネーム + 旧アプリケーションリリースで即時復旧可能。
- Firestore: 新規コレクションは prefix 付きで作成されるため、旧 prefix に戻すだけで復旧可能。
- アプリケーション: Cloud Run の旧リビジョンへ即時ロールバック。予約 API はメンテ中に停止していたため、整合性は維持される。

## 追加 TODO
- Terraform に Firestore インデックス管理を完全統合する（Phase 1 でモジュールを作成済み、環境に適用する）。
- Phase 2 でバックエンド同期ジョブ（MySQL → Firestore）を実装し、マルチデータストアの整合性テストを自動化する。
