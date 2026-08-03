package gapcoverage

import (
	"go.uber.org/zap"

	domain "vedo-edutrack/backend/internal/modules/gapcoverage/domain"
)

// CoverageService exposes FGOS coverage and attestation report use cases.
type CoverageService struct {
	logger *zap.Logger
}

// NewCoverageService creates a coverage service.
func NewCoverageService(logger *zap.Logger) *CoverageService {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &CoverageService{logger: logger.Named("coverage")}
}

// Coverage computes live FGOS coverage and deficits.
func (s *CoverageService) Coverage(bindings []domain.FgosBinding, mastery domain.Mastery) domain.CoverageReport {
	report := domain.ComputeCoverage(bindings, mastery)
	s.logger.Info("FGOS coverage", zap.Float64("percentage", report.Percent), zap.Int("covered", report.Covered), zap.Int("total", report.Total))
	return report
}

// Attestation builds an attestation readiness report.
func (s *CoverageService) Attestation(bindings []domain.FgosBinding, mastery domain.Mastery, thresholdPercent float64) domain.AttestationReport {
	report := domain.BuildAttestationReport(s.Coverage(bindings, mastery), thresholdPercent)
	s.logger.Info("attestation readiness", zap.Float64("percentage", report.Coverage.Percent), zap.Bool("ready", report.Ready))
	return report
}
