# 目的
- 本番運用を見据えて、Go API・React SPA・MySQL・GCP 環境の可観測性、アラート、バックアップ、および障害対応プロセスを定義し、運用チームが一貫した基準で監視・保守できる状態を整える。
- 障害復旧・改善サイクルを高速化するための Runbook、スケジュール、ロードマップを明文化し、プロジェクト全体の信頼性と継続的改善を担保する。

> **依存タスク**：インフラ IaC プロンプトで定義される Terraform モジュール／GCP プロジェクト構成が前提。

# 実装内容
## ログ収集・可観測性設計
- **構造化ログ標準化**：Go API は `slog`、React SPA は `@sentry/browser` + `console.error` Hook で JSON 形式に統一。全リクエストに `trace_id` / `span_id` / `request_id` を付与し、Cloud Trace に紐付ける。
- **収集経路**：
  - Cloud Run 標準ログ → Cloud Logging → Log Router → BigQuery（長期保管）と Cloud Storage（90 日ローテーション）へエクスポート。
  - フロントエンドのブラウザログは Cloud Logging の HTTP 受信エンドポイント（`/v1/log`）経由で API が受け取り、バックエンドログと同一スキーマで記録。
- **メトリクス**：
  - Cloud Monitoring（Managed） + OpenTelemetry Collector（Optional）で以下を収集：`request_count`, `error_count`, `latency_p95`, `cpu_utilization`, `memory_utilization`, `calendar_api_failures`, `sql_connections`, `cache_hit_ratio`, `queue_backlog`.
  - ビジネス KPI として `reservation_success_rate`, `blacklist_rejection_count`, `admin_login_count` をカスタムメトリクスに定義。
- **ダッシュボード**：Cloud Monitoring で運用/経営向け 2 種類のダッシュボードを作成。可視化ウィジェットはトラフィック傾向、レイテンシ、エラーレート、予約 KPI、GCP コストを含む。

## アラート／通知設計
- **通知チャネル**：PagerDuty（重大障害）、Slack #alerts（注意）、メールサマリー（日次）。発火時は Cloud Function でチケットシステム（Jira）の自動起票。
- **閾値とルール**：
  - `error_count / request_count > 5%` を 5 分継続 → Major（PagerDuty）。
  - `latency_p95 > 800ms` を 10 分継続 → Minor（Slack）。
  - `calendar_api_failures > 10/15m` または Circuit Breaker Open → Major。
  - `sql_backup_status != SUCCESS` → Major。
  - `admin_login_count = 0` for 24h（利用停止検知） → Info（メール）。
- **ランブック連携**：各アラートに Runbook URL を付与し、一次対応者が原因切り分けと復旧手順を即参照できるようにする。

## バックアップ／リストア手順書
- **Cloud SQL（MySQL）**：
  - 自動バックアップ：日次（02:00 JST）、7 世代保持。Point-in-Time Recovery（PITR）を有効化。
  - 手動バックアップ：重大リリース前／スキーマ変更前に `gcloud sql backups create` で取得。
  - リストア手順：
    1. `gcloud sql backups list` で対象バックアップ ID を特定。
    2. ステージング新インスタンスに `gcloud sql backups restore` 実行。
    3. 付属スクリプト `scripts/db-checksum.sh` でテーブル整合性を確認。
    4. 本番切替が必要な場合、Cloud Run の接続設定を新インスタンスに更新後、ヘルスチェック合格を確認してトラフィックを戻す。
- **Cloud Storage（静的アセット）**：
  - バージョニング有効化 + 週次で `gsutil rsync` による別バケット（`gs://<project>-backup-storage`）へコピー。
  - リストア：ロールバック対象オブジェクトのバージョン ID を指定し `gsutil cp` で復元。
- **Terraform State**：Cloud Storage バケット + バージョニング + Lock（DynamoDB 互換）で管理。週次でローカル暗号化アーカイブを Cloud Storage Coldline に保管。

