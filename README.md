# Personal Website

## 概要
Go (Gin) 製の API と React + TypeScript 製 SPA（公開サイト / 管理画面）を中心に、Firestore・Google Calendar / Gmail・Terraform・GCP（Cloud Run / Cloud Build）までを一貫して扱う個人ポートフォリオサイトの基盤です。ローカル開発から CI/CD、自動デプロイ、セキュリティ／可観測性までを含むプロダクション志向のテンプレートとして利用できます。

## 目次
- [実装ステータス](#実装ステータス)
- [技術スタック](#技術スタック)
- [リポジトリ構成](#リポジトリ構成)
- [セットアップ](#セットアップ)
  - [前提条件](#前提条件)
  - [環境変数](#環境変数)
  - [初期セットアップ手順](#初期セットアップ手順)
- [ローカル開発](#ローカル開発)
- [フロントエンド](#フロントエンド)
- [バックエンド API](#バックエンド-api)
- [データ永続化](#データ永続化)
- [テストと品質保証](#テストと品質保証)
- [CI/CD](#cicd)
- [インフラ / Terraform](#インフラ--terraform)
- [セキュリティと可観測性](#セキュリティと可観測性)
- [ドキュメント](#ドキュメント)
- [トラブルシューティング](#トラブルシューティング)
- [今後の改善アイデア](#今後の改善アイデア)

## 実装ステータス
| 領域 | 実装状況 | 補足 |
| --- | --- | --- |
| フロントエンド（公開 SPA） | 完了 | Home / Profile / Research / Projects / Contact を実装。i18n、フォーム検証、予約枠表示、reCAPTCHA 連携を含む。|
| フロントエンド（管理 SPA） | MVP | プロジェクト・研究・ブログ・予約・ブラックリスト CRUD を提供。UI/UX と e2e 自動テストは今後拡充予定。|
| バックエンド API | 完了 | 公開 / 管理エンドポイント、Google OAuth + JWT、予約（カレンダー・通知）処理、Clean Architecture 風レイヤ分離済み。|
| 認証・権限 | MVP | Google OAuth, ドメイン/メール許可リスト, JWT, Admin Guard を実装。Secret 管理と本番向け WIF 設定は環境依存。|
| 予約・外部連携 | MVP | Google Calendar への予定挿入、Gmail API 経由の通知をサポート。トークン更新やバックアップ導線は継続改善対象。|
| テスト | 進行中 | Go/Vitest 単体テストと ESLint を整備。Playwright など E2E、自動ビジュアル回帰は未導入。|
| CI/CD | 完了 | GitHub Actions → Cloud Build → Cloud Run で lint/test/build/deploy を自動化。Workload Identity Federation を使用。|
| インフラ (Terraform) | ベース構築済み | `terraform/environments/dev` で VPC / Cloud Run モジュールを提供。Firestore や本番環境差分は追加実装の余地あり。|

## 技術スタック
- **バックエンド**: Go 1.22, Gin, Uber Fx, sqlx, Viper, Google API (OAuth / Calendar / Gmail), Prometheus instrumentation
- **フロントエンド**: React 18, TypeScript, Vite, pnpm ワークスペース, Tailwind CSS, React Router, React Query, i18next
- **テスト / 品質**: Go test, Vitest, Jest DOM, ESLint, Prettier, Husky + lint-staged, Playwright（E2E 準備済み）
- **データベース**: Cloud Firestore (Native mode)
- **インフラ**: Docker Compose, Terraform 1.5+, GCP（Cloud Run, Cloud Build, Artifact Registry, Secret Manager, Firestore, VPC Connector）

## リポジトリ構成
```
.
├── backend/                  # Go API
│   ├── cmd/server            # エントリポイント
│   ├── internal              # ハンドラ／サービス／リポジトリ／ミドルウェアなど
│   ├── config                # 設定ファイル (config.yaml / example)
│   ├── scripts               # 補助スクリプト
│   └── Dockerfile
├── frontend/                 # React + Vite ワークスペース
│   ├── apps/public           # 公開 SPA
│   ├── apps/admin            # 管理者 SPA
│   ├── packages/shared       # 共有 UI / API クライアント
│   ├── package.json / pnpm-workspace.yaml
│   └── Dockerfile
├── deploy/
│   ├── cloudbuild/           # Cloud Build 設定
│   └── docker/nginx          # ローカル用 nginx 設定
├── terraform/                # IaC (環境別構成 / modules)
├── docker-compose.yml        # ローカル開発 (API + SPA)
├── Makefile                  # lint/test/build 等のコマンド
└── docs/architecture-design.md
```

## セットアップ

### 前提条件
- Go 1.22 以上
- Node.js 20 以上（Corepack 利用可）
- pnpm 8.x（`corepack enable` で自動インストール可能）
- Docker / Docker Compose v2
- Terraform 1.5 以上
- GCP プロジェクト（Cloud Run / Artifact Registry / Secret Manager / Cloud Build を有効化済み）

### 環境変数

#### `backend/.env`
| 変数 | 説明 |
| --- | --- |
| `APP_SERVER_PORT` | ローカル起動時のポート。Docker Compose では 8100 を使用。 |
| `APP_SERVER_MODE` | Gin のモード (`debug`/`release`)。 |
| `APP_FIRESTORE_PROJECT_ID` | Firestore を使用する GCP プロジェクト ID。設定されていない場合は永続化を無効化。 |
| `APP_FIRESTORE_DATABASE_ID` | Firestore のデータベース ID（通常は `(default)`）。 |
| `APP_FIRESTORE_COLLECTION_PREFIX` | コレクションに付与するプレフィックス（環境ごとの名前空間分離に利用）。 |
| `APP_FIRESTORE_EMULATOR_HOST` | Firestore Emulator に接続する場合のホスト（例: `localhost:8080`）。 |
| `APP_AUTH_JWT_SECRET` | 管理者 JWT 用のシークレット。 |
| `APP_AUTH_STATE_SECRET` | Google OAuth の state/トークン暗号化に使用するシークレット。 |
| `APP_SECURITY_CSRF_SIGNING_KEY` | CSRF トークン署名キー。 |
| `APP_GOOGLE_CLIENT_ID` / `APP_GOOGLE_CLIENT_SECRET` | Google OAuth クライアント情報。 |
| `APP_GOOGLE_REDIRECT_URL` | Google OAuth のリダイレクト URL。Cloud Run 公開 URL の `/api/auth/callback` を指定。 |

追加の詳細設定は `backend/config/config.yaml` または環境変数 `APP_*` で上書きします（例: `APP_SECURITY_ENABLE_CSRF=false`）。予約ワークフロー向けには `APP_BOOKING_CALENDAR_ID` や `APP_BOOKING_NOTIFICATION_SENDER` なども利用できます。

#### `frontend/.env`
| 変数 | 説明 |
| --- | --- |
| `VITE_API_BASE_URL` | フロントエンドからの API ベース URL（ローカルは `/api`、本番は Cloud Run の公開 URL）。 |
| `VITE_I18N_FALLBACK` | i18next のフォールバック言語。 |

reCAPTCHA を利用する場合は GitHub Actions / Cloud Build 側で `VITE_RECAPTCHA_SITE_KEY` / `VITE_RECAPTCHA_SECRET` をシークレットとして設定します。

#### その他
- Terraform 変数は `terraform/environments/<env>/terraform.tfvars` で管理。
- Cloud Build は Secret Manager に保管した JWT / OAuth / CSRF / reCAPTCHA シークレットを参照します（後述）。

### 初期セットアップ手順
1. `.env` を作成:
   ```bash
   cp backend/.env.example backend/.env
   cp frontend/.env.example frontend/.env
   ```
2. 依存関係を取得:
   ```bash
   make deps
   ```
   - `make deps-backend`: `go mod tidy`
   - `make deps-frontend`: `pnpm install`（corepack 未導入の場合は `npx pnpm@8.15.4 install`）
3. 必要に応じてフックを初期化:
   ```bash
   cd frontend && pnpm prepare
   ```

## ローカル開発
- **バックエンドのみ起動**: `cd backend && go run ./cmd/server`
- **フロントエンド（公開）**: `cd frontend && pnpm --filter @personal-website/public dev`
- **フロントエンド（管理）**: `cd frontend && pnpm --filter @personal-website/admin dev`
- **フルスタック（Docker Compose）**:
  ```bash
  make up           # フォアグラウンド
  make up-detached  # バックグラウンド
  make down         # 停止 & 後片付け
  ```
  - API: http://localhost:8100
  - フロント (nginx): http://localhost:3000
  - Firestore が必要な場合は `gcloud beta emulators firestore start --host-port=localhost:8080` などでエミュレータを併用してください（`.env` の `APP_FIRESTORE_EMULATOR_HOST` を設定）。
- **ユーティリティ**:
  - `make build`: バックエンドバイナリ / フロント dist を生成
  - `make fmt`: Go / TypeScript のフォーマッタ実行
  - `make smoke-backend`: API スモークテスト（`BASE_URL` や `TOKEN` でカスタマイズ可）

## フロントエンド

### 公開 SPA (`frontend/apps/public`)
- ルーティング: Home / Profile / Research / Projects / Contact。React Router v6.30 を使用。
- データ取得: `packages/shared` の API クライアント経由で `/api/v1/public/*` エンドポイントを叩く。
- 予約・問い合わせ: reCAPTCHA v3 トークンを取得し、`/v1/public/contact/bookings` へ送信。レスポンス検証＆エラーハンドリングを実装。
- 国際化: `packages/shared/src/i18n` の設定で ja/en をサポート。
- テスト: `pnpm --filter @personal-website/public test`。React Router の v7 transition 警告は既知（React Router ドキュメント参照）。
- ビルド: `pnpm --filter @personal-website/public build`。

### 管理 SPA (`frontend/apps/admin`)
- 認証済み管理者向けの CRUD UI。ダッシュボードサマリ、プロジェクト / 研究 / ブログ / 予約 / ブラックリスト管理を実装。
- API: `/api/admin/*`。JWT を `Authorization: Bearer` ヘッダで付与。
- テスト: `pnpm --filter @personal-website/admin test`。現状はユニットレベル中心で、E2E は今後 Playwright 導入予定。
- スタイル: Tailwind + Headless UI コンポーネントをベース。

## バックエンド API

### 公開エンドポイント（`/api` および `/api/v1/public`）
| メソッド | パス | 説明 |
| --- | --- | --- |
| GET /api/health | ヘルスチェック。HEAD も対応。 |
| GET /api/profile | プロフィール情報の取得。 |
| GET /api/projects | 公開プロジェクト一覧。 |
| GET /api/research | 研究コンテンツ一覧。 |
| GET /api/contact/availability | 予約可能枠の一覧（Google Calendar + DB を考慮）。 |
| GET /api/contact/config | フォーム設定（トピック、リードタイム等）。 |
| POST /api/contact | お問い合わせ送信（メール通知を想定）。 |
| POST /api/contact/bookings | 予約作成（Calendar イベント作成、メール通知、DB 永続化）。 |
| GET /api/auth/login | Google OAuth URL を発行。 |
| GET /api/auth/callback | OAuth コールバックで JWT を発行。 |
| GET /api/security/csrf | CSRF トークン / ダブルサブミット Cookie を発行。 |

### 管理エンドポイント（`/api/admin/*`、JWT + AdminGuard 必須）
- サマリ: `GET /summary`
- プロジェクト: `GET/POST/PUT/DELETE /projects` (+ `/projects/:id`)
- 研究: 同上（`/research`）
- ブログ: `GET/POST/PUT/DELETE /blogs`
- 予約: `GET/POST/PUT/DELETE /meetings`
- ブラックリスト: `GET/POST /blacklist`, `DELETE /blacklist/:id`
- ヘルス: `GET /health`

### 認証・セキュリティ
- Google OAuth 2.0 + JWT（HS256）。ドメイン / メールの許可リストを設定可能。
- CSRF: ダブルサブミットトークン（`ps_csrf` Cookie + `X-CSRF-Token` ヘッダ）。
- レートリミット: デフォルト 120req/min（`APP_SECURITY_RATE_LIMIT_*` で調整）。
- セキュリティヘッダ: CSP / HSTS / Referrer-Policy / X-Content-Type-Options / X-Frame-Options。
- HTTPS リダイレクト、CORS 設定、リクエスト ID、構造化ログ、Prometheus メトリクス (`/metrics`)。
- 予約時: Google Calendar API への挿入、Gmail API 経由のメール送信。Circuit Breaker + Retry + Timeout を実装。

## データ永続化
- DB スキーマは `deploy/mysql/init` の SQL で初期化（コンテナ起動時に自動適用）。
- エンティティ例:
  - `profile`: プロフィール情報
  - `projects`, `research`: 公開コンテンツ
  - `meetings`: 予約（`status`, `calendar_event_id` を保持）
  - `blacklist`: 予約を拒否するメールアドレス
  - `google_oauth_tokens`: Google API 用トークンの暗号化保存
- リポジトリ実装: Firestore / In-memory の両方を実装し、テスト容易性を確保。

## テストと品質保証
- `make lint`: gofmt チェック + `go vet` + ESLint
- `make test`: `go test ./...` + `pnpm -r test`
- `pnpm --filter @personal-website/public test --watch`: 公開 SPA のウォッチモード
- `pnpm test:e2e`: Playwright（セットアップ後に有効。CI には未組み込み）
- `pnpm test:perf`: Lighthouse CI（ビルド後に実行。Node 18+ が必要）
- `make ci`: lint / test / build を一括実行（CI と同構成）
- 目標: Go カバレッジ 80% 以上・Playwright E2E・負荷テスト（k6）を今後整備

## CI/CD

### GitHub Actions (`.github/workflows/ci.yml`)
- トリガ: PR（main / develop）、push（main / develop）、手動 (`workflow_dispatch`)
- quality ジョブ: `make deps` → `make lint` → `make test` → `make build`
- deploy ジョブ: main への push または `workflow_dispatch`。Workload Identity Federation で GCP 認証後、Cloud Build を起動。

### Cloud Build (`deploy/cloudbuild/cloudbuild.yaml`)
1. Backend / Frontend の Docker イメージをビルドし、Artifact Registry に `:${SHORT_SHA}` で push。
2. Cloud Run (`<service>-<env>`) へデプロイ。Secret Manager から JWT / OAuth / CSRF シークレットを注入し、必要に応じて VPC Connector を接続。
3. トラフィック割合（`_BACKEND_TRAFFIC`, `_FRONTEND_TRAFFIC`）で段階的リリースが可能。

### GitHub Environments / Secrets
1. `staging`, `production` 環境を作成し、下記を Environment Variables に登録:
   - `CLOUD_RUN_REGION`, `CLOUD_BUILD_ARTIFACT_REPO`, `CLOUD_BUILD_ARTIFACT_LOCATION`
   - `CLOUD_RUN_BACKEND_SERVICE`, `CLOUD_RUN_FRONTEND_SERVICE`
   - `CLOUD_RUN_VPC_CONNECTOR`（必要な場合のみ）
   - `BACKEND_SERVICE_ACCOUNT_EMAIL`, `FRONTEND_SERVICE_ACCOUNT_EMAIL`
   - `FRONTEND_API_BASE_URL`
   - `BACKEND_GOOGLE_CLIENT_ID`
   - `FIRESTORE_DATABASE_ID`（必要な場合のみ）
   - `FIRESTORE_COLLECTION_PREFIX`（必要な場合のみ）
   - `BACKEND_TRAFFIC_PERCENT`, `FRONTEND_TRAFFIC_PERCENT`
2. 同環境の Secrets:
   - `BACKEND_SECRET_JWT`
   - `BACKEND_SECRET_STATE`
   - `BACKEND_SECRET_CSRF`
   - `BACKEND_SECRET_GOOGLE_CLIENT_SECRET`
   - `BACKEND_SECRET_RECAPTCHA`（任意）
   - `FRONTEND_SECRET_RECAPTCHA`（任意）
3. リポジトリ全体の Actions Secrets:
   - `GCP_PROJECT_ID`
   - `GCP_WORKLOAD_IDENTITY_PROVIDER`
   - `GCP_SERVICE_ACCOUNT_EMAIL`
4. 上記サービスアカウントに必要なロールを付与:
   - Cloud Build Editor / Cloud Run Admin / Artifact Registry Administrator
   - Service Account User / Secret Manager Secret Accessor
   - Artifact Registry リポジトリ（`${CLOUD_BUILD_ARTIFACT_LOCATION}`）の作成忘れに注意。未作成の場合は下記コマンドで事前に作成する:

     ```sh
     gcloud artifacts repositories create "${CLOUD_BUILD_ARTIFACT_REPO}" \
       --repository-format=docker \
       --location="${CLOUD_BUILD_ARTIFACT_LOCATION}"
     ```

### デプロイ検証
- `workflow_dispatch` で `environment=staging` を指定 → Cloud Build 実行
- Cloud Build ログでイメージ push や Cloud Run デプロイを確認
- Cloud Run のリビジョン / トラフィック配分をチェックし、ステージング URL で動作確認
- 本番リリース時は `environment=production` を選択し、必要であれば段階的トラフィック配分を設定

## インフラ / Terraform
- `terraform/environments/dev`: 開発環境向け構成（VPC, Subnet, Cloud Run サービスなど）
- `terraform/modules/*`:
  - `network`: VPC / サブネット / VPC Connector
  - `cloudrun`: デフォルトサービス設定（CPU/メモリ、最小/最大インスタンス、IAM 付与など）
- 初回実行例:
  ```bash
  cd terraform/environments/dev
  terraform init
  terraform plan \
    -var="project_id=<your-gcp-project>" \
    -var="region=asia-northeast1" \
    -var="api_image=asia-northeast1-docker.pkg.dev/<project>/<repo>/backend:latest" \
    -var="frontend_image=asia-northeast1-docker.pkg.dev/<project>/<repo>/frontend:latest"
  ```
- Firestore や Secret Manager、モニタリングの IaC 化は拡張予定です。

## セキュリティと可観測性
- Secrets: GitHub Actions → Secret Manager → Cloud Run 環境変数で運用。コードへの直書きは禁止。
- ミドルウェア: CSRF / JWT / 管理者ガード / HTTPS リダイレクト / Security Headers / Rate Limiter / Request ID / 構造化ログ。
- 可観測性: Prometheus メトリクス (`/metrics`)、構造化ログ（`internal/logging`）、Google Calendar / Gmail API の障害を Circuit Breaker で緩和。
- リトライ / バックオフ: 予約処理で指数バックオフ + Circuit Breaker。カレンダー API 障害時も即座に失敗せず再試行。
- 監視連携: Cloud Monitoring / Logging と組み合わせてダッシュボードやアラートを設定することを推奨。

## ドキュメント
- アーキテクチャ全体の意図や設計判断: `docs/architecture-design.md`
- API 仕様（今後）: OpenAPI 化を計画
- 運用 Runbook や Terraform 拡張は `docs/` 以下に追記予定

## トラブルシューティング
- **OAuth コールバックが失敗する**: `GOOGLE_OAUTH_CLIENT_ID/SECRET` と Redirect URI を再確認。許可ドメイン / メールリストも要チェック。
- **予約作成時に 500**: `GOOGLE_GMAIL_TOKEN`（または Secret Manager で指定したトークン）が Cloud Build / Cloud Run に伝搬しているか確認。
- **CSRF 403**: `GET /api/security/csrf` で取得した Cookie を変更せずに送信しているかを確認。ダブルサブミット方式のため Cookie とヘッダが一致している必要があります。
- **React Router の `startTransition` 警告**: ランタイムには影響なし。React Router v7 移行時に `future` フラグを有効化予定。
- **Docker Compose で DB スキーマを更新したい**: `docker compose down -v` でボリュームを破棄して再起動。

## 今後の改善アイデア
1. Go / React 双方で統合テスト・E2E テスト（Playwright）を CI に組み込み、カバレッジ 80% 以上を目指す。
2. Terraform に Firestore / Secret Manager / Monitoring リソースを追加し、環境差分を完全 IaC 化する。
3. Cloud Build 後に自動でライトウェイトな Smoke テスト（`make smoke-backend`）や Lighthouse CI を走らせる。
4. Cloud Monitoring にダッシュボード・アラートポリシーを定義し、インシデント対応の基盤を整える。
5. 管理 SPA の UX 改善（並列編集、ドラフト機能）とアクセス制御細分化（ロールベース）を検討する。
