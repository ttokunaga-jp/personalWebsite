# 管理 API CRUD 仕様（v2 コンテンツモデル）

> 更新日時: 2024-03-15  
> 基準スキーマ: `deploy/mysql/schema.sql` / `docs/content-model-refactor.md`

## 共通方針
- エンドポイントは `/api/admin/*` 配下。認証は HttpOnly/Admin セッション Cookie（`ps_admin_session`）必須。
- リクエスト/レスポンスは JSON。レスポンスは `{ "data": <payload> }` フォーマットで統一し、エラーは `errs.AppError` で定義された構造を返す。
- 時刻はすべて ISO8601（UTC）。多言語フィールドは `LocalizedText`（`{ "ja"?: string, "en"?: string }`）に統一。
- バリデーション違反は `400 Bad Request`（`CodeInvalidInput`）、競合は `409 Conflict`、存在しないリソースは `404 Not Found` を返す。
- 書き込み系 API は `If-Match` ヘッダで `updatedAt` の ETag を受け取り、競合を検出できるようにする（後方互換のため省略可だが v2 実装で対応）。

---

## 1. プロフィール管理
エンドポイント: `/api/admin/profile`

### DTO
```jsonc
type ProfileDocument = {
  id: number
  displayName: string
  headline: LocalizedText
  summary: LocalizedText
  avatarUrl?: string
  location: LocalizedText
  theme: {
    mode: "light" | "dark" | "system"
    accentColor?: string
  }
  lab: {
    name: LocalizedText
    advisor?: LocalizedText
    room?: LocalizedText
    url?: string
  }
  affiliations: ProfileAffiliation[]
  communities: ProfileAffiliation[]
  workHistory: ProfileWorkExperience[]
  techSections: ProfileTechSection[]
  socialLinks: ProfileSocialLink[]
  updatedAt: string // UTC ISO8601
}
```

バリデーション要点:
- `displayName` は必須（1〜128文字）。
- `theme.mode` は列挙値。`accentColor` は `^#?[0-9a-fA-F]{6}$`。
- `affiliations/communities` の `startedAt` は現在時刻以前。`kind` はテーブルと整合。
- `socialLinks` は GitHub/Zenn/LinkedIn を `isFooter=true` で必須。URL は HTTPS。

### API
| Method | Path                    | 説明                            |
| ------ | ----------------------- | ------------------------------- |
| GET    | `/api/admin/profile`    | 現在のプロフィール文書を取得。 |
| PUT    | `/api/admin/profile`    | 全体更新（リクエストは ProfileDocument）。|

PUT 成功時は更新後文書と新しい `updatedAt` を返却。`If-Match` が不一致の場合は `409`.

---

## 2. プロジェクト管理
ベース URL: `/api/admin/projects`

### DTO
```jsonc
type ProjectDocument = {
  id: number
  slug: string
  title: LocalizedText
  summary: LocalizedText
  description?: LocalizedText
  coverImageUrl?: string
  primaryLink?: string
  links: ProjectLink[]
  period: { start?: string; end?: string } // YYYY-MM-DD
  tech: TechMembership[]
  highlight: boolean
  published: boolean
  sortOrder: number
  createdAt: string
  updatedAt: string
}
```

バリデーション:
- `slug` は英小文字/数字/ハイフン (`^[a-z0-9-]{3,64}$`)、一意。
- `period.start <= period.end`。終了日は `null` で進行中。
- `links[].type` は `repo|demo|article|slides|other`。URL は HTTPS。
- `tech` は `tech_catalog` に存在する ID のみ許可。

### API
| Method | Path                         | 説明 |
| ------ | ---------------------------- | ---- |
| GET    | `/api/admin/projects`        | プロジェクト一覧（ドラフト含む）。|
| POST   | `/api/admin/projects`        | 新規作成（body は ProjectDocument から `id/createdAt/updatedAt` を除いたもの）。|
| GET    | `/api/admin/projects/:id`    | 個別取得。|
| PUT    | `/api/admin/projects/:id`    | 上書き更新。|
| DELETE | `/api/admin/projects/:id`    | 論理削除（v2は `deleted_at` 列導入予定。暫定では物理削除）。|

---

## 3. 研究・ブログ管理
ベース URL: `/api/admin/research-blog`

### DTO
```jsonc
type ResearchEntry = {
  id: number
  slug: string
  kind: "research" | "blog"
  title: LocalizedText
  overview?: LocalizedText
  outcome?: LocalizedText
  outlook?: LocalizedText
  externalUrl: string
  publishedAt: string
  updatedAt: string
  highlightImageUrl?: string
  imageAlt?: LocalizedText
  isDraft: boolean
  tags: string[]
  links: ResearchLink[]
  assets: ResearchAsset[]
  tech: TechMembership[]
}
```

バリデーション:
- `slug` はプロジェクトと同様の命名制約かつ一意。
- `externalUrl` は HTTPS。`publishedAt` は過去〜現在。
- `links[].type` は `paper|slides|video|code|external`。
- `assets[].url` は HTTPS。`tags` は 32 文字以内。

