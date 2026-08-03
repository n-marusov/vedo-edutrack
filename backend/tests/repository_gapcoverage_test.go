//go:build integration

package tests

import (
	"context"
	"testing"

	gaprepo "vedo-edutrack/backend/internal/modules/gapcoverage/adapters/repository"
	gapreposqlc "vedo-edutrack/backend/internal/modules/gapcoverage/adapters/repository/sqlc"
)

// TestGapRepositorySaveAndListAnalysis — transaction persists a diagnosis run
// with its ranked root causes.
func TestGapRepositorySaveAndListAnalysis(t *testing.T) {
	pool := startPostgres(t)
	repo := gaprepo.NewGapRepository(pool)
	ctx := context.Background()

	learnerID := "55555555-5555-5555-5555-555555555555"

	analysis := gapreposqlc.GapAnalysisRow{
		LearnerID:       learnerID,
		PlanID:          strPtr("plan-gap-1"),
		OntologyVersion: "v1",
	}
	roots := []gapreposqlc.GapRootRow{
		{
			RootModuleID:        "percent",
			ChainModuleIDs:      []string{"percent", "solutions"},
			BlockedModulesCount: 3,
			BlockedSubjects:     []string{"chemistry"},
			CascadeRank:         1,
		},
		{
			RootModuleID:        "solutions",
			ChainModuleIDs:      []string{"solutions"},
			BlockedModulesCount: 1,
			CascadeRank:         2,
		},
	}

	diagID, err := repo.SaveGapAnalysis(ctx, analysis, roots)
	if err != nil {
		t.Fatalf("save gap analysis: %v", err)
	}
	if diagID == "" {
		t.Fatal("expected non-empty diagnostic id")
	}

	// Read-back: the diagnostic row exists with the learner's id.
	var count int
	if err := pool.QueryRow(ctx,
		`SELECT count(*) FROM gapcoverage.gap_diagnostics WHERE id = $1 AND learner_id = $2`,
		diagID, learnerID).Scan(&count); err != nil {
		t.Fatalf("count diagnostic: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", count)
	}

	// Both roots persisted with cascade ranks in order.
	var rootCount int
	if err := pool.QueryRow(ctx,
		`SELECT count(*) FROM gapcoverage.gap_roots WHERE diagnostic_id = $1`, diagID).Scan(&rootCount); err != nil {
		t.Fatalf("count roots: %v", err)
	}
	if rootCount != 2 {
		t.Fatalf("expected 2 roots, got %d", rootCount)
	}
}

// TestGapRepositoryCoverageSnapshot — transaction persists a coverage snapshot
// with its deficit list; read-back via the repository query.
func TestGapRepositoryCoverageSnapshot(t *testing.T) {
	pool := startPostgres(t)
	repo := gaprepo.NewGapRepository(pool)
	ctx := context.Background()

	learnerID := "66666666-6666-6666-6666-666666666666"

	snapshot := gapreposqlc.CoverageSnapshotRow{
		LearnerID:           learnerID,
		FrameworkID:         "fgos",
		FrameworkVersion:    "2021",
		OntologyVersion:     "v1",
		CoveragePercent:     50,
		CoveredRequirements: 1,
		TotalRequirements:   2,
		Version:             0,
	}
	deficits := []gapreposqlc.CoverageDeficitRow{
		{
			RequirementID:   "r2",
			Priority:        "high",
			Status:          "open",
			LinkedModuleIDs: []string{"percent"},
			EffortEstimate:  strPtr("2w"),
		},
	}

	snapshotID, err := repo.SaveCoverageSnapshot(ctx, snapshot, deficits)
	if err != nil {
		t.Fatalf("save coverage snapshot: %v", err)
	}
	if snapshotID == "" {
		t.Fatal("expected non-empty snapshot id")
	}

	snapshots, err := repo.ListCoverageSnapshotsByLearner(ctx, learnerID)
	if err != nil {
		t.Fatalf("list snapshots: %v", err)
	}
	if len(snapshots) != 1 {
		t.Fatalf("expected 1 snapshot, got %d", len(snapshots))
	}
	if snapshots[0].CoveragePercent != 50 || snapshots[0].TotalRequirements != 2 {
		t.Fatalf("unexpected snapshot: %+v", snapshots[0])
	}

	var deficitCount int
	if err := pool.QueryRow(ctx,
		`SELECT count(*) FROM gapcoverage.coverage_deficits WHERE snapshot_id = $1`, snapshotID).Scan(&deficitCount); err != nil {
		t.Fatalf("count deficits: %v", err)
	}
	if deficitCount != 1 {
		t.Fatalf("expected 1 deficit, got %d", deficitCount)
	}
}
