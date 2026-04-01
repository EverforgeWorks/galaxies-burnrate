// Package engine contains the core simulation mechanics and state mutation logic.
package engine

import (
	"fmt"
	"math"

	"galaxies-burnrate/internal/config"
	"galaxies-burnrate/internal/game"
	"galaxies-burnrate/internal/market"
	"galaxies-burnrate/internal/passenger"
	"galaxies-burnrate/internal/ship"
)

// processArrival handles the spatial transition of a vessel arriving at its destination, including port authority docking fees.
func processArrival(state *game.State, s *ship.Ship, settings *config.GameSettings) []string {
	var logs []string

	s.LocationID = s.DestinationID
	s.DestinationID = ""
	s.ArrivalDay = 0
	s.Status = "Idle"

	destPlanet := state.Planets[s.LocationID]
	polStatus := settings.PoliticalStatuses[destPlanet.PoliticalStatus]

	// Apply the trait-modified planetary docking fee.
	dockingFee := int64(math.Round(float64(destPlanet.DockingFee) * polStatus.DockingFeeMultiplier))
	if dockingFee > 0 {
		state.Player.Credits -= dockingFee
		logs = append(logs, fmt.Sprintf(">>> %s has arrived at %s. Port authority deducted a %d C docking fee.", s.Name, destPlanet.Name, dockingFee))
	} else {
		logs = append(logs, fmt.Sprintf(">>> %s has arrived at %s.", s.Name, destPlanet.Name))
	}

	return logs
}

// processContracts iterates through a newly arrived vessel's cargo hold and resolves any freight deliveries destined for the current port.
func processContracts(state *game.State, s *ship.Ship, settings *config.GameSettings) []string {
	var logs []string
	var remainingContracts []market.Contract

	for _, c := range s.Contracts {
		if c.DestinationID == s.LocationID {
			state.Player.Credits += c.Payout
			commodity := settings.Commodities[c.CommodityID]

			if c.Payout >= 0 {
				logs = append(logs, fmt.Sprintf("    - Contract Fulfilled: Delivered %d units of %s. Payout: %d C.", c.Quantity, commodity.Name, c.Payout))
			} else {
				logs = append(logs, fmt.Sprintf("    - Contract Fulfilled: Delivered %d units of %s. Late penalties incurred: %d C.", c.Quantity, commodity.Name, c.Payout))
			}

			// If severe late penalties push the corporate account into the negative, trigger fleet-wide morale degradation.
			if state.Player.Credits < 0 {
				logs = append(logs, "    - [WARNING] Contract penalties have overdrawn the corporate account!")
				logs = append(logs, applyBankruptcyPenalties(state)...)
			}

		} else {
			remainingContracts = append(remainingContracts, c)
		}
	}

	s.Contracts = remainingContracts
	return logs
}

// processPassengers handles the disembarkation and ticket payouts for commuters reaching their destination.
func processPassengers(state *game.State, s *ship.Ship) []string {
	var logs []string
	var remainingPassengers []passenger.Passenger

	for _, p := range s.Passengers {
		if p.DestinationID == s.LocationID {
			state.Player.Credits += p.Payout
			logs = append(logs, fmt.Sprintf("    - Passenger Disembarked: %s [%s]. Payout: %d C.", p.Name, p.Class.Name, p.Payout))
		} else {
			remainingPassengers = append(remainingPassengers, p)
		}
	}

	s.Passengers = remainingPassengers
	return logs
}

// processSpoilageAndLateFees checks active freight and passengers for missed deadlines.
func processSpoilageAndLateFees(state *game.State, s *ship.Ship, settings *config.GameSettings) []string {
	var logs []string
	var activeContracts []market.Contract

	// 1. Process Cargo Spoilage & Late Fees
	for _, c := range s.Contracts {
		if state.CurrentDay > c.ExpirationDay {
			absPayout := float64(c.Payout)
			if absPayout < 0 {
				absPayout = -absPayout
			}

			// Calculate and apply daily compounding late fees.
			penalty := int64(math.Round(absPayout * (float64(settings.CargoRisks.LateFeePercentPerDay) / 100.0)))
			if penalty < 1 {
				penalty = 1
			}
			c.Payout -= penalty
			logs = append(logs, fmt.Sprintf("    - [OVERDUE] %s contract overdue. %d C late fee applied.", settings.Commodities[c.CommodityID].Name, penalty))

			// Check for unprotected perishable cargo rotting in the hold.
			comm := settings.Commodities[c.CommodityID]
			if comm.HasTag("perishable") && s.HasUnprotectedCargo("perishable", settings.Commodities) {
				loss := int(math.Ceil(float64(c.Quantity) * 0.20))
				if loss < 1 {
					loss = 1
				}
				c.Quantity -= loss
				logs = append(logs, fmt.Sprintf("    - [WARNING] Unrefrigerated %s is spoiling. Lost %d units.", comm.Name, loss))
			}
		}

		if c.Quantity > 0 {
			activeContracts = append(activeContracts, c)
		} else {
			comm := settings.Commodities[c.CommodityID]
			logs = append(logs, fmt.Sprintf("    - [CONTRACT FAILED] Shipment of %s has completely rotted away.", comm.Name))
		}
	}
	s.Contracts = activeContracts

	// 2. Process Passenger Delays
	for i := range s.Passengers {
		p := &s.Passengers[i]
		if state.CurrentDay > p.ExpirationDay {
			absPayout := float64(p.Payout)
			if absPayout < 0 {
				absPayout = -absPayout
			}

			penaltyRate := float64(settings.CargoRisks.LateFeePercentPerDay)
			if p.Class.HasTag("vip") {
				penaltyRate = float64(settings.Multipliers.VipDelayPenaltyPercent)
			}

			penalty := int64(math.Round(absPayout * (penaltyRate / 100.0)))
			if penalty < 1 {
				penalty = 1
			}
			p.Payout -= penalty
			logs = append(logs, fmt.Sprintf("    - [COMPLAINT] %s travel delayed. Compensated %d C for itinerary failure.", p.Class.Name, penalty))
		}
	}

	return logs
}
