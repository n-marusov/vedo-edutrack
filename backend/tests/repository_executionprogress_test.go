//go:build integration

package tests

import (
	"context"
	"testing"
	"time"

	execrepo "vedo-edutrack/backend/internal/modules/executionprogress/adapters/repository"
	execreposqlc "vedo-edutrack/backend/internal/modules/executionprogress/adapters/repository/sqlc"
)

func ptrTime(t time.Time) *time.Time { return &t }

// TestProgressRepositoryUpsertAndList — real PostgreSQL via testcontainers:
// insert → conflict update (version bump) → read-back by learner and plan.
func TestProgressRepositoryUpsertAndList(t *testing.T) {
	pool := startPostgres(t)
	repo := execrepo.NewProgressRepository(pool)
	ctx := context.Background()

	learnerID := "22222222-2222-2222-2222-222222222222"
	planID := "plan-progress-1"
	now := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)

	// 1. First insert.
	row := execreposqlc.ModuleProgressRow{
		LearnerID: learnerID, ModuleID: "percent", PlanID: &planID,
		Status: "in_progress", Version: 0,
	}
	if err := repo.UpsertModuleProgress(ctx, row); err != nil {
		t.Fatalf("first upsert: %v", err)
	}

	// 2. Read-back by learner.
	got, err := repo.ListProgressByLearner(ctx, learnerID)
	if err != nil {
		t.Fatalf("list by learner: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 progress row, got %d", len(got))
	}
	if got[0].ModuleID != "percent" || got[0].Status != "in_progress" {
		t.Fatalf("unexpected row: %+v", got[0])
	}
	if got[0].Version != 0 {
		t.Fatalf("expected version 0 after first insert, got %d", got[0].Version)
	}

	// 3. Conflict update → version must bump (optimistic lock, SQL side).
	row.Status = "mastered"
	row.MasteredAt = ptrTime(now)
	if err := repo.UpsertModuleProgress(ctx, row); err != nil {
		t.Fatalf("second upsert: %v", err)
	}
	got, err = repo.ListProgressByLearner(ctx, learnerID)
	if err != nil {
		t.Fatalf("list by learner after update: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected still 1 row, got %d", len(got))
	}
	if got[0].Version != 1 {
		t.Fatalf("expected version 1 after update, got %d", got[0].Version)
	}
	if got[0].Status != "mastered" {
		t.Fatalf("expected status mastered, got %q", got[0].Status)
	}

	// 4. List by plan sees the same row.
	byPlan, err := repo.ListProgressByPlan(ctx, planID)
	if err != nil {
		t.Fatalf("list by plan: %v", err)
	}
	if len(byPlan) != 1 || byPlan[0].ModuleID != "percent" {
		t.Fatalf("unexpected plan rows: %+v", byPlan)
	}
}

// TestProgressRepositoryDeviationEventDedup — idempotent insert on the dedup
// key (learner_id, plan_id, deviation_type, detected_at).
func TestProgressRepositoryDeviationEventDedup(t *testing.T) {
	pool := startPostgres(t)
	repo := execrepo.NewProgressRepository(pool)
	ctx := context.Background()

	learnerID := "33333333-3333-3333-3333-333333333333"
	planID := "plan-progress-2"
	detectedAt := time.Date(2026, 9, 5, 8, 0, 0, 0, time.UTC)

	event := execreposqlc.DeviationEventRow{
		EventID:        "d4c1e2a5-1f0e-4c3a-9b2f-6a3c8e7d1a2b",
		LearnerID:      learnerID,
		PlanID:         &planID,
		ModuleID:       strPtr("percent"),
		DeviationType:  "late",
		DeviationValue: 3.0,
		Threshold:      0.15,
		DetectedAt:     detectedAt,
	}

	id, inserted, err := repo.InsertDeviationEvent(ctx, event)
	if err != nil {
		t.Fatalf("first insert: %v", err)
	}
	if !inserted || id == "" {
		t.Fatalf("expected insert, got id=%q inserted=%t", id, inserted)
	}

	// Duplicate (same dedup key, different event_id) must be ignored.
	dup := event
	dup.EventID = "e5f2b3c6-2a1d-4b4e-8c3f-7b4d9e8f2a3b"
	id2, inserted2, err := repo.InsertDeviationEvent(ctx, dup)
	if err != nil {
		t.Fatalf("duplicate insert: %v", err)
	}
	if inserted2 || id2 != "" {
		t.Fatalf("expected dedup, got id=%q inserted=%t", id2, inserted2)
	}
}

// TestProgressRepositoryReadinessForecast — insert and read-back.
func TestProgressRepositoryReadinessForecast(t *testing.T) {
	pool := startPostgres(t)
	repo := execrepo.NewProgressRepository(pool)
	ctx := context.Background()

	learnerID := "44444444-4444-4444-4444-444444444444"
	expected := time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)

	forecast := execreposqlc.ReadinessForecastRow{
		LearnerID:        learnerID,
		PlanID:           strPtr("plan-progress-3"),
		Status:           "on_track",
		ExpectedDate:     &expected,
		RemainingModules: 4,
		Velocity:         1.5,
		KeyRisks:         []string{"percent", "chemistry"},
		DataConfidence:   "low",
	}
	id, err := repo.SaveReadinessForecast(ctx, forecast)
	if err != nil {
		t.Fatalf("save forecast: %v", err)
	}
	if id == "" {
		t.Fatal("expected non-empty forecast id")
	}
}

func strPtr(s string) *string { return &s }
