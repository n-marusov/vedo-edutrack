package ontologyport

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"go.uber.org/zap"

	ontology "vedo-edutrack/backend/internal/modules/ontologyport/domain"
)

// GraphClient is the read-only Hub boundary used by ontology sync.
type GraphClient interface {
	GraphNeighborhood(ctx context.Context, ontologyID, conceptID string, depth int) (ontology.Subgraph, error)
}

// SyncOptions configures a full or incremental subgraph sync.
type SyncOptions struct {
	OntologyID  string
	ConceptIDs  []string
	Depth       int
	Incremental bool
}

// SyncResult describes the cached ontology snapshot after sync.
type SyncResult struct {
	OntologyID     string
	Version        string
	ModuleCount    int
	LinkCount      int
	ChangedModules int
	Incremental    bool
	Duration       time.Duration
	CachedAt       time.Time
}

// VersionedCache stores the latest copied subgraph by ontology id.
type VersionedCache struct {
	mu     sync.RWMutex
	graphs map[string]ontology.Subgraph
	cached map[string]time.Time
}

// NewVersionedCache creates an in-memory cache for copied Hub subgraphs.
func NewVersionedCache() *VersionedCache {
	return &VersionedCache{graphs: map[string]ontology.Subgraph{}, cached: map[string]time.Time{}}
}

// Get returns a cached graph snapshot if present.
func (c *VersionedCache) Get(ontologyID string) (ontology.Subgraph, time.Time, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	graph, ok := c.graphs[ontologyID]
	return graph, c.cached[ontologyID], ok
}

// Put replaces the cached graph snapshot.
func (c *VersionedCache) Put(graph ontology.Subgraph) time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	cachedAt := time.Now().UTC()
	c.graphs[graph.OntologyID] = graph
	c.cached[graph.OntologyID] = cachedAt
	return cachedAt
}

var defaultCache = NewVersionedCache()

// DefaultCache returns the process-local ontology sync cache.
func DefaultCache() *VersionedCache {
	return defaultCache
}

// SyncService copies relevant ontology subgraphs from VEDO Hub into an in-memory cache.
type SyncService struct {
	client GraphClient
	cache  *VersionedCache
	logger *zap.Logger
}

// NewSyncService creates an ontology sync application service.
func NewSyncService(client GraphClient, cache *VersionedCache, logger *zap.Logger) *SyncService {
	if cache == nil {
		cache = DefaultCache()
	}
	if logger == nil {
		logger = zap.NewNop()
	}
	return &SyncService{client: client, cache: cache, logger: logger.Named("ontology.sync")}
}

// Sync copies a full or incremental ontology subgraph and stores a versioned snapshot.
func (s *SyncService) Sync(ctx context.Context, opts SyncOptions) (SyncResult, error) {
	if err := validateSyncOptions(s.client, opts); err != nil {
		return SyncResult{}, err
	}

	started := time.Now()
	s.logger.Info("sync started for ontology", zap.String("ontology_id", opts.OntologyID), zap.Bool("incremental", opts.Incremental), zap.Strings("concept_ids", normalizeConceptIDs(opts.ConceptIDs)))

	merged, err := s.fetchMerged(ctx, opts)
	if err != nil {
		return SyncResult{}, err
	}

	previous, _, hadPrevious := s.cache.Get(opts.OntologyID)
	if result, unchanged := s.incrementalShortCircuit(opts, previous, merged, hadPrevious, started); unchanged {
		return result, nil
	}

	cachedAt := s.cache.Put(merged)
	result := SyncResult{OntologyID: merged.OntologyID, Version: merged.Version, ModuleCount: len(merged.Modules), LinkCount: len(merged.Links), ChangedModules: mergedChangedCount(opts, previous, merged, hadPrevious), Incremental: opts.Incremental, Duration: time.Since(started), CachedAt: cachedAt}
	s.logger.Info("fetched modules and links", zap.Int("count", result.ModuleCount), zap.Int("linkCount", result.LinkCount), zap.Duration("duration", result.Duration))
	s.logger.Info("subgraph cached", zap.String("version", result.Version), zap.Time("cached_at", result.CachedAt))
	if opts.Incremental {
		s.logger.Warn("incremental sync", zap.Int("changedCount", result.ChangedModules))
	}
	return result, nil
}

// validateSyncOptions checks the sync inputs before any network work.
func validateSyncOptions(client GraphClient, opts SyncOptions) error {
	if client == nil {
		return fmt.Errorf("ontology sync requires a Hub graph client")
	}
	if opts.OntologyID == "" {
		return fmt.Errorf("ontology id is required")
	}
	if len(normalizeConceptIDs(opts.ConceptIDs)) == 0 {
		return fmt.Errorf("at least one concept or subject id is required")
	}
	return nil
}

