package inmemory

import (
	"time"

	"github.com/takumi/personal-website/internal/model"
)

var (
	defaultAdminProfile = &model.AdminProfile{
		Name:        model.NewLocalizedText("高見 拓実", "Takumi Takami"),
		Title:       model.NewLocalizedText("ソフトウェアエンジニア / リサーチャー", "Software Engineer / Researcher"),
		Affiliation: model.NewLocalizedText("架空大学", "Example University"),
		Lab:         model.NewLocalizedText("ヒューマンコンピュータインタラクション研究室", "Human-Computer Interaction Lab"),
		Summary: model.NewLocalizedText(
			"人に寄り添う体験と堅牢なインフラを両立させるプロダクト開発に取り組んでいます。",
			"Building delightful experiences backed by resilient infrastructure.",
		),
		Skills: []model.LocalizedText{
			model.NewLocalizedText("Go", "Go"),
			model.NewLocalizedText("React", "React"),
			model.NewLocalizedText("GCP", "GCP"),
			model.NewLocalizedText("機械学習", "Machine Learning"),
		},
		FocusAreas: []model.LocalizedText{
			model.NewLocalizedText("AI支援開発", "AI-assisted development"),
			model.NewLocalizedText("分散システム", "Distributed systems"),
			model.NewLocalizedText("開発プロセス改善", "Development workflow improvement"),
		},
		UpdatedAt: func() *time.Time {
			t := time.Date(2024, 5, 1, 9, 0, 0, 0, time.UTC)
			return &t
		}(),
	}

	defaultProfile = adminProfileToPublic(defaultAdminProfile)

	defaultAdminProjects = []model.AdminProject{
		{
			ID: 1,
			Title: model.NewLocalizedText(
				"AI支援ポートフォリオ",
				"AI Assisted Portfolio",
			),
			Description: model.NewLocalizedText(
				"AIを活用してコンテンツ編集を支援する実験的なプラットフォームです。",
				"An experimental platform for AI assisted content authoring.",
			),
			TechStack: []string{"Go", "React", "GCP"},
			LinkURL:   "https://example.dev/projects/ai-portfolio",
			Year:      2023,
			Published: true,
			CreatedAt: time.Date(2023, 6, 1, 9, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2023, 7, 15, 9, 0, 0, 0, time.UTC),
		},
		{
			ID: 2,
			Title: model.NewLocalizedText(
				"分散トレーシング可視化ツール",
				"Distributed Tracing Visualiser",
			),
			Description: model.NewLocalizedText(
				"マイクロサービスのトレーシングデータを可視化し、ボトルネックを検知するダッシュボードです。",
				"A dashboard that visualises tracing spans to detect bottlenecks across microservices.",
			),
			TechStack: []string{"Go", "TypeScript", "OpenTelemetry"},
			LinkURL:   "",
			Year:      2024,
			Published: false,
			CreatedAt: time.Date(2024, 1, 5, 9, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 2, 12, 9, 0, 0, 0, time.UTC),
		},
	}

	defaultAdminResearch = []model.AdminResearch{
		{
			ID:    1,
			Title: model.NewLocalizedText("リモートチームの適応的スケジューリング", "Adaptive Scheduling for Remote Teams"),
			Summary: model.NewLocalizedText(
				"遠隔チームの会議体験を最適化するための動的スケジューリング手法。",
				"Dynamic scheduling strategies to optimise remote meeting experiences.",
			),
			ContentMD: model.NewLocalizedText(
				"### 概要\n\nカレンダーのコンテキストを考慮したスケジューリングモデルを提案。",
				"### Summary\n\nIntroduces a scheduling model that incorporates calendar context signals.",
			),
			Year:      2022,
			Published: true,
			CreatedAt: time.Date(2022, 3, 10, 9, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2022, 5, 18, 9, 0, 0, 0, time.UTC),
		},
		{
			ID:    2,
			Title: model.NewLocalizedText("生成AIによる会議議事録支援", "Generative AI Meeting Minutes"),
			Summary: model.NewLocalizedText(
				"生成AIを活用した議事録作成アシスタントの評価実験。",
				"Evaluation of AI-assisted minute taking for distributed teams.",
			),
			ContentMD: model.NewLocalizedText(
				"### メソッド\n\nLLM を用いた要約とタスク抽出を比較しました。",
				"### Method\n\nCompared LLM-based summarisation approaches for action item extraction.",
			),
			Year:      2024,
			Published: false,
			CreatedAt: time.Date(2024, 4, 21, 9, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 4, 25, 9, 0, 0, 0, time.UTC),
		},
	}

	defaultBlogPosts = []model.BlogPost{
		{
			ID:        1,
			Title:     model.NewLocalizedText("AI駆動開発の始め方", "Getting Started with AI-driven Development"),
			Summary:   model.NewLocalizedText("AIによる開発支援のための基本的なフレームワークを紹介します。", "Introducing foundational practices for AI-assisted development."),
			ContentMD: model.NewLocalizedText("## はじめに\n\nAI を活用した開発体験の設計指針についてまとめました。", "## Introduction\n\nA quick primer on designing workflows with AI copilots."),
			Tags:      []string{"ai", "productivity"},
			Published: true,
			PublishedAt: func() *time.Time {
				t := time.Date(2023, 8, 1, 9, 0, 0, 0, time.UTC)
				return &t
			}(),
			CreatedAt: time.Date(2023, 7, 25, 9, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2023, 7, 30, 9, 0, 0, 0, time.UTC),
		},
		{
			ID:        2,
			Title:     model.NewLocalizedText("クリーンアーキテクチャ再考", "Revisiting Clean Architecture"),
			Summary:   model.NewLocalizedText("フロントとバックの協調を意識した設計パターンを考察します。", "Discussing patterns that harmonise frontend and backend concerns."),
			ContentMD: model.NewLocalizedText("## モジュール設計\n\n境界づけられたコンテキストを定義する重要性について。", "## Modular design\n\nOn the importance of defining bounded contexts clearly."),
			Tags:      []string{"architecture"},
			Published: false,
			CreatedAt: time.Date(2024, 2, 11, 9, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 2, 20, 9, 0, 0, 0, time.UTC),
		},
	}

	defaultContactMessages = []model.ContactMessage{
		{
			ID:        "contact-1",
			Name:      "Akari Yamada",
			Email:     "akari@example.com",
			Topic:     "プロジェクト相談",
			Message:   "最新のリサーチ成果について詳しく伺いたいです。",
			Status:    model.ContactStatusPending,
			AdminNote: "",
			CreatedAt: time.Date(2024, 5, 10, 10, 29, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 5, 10, 10, 29, 0, 0, time.UTC),
		},
		{
			ID:        "contact-2",
			Name:      "Lucas Chen",
			Email:     "lucas@example.com",
			Topic:     "AIプロダクト連携",
			Message:   "AI駆動ポートフォリオに関してAPI連携を検討しています。",
			Status:    model.ContactStatusInReview,
			AdminNote: "メールで詳細ヒアリング中",
			CreatedAt: time.Date(2024, 5, 12, 22, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 5, 13, 9, 0, 0, 0, time.UTC),
		},
	}

	defaultMeetings = []model.Meeting{
		{
			ID:              1,
			Name:            "Akari Yamada",
			Email:           "akari@example.com",
			Datetime:        time.Date(2024, 5, 10, 10, 0, 0, 0, time.UTC),
			DurationMinutes: 30,
			MeetURL:         "https://meet.example.com/session/1",
			Status:          model.MeetingStatusPending,
			Notes:           "",
			CreatedAt:       time.Date(2024, 5, 1, 9, 0, 0, 0, time.UTC),
			UpdatedAt:       time.Date(2024, 5, 1, 9, 0, 0, 0, time.UTC),
		},
		{
			ID:              2,
			Name:            "Lucas Chen",
			Email:           "lucas@example.com",
			Datetime:        time.Date(2024, 5, 12, 13, 0, 0, 0, time.UTC),
			DurationMinutes: 45,
			MeetURL:         "https://meet.example.com/session/2",
			Status:          model.MeetingStatusConfirmed,
			Notes:           "Confirmed via email",
			CreatedAt:       time.Date(2024, 5, 2, 9, 0, 0, 0, time.UTC),
			UpdatedAt:       time.Date(2024, 5, 3, 9, 0, 0, 0, time.UTC),
		},
	}

	defaultBlacklist = []model.BlacklistEntry{
		{
			ID:        1,
			Email:     "spam@example.com",
			Reason:    "Repeated spam submissions",
			CreatedAt: time.Date(2023, 12, 25, 9, 0, 0, 0, time.UTC),
		},
	}
)
