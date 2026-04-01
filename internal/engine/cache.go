// Package engine contains the core simulation mechanics and state mutation logic.
package engine

import (
	"galaxies-burnrate/internal/config"
	"galaxies-burnrate/internal/crew"
	"galaxies-burnrate/internal/game"
)

// LookupCache holds pre-allocated slices of map keys and values to eliminate GC thrashing during simulation loops.
type LookupCache struct {
	PlanetIDs         []string
	EncounterIDsByTag map[string][]string
	PassengerClassIDs []string
	CrewRoles         []crew.RoleTemplate
	Origins           []crew.Origin
	Specialties       []crew.Specialty
}

// cache provides global read-only access to static data slices.
var cache *LookupCache

// Initialize builds the static lookup tables required for high-performance engine loops.
func Initialize(state *game.State, settings *config.GameSettings) {
	cache = &LookupCache{
		EncounterIDsByTag: make(map[string][]string),
	}

	for id := range state.Planets {
		cache.PlanetIDs = append(cache.PlanetIDs, id)
	}

	for id, enc := range settings.Encounters {
		cache.EncounterIDsByTag[enc.RequiredTag] = append(cache.EncounterIDsByTag[enc.RequiredTag], id)
	}

	for id := range settings.PassengerClasses {
		cache.PassengerClassIDs = append(cache.PassengerClassIDs, id)
	}

	for _, r := range settings.CrewRoles {
		cache.CrewRoles = append(cache.CrewRoles, r)
	}

	for _, o := range settings.Origins {
		cache.Origins = append(cache.Origins, o)
	}

	for _, s := range settings.Specialties {
		cache.Specialties = append(cache.Specialties, s)
	}
}