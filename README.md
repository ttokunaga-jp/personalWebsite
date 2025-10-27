# Personal Website

## 概要
このリポジトリは Go (Gin) 製バックエンド API、React + TypeScript + Tailwind による SPA（公開用・管理用）および GCP（Cloud Run / Cloud SQL / VPC）向け Terraform 構成を含む環境セットアップ用テンプレートです。品質・可観測性・セキュリティで妥当な初期設定を整え、すぐに開発を開始できる状態を提供します。

## リポジトリ構成
```
.
├── backend/                 # Go API スケルトン (Gin + Fx + Viper)
│   ├── cmd/server           # アプリケーションエントリポイント
│   ├── internal             # ハンドラや DI を含む内部レイヤ
│   ├── config               # 設定ファイル例
│   ├── Dockerfile           # バックエンド用 Dockerfile
│   └── .env.example         # 環境変数テンプレート
├── frontend/                # React ワークスペース (公開 SPA + 管理 SPA)
│   ├── apps/public          # 公開向け SPA
│   ├── apps/admin           # 管理者向け SPA
│   ├── packages/shared      # 共有 UI / API クライアント
│   ├── package.json         # pnpm ワークスペース定義
│   ├── pnpm-workspace.yaml  # ワークスペースマッピング
│   └── Dockerfile           # フロントエンドビルド用 Dockerfile
├── deploy/
│   └── docker/nginx         # ローカル用 nginx リバースプロキシ設定
├── terraform/               # Terraform 初期構成 (ネットワーク + Cloud Run)
│   ├── environments/dev     # 環境別コンポジション
│   └── modules              # Cloud Run / Network モジュール
├── docker-compose.yml       # ローカル開発用 (API + SPA + MySQL)
└── Makefile                 # lint/test/build/up/down 用のタスク定義
```

## 前提条件
- Go 1.22 以上
- Node.js 20 以上（Corepack / pnpm 利用）
- Docker & Docker Compose v2
- Terraform 1.5 以上

## 初期設定手順
1. 環境変数テンプレートをコピーして秘密情報を設定します。
   ```bash
   cp backend/.env.example backend/.env
   cp frontend/.env.example frontend/.env
   ```
2. 依存関係をインストールします。
   ```bash
   make deps
   ```
   - `make deps-backend`: `go mod tidy` を実行し `go.sum` を生成します。
   - `make deps-frontend`: Corepack が利用可能なら有効化し、無い場合は `npx pnpm@8.15.4 install` を自動利用します (`pnpm-lock.yaml` を生成)。
3. (任意) Husky のフックを初期化します。
   ```bash
   cd frontend && pnpm prepare
   ```

## 開発フロー
- **Lint**: `make lint`
- **Test**: `make test`
- **Build**: `make build`
- **Docker compose up**: `make up` または `make up-detached`
- **Docker compose down**: `make down`

フロントエンドの開発サーバーを個別に立ち上げる場合:
```bash
cd frontend
pnpm --filter @personal-website/public dev
pnpm --filter @personal-website/admin dev
```

## 公開 SPA 実装メモ
- **実装済みページ**: Home / Profile / Research / Projects / Contact を詳細化し、バックエンドの `/v1/public/*` API からデータを取得します。`apps/public/src/modules/public-api` に型付きクライアントと `useApiResource` フックを追加しました。
- **主要機能**:
  - Home: プロフィール概要・所属・SNS を API から描画しつつ、Go API `/health` のステータスをヘルス表示。
  - Profile: 所属・研究室・職歴・スキル・コミュニティをカード表示し、ローディング/空データ時のプレースホルダを備えています。
  - Research: タグフィルタと Markdown/HTML 表示 (`MarkdownRenderer`) を実装し、画像やリンクのサニタイズを実施。
  - Projects: 技術スタックフィルタとカード UI（期間整形・リンクバッジ付き）でプロジェクトを一覧化。
  - Contact: 予約枠表示、フォームバリデーション、reCAPTCHA v3 トークン取得、`/v1/public/contact/bookings` 送信までをサポート。
- **セットアップの注意**:
  - `.env` に `VITE_API_BASE_URL` と `VITE_RECAPTCHA_SITE_KEY` を設定してください。未設定の場合、Contact ページの送信が失敗します。
  - 研究コンテンツの HTML を提供する場合でも、`http/https/mailto` 以外のプロトコルは描画されません（`MarkdownRenderer` で制御）。
- **テスト/ビルド**:
  - `pnpm --filter @personal-website/public test` で UI テスト（Projects フィルタ、Contact フォーム検証など）を実行。React Router v7 への移行警告（`startTransition`）は既知のものです。
  - `pnpm --filter @personal-website/public lint` / `pnpm --filter @personal-website/public build` で静的解析・ビルド確認。
  - Docker でのビルド検証は `docker compose build --no-cache frontend` を推奨します。