// fetchMerged fetches each concept neighborhood and merges modules/links by key.
func (s *SyncService) fetchMerged(ctx context.Context, opts SyncOptions) (ontology.Subgraph, error) {
	depth := opts.Depth
	if depth <= 0 {
		depth = 1
	}
	merged := ontology.Subgraph{OntologyID: opts.OntologyID}
	modulesByID := map[string]ontology.OntologyModule{}
	linksByKey := map[string]ontology.OntologyLink{}
	for _, conceptID := range normalizeConceptIDs(opts.ConceptIDs) {
		subgraph, err := s.client.GraphNeighborhood(ctx, opts.OntologyID, conceptID, depth)
		if err != nil {
			s.logger.Error("sync failed", zap.String("ontology_id", opts.OntologyID), zap.String("concept_id", conceptID), zap.Error(err))
			return ontology.Subgraph{}, fmt.Errorf("fetch graph neighborhood for %s: %w", conceptID, err)
		}
		mergeSubgraph(&merged, subgraph, modulesByID, linksByKey)
	}
	merged.Modules = sortedModules(modulesByID)
	merged.Links = sortedLinks(linksByKey)
	return merged, nil
}

// mergeSubgraph folds a fetched subgraph into the shared maps.
func mergeSubgraph(merged *ontology.Subgraph, subgraph ontology.Subgraph, modulesByID map[string]ontology.OntologyModule, linksByKey map[string]ontology.OntologyLink) {
	if subgraph.OntologyID != "" {
		merged.OntologyID = subgraph.OntologyID
	}
	if subgraph.Version != "" {
		merged.Version = subgraph.Version
	}
	for _, module := range subgraph.Modules {
		modulesByID[module.ID] = module
	}
	for _, link := range subgraph.Links {
		linksByKey[linkKey(link)] = link
	}
}

// incrementalShortCircuit detects an incremental sync with no changes and
// returns the early result (true) so the caller can skip re-caching.
func (s *SyncService) incrementalShortCircuit(opts SyncOptions, previous, merged ontology.Subgraph, hadPrevious bool, started time.Time) (SyncResult, bool) {
	if !opts.Incremental || !hadPrevious {
		return SyncResult{}, false
	}
	if changedModuleCount(previous, merged) == 0 && previous.Version == merged.Version {
		cachedAt := s.cache.Put(previous)
		result := SyncResult{OntologyID: previous.OntologyID, Version: previous.Version, ModuleCount: len(previous.Modules), LinkCount: len(previous.Links), ChangedModules: 0, Incremental: true, Duration: time.Since(started), CachedAt: cachedAt}
		s.logger.Warn("incremental sync", zap.Int("changedCount", 0), zap.String("version", previous.Version))
		return result, true
	}
	return SyncResult{}, false
}

// mergedChangedCount computes the changed-module count for the result.
func mergedChangedCount(opts SyncOptions, previous, merged ontology.Subgraph, hadPrevious bool) int {
	if opts.Incremental && hadPrevious {
		return changedModuleCount(previous, merged)
	}
	return len(merged.Modules)
}

func normalizeConceptIDs(ids []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		if id == "" || seen[id] {
			continue
		}
		seen[id] = true
		out = append(out, id)
	}
	sort.Strings(out)
	return out
}

func linkKey(link ontology.OntologyLink) string {
	return string(link.Type) + "\x00" + link.SourceModuleID + "\x00" + link.TargetModuleID
}

func sortedModules(modules map[string]ontology.OntologyModule) []ontology.OntologyModule {
	ids := make([]string, 0, len(modules))
	for id := range modules {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	out := make([]ontology.OntologyModule, 0, len(ids))
	for _, id := range ids {
		out = append(out, modules[id])
	}
	return out
}

func sortedLinks(links map[string]ontology.OntologyLink) []ontology.OntologyLink {
	keys := make([]string, 0, len(links))
	for key := range links {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	out := make([]ontology.OntologyLink, 0, len(keys))
	for _, key := range keys {
		out = append(out, links[key])
	}
	return out
}

func changedModuleCount(previous, next ontology.Subgraph) int {
	prevVersion := map[string]string{}
	for _, module := range previous.Modules {
		prevVersion[module.ID] = module.Version
	}
	changed := 0
	for _, module := range next.Modules {
		if prevVersion[module.ID] != module.Version {
			changed++
		}
	}
	return changed
}
