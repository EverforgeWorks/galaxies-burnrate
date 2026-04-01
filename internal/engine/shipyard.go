// Package engine contains the core simulation mechanics and state mutation logic.
package engine

import (
	"errors"
	"math"

	"galaxies-burnrate/internal/config"
	"galaxies-burnrate/internal/game"
	"galaxies-burnrate/internal/planet"
	"galaxies-burnrate/internal/ship"
)

// CalculateNewShipCost determines the localized market price of a new vessel based on the planet's economic status.
func CalculateNewShipCost(chassis ship.ChassisTemplate, p planet.Planet, settings *config.GameSettings) int64 {
	ecoStatus := settings.EconomicStatuses[p.EconomicStatus]
	return int64(math.Round(float64(chassis.BaseValue) * ecoStatus.ShipyardCostMultiplier))
}

// CalculateModuleCost determines the localized market price of a new module, factoring in planetary traits and crew specialties.
func CalculateModuleCost(module ship.ModuleTemplate, p planet.Planet, s *ship.Ship, settings *config.GameSettings) int64 {
	ecoStatus := settings.EconomicStatuses[p.EconomicStatus]
	baseLocalizedCost := float64(module.Value) * ecoStatus.ShipyardCostMultiplier

	highestDiscount := 0.0
	for _, member := range s.Roster {
		if member.Specialty.ShipyardDiscount > highestDiscount {
			highestDiscount = member.Specialty.ShipyardDiscount
		}
	}

	finalCost := baseLocalizedCost * (1.0 - highestDiscount)
	return int64(math.Round(finalCost))
}

// PurchaseShip deducts corporate funds and adds a newly provisioned vessel to the player's fleet.
func PurchaseShip(state *game.State, settings *config.GameSettings, chassis ship.ChassisTemplate, p planet.Planet, name string) error {
	cost := CalculateNewShipCost(chassis, p, settings)

	if state.Player.Credits < cost {
		return errors.New("insufficient corporate funds")
	}

	newShip := ship.New(name, chassis, p.ID)

	for _, modID := range chassis.DefaultModules {
		mod, exists := settings.Modules[modID]
		if exists {
			_ = newShip.InstallModule(mod)
		}
	}

	newShip.Refuel()

	state.Player.Credits -= cost
	state.Player.Fleet[newShip.ID] = newShip

	return nil
}

// PurchaseModule deducts corporate funds and installs a new component onto the active vessel.
func PurchaseModule(state *game.State, settings *config.GameSettings, s *ship.Ship, p planet.Planet, module ship.ModuleTemplate) error {
	if len(s.Modules) >= s.Chassis.MaxModules {
		return errors.New("maximum chassis module capacity reached")
	}

	cost := CalculateModuleCost(module, p, s, settings)

	if state.Player.Credits < cost {
		return errors.New("insufficient corporate funds")
	}

	state.Player.Credits -= cost
	return s.InstallModule(module)
}

// SellModule removes a component from the vessel and refunds a percentage of its base value to the corporate account.
func SellModule(state *game.State, settings *config.GameSettings, s *ship.Ship, moduleIndex int) error {
	if moduleIndex < 0 || moduleIndex >= len(s.Modules) {
		return errors.New("invalid module selection")
	}

	mod := s.Modules[moduleIndex]
	refund := float64(mod.Template.Value) * settings.Multipliers.ModuleRefundRate

	// Damaged modules sell for significantly less scrap value (50% penalty)
	if mod.IsDamaged {
		refund *= 0.5
	}

	err := s.UninstallModule(moduleIndex)
	if err != nil {
		return err
	}

	state.Player.Credits += int64(math.Round(refund))
	return nil
}

// CalculateRepairCost determines the localized market price to fix a disabled module.
func CalculateRepairCost(module ship.ModuleTemplate, p planet.Planet, s *ship.Ship, settings *config.GameSettings) int64 {
	ecoStatus := settings.EconomicStatuses[p.EconomicStatus]

	// Base repair fee is the module's value * the configured repair rate (e.g., 25%)
	baseRepairFee := float64(module.Value) * settings.Multipliers.ModuleRepairRate

	// Apply the planet's economic repair multiplier
	localizedRepairFee := baseRepairFee * ecoStatus.RepairCostMultiplier

	// Apply crew shipyard discounts
	highestDiscount := 0.0
	for _, member := range s.Roster {
		if member.Specialty.ShipyardDiscount > highestDiscount {
			highestDiscount = member.Specialty.ShipyardDiscount
		}
	}

	finalCost := localizedRepairFee * (1.0 - highestDiscount)
	return int64(math.Round(finalCost))
}

// RepairModule deducts corporate funds and restores a damaged module to operational status.
func RepairModule(state *game.State, settings *config.GameSettings, s *ship.Ship, p planet.Planet, moduleIndex int) error {
	if moduleIndex < 0 || moduleIndex >= len(s.Modules) {
		return errors.New("invalid module selection")
	}

	if !s.Modules[moduleIndex].IsDamaged {
		return errors.New("module is already fully operational")
	}

	cost := CalculateRepairCost(s.Modules[moduleIndex].Template, p, s, settings)

	if state.Player.Credits < cost {
		return errors.New("insufficient corporate funds to authorize repairs")
	}

	state.Player.Credits -= cost
	s.Modules[moduleIndex].IsDamaged = false
	return nil
}
