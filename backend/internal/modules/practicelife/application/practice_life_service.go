// Package practicelife provides the practicelife bounded context application layer.
package practicelife

import (
	"go.uber.org/zap"

	domain "vedo-edutrack/backend/internal/modules/practicelife/domain"
)

// PracticeLifeService exposes story recommendation and project suggestion use cases.
type PracticeLifeService struct {
	stories *domain.StoryCatalog
	ideas   *domain.ProjectIdeaCatalog
	logger  *zap.Logger
}

// NewPracticeLifeService creates a practice-life application service backed
// by ontology-sourced story and project idea catalogs.
func NewPracticeLifeService(logger *zap.Logger, stories []domain.Story, ideas []domain.ProjectIdea) *PracticeLifeService {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &PracticeLifeService{
		stories: domain.NewStoryCatalog(stories),
		ideas:   domain.NewProjectIdeaCatalog(ideas),
		logger:  logger.Named("practice.life"),
	}
}

// StoriesForModule returns stories linked to a module via appliesTo/enriches graph links.
func (s *PracticeLifeService) StoriesForModule(moduleID string) []domain.Story {
	result := s.stories.ByModule(moduleID)
	s.logger.Info("stories for module",
		zap.String("module_id", moduleID),
		zap.Int("count", len(result)),
	)
	return result
}

// ProjectsForModule returns project ideas linked to a module.
func (s *PracticeLifeService) ProjectsForModule(moduleID string) []domain.ProjectIdea {
	result := s.ideas.ByModule(moduleID)
	s.logger.Info("projects for module",
		zap.String("module_id", moduleID),
		zap.Int("count", len(result)),
	)
	return result
}

// RecommendStories returns stories recommended at module mastery.
// Stories are linked via appliesTo/enriches graph edges in the ontology.
func (s *PracticeLifeService) RecommendStories(masteredModuleIDs []string) []domain.Story {
	seen := map[string]bool{}
	var result []domain.Story
	for _, modID := range masteredModuleIDs {
		for _, story := range s.stories.ByModule(modID) {
			if !seen[story.ID] {
				seen[story.ID] = true
				result = append(result, story)
			}
		}
	}
	s.logger.Info("story recommendations generated",
		zap.Int("mastered_modules", len(masteredModuleIDs)),
		zap.Int("recommended", len(result)),
	)
	return result
}

// SuggestProjects returns project ideas where ≥ 80% of required modules
// are in the mastered or available set (cross-subject readiness gate).
func (s *PracticeLifeService) SuggestProjects(masteredOrAvailable map[string]bool) []domain.ProjectIdea {
	result := s.ideas.SuggestEligible(masteredOrAvailable)
	s.logger.Info("project suggestions generated",
		zap.Int("eligible", len(result)),
		zap.Int("total_ideas", s.ideas.Count()),
	)
	return result
}

// StoryCount returns the total number of stories in the catalog.
func (s *PracticeLifeService) StoryCount() int { return s.stories.Count() }

// ProjectCount returns the total number of project ideas in the catalog.
func (s *PracticeLifeService) ProjectCount() int { return s.ideas.Count() }
