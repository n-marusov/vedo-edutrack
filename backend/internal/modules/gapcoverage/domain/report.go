package gapcoverage

// AttestationReport combines coverage and forecast-like readiness status.
type AttestationReport struct {
	Coverage CoverageReport
	Ready    bool
}

// BuildAttestationReport marks a learner ready when coverage meets the threshold.
func BuildAttestationReport(coverage CoverageReport, thresholdPercent float64) AttestationReport {
	if thresholdPercent <= 0 {
		thresholdPercent = 80
	}
	return AttestationReport{Coverage: coverage, Ready: coverage.Percent >= thresholdPercent}
}