## 障害対応 Runbook（抜粋）
- **API レスポンス 5xx 増加**：
  1. Cloud Monitoring ダッシュボードでレイテンシ・エラーメッセージを確認。
  2. Cloud Logging で `trace_id` を辿り、Go Panic / SQL 接続エラーを特定。
  3. 必要であれば直近デプロイを Cloud Run リビジョンロールバック。
  4. カナリアテスト用 `make smoke-test` を実行し、機能回復を検証。
  5. 事後対応：根本原因分析（RCA）を Notion に記録し、再発防止策を作成。
- **Google Calendar API 障害**：
  1. Circuit Breaker メトリクスを確認して Open 状態ならバックオフ中であることを周知。
  2. Retry キュー（Pub/Sub）内のイベント数を確認し、溢れそうならバッチ処理を実行。
  3. SLA 超過が見込まれる場合、管理画面にバナー表示 + 管理者へ手動登録を案内。
  4. 復旧後、Pending 予約を `scripts/replay-calendar-events.sh` で再送。
- **Cloud SQL 接続枯渇**：
  1. MySQL Performance Schema でアクティブ接続数を調査。
  2. コネクションプール設定（Go `database/sql`）を一時的に引き下げ、長時間トランザクションを停止。
  3. バックアップ・リストア状況を確認し、必要に応じて只読レプリカへ切替。

## 運用スケジュールとタスク一覧
- **日次**：バックアップ結果確認、アラートレビュー、予約失敗ログ確認、セキュリティログ（Cloud Armor）チェック。
- **週次（火曜 10:00 JST）**：依存ライブラリアップデート確認、エラーレポート集計、SLO（稼働率・レイテンシ）レビュー、Terraform Plan のドライラン。
- **月次（第 1 水曜 15:00 JST）**：DR ドリル（PITR リストア演習）、コストレポート分析、アクセス権棚卸し、ロードマップ更新。
- **四半期**：総合監査（監視ルール、プレイブック、セキュリティポリシーの棚卸し）、性能テスト（k6 + Lighthouse）。

## 拡張・改善ロードマップ
1. **短期（1-2 か月）**：Synthetic Monitoring（Cloud Monitoring Uptime Check）導入、アドミン画面にリアルタイムメトリクス表示、ログクエリのプリセット共有。
2. **中期（3-6 か月）**：Service Level Objective (SLO) & Error Budget 運用開始、AIOps（Cloud Operations Recommender）活用、Cloud DLP を用いたログの個人情報マスキング。
3. **長期（6 か月以降）**：BigQuery + Looker Studio でカスタマー分析ダッシュボード、Blue/Green デプロイ自動化、マルチリージョン DR（セカンダリ Cloud Run & Cloud SQL リードレプリカ）。

# コード
```hcl
# terraform/modules/monitoring/main.tf
resource "google_monitoring_dashboard" "service_observability" {
  dashboard_json = templatefile("${path.module}/dashboards/service.json", {
    project_id = var.project_id
  })
}

resource "google_monitoring_alert_policy" "api_error_rate" {
  display_name = "API 5xx rate > 5%"
  combiner     = "OR"

  conditions {
    display_name = "High 5xx ratio"
    condition_threshold {
      filter = <<-EOT
        resource.type="cloud_run_revision"
        AND metric.type="run.googleapis.com/request_count"
        AND metric.label."response_code" =~ "^5"
      EOT
      aggregations {
        alignment_period     = "300s"
        per_series_aligner   = "ALIGN_RATE"
        cross_series_reducer = "REDUCE_SUM"
        group_by_fields      = ["resource.label.service_name"]
      }
      denominator_filter = <<-EOT
        resource.type="cloud_run_revision"
        AND metric.type="run.googleapis.com/request_count"
      EOT
      denominator_aggregations {
        alignment_period   = "300s"
        per_series_aligner = "ALIGN_RATE"
      }
      comparison = "COMPARISON_GT"
      threshold_value = 0.05
      duration = "300s"
      trigger { count = 1 }
    }
  }

  notification_channels = [
    google_monitoring_notification_channel.pagerduty.id,
    google_monitoring_notification_channel.slack.id
  ]
  documentation {
    content = "Runbook: https://docs.example.com/runbooks/api-5xx"
  }
}

resource "google_sql_backup_run" "manual_pre_release" {
  instance = google_sql_database_instance.primary.name
  depends_on = [time_rotating.cron_pre_release]
}
```

