# コンテンツモデル v2 リポジトリ刷新メモ

## 背景
- `deploy/mysql/schema.sql` で定義した新テーブル群 (tech_catalog, profile_* 系, research_blog_entries など) を取り扱うため、旧来のシンプルなドメイン/リポジトリ構造では表現力が不足していた。
- 既存 API は v1 データを前提としているため、互換性を保ちつつ新スキーマを扱うクリーンな API を段階的に導入する必要がある。

## ドメインモデル整備
- `internal/model/content_v2.go` に新しい集約モデルを追加。
  - 技術カタログ (`TechCatalogEntry`, `TechMembership`) とコンテキスト情報を一元管理。
  - プロフィール (`ProfileDocument`) は所属/コミュニティ/経歴/技術セクション/ソーシャルリンクをフルで保持。
  - プロジェクト・研究ブログ・ホーム設定・お問い合わせ設定・予約情報をそれぞれ新しいスキーマに対応した構造で表現。
- 既存の `model.Project`, `model.Profile` などは旧 API 互換のために残しており、段階的に置き換えていく。

## リポジトリ実装
- **MySQL**: `internal/repository/mysql/content_profile.go`
  - 新テーブル群からプロフィールドキュメントを構築。
  - 付随テーブル (`profile_affiliations`, `profile_work_history`, `profile_social_links`, `tech_relationships`) を JOIN/複合クエリで取得し、`ProfileDocument` にマッピング。
  - `tech_relationships` を `sqlx.In` で展開し、技術カタログを同時に取り込む。
- **Firestore**: `internal/repository/firestore/content_profile.go`
  - `profiles/<primary>` ドキュメントを `ProfileDocument` に変換。
  - 既存の firestore util (`localizedDoc`) を流用して LocalizedText をマッピング。
- **In-memory**: `internal/repository/inmemory/content_profile.go`
  - 既存フィクスチャから必要情報を再構成し、新モデルを返す。
- **DI**: `repository/provider/repositories.go` に `NewContentProfileRepository` を追加し、 MySQL/Firestore/In-memory を切り替えられるようにした。

## 互換性と移行指針
- 現行 API/サービス層は旧モデル (`model.Profile`, `repository.ProfileRepository`) を利用し続けている。並行して `ContentProfileRepository` を導入し、新 API 実装時に差し替える。
- 仕様差分:
  - 旧 `Profile` はスキル配列やフォーカスエリアのみ。新 `ProfileDocument` では技術カタログ・所属情報などが増えている。
  - 技術タグは `tech_relationships` を通じて `TechCatalogEntry` と結合しており、旧 `TechStack []string` では表現できない。旧 API から新 API へ移行するフローで、`[]string` を `[]TechMembership` へマッピングするアダプタを用意する。
- 移行ステップ:
  1. 新リポジトリ/モデルをサービス層に組み込み、新 API (`/api/v1/public/*`) で `ProfileDocument` を返却するようにする。
  2. 旧 API は暫定的にアダプタで `ProfileDocument` → `model.Profile` にダウングレードして互換性を維持。
  3. 管理画面・公開画面のクライアント実装が新 JSON を消費するようになったら旧モデル/リポジトリを撤去。

## 今後のタスク
- プロジェクト・研究・お問い合わせ設定など、残りのエンティティについても `*_Document` 系モデルとリポジトリ実装を追加する。
- サービス層/ハンドラを `ContentProfileRepository` へ差し替え、レスポンススキーマを v2 に更新。
- Google カレンダー連携・通知ワークフローが `MeetingReservationV2` を扱えるようリファクタリング。
