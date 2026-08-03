package resources

import (
	"go.uber.org/zap"

	domain "vedo-edutrack/backend/internal/modules/resources/domain"
)

// CatalogQuery is a catalog read use case input.
type CatalogQuery struct {
	Filters  domain.ResourceFilter
	ModuleID string
}

// CatalogResult contains catalog items or module-bound resources.
type CatalogResult struct {
	Items []domain.Resource
	Total int
}

// CatalogService exposes resource catalog read use cases.
type CatalogService struct {
	catalog *domain.Catalog
	logger  *zap.Logger
}

// NewCatalogService creates a resource catalog application service.
func NewCatalogService(catalog *domain.Catalog, logger *zap.Logger) *CatalogService {
	if catalog == nil {
		catalog = &domain.Catalog{}
	}
	if logger == nil {
		logger = zap.NewNop()
	}
	return &CatalogService{catalog: catalog, logger: logger.Named("resources.catalog")}
}

// Query returns filtered resources or resources bound to a module.
func (s *CatalogService) Query(query CatalogQuery) CatalogResult {
	if query.ModuleID != "" {
		items := s.catalog.ResourcesForModule(query.ModuleID)
		return CatalogResult{Items: items, Total: len(items)}
	}
	page := s.catalog.Filter(query.Filters)
	s.logger.Info("query", zap.Any("filters", query.Filters), zap.Int("results", len(page.Items)))
	return CatalogResult{Items: page.Items, Total: page.Total}
}
