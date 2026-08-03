package executionprogress

import (
	"context"

	"go.uber.org/zap"

	domain "vedo-edutrack/backend/internal/modules/executionprogress/domain"
)

// DeviationAlertService evaluates plan deviations and decides whether to publish alerts.
type DeviationAlertService struct {
	logger    *zap.Logger
	threshold int // N days threshold before alerting
}

// NewDeviationAlertService creates a deviation alert service.
// thresholdDays is the number of days of deviation before an alert is published.
func NewDeviationAlertService(logger *zap.Logger, thresholdDays int) *DeviationAlertService {
	if logger == nil {
		logger = zap.NewNop()
	}
	if thresholdDays <= 0 {
		thresholdDays = 7 // default: alert after 7 days of deviation
	}
	return &DeviationAlertService{
		logger:    logger.Named("execution.alert"),
		threshold: thresholdDays,
	}
}

// AlertResult is the outcome of evaluating a deviation.
type AlertResult struct {
	ShouldAlert bool   `json:"should_alert"`
	ModuleID    string `json:"module_id"`
	DeltaDays   int    `json:"delta_days"`
	Threshold   int    `json:"threshold"`
	Reason      string `json:"reason,omitempty"`
}

// Evaluate checks whether plan-vs-actual deviations exceed the alert threshold.
func (s *DeviationAlertService) Evaluate(
	_ context.Context,
	result domain.ComparisonResult,
) []AlertResult {
	var alerts []AlertResult
	for _, event := range result.Events {
		deltaDays := int(event.DeltaPct * 30) // approximate: pct of planned timeline → days
		if deltaDays > s.threshold {
			alert := AlertResult{
				ShouldAlert: true,
				ModuleID:    event.ModuleID,
				DeltaDays:   deltaDays,
				Threshold:   s.threshold,
				Reason:      "deviation exceeds configured threshold",
			}
			alerts = append(alerts, alert)
			s.logger.Warn("deviation alert triggered",
				zap.String("module", event.ModuleID),
				zap.Float64("delta_pct", event.DeltaPct),
				zap.Int("delta_days", deltaDays),
				zap.Int("threshold", s.threshold),
			)
		}
	}
	if len(alerts) == 0 {
		s.logger.Info("deviation evaluation complete, no alerts")
	}
	return alerts
}

// SetThreshold updates the alert threshold dynamically (without restart).
func (s *DeviationAlertService) SetThreshold(days int) {
	if days > 0 {
		s.threshold = days
		s.logger.Info("alert threshold updated", zap.Int("new_threshold", days))
	}
}