```yaml
# deploy/cloud-logging/log-router.yaml
metadata:
  name: personal-website-log-router
  description: Route structured logs to BigQuery and Storage
sink:
  destination:
    - bigquery.googleapis.com/projects/${PROJECT_ID}/datasets/app_logs
    - storage.googleapis.com/${PROJECT_ID}-logs-archive
  filter: |
    resource.type = ("cloud_run_revision" OR "cloudsql_database")
    AND severity >= DEFAULT
  includeChildren: true
```

# 設定
- **ディレクトリ構造**：
  ```
  docs/
    architecture-design.md
    operations-design.md
  deploy/
    mysql/
      schema.sql
    cloud-logging/
      log-router.yaml
    monitoring/
      service.json
    runbooks/            # TODO: Runbook テンプレートを追加予定
  backend/
    scripts/
      api_smoke.sh
  terraform/
    modules/
      monitoring/
        main.tf
        dashboards/
          service.json
    envs/
      prod/
        main.tf
        variables.auto.tfvars
  ```
- **環境変数／Secret 管理**：`OPS_PAGERDUTY_KEY`, `OPS_SLACK_WEBHOOK`, `OPS_JIRA_TOKEN`, `BACKUP_BUCKET`, `CALENDAR_RETRY_TOPIC` を Secret Manager で管理し、Cloud Run/Functions へ Workload Identity Federation 経由で注入。
- **アクセス権限**：運用担当ロール（`roles/monitoring.admin`, `roles/logging.admin`, `roles/cloudsql.editor`）を最小権限で付与。Runbook 編集は GitHub PR ベースでレビュー必須。
- **変更管理**：アラート閾値、Runbook 変更は `ops-change-request` テンプレートを用いた Pull Request で行い、SRE + 開発リードの承認を必須とする。

# テスト
- **障害通知テスト**：
  - Cloud Monitoring で `test_notification` を実行し、PagerDuty/Slack に到達することを確認。
  - `gcloud alpha monitoring policies enable --policy=<id> --notification-channel=<test>` を用いて条件を一時的にトリガーし、Runbook リンクの動作を確認。
- **バックアップ & リストア演習**：
  - ステージング環境で日次バックアップを取得後、`gcloud sql backups restore` で別インスタンスにリストアし、`deploy/mysql/schema.sql` に含まれる代表テーブルで `CHECKSUM TABLE` を実行して整合性を確認。
  - Cloud Storage のバージョン復元を `gsutil cp -a` で実施し、React ビルドアセットのハッシュ一致を `pnpm --filter @personal-website/public build` の出力と比較。API の基本機能は `backend/scripts/api_smoke.sh` で確認。
- **可観測性検証**：
  - k6 シナリオで 10 分間負荷を与え、`latency_p95` と `error_rate` のメトリクス更新をダッシュボードで確認。
  - 管理画面のメトリクス表示コンポーネントが Cloud Monitoring API から最新値を取得する E2E テスト（Playwright）を実行。

# 検証方法
- `make ops-validate`（新規追加予定）を実行し、Terraform のアラート定義が `terraform validate` をパスし、`yamllint deploy/cloud-logging/log-router.yaml` が成功することを確認。
- ステージング環境で意図的に API に障害フラグを設定（Feature Flag）し、PagerDuty・Slack 通知、Runbook フロー、ステータスページ更新までの一連をタイムスタンプ付きで記録。
- 月次 DR ドリルの結果（リストア所要時間、データ整合性チェックログ）を `docs/ops/dr-drill-YYYYMM.md` として保存し、レビュー済みであることを確認。
- 監視/バックアップ設定の変更は GitHub Actions の `ops-check` ワークフロー（`terraform plan` + `gcloud beta monitoring policies lint`）が成功し、コードレビューで承認されたことをもって検証完了とする。
