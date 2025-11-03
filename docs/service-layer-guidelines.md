# サービス層ガイドライン (v2 対応)

## 目的
- 新しいコンテンツモデル (ProfileDocument / ProjectDocument など) を利用するサービス層で、エラーハンドリングとトランザクション境界の扱いを統一する。
- 旧 API との互換性を維持しつつ段階的にレスポンス構造を移行するためのルールを明文化する。

## エラーハンドリング
- リポジトリ層で返却されるエラーは `service/support.MapRepositoryError` を通じて `errs.AppError` に変換する。
  - `ErrNotFound` → 404 / `CodeNotFound`
  - `ErrInvalidInput` → 400 / `CodeInvalidInput`
  - `ErrDuplicate` → 409 / `CodeConflict`
  - それ以外 → 500 / `CodeInternal`
- サービス層からは必ず `errs.AppError` を返し、ハンドラで `respondError` を呼べば既存のレスポンス構造と整合する。
- 旧 API 互換のためにアダプタでデータを縮約する場合も、失敗時は上記ルールでエラーを生成する。

## トランザクション境界
- 読み取り専用のサービスは原則としてトランザクション不要。
- 書き込み処理を伴うサービスでは、`sqlx.DB` を直接利用するのではなく、`service/support` に用意するトランザクションヘルパ (今後追加予定) を介して `BEGIN`, `COMMIT`, `ROLLBACK` を集中管理する。
- トランザクション内で送出されたエラーはそのまま `MapRepositoryError` で変換し、呼び出し側は再ラップしない。

## 互換アダプタ
- v2 モデルを旧 API の JSON に変換する場合は、`service` パッケージにアダプタを置く。
  - 例: `profileDocumentToLegacy` は `ProfileDocument` → `model.Profile` を安全に縮約する。
  - 旧構造に存在しないデータは無理に詰め込まず、必要に応じてクライアント側で新エンドポイントを利用するよう誘導する。
- アダプタはドメインロジックを持たず、単純なマッピングのみに留める。

## 今後の移行方針
1. 公開 API から順に新サービスへ差し替え (`Profile` → `Projects` → `Research` → `Contact` → `Home`)。
2. 管理 API も同じドキュメントモデルを返すように再構築し、フロントエンドに実装変更を段階的に行う。
3. 全エンドポイントが新スキーマを返し終えたら、旧 `model.Profile` など互換用の構造体とリポジトリを削除する。
