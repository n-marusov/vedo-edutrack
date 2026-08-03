package resources

import "fmt"

// ResourceFilter restricts catalog queries.
type ResourceFilter struct {
	Type       ResourceType
	Format     string
	Source     string
	Difficulty string
	Limit      int
	Offset     int
}

// ResourcePage is a paginated filter result.
type ResourcePage struct {
	Items  []Resource
	Total  int
	Limit  int
	Offset int
}

// Catalog is an in-memory domain catalog with module bindings.
type Catalog struct {
	resources map[string]Resource
	bindings  []ResourceBinding
}

// NewCatalog creates a catalog and validates resources.
func NewCatalog(resources []Resource) (*Catalog, error) {
	catalog := &Catalog{resources: map[string]Resource{}}
	for _, resource := range resources {
		if err := resource.Validate(); err != nil {
			return nil, err
		}
		catalog.resources[resource.ID] = resource
	}
	return catalog, nil
}

// Filter returns resources matching all provided filters.
func (c *Catalog) Filter(filter ResourceFilter) ResourcePage {
	matches := make([]Resource, 0, len(c.resources))
	for _, resource := range c.resources {
		if filter.matches(resource) {
			matches = append(matches, resource)
		}
	}
	return paginate(matches, filter.Limit, filter.Offset)
}

// matches reports whether a resource satisfies all set filter fields.
func (filter ResourceFilter) matches(resource Resource) bool {
	if filter.Type != "" && resource.Type != filter.Type {
		return false
	}
	if filter.Format != "" && resource.Format != filter.Format {
		return false
	}
	if filter.Source != "" && resource.Source != filter.Source {
		return false
	}
	if filter.Difficulty != "" && resource.Difficulty != filter.Difficulty {
		return false
	}
	return true
}

// paginate applies offset/limit bounds to a result slice.
func paginate(items []Resource, limit, offset int) ResourcePage {
	total := len(items)
	if offset < 0 {
		offset = 0
	}
	if offset > total {
		offset = total
	}
	if limit <= 0 || limit > total-offset {
		limit = total - offset
	}
	return ResourcePage{Items: items[offset : offset+limit], Total: total, Limit: limit, Offset: offset}
}

// BindToModule associates a resource with a module and rejects orphan resources.
func (c *Catalog) BindToModule(binding ResourceBinding) error {
	if err := binding.Validate(); err != nil {
		return err
	}
	if _, ok := c.resources[binding.ResourceID]; !ok {
		return fmt.Errorf("resource %q does not exist", binding.ResourceID)
	}
	for _, existing := range c.bindings {
		if existing == binding {
			return nil
		}
	}
	c.bindings = append(c.bindings, binding)
	return nil
}

// ResourcesForModule returns resources bound to a module.
func (c *Catalog) ResourcesForModule(moduleID string) []Resource {
	out := []Resource{}
	for _, binding := range c.bindings {
		if binding.ModuleID == moduleID {
			out = append(out, c.resources[binding.ResourceID])
		}
	}
	return out
}
