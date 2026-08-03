//go:build integration

package tests

import (
	"context"
	"testing"
	"time"

	planrepo "vedo-edutrack/backend/internal/modules/planmanagement/adapters/repository"
	plansqlc "vedo-edutrack/backend/internal/modules/planmanagement/adapters/repository/sqlc"
)

// TestPlanRepositorySaveAndReadBack — INSERT-only plan snapshot (plan + steps
// + checkpoints) round-trips through the repository.
func TestPlanRepositorySaveAndReadBack(t *testing.T) {
	pool := startPostgres(t)
	repo := planrepo.NewPlanRepository(pool)
	ctx := context.Background()

	learnerID := "77777777-7777-7777-7777-777777777777"
	start := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)

	plan := plansqlc.LearningPlanRow{
		LearnerID:       learnerID,
		GoalModuleID:    "chemistry",
		PedagogyConcept: strPtr("essentialism"),
		OntologyVersion: "v1",
		Status:          "active",
		TimelineStart:   start,
		TimelineEnd:     end,
		SnapshotHash:    "hash-1",
		Version:         1,
		FixedAt:         ptrTime(start),
	}
	steps := []plansqlc.PlanStepRow{
		{ModuleID: "percent", Position: 0, Horizon: "near", IsEssential: true, PlannedStart: start, PlannedEnd: start.AddDate(0, 1, 0)},
		{ModuleID: "solutions", Position: 1, Horizon: "mid", IsEssential: false, PlannedStart: start.AddDate(0, 1, 0), PlannedEnd: start.AddDate(0, 2, 0)},
	}
	checkpoints := []plansqlc.PlanCheckpointRow{
		{Name: "checkpoint-1", CheckpointDate: start.AddDate(0, 2, 0), FrameworkID: strPtr("fgos"), TargetCoveragePercent: 50},
	}

	planID, err := repo.SavePlan(ctx, plan, steps, checkpoints)
	if err != nil {
		t.Fatalf("save plan: %v", err)
	}
	if planID == "" {
		t.Fatal("expected non-empty plan id")
	}

	// Read-back: plan + steps + checkpoints.
	got, gotSteps, gotCheckpoints, err := repo.GetPlan(ctx, planID)
	if err != nil {
		t.Fatalf("get plan: %v", err)
	}
	if got.GoalModuleID != "chemistry" || got.Status != "active" {
		t.Fatalf("unexpected plan: %+v", got)
	}
	if len(gotSteps) != 2 {
		t.Fatalf("expected 2 steps, got %d", len(gotSteps))
	}
	if len(gotCheckpoints) != 1 {
		t.Fatalf("expected 1 checkpoint, got %d", len(gotCheckpoints))
	}
	// Steps are ordered by position.
	if gotSteps[0].ModuleID != "percent" || gotSteps[1].ModuleID != "solutions" {
		t.Fatalf("unexpected step order: %+v", gotSteps)
	}

	// List by learner returns the plan.
	plans, err := repo.ListPlansByLearner(ctx, learnerID)
	if err != nil {
		t.Fatalf("list plans: %v", err)
	}
	if len(plans) != 1 || plans[0].ID != planID {
		t.Fatalf("unexpected plan list: %+v", plans)
	}
}

// TestPlanRepositoryInsertOnlySemantics — a second SavePlan with the same id
// is impossible (id is DB-generated): plans are immutable snapshots.
func TestPlanRepositoryInsertOnlySemantics(t *testing.T) {
	pool := startPostgres(t)
	repo := planrepo.NewPlanRepository(pool)
	ctx := context.Background()

	learnerID := "88888888-8888-8888-8888-888888888888"
	plan := plansqlc.LearningPlanRow{
		LearnerID:       learnerID,
		GoalModuleID:    "percent",
		OntologyVersion: "v1",
		Status:          "active",
		TimelineStart:   time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
		TimelineEnd:     time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC),
		SnapshotHash:    "hash-2",
		Version:         1,
	}

	firstID, err := repo.SavePlan(ctx, plan, nil, nil)
	if err != nil {
		t.Fatalf("first save: %v", err)
	}
	secondID, err := repo.SavePlan(ctx, plan, nil, nil)
	if err != nil {
		t.Fatalf("second save: %v", err)
	}
	// Two distinct snapshots — the INSERT-only model does not overwrite.
	if firstID == secondID {
		t.Fatal("expected distinct plan ids (INSERT-only snapshots)")
	}

	plans, err := repo.ListPlansByLearner(ctx, learnerID)
	if err != nil {
		t.Fatalf("list plans: %v", err)
	}
	if len(plans) != 2 {
		t.Fatalf("expected 2 snapshots, got %d", len(plans))
	}
}
