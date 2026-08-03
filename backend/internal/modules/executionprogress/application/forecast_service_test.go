package executionprogress

import (
	"context"
	"testing"
	"time"

	"go.uber.org/zap"

	domain "vedo-edutrack/backend/internal/modules/executionprogress/domain"
)

// stubProgressRepo returns fixed progress data for forecast tests.
type stubProgressRepo struct {
	progress       []domain.ModuleProgress
	plannedModules int
	remaining      int
	progressErr    error
	remainingErr   error
}

func (s stubProgressRepo) GetProgress(_ context.Context, _ string) ([]domain.ModuleProgress, error) {
	return s.progress, s.progressErr
}

func (s stubProgressRepo) GetPlannedModuleCount(_ context.Context, _, _ string) (int, error) {
	return s.plannedModules, nil
}

func (s stubProgressRepo) GetRemainingModules(_ context.Context, _, _ string) (int, error) {
	return s.remaining, s.remainingErr
}

func TestForecastServiceOnTrack(t *testing.T) {
	now := time.Now()
	repo := stubProgressRepo{
		progress: []domain.ModuleProgress{
			{ModuleID: "m1", Status: domain.StatusMastered, MasteredAt: &now},
			{ModuleID: "m2", Status: domain.StatusMastered, MasteredAt: &now},
			{ModuleID: "m3", Status: domain.StatusMastered, MasteredAt: &now},
		},
		remaining: 3,
	}
	svc := NewForecastService(repo, zap.NewNop())

	result, err := svc.ForecastReadiness(context.Background(), "l1", "p1", 10, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Status != domain.ForecastOnTrack {
		t.Fatalf("status=%s, want on-track", result.Status)
	}
	if result.RemainingModules != 3 {
		t.Fatalf("remaining=%d, want 3", result.RemainingModules)
	}
}

func TestForecastServiceNotOnTrack(t *testing.T) {
	now := time.Now()
	repo := stubProgressRepo{
		progress: []domain.ModuleProgress{
			{ModuleID: "m1", Status: domain.StatusMastered, MasteredAt: &now},
		},
		remaining: 9,
	}
	svc := NewForecastService(repo, zap.NewNop())

	result, err := svc.ForecastReadiness(context.Background(), "l1", "p1", 10, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Status != domain.ForecastNotOnTrack {
		t.Fatalf("status=%s, want not-on-track", result.Status)
	}
	if len(result.KeyRisks) == 0 {
		t.Fatal("expected key risks for not-on-track forecast")
	}
}

func TestForecastServiceLowConfidence(t *testing.T) {
	repo := stubProgressRepo{
		progress:  []domain.ModuleProgress{},
		remaining: 10,
	}
	svc := NewForecastService(repo, zap.NewNop())

	result, err := svc.ForecastReadiness(context.Background(), "l1", "p1", 0, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.DataConfidence != "low" {
		t.Fatalf("confidence=%s, want low", result.DataConfidence)
	}
}
