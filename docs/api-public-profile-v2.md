# Public API v1 Profile (v2 Schema)

Endpoint: `GET /api/v1/public/profile`

```jsonc
{
  "data": {
    "id": 1,
    "displayName": "Takumi",
    "headline": {
      "ja": "ソフトウェアエンジニア",
      "en": "Software Engineer"
    },
    "summary": {
      "ja": "AI 駆動開発に取り組んでいます",
      "en": "Building AI-assisted workflows"
    },
    "avatarUrl": "https://example.dev/avatar.png",
    "location": {
      "ja": "東京",
      "en": "Tokyo"
    },
    "theme": {
      "mode": "light",
      "accentColor": "#3b82f6"
    },
    "lab": {
      "name": {"ja": "HCI 研究室", "en": "HCI Lab"},
      "advisor": {"ja": "指導教員", "en": "Advisor"},
      "room": {"ja": "4F", "en": "4F"},
      "url": "https://example.dev/lab"
    },
    "affiliations": [
      {
        "id": 1,
        "profileId": 1,
        "kind": "affiliation",
        "name": "Example University",
        "url": "https://example.dev",
        "description": {"ja": "研究員", "en": "Researcher"},
        "startedAt": "2021-04-01T00:00:00Z",
        "sortOrder": 1
      }
    ],
    "communities": [],
    "workHistory": [],
    "techSections": [
      {
        "id": 1,
        "profileId": 1,
        "title": {"ja": "スキル", "en": "Skills"},
        "layout": "chips",
        "breakpoint": "lg",
        "sortOrder": 1,
        "members": [
          {
            "membershipId": 1,
            "entityType": "profile_section",
            "entityId": 1,
            "tech": {
              "id": 1,
              "slug": "go",
              "displayName": "Go",
              "category": "backend",
              "level": "advanced",
              "icon": "",
              "sortOrder": 1,
              "active": true,
              "createdAt": "2023-01-01T00:00:00Z",
              "updatedAt": "2024-01-01T00:00:00Z"
            },
            "context": "primary",
            "note": "",
            "sortOrder": 1
          }
        ]
      }
    ],
    "socialLinks": [
      {
        "id": 1,
        "profileId": 1,
        "provider": "github",
        "label": {"ja": "GitHub", "en": "GitHub"},
        "url": "https://github.com/example",
        "isFooter": true,
        "sortOrder": 1
      }
    ],
    "updatedAt": "2024-05-01T00:00:00Z"
  }
}
```

## Notes
- 全ての datetime は ISO8601 (UTC) で返却。
- `techSections[].members[].tech` は `tech_catalog` 由来のカタログ情報をインライン展開。
- 旧 API が返していた `skills` / `focusAreas` は、`techSections` や `workHistory` を参照してフロントで再構成する。
