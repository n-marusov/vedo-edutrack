package routeplanning

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	domain "vedo-edutrack/backend/internal/modules/routeplanning/domain"
)

// OntologyGraphProvider supplies route-planning graph snapshots for a requested ontology version.
type OntologyGraphProvider interface {
	GraphForVersion(ctx context.Context, ontologyVersion string) (domain.OntologyGraph, error)
}

// ComputeRequest is the application-layer route computation request.
type ComputeRequest struct {
	LearnerID         string
	PositionID        string
	GoalID            string
	PedagogyConceptID string
	OntologyVersion   string
	MidHorizonModuleN int
}

// ComputeResult contains a computed route plus version metadata.
type ComputeResult struct {
	LearnerID       string
	OntologyVersion string
	Route           domain.Route
	ComputedAt      time.Time
}

// ComputeService orchestrates ontology graph loading and pure-domain pathfinding.
type ComputeService struct {
	graphs     OntologyGraphProvider
	pathfinder *domain.Pathfinder
	logger     *zap.Logger
}

// NewComputeService creates a route computation application service.
func NewComputeService(graphs OntologyGraphProvider, pathfinder *domain.Pathfinder, logger *zap.Logger) *ComputeService {
	if pathfinder == nil {
		pathfinder = domain.NewPathfinder(domain.DefaultWeightProfile())
	}
	if logger == nil {
		logger = zap.NewNop()
	}
	return &ComputeService{graphs: graphs, pathfinder: pathfinder, logger: logger.Named("route.compute")}
}

// Compute computes a route for a learner and goal using the requested ontology snapshot.
func (s *ComputeService) Compute(ctx context.Context, req ComputeRequest) (ComputeResult, error) {
	if req.LearnerID == "" {
		return ComputeResult{}, fmt.Errorf("learner id is required")
	}
	if s.graphs == nil {
		return ComputeResult{}, fmt.Errorf("ontology graph provider is required")
	}
	s.logger.Info("computation started", zap.String("learner", req.LearnerID), zap.String("goal", req.GoalID), zap.String("ontology_version", req.OntologyVersion))
	started := time.Now()

	graph, err := s.graphs.GraphForVersion(ctx, req.OntologyVersion)
	if err != nil {
		s.logger.Error("computation failed", zap.String("reason", "load ontology graph"), zap.Error(err))
		return ComputeResult{}, fmt.Errorf("load ontology graph: %w", err)
	}
	route, err := s.pathfinder.Compute(graph, domain.ComputeRequest{PositionID: req.PositionID, GoalID: req.GoalID, PedagogyConceptID: req.PedagogyConceptID, MidHorizonModuleN: req.MidHorizonModuleN})
	if err != nil {
		if _, ok := err.(*domain.UnreachableGoalError); ok {
			s.logger.Warn("path not found", zap.String("goal", req.GoalID), zap.String("position", req.PositionID), zap.Error(err))
		} else {
			s.logger.Error("computation failed", zap.String("reason", "pathfinding"), zap.Error(err))
		}
		return ComputeResult{}, err
	}
	s.logger.Info("path found", zap.Int("moduleCount", len(route.Steps)), zap.Duration("duration", time.Since(started)))
	return ComputeResult{LearnerID: req.LearnerID, OntologyVersion: req.OntologyVersion, Route: route, ComputedAt: time.Now().UTC()}, nil
}
