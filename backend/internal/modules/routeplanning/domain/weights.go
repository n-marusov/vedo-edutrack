package routeplanning

// WeightProfile defines route edge costs. Lower cost means stronger preference.
type WeightProfile struct {
	StrictPrerequisite int
	SoftPrerequisite   int
	AppliesTo          int
}

// DefaultWeightProfile returns MVP weights from the M1 plan.
func DefaultWeightProfile() WeightProfile {
	return WeightProfile{StrictPrerequisite: 1, SoftPrerequisite: 5, AppliesTo: 20}
}

func (p WeightProfile) weight(linkType LinkType) (int, bool) {
	switch linkType {
	case LinkStrictPrerequisite:
		return p.StrictPrerequisite, true
	case LinkSoftPrerequisite:
		return p.SoftPrerequisite, true
	case LinkAppliesTo:
		return p.AppliesTo, true
	case LinkEnriches:
		return 0, false
	default:
		return 0, false
	}
}
