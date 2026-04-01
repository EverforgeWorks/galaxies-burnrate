// Package engine contains the core simulation mechanics and state mutation logic.
package engine

import (
	"errors"
	"math"

	"galaxies-burnrate/internal/config"
	"galaxies-burnrate/internal/planet"
	"galaxies-burnrate/internal/player"
	"galaxies-burnrate/internal/ship"
)

// RefuelInvoice holds the calculated cost and fuel amount for a proposed refueling operation.
type RefuelInvoice struct {
	FuelToAdd int
	TotalCost int64
}

// CalculateRefuel determines the fuel required to fill the tank and the associated credit cost, factoring in global multipliers and local economic traits.
func CalculateRefuel(s *ship.Ship, p planet.Planet, credits int64, settings *config.GameSettings) (RefuelInvoice, error) {
	if s.Status == "In Transit" {
		return RefuelInvoice{}, errors.New("cannot refuel a ship while it is in transit")
	}

	if !p.HasTag("refuel") {
		return RefuelInvoice{}, errors.New("this location does not possess refueling facilities")
	}

	stats := s.CalculateStats()
	missingFuel := stats.FuelCapacity - s.CurrentFuel

	if missingFuel <= 0 {
		return RefuelInvoice{}, errors.New("fuel tanks are already at maximum capacity")
	}

	ecoStatus := settings.EconomicStatuses[p.EconomicStatus]
	adjustedUnitCost := float64(p.Refuel.Cost) * settings.Multipliers.FuelCost * ecoStatus.FuelCostMultiplier
	totalCost := int64(math.Round(float64(missingFuel) * adjustedUnitCost))

	// Process partial fill logic if corporate funds are insufficient for a full tank.
	if credits < totalCost {
		affordableFuel := int(math.Floor(float64(credits) / adjustedUnitCost))
		if affordableFuel <= 0 {
			return RefuelInvoice{}, errors.New("insufficient corporate credits to purchase any fuel")
		}
		return RefuelInvoice{
			FuelToAdd: affordableFuel,
			TotalCost: int64(math.Round(float64(affordableFuel) * adjustedUnitCost)),
		}, nil
	}

	return RefuelInvoice{
		FuelToAdd: missingFuel,
		TotalCost: totalCost,
	}, nil
}

// ExecuteRefuel processes the financial transaction and physically adds fuel to the vessel.
func ExecuteRefuel(s *ship.Ship, player *player.Player, p planet.Planet, settings *config.GameSettings) error {
	invoice, err := CalculateRefuel(s, p, player.Credits, settings)
	if err != nil {
		return err
	}

	player.Credits -= invoice.TotalCost
	s.CurrentFuel += invoice.FuelToAdd

	return nil
}
