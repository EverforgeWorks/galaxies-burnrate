// Package engine contains the core simulation mechanics and state mutation logic.
package engine

import (
	"fmt"

	"galaxies-burnrate/internal/config"
	"galaxies-burnrate/internal/game"
)

// AdvanceTime progresses the global simulation clock day-by-day, orchestrating the resolution of events, deadlines, payroll, and crew health.
func AdvanceTime(state *game.State, settings *config.GameSettings, days int) []string {
	var events []string

	// Execute the simulation loop for the requested number of days.
	for i := 0; i < days; i++ {
		state.CurrentDay++
		events = append(events, fmt.Sprintf("\n--- End of Day %d ---", state.CurrentDay))

		// Iterate through the player's fleet to process localized logistics and health mechanics.
		for _, s := range state.Player.Fleet {
			// Ignore vessels that have been catastrophically destroyed.
			if s.Status == "Destroyed" {
				continue
			}

			// 1. Biological and Medical Processing
			events = append(events, processDailyCrewHealth(s, settings)...)

			// 2. Transit and Arrival Logistics
			if s.Status == "In Transit" {
				if state.CurrentDay >= s.ArrivalDay {
					events = append(events, processArrival(state, s, settings)...)
					events = append(events, resolveCustomsInspection(state, s, settings)...)
					events = append(events, processContracts(state, s, settings)...)
					events = append(events, processPassengers(state, s)...)
				} else {
					events = append(events, resolveTransitEvent(state, s, settings)...)
				}
			}

			// 3. Active Cargo Degradation and Expiration
			events = append(events, processSpoilageAndLateFees(state, s, settings)...)
		}

		// 4. Global Corporate Processing (Payroll and Morale)
		if state.CurrentDay%settings.Crew.PayPeriodDays == 0 {
			events = append(events, processFleetPayroll(state, settings)...)
		}

		// 5. NPC AI Processing
		events = append(events, ProcessCompetitors(state, settings)...)
	}

	// 6. Market Board Refreshes
	RefreshContracts(state, settings)
	RefreshPassengers(state, settings)
	RefreshCrew(state, settings)

	return events
}
