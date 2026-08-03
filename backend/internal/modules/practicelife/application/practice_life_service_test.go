package practicelife

import (
	"testing"

	"go.uber.org/zap"

	domain "vedo-edutrack/backend/internal/modules/practicelife/domain"
)

func testService(t *testing.T) *PracticeLifeService {
	t.Helper()
	stories := []domain.Story{
		{ID: "s1", Title: "Проценты в жизни", LinkedModules: []string{"percent", "math-5-11"}, Subjects: []string{"Математика"}, RealWorld: "Скидки и банки.", ReadingMinutes: 3},
		{ID: "s2", Title: "Растворы", LinkedModules: []string{"solutions"}, Subjects: []string{"Химия"}, RealWorld: "Медицина.", ReadingMinutes: 4},
		{ID: "s3", Title: "Химия и экология", LinkedModules: []string{"chemistry"}, Subjects: []string{"Химия", "Биология"}, RealWorld: "Фотосинтез.", ReadingMinutes: 5},
	}
	projects := []domain.ProjectIdea{
		{ID: "p1", Title: "Лаборатория дома", Modules: []string{"solutions", "chemistry"}, DifficultyLevel: domain.DifficultyMedium, ExpectedOutcome: "Опыты."},
		{ID: "p2", Title: "Экология двора", Modules: []string{"chemistry", "percent"}, DifficultyLevel: domain.DifficultyBasic, ExpectedOutcome: "Расчёт."},
		{ID: "p3", Title: "Большой проект", Modules: []string{"percent", "solutions", "chemistry"}, DifficultyLevel: domain.DifficultyAdvanced, ExpectedOutcome: "Комплексное исследование."},
	}
	return NewPracticeLifeService(zap.NewNop(), stories, projects)
}

func TestStoriesForModule(t *testing.T) {
	svc := testService(t)
	stories := svc.StoriesForModule("percent")
	if len(stories) != 1 {
		t.Fatalf("expected 1 story for percent, got %d", len(stories))
	}
	if stories[0].ID != "s1" {
		t.Fatalf("story id=%s, want s1", stories[0].ID)
	}
}

func TestStoriesForUnknownModule(t *testing.T) {
	svc := testService(t)
	stories := svc.StoriesForModule("unknown")
	if len(stories) != 0 {
		t.Fatalf("expected 0 stories for unknown module, got %d", len(stories))
	}
}

func TestRecommendStoriesDeduplicates(t *testing.T) {
	svc := testService(t)
	// s1 is linked to both percent and math-5-11 — must appear once.
	stories := svc.RecommendStories([]string{"percent", "math-5-11"})
	if len(stories) != 1 {
		t.Fatalf("expected 1 deduplicated story, got %d", len(stories))
	}
}

func TestSuggestProjectsEligibility(t *testing.T) {
	svc := testService(t)
	ready := map[string]bool{"solutions": true, "chemistry": true}
	projects := svc.SuggestProjects(ready)
	// p1 (solutions+chemistry) eligible; p3 (percent+solutions+chemistry) = 66% < 80% — not eligible.
	if len(projects) != 1 || projects[0].ID != "p1" {
		t.Fatalf("expected only p1 eligible, got %+v", projects)
	}
}

func TestProjectCountAndStoryCount(t *testing.T) {
	svc := testService(t)
	if svc.StoryCount() != 3 {
		t.Fatalf("story count=%d, want 3", svc.StoryCount())
	}
	if svc.ProjectCount() != 3 {
		t.Fatalf("project count=%d, want 3", svc.ProjectCount())
	}
}
