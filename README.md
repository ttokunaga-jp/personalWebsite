# Personal Website Platform

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

## Docker Compose メモ
- バックエンド API: `http://localhost:8080`
- フロントエンド (nginx): `http://localhost:3000` で `/api` をバックエンドにプロキシ
- MySQL: `localhost:3306`（認証情報は `docker-compose.yml` を参照）

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

## 次のステップ
- クリーンアーキテクチャ準拠のユースケース・リポジトリ実装を追加
- Terraform モジュールに Cloud SQL / Secret Manager / IAM 設定を拡張
- CI (GitHub Actions / Cloud Build) を整備し make タスク・Terraform plan を自動実行
