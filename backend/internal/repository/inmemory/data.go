package inmemory

import (
	"time"

	"github.com/takumi/personal-website/internal/model"
)

var (
	defaultAdminProfile = &model.AdminProfile{
		ID:          1,
		DisplayName: "Takumi Takami",
		Headline: model.NewLocalizedText(
			"ソフトウェアエンジニア / リサーチャー",
			"Software Engineer / Researcher",
		),
		Summary: model.NewLocalizedText(
			"人に寄り添う体験と堅牢なインフラを両立させるプロダクト開発に取り組んでいます。",
			"Building delightful experiences backed by resilient infrastructure.",
		),
		AvatarURL: "https://example.dev/avatar.png",
		Location:  model.NewLocalizedText("東京", "Tokyo"),
		Theme: model.ProfileTheme{
			Mode:        model.ProfileThemeModeLight,
			AccentColor: "#3b82f6",
		},
		Lab: model.ProfileLab{
			Name:    model.NewLocalizedText("ヒューマンコンピュータインタラクション研究室", "Human-Computer Interaction Lab"),
			Advisor: model.NewLocalizedText("指導教員", "Advisor"),
			Room:    model.NewLocalizedText("4F 研究室", "Lab 4F"),
			URL:     "https://example.dev/lab",
		},
		Affiliations: []model.ProfileAffiliation{
			{
				ID:          1,
				ProfileID:   1,
				Kind:        model.ProfileAffiliationKindAffiliation,
				Name:        "Example University",
				URL:         "https://example.dev",
				Description: model.NewLocalizedText("研究員", "Researcher"),
				StartedAt:   time.Date(2021, 4, 1, 0, 0, 0, 0, time.UTC),
				SortOrder:   1,
			},
			{
				ID:          2,
				ProfileID:   1,
				Kind:        model.ProfileAffiliationKindAffiliation,
				Name:        "Example Graduate School",
				URL:         "https://grad.example",
				Description: model.NewLocalizedText("修士課程", "Graduate program"),
				StartedAt:   time.Date(2019, 4, 1, 0, 0, 0, 0, time.UTC),
				SortOrder:   2,
			},
		},
		Communities: []model.ProfileAffiliation{
			{
				ID:          3,
				ProfileID:   1,
				Kind:        model.ProfileAffiliationKindCommunity,
				Name:        "Open Source Guild",
				URL:         "https://oss.example",
				Description: model.NewLocalizedText("OSS コミュニティ", "OSS community"),
				StartedAt:   time.Date(2020, 9, 1, 0, 0, 0, 0, time.UTC),
				SortOrder:   1,
			},
		},
		WorkHistory: []model.ProfileWorkExperience{
			{
				ID:           1,
				ProfileID:    1,
				Organization: model.NewLocalizedText("Example Corp", "Example Corp"),
				Role:         model.NewLocalizedText("フルスタックエンジニア", "Full-stack Engineer"),
				Summary:      model.NewLocalizedText("AI と Web を横断するプロダクト開発をリード。", "Led AI and web product initiatives."),
				StartedAt:    time.Date(2019, 4, 1, 0, 0, 0, 0, time.UTC),
				ExternalURL:  "https://example.dev/company",
				SortOrder:    1,
			},
			{
				ID:           2,
				ProfileID:    1,
				Organization: model.NewLocalizedText("Example Labs", "Example Labs"),
				Role:         model.NewLocalizedText("リサーチエンジニア", "Research Engineer"),
				Summary:      model.NewLocalizedText("ユーザー行動分析とプロトタイピングを担当。", "Focused on user research and prototyping."),
				StartedAt:    time.Date(2022, 6, 1, 0, 0, 0, 0, time.UTC),
				SortOrder:    2,
			},
		},
		TechSections: []model.ProfileTechSection{
			{
				ID:         1,
				ProfileID:  1,
				Title:      model.NewLocalizedText("スキルセット", "Skill Set"),
				Layout:     "chips",
				Breakpoint: "lg",
				SortOrder:  1,
				Members: []model.TechMembership{
					{
						MembershipID: 1,
						EntityType:   "profile_section",
						EntityID:     1,
						Tech: model.TechCatalogEntry{
							ID:          1,
							Slug:        "go",
							DisplayName: "Go",
							Category:    "backend",
							Level:       model.TechLevelAdvanced,
							SortOrder:   1,
							Active:      true,
							CreatedAt:   time.Now().Add(-365 * 24 * time.Hour),
							UpdatedAt:   time.Now(),
						},
						Context:   model.TechContextPrimary,
						SortOrder: 1,
					},
					{
						MembershipID: 2,
						EntityType:   "profile_section",
						EntityID:     1,
						Tech: model.TechCatalogEntry{
							ID:          2,
							Slug:        "react",
							DisplayName: "React",
							Category:    "frontend",
							Level:       model.TechLevelAdvanced,
							SortOrder:   2,
							Active:      true,
							CreatedAt:   time.Now().Add(-400 * 24 * time.Hour),
							UpdatedAt:   time.Now(),
						},
						Context:   model.TechContextSupporting,
						SortOrder: 2,
					},
				},
			},
		},
		SocialLinks: []model.ProfileSocialLink{
			{
				ID:        1,
				ProfileID: 1,
				Provider:  model.ProfileSocialProviderGitHub,
				Label:     model.NewLocalizedText("GitHub", "GitHub"),
				URL:       "https://github.com/example",
				IsFooter:  true,
				SortOrder: 1,
			},
			{
				ID:        2,
				ProfileID: 1,
				Provider:  model.ProfileSocialProviderZenn,
				Label:     model.NewLocalizedText("Zenn", "Zenn"),
				URL:       "https://zenn.dev/example",
				IsFooter:  true,
				SortOrder: 2,
			},
			{
				ID:        3,
				ProfileID: 1,
				Provider:  model.ProfileSocialProviderLinkedIn,
				Label:     model.NewLocalizedText("LinkedIn", "LinkedIn"),
				URL:       "https://linkedin.com/in/example",
				IsFooter:  true,
				SortOrder: 3,
			},
		},
		UpdatedAt: time.Date(2024, 5, 1, 9, 0, 0, 0, time.UTC),
	}

	defaultProfile = documentToLegacy(defaultAdminProfile)

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
			Tech: []model.TechMembership{
				{
					MembershipID: 1,
					EntityType:   "project",
					EntityID:     1,
					Tech: model.TechCatalogEntry{
						ID:          1,
						Slug:        "go",
						DisplayName: "Go",
						Level:       model.TechLevelAdvanced,
						SortOrder:   1,
						Active:      true,
					},
					Context:   model.TechContextPrimary,
					SortOrder: 1,
				},
				{
					MembershipID: 2,
					EntityType:   "project",
					EntityID:     1,
					Tech: model.TechCatalogEntry{
						ID:          2,
						Slug:        "react",
						DisplayName: "React",
						Level:       model.TechLevelAdvanced,
						SortOrder:   2,
						Active:      true,
					},
					Context:   model.TechContextSupporting,
					SortOrder: 2,
				},
				{
					MembershipID: 3,
					EntityType:   "project",
					EntityID:     1,
					Tech: model.TechCatalogEntry{
						ID:          3,
						Slug:        "gcp",
						DisplayName: "GCP",
						Level:       model.TechLevelIntermediate,
						SortOrder:   3,
						Active:      true,
					},
					Context:   model.TechContextSupporting,
					SortOrder: 3,
				},
			},
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
			Tech: []model.TechMembership{
				{
					MembershipID: 4,
					EntityType:   "project",
					EntityID:     2,
					Tech: model.TechCatalogEntry{
						ID:          1,
						Slug:        "go",
						DisplayName: "Go",
						Level:       model.TechLevelAdvanced,
						SortOrder:   1,
						Active:      true,
					},
					Context:   model.TechContextPrimary,
					SortOrder: 1,
				},
				{
					MembershipID: 5,
					EntityType:   "project",
					EntityID:     2,
					Tech: model.TechCatalogEntry{
						ID:          4,
						Slug:        "typescript",
						DisplayName: "TypeScript",
						Level:       model.TechLevelAdvanced,
						SortOrder:   2,
						Active:      true,
					},
					Context:   model.TechContextSupporting,
					SortOrder: 2,
				},
				{
					MembershipID: 6,
					EntityType:   "project",
					EntityID:     2,
					Tech: model.TechCatalogEntry{
						ID:          5,
						Slug:        "opentelemetry",
						DisplayName: "OpenTelemetry",
						Level:       model.TechLevelIntermediate,
						SortOrder:   3,
						Active:      true,
					},
					Context:   model.TechContextSupporting,
					SortOrder: 3,
				},
			},
			LinkURL:   "",
			Year:      2024,
			Published: false,
			CreatedAt: time.Date(2024, 1, 5, 9, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 2, 12, 9, 0, 0, 0, time.UTC),
		},
	}

	defaultAdminResearch = []model.AdminResearch{
		{
			ID:   1,
			Slug: "adaptive-scheduling",
			Kind: model.ResearchKindResearch,
			Title: model.NewLocalizedText(
				"リモートチームの適応的スケジューリング",
				"Adaptive Scheduling for Remote Teams",
			),
			Overview: model.NewLocalizedText(
				"遠隔チームの会議体験を最適化するための動的スケジューリング手法を提案します。",
				"Introduces a dynamic scheduling approach tailored for distributed teams.",
			),
			Outcome: model.NewLocalizedText(
				"パイロット導入で会議時間を平均 14% 削減し、満足度スコアを 1.3 向上させました。",
				"Pilot adoption reduced average meeting time by 14% and improved satisfaction by 1.3 points.",
			),
			Outlook: model.NewLocalizedText(
				"今後は Google Calendar のインサイト連携による最適化を検証予定です。",
				"Validation with Google Calendar insights will drive further optimisation.",
			),
			ExternalURL:       "https://example.com/research/adaptive-scheduling",
			HighlightImageURL: "https://cdn.example.com/images/adaptive-scheduling.jpg",
			ImageAlt: model.NewLocalizedText(
				"会議スケジュール最適化の可視化パネル",
				"Visual panel illustrating scheduling optimisation",
			),
			PublishedAt: time.Date(2022, 3, 15, 9, 0, 0, 0, time.UTC),
			IsDraft:     false,
			CreatedAt:   time.Date(2022, 3, 10, 9, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2022, 5, 18, 9, 0, 0, 0, time.UTC),
			Tags: []model.ResearchTag{
				{ID: 1, EntryID: 1, Value: "scheduling", SortOrder: 1},
				{ID: 2, EntryID: 1, Value: "remote-work", SortOrder: 2},
			},
			Links: []model.ResearchLink{
				{
					ID:        1,
					EntryID:   1,
					Type:      model.ResearchLinkTypePaper,
					Label:     model.NewLocalizedText("論文 PDF", "Paper PDF"),
					URL:       "https://example.com/research/adaptive-scheduling-paper.pdf",
					SortOrder: 1,
				},
				{
					ID:        2,
					EntryID:   1,
					Type:      model.ResearchLinkTypeSlides,
					Label:     model.NewLocalizedText("発表資料", "Conference Slides"),
					URL:       "https://example.com/research/adaptive-slides",
					SortOrder: 2,
				},
			},
			Assets: []model.ResearchAsset{
				{
					ID:        1,
					EntryID:   1,
					URL:       "https://cdn.example.com/assets/adaptive-dashboard.png",
					Caption:   model.NewLocalizedText("ダッシュボード UI", "Dashboard UI"),
					SortOrder: 1,
				},
			},
			Tech: []model.TechMembership{
				{
					MembershipID: 1,
					EntityType:   researchEntityType,
					EntityID:     1,
					Tech: model.TechCatalogEntry{
						ID:          1,
						Slug:        "go",
						DisplayName: "Go",
						Level:       model.TechLevelAdvanced,
						SortOrder:   1,
						Active:      true,
						CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					Context:   model.TechContextPrimary,
					SortOrder: 1,
				},
				{
					MembershipID: 2,
					EntityType:   researchEntityType,
					EntityID:     1,
					Tech: model.TechCatalogEntry{
						ID:          2,
						Slug:        "react",
						DisplayName: "React",
						Level:       model.TechLevelIntermediate,
						SortOrder:   2,
						Active:      true,
						CreatedAt:   time.Date(2019, 7, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt:   time.Date(2022, 11, 1, 0, 0, 0, 0, time.UTC),
					},
					Context:   model.TechContextSupporting,
					SortOrder: 2,
				},
			},
		},
		{
			ID:   2,
			Slug: "llm-meeting-minutes",
			Kind: model.ResearchKindResearch,
			Title: model.NewLocalizedText(
				"生成 AI による会議議事録支援",
				"Generative AI Meeting Minutes",
			),
			Overview: model.NewLocalizedText(
				"LLM を活用して議事録作成を支援するワークフローを設計しました。",
				"Designed a workflow that leverages LLMs to streamline meeting minutes creation.",
			),
			Outcome: model.NewLocalizedText(
				"社内ユーザーテストで 86% が手動作業削減を実感しました。",
				"Internal tests showed 86% of participants reduced manual effort.",
			),
			Outlook: model.NewLocalizedText(
				"Azure OpenAI での多言語対応とプライバシー機能の検証を進めています。",
				"Working on multilingual support and privacy controls with Azure OpenAI.",
			),
			ExternalURL:       "https://example.com/research/llm-meeting-minutes",
			HighlightImageURL: "",
			ImageAlt: model.NewLocalizedText(
				"議事録支援アシスタントの画面",
				"Meeting minutes assistant interface",
			),
			PublishedAt: time.Date(2024, 4, 24, 9, 0, 0, 0, time.UTC),
			IsDraft:     true,
			CreatedAt:   time.Date(2024, 4, 21, 9, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 4, 25, 9, 0, 0, 0, time.UTC),
			Tags: []model.ResearchTag{
				{ID: 1, EntryID: 2, Value: "generative-ai", SortOrder: 1},
				{ID: 2, EntryID: 2, Value: "productivity", SortOrder: 2},
			},
			Links: []model.ResearchLink{
				{
					ID:        1,
					EntryID:   2,
					Type:      model.ResearchLinkTypeVideo,
					Label:     model.NewLocalizedText("デモ動画", "Demo Video"),
					URL:       "https://example.com/research/llm-demo",
					SortOrder: 1,
				},
			},
			Tech: []model.TechMembership{
				{
					MembershipID: 3,
					EntityType:   researchEntityType,
					EntityID:     2,
					Tech: model.TechCatalogEntry{
						ID:          3,
						Slug:        "python",
						DisplayName: "Python",
						Level:       model.TechLevelAdvanced,
						SortOrder:   3,
						Active:      true,
						CreatedAt:   time.Date(2018, 5, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt:   time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC),
					},
					Context:   model.TechContextPrimary,
					SortOrder: 1,
				},
			},
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

	defaultMeetingReservations = []model.MeetingReservation{
		{
			ID:              1,
			LookupHash:      "lookup-akari",
			Name:            "Akari Yamada",
			Email:           "akari@example.com",
			Topic:           "プロジェクト相談",
			Message:         "最新のリサーチ成果について詳しく伺いたいです。",
			StartAt:         time.Date(2024, 5, 10, 10, 0, 0, 0, time.UTC),
			EndAt:           time.Date(2024, 5, 10, 10, 30, 0, 0, time.UTC),
			DurationMinutes: 30,
			GoogleEventID:   "evt-akari",
			Status:          model.MeetingReservationStatusPending,
			CreatedAt:       time.Date(2024, 5, 1, 9, 0, 0, 0, time.UTC),
			UpdatedAt:       time.Date(2024, 5, 1, 9, 0, 0, 0, time.UTC),
		},
		{
			ID:              2,
			LookupHash:      "lookup-lucas",
			Name:            "Lucas Chen",
			Email:           "lucas@example.com",
			Topic:           "AIプロダクト連携",
			Message:         "AI駆動の活用について相談したいです。",
			StartAt:         time.Date(2024, 5, 12, 13, 0, 0, 0, time.UTC),
			EndAt:           time.Date(2024, 5, 12, 13, 45, 0, 0, time.UTC),
			DurationMinutes: 45,
			GoogleEventID:   "evt-lucas",
			Status:          model.MeetingReservationStatusConfirmed,
			ConfirmationSentAt: func() *time.Time {
				ts := time.Date(2024, 5, 3, 9, 0, 0, 0, time.UTC)
				return &ts
			}(),
			LastNotificationSentAt: func() *time.Time {
				ts := time.Date(2024, 5, 3, 9, 0, 0, 0, time.UTC)
				return &ts
			}(),
			CreatedAt: time.Date(2024, 5, 2, 9, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 5, 3, 9, 0, 0, 0, time.UTC),
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
