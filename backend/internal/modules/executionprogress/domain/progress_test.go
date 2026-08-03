package executionprogress

import (
	"testing"
	"time"
)

func TestPlanVsActualDeviationAndEventThreshold(t *testing.T) {
	plan := FixedPlan{ID: "p1", LearnerID: "l1", Modules: []PlannedModule{{ModuleID: "m1", PlannedStart: day(0), PlannedEnd: day(10)}}}
	progress := []ModuleProgress{{ModuleID: "m1", Status: StatusMastered, MasteredAt: ptrTime(day(12))}}
	comparison := ComparePlanActual(plan, progress)
	if len(comparison.Deviations) != 1 || comparison.Deviations[0].DeltaDays != 2 {
		t.Fatalf("unexpected deviations: %+v", comparison.Deviations)
	}
	if len(comparison.Events) != 1 {
		t.Fatalf("expected deviation event above ±15%% threshold, got %+v", comparison.Events)
	}
}

func TestModuleOutsidePlanFlaggedAsDivergence(t *testing.T) {
	comparison := ComparePlanActual(FixedPlan{ID: "p1"}, []ModuleProgress{{ModuleID: "extra", Status: StatusMastered}})
	if len(comparison.Divergences) != 1 || comparison.Divergences[0] != "extra" {
		t.Fatalf("expected outside-plan divergence, got %+v", comparison.Divergences)
	}
}

func TestBinaryForecast(t *testing.T) {
	onTrack := ForecastReadiness(ForecastInput{CompletedModules: 5, RemainingModules: 5, DaysElapsed: 5, DaysRemaining: 6})
	if onTrack.Status != ForecastOnTrack {
		t.Fatalf("status=%s, want on-track", onTrack.Status)
	}
	lowAccuracy := ForecastReadiness(ForecastInput{CompletedModules: 0, RemainingModules: 10, DaysElapsed: 0, DaysRemaining: 10})
	if !lowAccuracy.LowAccuracy {
		t.Fatalf("expected low accuracy label for insufficient data")
	}
}

func day(offset int) time.Time       { return time.Date(2026, 1, 1+offset, 0, 0, 0, 0, time.UTC) }
func ptrTime(t time.Time) *time.Time { return &t }
