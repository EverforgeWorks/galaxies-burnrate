// Package engine contains the core simulation mechanics and state mutation logic.
package engine

import (
	"errors"
	"math"

	"galaxies-burnrate/internal/planet"
	"galaxies-burnrate/internal/ship"
)

// TravelPlan holds the calculated logistical requirements for a proposed journey.
type TravelPlan struct {
	Distance int
	FuelCost int
	TimeCost int
}

// CalculateTravel computes the distance, fuel, and time required for a ship to travel between two planets.
func CalculateTravel(s *ship.Ship, origin, destination planet.Planet) (TravelPlan, error) {
	stats := s.CalculateStats()

	// Prevent division by zero and enforce mechanical logic: ships without engines cannot move.
	if stats.Speed <= 0 {
		return TravelPlan{}, errors.New("ship has zero speed and cannot execute travel")
	}

	// Calculate Euclidean distance between the two grid coordinates.
	dx := float64(destination.X - origin.X)
	dy := float64(destination.Y - origin.Y)
	distance := int(math.Round(math.Sqrt(dx*dx + dy*dy)))

	fuelCost := distance * stats.FuelConsumption

	// Calculate time cost. Math.Ceil ensures that even a short jump takes at least 1 unit of time.
	timeCost := int(math.Ceil(float64(distance) / float64(stats.Speed)))

	return TravelPlan{
		Distance: distance,
		FuelCost: fuelCost,
		TimeCost: timeCost,
	}, nil
}

// ExecuteTravel consumes resources and dispatches the ship on an asynchronous journey.
func ExecuteTravel(s *ship.Ship, origin, destination planet.Planet, currentDay int) error {
	plan, err := CalculateTravel(s, origin, destination)
	if err != nil {
		return err
	}

	if s.CurrentFuel < plan.FuelCost {
		return errors.New("insufficient fuel for this journey")
	}

	// Deduct fuel immediately upon departure and set the arrival parameters.
	s.CurrentFuel -= plan.FuelCost
	s.DestinationID = destination.ID
	s.ArrivalDay = currentDay + plan.TimeCost
	s.Status = "In Transit"

	return nil
}
