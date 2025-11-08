# API インシデント対応テンプレート / API Incident Response Template

## 1. 概要 / Summary
- **発生日時 / Start time**: <!-- YYYY-MM-DD HH:MM JST -->
- **検知方法 / Detection**: PagerDuty, Slack Alert, Manual, etc.
- **影響範囲 / Impact**: 例) 全ユーザーの 5xx が 10 分継続。予約フォーム送信不可。
- **対応担当 / On-call owner**: <!-- Name -->

## 2. 初動対応 / Immediate Actions
1. PagerDuty / Slack でアラートを確認し、担当者にアサインする。
2. Cloud Monitoring ダッシュボード `Personal Website Operations Overview` を開き、トラフィック・レイテンシ・エラーレートを確認する。
3. `make smoke-backend` を実行し、主要 API エンドポイントの応答を確認する。
4. 直近デプロイがある場合は Cloud Run リビジョンを確認し、健全なリビジョンへトラフィックを切り戻す。

## 3. 調査 / Investigation
- **ログ**: Cloud Logging で `trace_id` / `request_id` をキーにエラーを抽出。フィルタ例:  
  ```
  resource.type="cloud_run_revision"
  severity>=ERROR
  resource.labels.service_name="personal-website-api"
  ```
- **依存サービス**: Cloud SQL 接続数、Google Calendar API 呼び出し結果、外部 API ステータスを確認。
- **最近の変更**: GitHub のデプロイ履歴、Terraform 変更、DB マイグレーションを確認。

## 4. 回復手順 / Mitigation
- Cloud Run ロールバック、環境変数の再適用、必要であれば一時的にエッジキャッシュを有効化。
- Cloud SQL / Firestore のパフォーマンス劣化が原因の場合は、接続プール調整またはリトライ間隔の変更。
- 外部 API 障害時は機能制限モードに切り替え、管理画面にバナーを掲示。

## 5. エスカレーション / Escalation
- エラーが 15 分以内に収束しない場合、SRE リードに電話連絡。
- データ損失が疑われる場合はインシデントレスポンスチームへエスカレーション。

## 6. 事後対応 / Postmortem
- RCA（Root Cause Analysis）を 48 時間以内に作成し、`docs/ops/incidents/` に保存。
- 再発防止策・テスト・アラート閾値の更新をチームでレビュー。
- ユーザー影響が大きい場合はステータスページや SNS でアナウンス。

