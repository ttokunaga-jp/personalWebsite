# Runbook Directory

運用時の一次対応者が参照する Runbook を格納するディレクトリです。新しいアラートやサービスが追加された場合は以下のガイドラインに従って Runbook を作成してください。

## テンプレート
- `api-incident-template.md`：Cloud Run API で 5xx が増加した場合の対応手順テンプレート。

## 運用ルール
1. すべての Runbook は Pull Request 経由で更新し、SRE と開発リードのレビューを必須とします。
2. ファイル名は `service-name--incident-type.md` とし、ハイフンは単語区切り、英小文字で記載します。
3. 本番向け Runbook は日本語／英語併記を推奨します。少なくとも概要・影響範囲・一次対応・エスカレーションの 4 セクションを含めてください。
4. アラートポリシーや GitHub Actions からリンクする場合は `https://github.com/<org>/<repo>/blob/main/docs/runbooks/...` のフルパスを使用します。

