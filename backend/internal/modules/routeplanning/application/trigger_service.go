package routeplanning

import (
	"time"

	"go.uber.org/zap"
)

// RecomputeReason enumerates route recomputation triggers.
type RecomputeReason string

const (
	RecomputeReasonProgressChange RecomputeReason = "progressChange"
	RecomputeReasonGoalChange     RecomputeReason = "goalChange"
	RecomputeReasonOntologyUpdate RecomputeReason = "ontologyUpdate"
)

// RecomputeInput is the current route state to evaluate for recalculation.
type RecomputeInput struct {
	LearnerID            string
	PreviousGoalID       string
	CurrentGoalID        string
	PreviousOntology     string
	CurrentOntology      string
	CompletedModuleRatio float64
	PreviousModuleRatio  float64
}

// RouteRecalculationNeeded is emitted when a route must be recomputed.
type RouteRecalculationNeeded struct {
	LearnerID string
	Reason    RecomputeReason
	CreatedAt time.Time
}

// RecomputeTrigger detects route recomputation conditions.
type RecomputeTrigger struct {
	progressThreshold float64
	logger            *zap.Logger
}

// NewRecomputeTrigger creates a trigger detector. Default progress threshold is 15%.
func NewRecomputeTrigger(logger *zap.Logger) *RecomputeTrigger {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &RecomputeTrigger{progressThreshold: 0.15, logger: logger.Named("route.compute")}
}

// Detect returns a recalculation event when progress, goal, or ontology version changed enough.
func (t *RecomputeTrigger) Detect(input RecomputeInput) (RouteRecalculationNeeded, bool) {
	reason, ok := t.reason(input)
	if !ok {
		return RouteRecalculationNeeded{}, false
	}
	event := RouteRecalculationNeeded{LearnerID: input.LearnerID, Reason: reason, CreatedAt: time.Now().UTC()}
	t.logger.Info("recompute triggered", zap.String("learner", input.LearnerID), zap.String("reason", string(reason)))
	return event, true
}

func (t *RecomputeTrigger) reason(input RecomputeInput) (RecomputeReason, bool) {
	if input.PreviousGoalID != "" && input.CurrentGoalID != "" && input.PreviousGoalID != input.CurrentGoalID {
		return RecomputeReasonGoalChange, true
	}
	if input.PreviousOntology != "" && input.CurrentOntology != "" && input.PreviousOntology != input.CurrentOntology {
		return RecomputeReasonOntologyUpdate, true
	}
	progressDelta := input.CompletedModuleRatio - input.PreviousModuleRatio
	if progressDelta < 0 {
		progressDelta = -progressDelta
	}
	if progressDelta > t.progressThreshold {
		return RecomputeReasonProgressChange, true
	}
	return "", false
}
