package executionprogress

import (
	"context"

	"go.uber.org/zap"

	domain "vedo-edutrack/backend/internal/modules/executionprogress/domain"
)

// ForecastService computes binary readiness forecasts.
type ForecastService struct {
	progress ProgressRepository
	logger   *zap.Logger
}

// ProgressRepository abstracts progress data reads for forecasting.
type ProgressRepository interface {
	GetProgress(ctx context.Context, learnerID string) ([]domain.ModuleProgress, error)
	GetPlannedModuleCount(ctx context.Context, learnerID, planID string) (int, error)
	GetRemainingModules(ctx context.Context, learnerID, planID string) (int, error)
}

// NewForecastService creates a forecast service.
func NewForecastService(repo ProgressRepository, logger *zap.Logger) *ForecastService {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &ForecastService{
		progress: repo,
		logger:   logger.Named("execution.forecast"),
	}
}

// ForecastResult wraps the domain forecast with additional context.
type ForecastResult struct {
	domain.ForecastResult
	ExpectedDate     string   `json:"expected_date,omitempty"`
	RemainingModules int      `json:"remaining_modules"`
	KeyRisks         []string `json:"key_risks,omitempty"`
	DataConfidence   string   `json:"data_confidence"` // "high" or "low"
}

// ForecastReadiness computes a binary forecast for a learner's plan.
func (s *ForecastService) ForecastReadiness(
	ctx context.Context,
	learnerID string,
	planID string,
	daysElapsed int,
	daysRemaining int,
) (*ForecastResult, error) {
	progress, err := s.progress.GetProgress(ctx, learnerID)
	if err != nil {
		s.logger.Error("failed to get progress", zap.Error(err))
		return nil, err
	}

	completedModules := 0
	remainingModules := 0
	for _, p := range progress {
		if p.Status == domain.StatusMastered {
			completedModules++
		}
	}

	remainingModules, err = s.progress.GetRemainingModules(ctx, learnerID, planID)
	if err != nil {
		s.logger.Warn("could not get remaining modules, counting from progress",
			zap.Error(err),
		)
		for _, p := range progress {
			if p.Status != domain.StatusMastered && p.Status != domain.StatusSkipped {
				remainingModules++
			}
		}
	}

	input := domain.ForecastInput{
		CompletedModules: completedModules,
		RemainingModules: remainingModules,
		DaysElapsed:      daysElapsed,
		DaysRemaining:    daysRemaining,
	}

	result := domain.ForecastReadiness(input)

	confidence := "high"
	var risks []string
	if result.LowAccuracy {
		confidence = "low"
		risks = append(risks, "insufficient data for accurate forecast")
	}
	if result.Status == domain.ForecastNotOnTrack {
		risks = append(risks, "current pace insufficient to complete on time")
	}

	s.logger.Info("forecast computed",
		zap.String("learner_id", learnerID),
		zap.String("plan_id", planID),
		zap.String("status", string(result.Status)),
		zap.Float64("velocity", result.Velocity),
		zap.String("confidence", confidence),
	)

	return &ForecastResult{
		ForecastResult:   result,
		RemainingModules: remainingModules,
		KeyRisks:         risks,
		DataConfidence:   confidence,
	}, nil
}
