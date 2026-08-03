package routeplanning

// Horizon is a named route slice for far/mid/near planning views.
type Horizon struct {
	Name      string
	ModuleIDs []string
}

func buildHorizons(moduleIDs []string, midLimit int) (Horizon, Horizon, Horizon) {
	if midLimit <= 0 {
		midLimit = 5
	}
	farIDs := append([]string(nil), moduleIDs...)
	midEnd := len(moduleIDs)
	if midEnd > midLimit {
		midEnd = midLimit
	}
	midIDs := append([]string(nil), moduleIDs[:midEnd]...)
	nearIDs := []string{}
	if len(moduleIDs) > 0 {
		nearIDs = []string{moduleIDs[0]}
	}
	return Horizon{Name: "far", ModuleIDs: farIDs}, Horizon{Name: "mid", ModuleIDs: midIDs}, Horizon{Name: "near", ModuleIDs: nearIDs}
}
