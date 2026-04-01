// Package ship manages vessel architecture, module loadouts, and stat aggregation.
package ship

import (
	crand "crypto/rand"
	"fmt"
	"math/rand"
	"time"

	"galaxies-burnrate/internal/crew"
	"galaxies-burnrate/internal/market"
	"galaxies-burnrate/internal/passenger"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// ShipStats represents the aggregated capabilities of a vessel.
type ShipStats struct {
	CargoSpace             int
	FuelCapacity           int
	Speed                  int
	PassengerSpace         int
	CrewCapacity           int
	HazardousProtection    int
	FragileProtection      int
	IllicitConcealment     int
	PerishablePreservation int
	FuelConsumption        int
}

// Ship represents an active, physical vessel in the simulation.
type Ship struct {
	ID            string
	Name          string
	Chassis       ChassisTemplate
	Modules       []InstalledModule
	Roster        []crew.CrewMember
	Contracts     []market.Contract
	Passengers    []passenger.Passenger
	CurrentFuel   int
	LocationID    string
	DestinationID string
	ArrivalDay    int
	Status        string
}

// generateShipID creates a secure, unique hexadecimal identifier.
func generateShipID() string {
	b := make([]byte, 8)
	_, _ = crand.Read(b)
	return fmt.Sprintf("%x", b)
}

// New constructs a fresh Ship instance with a unique ID.
func New(name string, chassis ChassisTemplate, locationID string) *Ship {
	return &Ship{
		ID:         generateShipID(),
		Name:       name,
		Chassis:    chassis,
		Modules:    make([]InstalledModule, 0),
		LocationID: locationID,
		Status:     "Idle",
	}
}

// InstallModule adds a pristine component to the ship.
func (s *Ship) InstallModule(m ModuleTemplate) error {
	s.Modules = append(s.Modules, InstalledModule{
		Template:  m,
		IsDamaged: false,
	})
	return nil
}

// UninstallModule removes a component at the specified index.
func (s *Ship) UninstallModule(index int) error {
	if index < 0 || index >= len(s.Modules) {
		return fmt.Errorf("module index out of bounds")
	}
	s.Modules = append(s.Modules[:index], s.Modules[index+1:]...)
	return nil
}

// DamageRandomModule selects a random operational module and disables it. Returns the name of the module, and true if successful.
func (s *Ship) DamageRandomModule() (string, bool) {
	var operational []int
	for i, mod := range s.Modules {
		if !mod.IsDamaged {
			operational = append(operational, i)
		}
	}

	if len(operational) == 0 {
		return "", false
	}

	targetIdx := operational[rng.Intn(len(operational))]
	s.Modules[targetIdx].IsDamaged = true
	return s.Modules[targetIdx].Template.Name, true
}

// CalculateStats aggregates the baseline chassis capabilities with all fully operational modules.
func (s *Ship) CalculateStats() ShipStats {
	stats := ShipStats{
		CargoSpace:      s.Chassis.BaseCargo,
		FuelCapacity:    s.Chassis.BaseFuel,
		Speed:           s.Chassis.BaseSpeed,
		PassengerSpace:  s.Chassis.BaseCabins,
		CrewCapacity:    s.Chassis.BaseQuarters,
		FuelConsumption: s.Chassis.FuelConsumption,
	}

	for _, mod := range s.Modules {
		if !mod.IsDamaged {
			stats.CargoSpace += mod.Template.CargoBonus
			stats.FuelCapacity += mod.Template.FuelBonus
			stats.Speed += mod.Template.SpeedBonus
			stats.PassengerSpace += mod.Template.PassengerSpace
			stats.CrewCapacity += mod.Template.QuartersBonus
			stats.HazardousProtection += mod.Template.HazardousProtection
			stats.FragileProtection += mod.Template.FragileProtection
			stats.IllicitConcealment += mod.Template.IllicitConcealment
			stats.PerishablePreservation += mod.Template.PerishablePreservation
		}
	}

	return stats
}

// Refuel tops off the vessel's fuel tank to its maximum calculated capacity.
func (s *Ship) Refuel() {
	stats := s.CalculateStats()
	s.CurrentFuel = stats.FuelCapacity
}

// AvailableCargoSpace calculates remaining volumetric freight capacity.
func (s *Ship) AvailableCargoSpace() int {
	stats := s.CalculateStats()
	used := 0
	for _, c := range s.Contracts {
		used += c.Quantity
	}
	return stats.CargoSpace - used
}

// AvailablePassengerSpace calculates unoccupied cabin accommodations.
func (s *Ship) AvailablePassengerSpace() int {
	stats := s.CalculateStats()
	return stats.PassengerSpace - len(s.Passengers)
}

// AvailableCrewSpace calculates unoccupied crew quarters.
func (s *Ship) AvailableCrewSpace() int {
	stats := s.CalculateStats()
	return stats.CrewCapacity - len(s.Roster)
}

// UnprotectedCargoAmount determines the volume of specific tagged cargo lacking proper module protections.
func (s *Ship) UnprotectedCargoAmount(tag string, commodities map[string]market.Commodity) int {
	total := 0
	for _, c := range s.Contracts {
		if comm, exists := commodities[c.CommodityID]; exists && comm.HasTag(tag) {
			total += c.Quantity
		}
	}

	stats := s.CalculateStats()
	protection := 0

	switch tag {
	case "hazardous":
		protection = stats.HazardousProtection
	case "fragile":
		protection = stats.FragileProtection
	case "illicit":
		protection = stats.IllicitConcealment
	case "perishable":
		protection = stats.PerishablePreservation
	}

	unprotected := total - protection
	if unprotected < 0 {
		return 0
	}
	return unprotected
}

// HasUnprotectedCargo is a boolean convenience wrapper for UnprotectedCargoAmount.
func (s *Ship) HasUnprotectedCargo(tag string, commodities map[string]market.Commodity) bool {
	return s.UnprotectedCargoAmount(tag, commodities) > 0
}
