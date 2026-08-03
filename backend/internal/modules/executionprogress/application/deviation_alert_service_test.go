package executionprogress

import (
	"context"
	"testing"

	"go.uber.org/zap"

	domain "vedo-edutrack/backend/internal/modules/executionprogress/domain"
)

func TestDeviationAlertThresholdNotExceeded(t *testing.T) {
	svc := NewDeviationAlertService(zap.NewNop(), 7)

	result := domain.ComparisonResult{
		Events: []domain.PlanDeviationDetected{
			{ModuleID: "m1", DeltaPct: 0.1}, // ~3 days — below 7-day threshold
		},
	}
	alerts := svc.Evaluate(context.Background(), result)
	if len(alerts) != 0 {
		t.Fatalf("expected no alerts for small deviation, got %+v", alerts)
	}
}

func TestDeviationAlertThresholdExceeded(t *testing.T) {
	svc := NewDeviationAlertService(zap.NewNop(), 7)

	result := domain.ComparisonResult{
		Events: []domain.PlanDeviationDetected{
			{ModuleID: "m1", DeltaPct: 0.5}, // ~15 days — above 7-day threshold
		},
	}
	alerts := svc.Evaluate(context.Background(), result)
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if !alerts[0].ShouldAlert {
		t.Fatal("expected ShouldAlert=true")
	}
	if alerts[0].ModuleID != "m1" {
		t.Fatalf("module=%s, want m1", alerts[0].ModuleID)
	}
}

func TestDeviationAlertDefaultThreshold(t *testing.T) {
	svc := NewDeviationAlertService(zap.NewNop(), 0) // default 7
	result := domain.ComparisonResult{
		Events: []domain.PlanDeviationDetected{
			{ModuleID: "m1", DeltaPct: 0.4}, // ~12 days — above default 7-day threshold
		},
	}
	alerts := svc.Evaluate(context.Background(), result)
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert with default threshold, got %d", len(alerts))
	}
}

func TestDeviationAlertThresholdUpdate(t *testing.T) {
	svc := NewDeviationAlertService(zap.NewNop(), 7)
	svc.SetThreshold(20)

	result := domain.ComparisonResult{
		Events: []domain.PlanDeviationDetected{
			{ModuleID: "m1", DeltaPct: 0.5}, // ~15 days — below updated 20-day threshold
		},
	}
	alerts := svc.Evaluate(context.Background(), result)
	if len(alerts) != 0 {
		t.Fatalf("expected no alerts after threshold update, got %+v", alerts)
	}
}
