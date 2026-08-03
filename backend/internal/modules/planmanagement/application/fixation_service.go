package planmanagement

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	routedomain "vedo-edutrack/backend/internal/modules/routeplanning/domain"
)

// LearningPlan is an immutable snapshot of a computed route with timeline and checkpoints.
type LearningPlan struct {
	ID              string
	LearnerID       string
	GoalID          string
	OntologyVersion string
	Modules         []PlanModule
	Checkpoints     []Checkpoint
	StartAt         time.Time
	EndAt           time.Time
	CreatedAt       time.Time
	Version         int
}

// PlanModule is one scheduled module in a fixed learning plan.
type PlanModule struct {
	ModuleID string
	Order    int
	StartAt  time.Time
	DueAt    time.Time
	LinkType string
}

// Checkpoint is a timeline control point for plan-vs-actual tracking.
type Checkpoint struct {
	ID        string
	ModuleID  string
	DueAt     time.Time
	Completed bool
}

// PlanRepository persists immutable learning plans. Implementations must be insert-only.
type PlanRepository interface {
	CreatePlan(ctx context.Context, plan LearningPlan) error
}

// FixationRequest asks the service to snapshot a computed route.
type FixationRequest struct {
	PlanID          string
	LearnerID       string
	GoalID          string
	OntologyVersion string
	Route           routedomain.Route
	StartAt         time.Time
	ModuleDuration  time.Duration
}

// PlanFixationService creates immutable plan snapshots from computed routes.
type PlanFixationService struct {
	repository PlanRepository
	logger     *zap.Logger
}

// NewPlanFixationService creates a plan fixation application service.
func NewPlanFixationService(repository PlanRepository, logger *zap.Logger) *PlanFixationService {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &PlanFixationService{repository: repository, logger: logger.Named("plan.fixation")}
}

// Fix creates and optionally persists an immutable learning plan snapshot.
func (s *PlanFixationService) Fix(ctx context.Context, req FixationRequest) (LearningPlan, error) {
	if req.PlanID == "" || req.LearnerID == "" || req.GoalID == "" {
		return LearningPlan{}, fmt.Errorf("plan id, learner id, and goal id are required")
	}
	if len(req.Route.Steps) == 0 {
		return LearningPlan{}, fmt.Errorf("route must contain at least one module")
	}
	startAt := req.StartAt
	if startAt.IsZero() {
		startAt = time.Now().UTC()
	}
	moduleDuration := req.ModuleDuration
	if moduleDuration <= 0 {
		moduleDuration = 24 * time.Hour
	}

	modules := make([]PlanModule, 0, len(req.Route.Steps))
	checkpoints := make([]Checkpoint, 0, len(req.Route.Steps))
	for _, step := range req.Route.Steps {
		moduleStart := startAt.Add(time.Duration(step.Order) * moduleDuration)
		moduleDue := moduleStart.Add(moduleDuration)
		modules = append(modules, PlanModule{ModuleID: step.ModuleID, Order: step.Order, StartAt: moduleStart, DueAt: moduleDue, LinkType: string(step.Via)})
		checkpoints = append(checkpoints, Checkpoint{ID: fmt.Sprintf("%s-%03d", req.PlanID, step.Order+1), ModuleID: step.ModuleID, DueAt: moduleDue})
	}
	plan := LearningPlan{ID: req.PlanID, LearnerID: req.LearnerID, GoalID: req.GoalID, OntologyVersion: req.OntologyVersion, Modules: modules, Checkpoints: checkpoints, StartAt: startAt, EndAt: modules[len(modules)-1].DueAt, CreatedAt: time.Now().UTC(), Version: 1}
	if s.repository != nil {
		if err := s.repository.CreatePlan(ctx, plan); err != nil {
			s.logger.Error("plan fixation failed", zap.String("planId", req.PlanID), zap.Error(err))
			return LearningPlan{}, fmt.Errorf("create immutable plan: %w", err)
		}
	}
	s.logger.Info("plan fixed", zap.String("planId", plan.ID), zap.Int("modules", len(plan.Modules)), zap.Time("timeline_start", plan.StartAt), zap.Time("timeline_end", plan.EndAt))
	return plan, nil
}