API エンドポイント:
| Method | Path                                   | 説明 |
| ------ | -------------------------------------- | ---- |
| GET    | `/api/admin/research-blog`             | 全件取得。クエリ `kind`, `isDraft` フィルタ対応。|
| POST   | `/api/admin/research-blog`             | 新規作成。|
| GET    | `/api/admin/research-blog/:id`         | 個別取得。|
| PUT    | `/api/admin/research-blog/:id`         | 更新。|
| DELETE | `/api/admin/research-blog/:id`         | 削除。|

---

## 4. ホーム設定管理
ベース URL: `/api/admin/home-config`

### DTO
```jsonc
type HomePageConfig = {
  id: number
  profileId: number
  heroSubtitle?: LocalizedText
  quickLinks: HomeQuickLink[]
  chipSources: HomeChipSource[]
  updatedAt: string
}
```

バリデーション:
- `profileId` は `profiles.id` の参照。存在チェック必須。
- `quickLinks[].section` は `profile|research_blog|projects|contact`。
- `chipSources[].source` は `affiliation|community|tech`。`limit` は `1-12`。

API:
| Method | Path                    | 説明 |
| ------ | ----------------------- | ---- |
| GET    | `/api/admin/home-config`| 設定取得（単一レコード）。|
| PUT    | `/api/admin/home-config`| 更新。|

---

## 5. お問い合わせ設定 / 予約

### 5.1 Contact Form 設定
ベース URL: `/api/admin/contact-settings`

DTO:
```jsonc
type ContactFormSettings = {
  id: number
  heroTitle?: LocalizedText
  heroDescription?: LocalizedText
  topics: ContactTopic[]
  consentText: LocalizedText
  minimumLeadHours: number
  recaptchaPublicKey?: string
  supportEmail: string
  calendarTimezone: string
  googleCalendarId?: string
  bookingWindowDays: number
  createdAt: string
  updatedAt: string
}
```

バリデーション:
- `supportEmail` はメール形式。`calendarTimezone` は IANA タイムゾーン。
- `topics[].id` は 32 文字以内で一意。`topics[].label` いずれかの言語が必須。
- `minimumLeadHours` ≥ 0、`bookingWindowDays` ≥ 1。

API:
| Method | Path                            | 説明 |
| ------ | ------------------------------- | ---- |
| GET    | `/api/admin/contact-settings`   | 設定取得。|
| PUT    | `/api/admin/contact-settings`   | 更新。|

### 5.2 予約管理
ベース URL: `/api/admin/reservations`

DTO:
```jsonc
type MeetingReservation = {
  id: number
  name: string
  email: string
  topic?: string
  message?: string
  startAt: string
  endAt: string
  durationMinutes: number
  googleEventId?: string
  googleCalendarStatus: "pending" | "confirmed" | "declined" | "cancelled"
  status: "pending" | "confirmed" | "cancelled"
  confirmationSentAt?: string
  lastNotificationSentAt?: string
  lookupHash: string
  cancellationReason?: string
  createdAt: string
  updatedAt: string
  notifications: MeetingNotification[]
}
```

バリデーション:
- `startAt < endAt`。`durationMinutes` は 15 の倍数かつ 15〜180。
- `lookupHash` は `SHA256(name:email)` を保持。更新時に再計算。
- `notifications[].status` は `pending|sent|failed`。

API:
| Method | Path                                | 説明 |
| ------ | ----------------------------------- | ---- |
| GET    | `/api/admin/reservations`           | 予約一覧（クエリ `status`, `email`, `date` フィルタ）。|
| GET    | `/api/admin/reservations/:id`       | 個別取得。|
| PUT    | `/api/admin/reservations/:id`       | 状態更新・キャンセル処理。|
| POST   | `/api/admin/reservations/:id/retry` | 通知再送。|

---

## 6. 技術カタログ管理（補足）
ベース URL: `/api/admin/tech-catalog`

DTO:
```jsonc
type TechCatalogEntry = {
  id: number
  slug: string
  displayName: string
  category?: string
  level: "beginner" | "intermediate" | "advanced"
  icon?: string
  sortOrder: number
  active: boolean
  createdAt: string
  updatedAt: string
}
```

バリデーション:
- `slug` は英数字/ハイフン、既存利用状況を考慮して重複不可。
- `displayName` 1〜128 文字。`icon` は emoji も許可。

API:
| Method | Path                           | 説明 |
| ------ | ------------------------------ | ---- |
| GET    | `/api/admin/tech-catalog`      | 一覧。クエリ `active`、`q` (検索) 対応。|
| POST   | `/api/admin/tech-catalog`      | 追加。|
| PUT    | `/api/admin/tech-catalog/:id`  | 更新。|
| DELETE | `/api/admin/tech-catalog/:id`  | 非アクティブ化（`active=false` に更新）。|

---

## 7. 共通 DTO
```jsonc
type LocalizedText = {
  ja?: string
  en?: string
}

type TechMembership = {
  membershipId?: number
  entityType: "profile_section" | "project" | "research_blog"
  entityId?: number
  techId: number
  context: "primary" | "supporting"
  note?: string
  sortOrder: number
}
```

---

## 8. ドキュメント管理
- 本仕様は管理 API 実装のソース (handler/service) と同リポジトリ内で同期する。
- JSON スキーマの自動生成（`openapi.yaml` 更新）はフェーズ 2 にて対応予定。
- 管理 SPA では上記 DTO を TypeScript 型として取り込み、FieldArray / バリデーションは Yup/React Hook Form 等で実装する。
