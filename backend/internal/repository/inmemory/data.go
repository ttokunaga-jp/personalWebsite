package inmemory

import "github.com/takumi/personal-website/internal/model"

var (
	defaultProfile = &model.Profile{
		Name:        "Takumi Example",
		Title:       "Software Engineer / Researcher",
		Affiliation: "Example University",
		Lab:         "Human-Computer Interaction Lab",
		Summary:     "Building delightful experiences backed by resilient infrastructure.",
		Skills:      []string{"Go", "React", "GCP", "Machine Learning"},
	}

	defaultProjects = []model.Project{
		{
			ID:          1,
			Title:       "AI Assisted Portfolio",
			Description: "An experimental platform for AI assisted content authoring.",
			TechStack:   []string{"Go", "React", "GCP"},
			LinkURL:     "https://example.dev/projects/ai-portfolio",
			Year:        2023,
		},
	}

	defaultResearch = []model.Research{
		{
			ID:        1,
			Title:     "Adaptive Scheduling for Remote Teams",
			ContentMD: "Exploration of reinforcement learning to optimise meeting schedules.",
			Year:      2022,
		},
	}
)
