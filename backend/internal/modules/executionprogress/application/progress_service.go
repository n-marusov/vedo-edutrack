package executionprogress

import (
	"go.uber.org/zap"

	domain "vedo-edutrack/backend/internal/modules/executionprogress/domain"
)

// ProgressService exposes plan-vs-actual and forecast use cases.
type ProgressService struct {
	logger *zap.Logger
}

// NewProgressService creates an execution-progress application service.
func NewProgressService(logger *zap.Logger) *ProgressService {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &ProgressService{logger: logger.Named("execution.progress")}
}

// Compare compares fixed plan modules to actual progress.
func (s *ProgressService) Compare(plan domain.FixedPlan, progress []domain.ModuleProgress) domain.ComparisonResult {
	result := domain.ComparePlanActual(plan, progress)
	for _, event := range result.Events {
		s.logger.Info("deviation detected", zap.String("module", event.ModuleID), zap.Float64("delta", event.DeltaPct))
	}
	return result
}

// Forecast computes binary readiness.
func (s *ProgressService) Forecast(input domain.ForecastInput) domain.ForecastResult {
	result := domain.ForecastReadiness(input)
	if result.Status == domain.ForecastNotOnTrack {
		s.logger.Warn("at-risk forecast", zap.Int("remaining", input.RemainingModules), zap.Float64("velocity", result.Velocity))
	}
	return result
}
