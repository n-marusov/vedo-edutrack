package practicelife

import (
	"testing"
)

// TestLaunchContentVolume enforces the M2 launch content exit criterion:
// ≥ 50 stories and ≥ 30 project ideas (REQ-FR-practice.stories.context-delivery
// AC: story linked to 1-3 topics; REQ-FR-practice.projects.suggest-cross-subject
// AC: project requires modules from ≥ 2 subjects).
func TestLaunchContentVolume(t *testing.T) {
	stories := LaunchStories()
	if len(stories) < 50 {
		t.Fatalf("launch stories=%d, want ≥ 50", len(stories))
	}
	projects := LaunchProjects()
	if len(projects) < 30 {
		t.Fatalf("launch projects=%d, want ≥ 30", len(projects))
	}
}

// TestLaunchStoriesValid verifies every story has 1-3 linked modules, a
// mandatory real-world section, and ≤ 5 min reading time.
func TestLaunchStoriesValid(t *testing.T) {
	stories := LaunchStories()
	seen := map[string]bool{}
	for _, s := range stories {
		if seen[s.ID] {
			t.Fatalf("duplicate story id %q", s.ID)
		}
		seen[s.ID] = true
		if s.Title == "" {
			t.Fatalf("story %q: empty title", s.ID)
		}
		if len(s.LinkedModules) == 0 || len(s.LinkedModules) > 3 {
			t.Fatalf("story %q: linked modules=%d, want 1-3", s.ID, len(s.LinkedModules))
		}
		if s.RealWorld == "" {
			t.Fatalf("story %q: missing real-world section", s.ID)
		}
		if s.ReadingMinutes <= 0 || s.ReadingMinutes > 5 {
			t.Fatalf("story %q: reading minutes=%d, want 1-5", s.ID, s.ReadingMinutes)
		}
	}
}

// TestLaunchProjectsValid verifies every project has a graded difficulty,
// an expected outcome, and requires modules from ≥ 2 distinct subjects.
func TestLaunchProjectsValid(t *testing.T) {
	projects := LaunchProjects()
	seen := map[string]bool{}
	for _, p := range projects {
		if seen[p.ID] {
			t.Fatalf("duplicate project id %q", p.ID)
		}
		seen[p.ID] = true
		if p.Title == "" {
			t.Fatalf("project %q: empty title", p.ID)
		}
		if len(p.Modules) < 2 {
			t.Fatalf("project %q: modules=%d, want ≥ 2", p.ID, len(p.Modules))
		}
		if p.DifficultyLevel != "basic" && p.DifficultyLevel != "medium" && p.DifficultyLevel != "advanced" {
			t.Fatalf("project %q: invalid difficulty %q", p.ID, p.DifficultyLevel)
		}
		if p.ExpectedOutcome == "" {
			t.Fatalf("project %q: missing expected outcome", p.ID)
		}
	}
}

// TestLaunchContentUniqueIDs verifies story and project IDs do not collide.
func TestLaunchContentUniqueIDs(t *testing.T) {
	stories := LaunchStories()
	projects := LaunchProjects()
	projectIDs := map[string]bool{}
	for _, p := range projects {
		projectIDs[p.ID] = true
	}
	for _, s := range stories {
		if projectIDs[s.ID] {
			t.Fatalf("story id %q collides with a project id", s.ID)
		}
	}
}