- **UX メモ**: ARIA ラベルとキーボード操作に対応済みです。`pnpm --filter @personal-website/public dev` で起動し、モバイル幅も合わせて確認してください。

## Docker Compose メモ
- バックエンド API: `http://localhost:8100`
- フロントエンド (nginx): `http://localhost:3000` で `/api` をバックエンドにプロキシ
- MySQL: `localhost:23306`（認証情報は `docker-compose.yml` を参照）

### 初期データベース
`deploy/mysql/init` 配下の SQL がコンテナ起動時に順番に実行され、`google_oauth_tokens` や `blacklist` など管理機能で利用するテーブルを自動生成します。すでに `mysql_data` ボリュームがある状態でスキーマを更新したい場合は、以下でボリュームを一度破棄してください（永続化データは消えます）。

```bash
docker compose down -v
docker compose up -d
```

## 管理 API / GUI

認証済みの管理者のみがアクセスできる `/api/admin` 配下のエンドポイントを実装しました。主な REST エンドポイントは次の通りです。

| メソッド | パス | 用途 |
| --- | --- | --- |
| `GET` | `/api/admin/summary` | 公開済み/下書き件数や予約状況のサマリー取得 |
| `GET` / `POST` / `PUT` / `DELETE` | `/api/admin/projects[:id]` | プロジェクトの CRUD |
| `GET` / `POST` / `PUT` / `DELETE` | `/api/admin/research[:id]` | 研究コンテンツの CRUD |
| `GET` / `POST` / `PUT` / `DELETE` | `/api/admin/blogs[:id]` | ブログ投稿の CRUD |
| `GET` / `POST` / `PUT` / `DELETE` | `/api/admin/meetings[:id]` | 予約（面談）情報の CRUD |
| `GET` / `POST` / `DELETE` | `/api/admin/blacklist[:id]` | ブラックリスト管理 |

管理者 SPA (`frontend/apps/admin`) は上記 API を利用してコンテンツや予約・ブラックリストを操作します。ローカル開発時は次のコマンドで起動できます。

```bash
cd frontend
pnpm --filter @personal-website/admin dev
```

`.env` で `VITE_API_BASE_URL` を指定している場合は、Cloud Run やリバースプロキシ環境に合わせて更新してください。

## Terraform 初期構成
最小の Terraform スタックが `terraform/environments/dev` に用意されています。
```bash
cd terraform/environments/dev
terraform init
terraform plan -var="project_id=<your-project>" -var="api_image=<artifact-registry-image>" -var="frontend_image=<artifact-registry-image>"
```
モジュール内容:
- 専用 VPC / サブネット / VPC Connector (`modules/network`)
- API とフロントエンド用 Cloud Run サービス (`modules/cloudrun/*`)

## テスト戦略スナップショット
- バックエンド: `go test ./...`（Testify ベースのテストを追加可能）
- フロントエンド: `vitest` + Testing Library (各ワークスペースで設定済み)
- Lint: Go vet + フォーマッタ整合性、ESLint (React / TypeScript プリセット)
- API スモークテスト: `make smoke-backend`（`TOKEN` または `ADMIN_TOKEN` を設定すると管理 API も検証）。別コンテナから実行する場合は `BASE_URL=http://backend:8100 make smoke-backend` のように明示的にエンドポイントを指定してください。

### CSRF トークンの手動確認
`/api/security/csrf` はダブルサブミットトークン方式を採用しています。Cookie には `ランダム値:有効期限UNIX秒:署名` の形式で保存されるため、手動テスト時は以下の点に注意してください。

```bash
# 1. トークンと Cookie を取得
curl -i http://localhost:8100/api/security/csrf

# 2. 応答ヘッダーの Set-Cookie をそのまま利用してリクエストを送る
curl -X POST http://localhost:8100/api/contact \
  -H "Content-Type: application/json" \
  -H "X-CSRF-Token: <Set-Cookie の先頭にあるランダム文字列>" \
  --cookie "ps_csrf=<Set-Cookie の値を完全に貼り付ける>" \
  -d '{"name":"Tester","email":"tester@example.com","message":"Hello"}'
```

Cookie 値から署名部分（`:<timestamp>:<signature>`）を削除すると検証に失敗し 403 が返るため、コピー漏れに注意してください。

## 次のステップ
- クリーンアーキテクチャ準拠のユースケース・リポジトリ実装を追加
- Terraform モジュールに Cloud SQL / Secret Manager / IAM 設定を拡張
- CI (GitHub Actions / Cloud Build) を整備し make タスク・Terraform plan を自動実行
