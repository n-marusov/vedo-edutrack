package executionprogress

// DeviationCause records why a learner deviated from the fixed plan.
type DeviationCause string

const (
	DeviationCauseUnknown  DeviationCause = "unknown"
	DeviationCauseSkipped  DeviationCause = "skipped"
	DeviationCauseExternal DeviationCause = "external"
)
