package executionprogress

// ForecastStatus is the binary MVP readiness forecast.
type ForecastStatus string

const (
	ForecastOnTrack    ForecastStatus = "on-track"
	ForecastNotOnTrack ForecastStatus = "not-on-track"
)

// ForecastInput contains current velocity and remaining workload.
type ForecastInput struct {
	CompletedModules int
	RemainingModules int
	DaysElapsed      int
	DaysRemaining    int
}

// ForecastResult is a binary readiness forecast with a low-accuracy label.
type ForecastResult struct {
	Status      ForecastStatus
	LowAccuracy bool
	Velocity    float64
}

// ForecastReadiness estimates whether the learner can complete remaining modules on time.
func ForecastReadiness(input ForecastInput) ForecastResult {
	if input.CompletedModules <= 0 || input.DaysElapsed <= 0 {
		return ForecastResult{Status: ForecastNotOnTrack, LowAccuracy: true}
	}
	velocity := float64(input.CompletedModules) / float64(input.DaysElapsed)
	required := float64(input.RemainingModules)
	if input.DaysRemaining > 0 {
		required = required / float64(input.DaysRemaining)
	}
	status := ForecastNotOnTrack
	if velocity >= required {
		status = ForecastOnTrack
	}
	return ForecastResult{Status: status, Velocity: velocity}
}
