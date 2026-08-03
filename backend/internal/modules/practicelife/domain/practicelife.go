// Package practicelife provides the domain model for the practicelife bounded context.
//
// Stories, project ideas, and quality maps are ontology-sourced content
// delivered to learners at the moment of module mastery via appliesTo/enriches
// graph links. The practicelife context is a read-only supporting context
// that caches projections from the ontology-port upstream.
package practicelife

// Story is a short contextual material (≤ 5 min reading) linked to 1–3 topics
// explaining why the knowledge matters in the real world.
type Story struct {
	ID             string
	Title          string
	Text           string   // full story content
	LinkedModules  []string // 1–3 module IDs from ontology
	Subjects       []string // subjects involved (for cross-subject display)
	RealWorld      string   // mandatory real-world application section
	ReadingMinutes int      // estimated reading time (≤ 5)
}

// Difficulty is the graded complexity level for project ideas.
type Difficulty string

const (
	DifficultyBasic    Difficulty = "basic"
	DifficultyMedium   Difficulty = "medium"
	DifficultyAdvanced Difficulty = "advanced"
)

// ProjectIdea is a cross-subject project requiring modules from ≥ 2 subjects.
type ProjectIdea struct {
	ID              string
	Title           string
	Modules         []string // required module IDs (≥ 2 subjects)
	DifficultyLevel Difficulty
	ExpectedOutcome string // what the learner will produce or demonstrate
}

// StoryCatalog is a read-only projection of stories cached from the ontology.
type StoryCatalog struct {
	stories       []Story
	moduleToIndex map[string][]int // module ID → story indices
}

// NewStoryCatalog builds a story catalog and indexes by module links.
func NewStoryCatalog(stories []Story) *StoryCatalog {
	c := &StoryCatalog{
		stories:       stories,
		moduleToIndex: make(map[string][]int),
	}
	for i, s := range stories {
		for _, modID := range s.LinkedModules {
			c.moduleToIndex[modID] = append(c.moduleToIndex[modID], i)
		}
	}
	return c
}

// ByModule returns stories linked to a specific module.
func (c *StoryCatalog) ByModule(moduleID string) []Story {
	indices := c.moduleToIndex[moduleID]
	if len(indices) == 0 {
		return nil
	}
	result := make([]Story, len(indices))
	for i, idx := range indices {
		result[i] = c.stories[idx]
	}
	return result
}

// Count returns the total number of stories in the catalog.
func (c *StoryCatalog) Count() int {
	return len(c.stories)
}

// ProjectIdeaCatalog is a read-only projection of project ideas from the ontology.
type ProjectIdeaCatalog struct {
	ideas         []ProjectIdea
	moduleToIndex map[string][]int // module ID → idea indices
}

// NewProjectIdeaCatalog builds a project idea catalog and indexes by module links.
func NewProjectIdeaCatalog(ideas []ProjectIdea) *ProjectIdeaCatalog {
	pc := &ProjectIdeaCatalog{
		ideas:         ideas,
		moduleToIndex: make(map[string][]int),
	}
	for i, idea := range ideas {
		for _, modID := range idea.Modules {
			pc.moduleToIndex[modID] = append(pc.moduleToIndex[modID], i)
		}
	}
	return pc
}

// ByModule returns project ideas linked to a specific module.
func (pc *ProjectIdeaCatalog) ByModule(moduleID string) []ProjectIdea {
	indices := pc.moduleToIndex[moduleID]
	if len(indices) == 0 {
		return nil
	}
	result := make([]ProjectIdea, len(indices))
	for i, idx := range indices {
		result[i] = pc.ideas[idx]
	}
	return result
}

// SuggestEligible returns project ideas where ≥ 80% of required modules
// are in the mastered or available set.
func (pc *ProjectIdeaCatalog) SuggestEligible(masteredOrAvailable map[string]bool) []ProjectIdea {
	var result []ProjectIdea
	for _, idea := range pc.ideas {
		if len(idea.Modules) == 0 {
			continue
		}
		ready := 0
		for _, modID := range idea.Modules {
			if masteredOrAvailable[modID] {
				ready++
			}
		}
		if float64(ready)/float64(len(idea.Modules)) >= 0.80 {
			result = append(result, idea)
		}
	}
	return result
}

// Count returns the total number of project ideas.
func (pc *ProjectIdeaCatalog) Count() int {
	return len(pc.ideas)
}
