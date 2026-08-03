// Package adapters provides infrastructure adapters for the practicelife bounded context.
package adapters

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	app "vedo-edutrack/backend/internal/modules/practicelife/application"
	domain "vedo-edutrack/backend/internal/modules/practicelife/domain"
)

// Handler serves practice-life REST endpoints.
type Handler struct {
	svc    *app.PracticeLifeService
	logger *zap.Logger
}

// NewHandler creates a practice-life HTTP handler.
func NewHandler(svc *app.PracticeLifeService, logger *zap.Logger) *Handler {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &Handler{svc: svc, logger: logger.Named("practice.handler")}
}

// Register routes onto the given chi router.
func (h *Handler) Register(r chi.Router) {
	r.Get("/modules/{moduleID}/stories", h.handleStoriesForModule)
	r.Get("/modules/{moduleID}/projects", h.handleProjectsForModule)
	r.Get("/learners/{learnerID}/recommended-stories", h.handleRecommendStories)
	r.Get("/learners/{learnerID}/recommended-projects", h.handleSuggestProjects)
}

// handleStoriesForModule returns stories linked to a module.
func (h *Handler) handleStoriesForModule(w http.ResponseWriter, r *http.Request) {
	moduleID := chi.URLParam(r, "moduleID")
	stories := h.svc.StoriesForModule(moduleID)
	h.logger.Info("stories for module",
		zap.String("module_id", moduleID),
		zap.Int("count", len(stories)),
	)
	writeJSON(w, http.StatusOK, mapStoriesToResponse(stories))
}

// handleProjectsForModule returns project ideas linked to a module.
func (h *Handler) handleProjectsForModule(w http.ResponseWriter, r *http.Request) {
	moduleID := chi.URLParam(r, "moduleID")
	projects := h.svc.ProjectsForModule(moduleID)
	h.logger.Info("projects for module",
		zap.String("module_id", moduleID),
		zap.Int("count", len(projects)),
	)
	writeJSON(w, http.StatusOK, mapProjectsToResponse(projects))
}

// handleRecommendStories returns story recommendations for a learner.
func (h *Handler) handleRecommendStories(w http.ResponseWriter, r *http.Request) {
	learnerID := chi.URLParam(r, "learnerID")
	masteredModules := parseModuleIDs(r.URL.Query().Get("mastered_modules"))
	stories := h.svc.RecommendStories(masteredModules)
	h.logger.Info("recommended stories",
		zap.String("learner_id", learnerID),
		zap.Int("count", len(stories)),
	)
	writeJSON(w, http.StatusOK, mapStoriesToResponse(stories))
}

// handleSuggestProjects returns project suggestions for a learner.
func (h *Handler) handleSuggestProjects(w http.ResponseWriter, r *http.Request) {
	learnerID := chi.URLParam(r, "learnerID")
	masteredModules := parseModuleIDs(r.URL.Query().Get("mastered_modules"))
	availableModules := parseModuleIDs(r.URL.Query().Get("available_modules"))

	readySet := make(map[string]bool, len(masteredModules)+len(availableModules))
	for _, m := range masteredModules {
		readySet[m] = true
	}
	for _, m := range availableModules {
		readySet[m] = true
	}

	projects := h.svc.SuggestProjects(readySet)
	h.logger.Info("suggested projects",
		zap.String("learner_id", learnerID),
		zap.Int("count", len(projects)),
	)
	writeJSON(w, http.StatusOK, mapProjectsToResponse(projects))
}

// --- helper types and functions ---

type storyResponse struct {
	ID             string   `json:"id"`
	Title          string   `json:"title"`
	Text           string   `json:"text"`
	LinkedModules  []string `json:"linked_modules"`
	Subjects       []string `json:"subjects"`
	RealWorld      string   `json:"real_world"`
	ReadingMinutes int      `json:"reading_minutes"`
}

type projectResponse struct {
	ID              string   `json:"id"`
	Title           string   `json:"title"`
	Modules         []string `json:"modules"`
	DifficultyLevel string   `json:"difficulty_level"`
	ExpectedOutcome string   `json:"expected_outcome"`
}

func mapStoriesToResponse(stories []domain.Story) []storyResponse {
	result := make([]storyResponse, len(stories))
	for i, s := range stories {
		result[i] = storyResponse{
			ID:             s.ID,
			Title:          s.Title,
			Text:           s.Text,
			LinkedModules:  s.LinkedModules,
			Subjects:       s.Subjects,
			RealWorld:      s.RealWorld,
			ReadingMinutes: s.ReadingMinutes,
		}
	}
	return result
}

func mapProjectsToResponse(projects []domain.ProjectIdea) []projectResponse {
	result := make([]projectResponse, len(projects))
	for i, p := range projects {
		result[i] = projectResponse{
			ID:              p.ID,
			Title:           p.Title,
			Modules:         p.Modules,
			DifficultyLevel: string(p.DifficultyLevel),
			ExpectedOutcome: p.ExpectedOutcome,
		}
	}
	return result
}

// parseModuleIDs splits a comma-separated parameter into a slice.
func parseModuleIDs(raw string) []string {
	if raw == "" {
		return nil
	}
	var result []string
	start := 0
	for i := 0; i < len(raw); i++ {
		if raw[i] == ',' {
			if trimmed := trimSpaces(raw[start:i]); trimmed != "" {
				result = append(result, trimmed)
			}
			start = i + 1
		}
	}
	if trimmed := trimSpaces(raw[start:]); trimmed != "" {
		result = append(result, trimmed)
	}
	return result
}

func trimSpaces(s string) string {
	for len(s) > 0 && s[0] == ' ' {
		s = s[1:]
	}
	for len(s) > 0 && s[len(s)-1] == ' ' {
		s = s[:len(s)-1]
	}
	return s
}

// writeJSON encodes any value as JSON with the given status code.
func writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
