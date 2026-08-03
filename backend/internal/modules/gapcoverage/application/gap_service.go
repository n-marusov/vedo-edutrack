package gapcoverage

import (
	"go.uber.org/zap"

	domain "vedo-edutrack/backend/internal/modules/gapcoverage/domain"
)

// GapService exposes root-cause diagnosis use cases.
type GapService struct {
	logger *zap.Logger
}

// NewGapService creates a gap diagnosis service.
func NewGapService(logger *zap.Logger) *GapService {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &GapService{logger: logger.Named("gap.diagnose")}
}

// Diagnose finds ranked root causes for a learner lag module.
func (s *GapService) Diagnose(graph domain.Graph, mastery domain.Mastery, lagModuleID string) domain.DiagnosisResult {
	s.logger.Info("diagnosis started", zap.String("lag_module", lagModuleID))
	result := domain.DiagnoseRootCause(graph, mastery, lagModuleID)
	s.logger.Info("root causes found", zap.Int("count", len(result.RootCauses)))
	return result
}
