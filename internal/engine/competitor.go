// Package engine contains the core simulation mechanics and state mutation logic.
package engine

import (
	"fmt"
	"math"
	mathrand "math/rand"

	"galaxies-burnrate/internal/config"
	"galaxies-burnrate/internal/game"
	"galaxies-burnrate/internal/npc"
)

// GenerateInitialCompetitors populates the galaxy with independent AI ships based on the configured population density.
func GenerateInitialCompetitors(state *game.State, settings *config.GameSettings) {
	if len(cache.PlanetIDs) == 0 {
		return
	}

	min := settings.Competitors.MinActive
	max := settings.Competitors.MaxActive

	count := min
	if max > min {
		count = mathrand.Intn(max-min+1) + min
	}

	for i := 0; i < count; i++ {
		startID := cache.PlanetIDs[mathrand.Intn(len(cache.PlanetIDs))]
		newComp := npc.NewCompetitor(startID)
		state.Competitors = append(state.Competitors, &newComp)
	}
}

// ProcessCompetitors iterates through all rival ships, resolving their arrivals and assigning them new contracts.
func ProcessCompetitors(state *game.State, settings *config.GameSettings) []string {
	var logs []string

	for _, c := range state.Competitors {
		switch c.Status {
		case "In Transit":
			if state.CurrentDay >= c.ArrivalDay {
				destName := state.Planets[c.DestinationID].Name
				commName := "Unknown Cargo"
				if c.ActiveContract != nil {
					if comm, exists := settings.Commodities[c.ActiveContract.CommodityID]; exists {
						commName = comm.Name
					}
				}

				logs = append(logs, fmt.Sprintf("    - [RIVAL] %s arrived at %s and delivered %s.", c.Name, destName, commName))

				c.LocationID = c.DestinationID
				c.DestinationID = ""
				c.Status = "Idle"
				c.ActiveContract = nil
			}
		case "Idle":
			contracts := state.Contracts[c.LocationID]
			if len(contracts) > 0 {
				idx := mathrand.Intn(len(contracts))
				selected := contracts[idx]
				state.Contracts[c.LocationID] = append(contracts[:idx], contracts[idx+1:]...)

				c.ActiveContract = &selected
				c.DestinationID = selected.DestinationID

				p1 := state.Planets[c.LocationID]
				p2 := state.Planets[c.DestinationID]
				dx := float64(p2.X - p1.X)
				dy := float64(p2.Y - p1.Y)
				distance := int(math.Round(math.Sqrt(dx*dx + dy*dy)))
				if distance < 1 {
					distance = 1
				}

				timeCost := int(math.Ceil(float64(distance) / 10.0))
				if timeCost < 1 {
					timeCost = 1
				}

				c.ArrivalDay = state.CurrentDay + timeCost
				c.Status = "In Transit"

				originName := p1.Name
				destName := p2.Name
				commName := "Unknown Cargo"
				if comm, exists := settings.Commodities[selected.CommodityID]; exists {
					commName = comm.Name
				}

				logs = append(logs, fmt.Sprintf("    - [RIVAL] %s accepted a contract for %s and departed %s for %s.", c.Name, commName, originName, destName))
			}
		}
	}

	return logs
}
