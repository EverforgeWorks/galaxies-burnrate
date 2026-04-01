// Package cli provides the command-line interface and menu handling for the game.
package cli

import (
	"fmt"
	"strconv"
	"strings"

	"galaxies-burnrate/internal/engine"
	"galaxies-burnrate/internal/planet"
	"galaxies-burnrate/internal/ship"
)

// handleShipyard routes the player to the local shipyard interface for the active ship.
func (c *CLI) handleShipyard() {
	activeShip := c.getActiveShip()
	if activeShip == nil {
		fmt.Println("\nYou have no ships available.")
		return
	}

	if activeShip.Status != "Idle" {
		fmt.Printf("\n%s is currently in transit and cannot access shipyard facilities.\n", activeShip.Name)
		return
	}

	p := c.gameState.Planets[activeShip.LocationID]
	if !p.HasTag("shipyard") {
		fmt.Printf("\nThere are no recognized Shipyard facilities at %s.\n", p.Name)
		return
	}

	c.runShipyardMenu(activeShip, p)
}

// runShipyardMenu traps the user in the planetary shipyard loop for the active ship.
func (c *CLI) runShipyardMenu(s *ship.Ship, p planet.Planet) {
	inventory := p.Shipyard

	for {
		fmt.Printf("\n=== %s SHIPYARD ===\n", strings.ToUpper(p.Name))
		fmt.Printf("Active Vessel: %s | Corporate Funds: %d C\n", s.Name, c.gameState.Player.Credits)
		fmt.Println("\nServices:")
		fmt.Println("[1] Purchase New Vessel")
		fmt.Println("[2] Modify Active Vessel (Purchase Module)")
		fmt.Println("[3] Modify Active Vessel (Sell Module)")
		fmt.Println("[4] Repair Damaged Modules")
		fmt.Println("[0] Return to Fleet Operations")
		fmt.Print("\nSelect a service > ")

		if !c.scanner.Scan() {
			return
		}
		input := strings.TrimSpace(c.scanner.Text())

		switch input {
		case "1":
			c.purchaseVesselFlow(p, inventory.Chassis)
		case "2":
			c.purchaseModuleFlow(s, p, inventory.Modules)
		case "3":
			c.sellModuleFlow(s)
		case "4":
			c.repairModuleFlow(s, p)
		case "0":
			return
		default:
			fmt.Println("\nInvalid selection.")
		}
	}
}

// purchaseVesselFlow handles the listing and purchasing of new ships at the current planetary shipyard.
func (c *CLI) purchaseVesselFlow(p planet.Planet, availableChassis []string) {
	if len(availableChassis) == 0 {
		fmt.Println("\nNo chassis currently available at this shipyard.")
		return
	}

	fmt.Println("\n--- Available Vessels ---")
	for i, cID := range availableChassis {
		chassis := c.settings.Chassis[cID]
		cost := engine.CalculateNewShipCost(chassis, p, c.settings)
		fmt.Printf("[%d] %s | Cost: %d C | Max Modules: %d\n", i+1, chassis.Name, cost, chassis.MaxModules)
		fmt.Printf("    Description: %s\n", chassis.Description)
	}
	fmt.Println("[0] Cancel")
	fmt.Print("Select chassis number to purchase > ")

	c.scanner.Scan()
	chassisInput := strings.TrimSpace(c.scanner.Text())
	chassisIndex, err := strconv.Atoi(chassisInput)
	if err != nil || chassisIndex < 0 || chassisIndex > len(availableChassis) {
		fmt.Println("Invalid selection. Transaction aborted.")
		return
	}
	if chassisIndex == 0 {
		return
	}

	selectedChassis := c.settings.Chassis[availableChassis[chassisIndex-1]]
	cost := engine.CalculateNewShipCost(selectedChassis, p, c.settings)

	if c.gameState.Player.Credits < cost {
		fmt.Println("\nWARNING: Insufficient corporate funds for this transaction.")
		return
	}

	fmt.Print("\nEnter a designation for the new vessel > ")
	c.scanner.Scan()
	shipName := strings.TrimSpace(c.scanner.Text())
	if shipName == "" {
		fmt.Println("Designation required. Transaction aborted.")
		return
	}

	err = engine.PurchaseShip(c.gameState, c.settings, selectedChassis, p, shipName)
	if err != nil {
		fmt.Printf("\nTransaction Error: %v\n", err)
		return
	}

	fmt.Printf("\n>>> Transaction complete. %s (%s) added to fleet.\n", shipName, selectedChassis.Name)
}

// purchaseModuleFlow handles the listing and installation of new modules onto the active vessel.
func (c *CLI) purchaseModuleFlow(s *ship.Ship, p planet.Planet, availableModules []string) {
	if len(availableModules) == 0 {
		fmt.Println("\nNo modules currently available at this shipyard.")
		return
	}

	fmt.Printf("\n--- Available Modules for %s (Capacity: %d/%d) ---\n", s.Name, len(s.Modules), s.Chassis.MaxModules)
	for i, mID := range availableModules {
		mod := c.settings.Modules[mID]
		cost := engine.CalculateModuleCost(mod, p, s, c.settings)
		fmt.Printf("[%d] %s | Cost: %d C\n", i+1, mod.Name, cost)
	}
	fmt.Println("[0] Cancel")
	fmt.Print("Select module number to purchase > ")

	c.scanner.Scan()
	modInput := strings.TrimSpace(c.scanner.Text())
	modIndex, err := strconv.Atoi(modInput)
	if err != nil || modIndex < 0 || modIndex > len(availableModules) {
		fmt.Println("Invalid selection. Transaction aborted.")
		return
	}
	if modIndex == 0 {
		return
	}

	selectedMod := c.settings.Modules[availableModules[modIndex-1]]

	err = engine.PurchaseModule(c.gameState, c.settings, s, p, selectedMod)
	if err != nil {
		fmt.Printf("\nInstallation Error: %v\n", err)
		return
	}

	fmt.Printf("\n>>> Transaction complete. %s installed on %s.\n", selectedMod.Name, s.Name)
}

// sellModuleFlow handles the uninstallation and liquidation of existing modules on the active vessel.
func (c *CLI) sellModuleFlow(s *ship.Ship) {
	if len(s.Modules) == 0 {
		fmt.Println("\nThis vessel currently has no removable modules installed.")
		return
	}

	fmt.Printf("\n--- Installed Modules on %s ---\n", s.Name)
	for i, mod := range s.Modules {
		refund := int64(float64(mod.Template.Value) * c.settings.Multipliers.ModuleRefundRate)

		statusStr := ""
		if mod.IsDamaged {
			refund = int64(float64(refund) * 0.5)
			statusStr = "[DAMAGED - Scrap Value] "
		}

		fmt.Printf("[%d] %s %s(Refund Value: %d C)\n", i+1, mod.Template.Name, statusStr, refund)
	}
	fmt.Println("[0] Cancel")
	fmt.Print("Select module number to sell/uninstall > ")

	c.scanner.Scan()
	modInput := strings.TrimSpace(c.scanner.Text())
	modIndex, err := strconv.Atoi(modInput)
	if err != nil || modIndex < 0 || modIndex > len(s.Modules) {
		fmt.Println("Invalid selection. Operation aborted.")
		return
	}
	if modIndex == 0 {
		return
	}

	modToSell := s.Modules[modIndex-1]
	actualRefund := int64(float64(modToSell.Template.Value) * c.settings.Multipliers.ModuleRefundRate)
	if modToSell.IsDamaged {
		actualRefund = int64(float64(actualRefund) * 0.5)
	}

	err = engine.SellModule(c.gameState, c.settings, s, modIndex-1)
	if err != nil {
		fmt.Printf("\nUninstallation Error: %v\n", err)
		return
	}

	fmt.Printf("\n>>> Asset Liquidated. %s removed from %s and %d C refunded.\n", modToSell.Template.Name, s.Name, actualRefund)
}

// repairModuleFlow handles restoring damaged modules to operational status.
func (c *CLI) repairModuleFlow(s *ship.Ship, p planet.Planet) {
	damagedCount := 0
	fmt.Printf("\n--- Damaged Modules on %s ---\n", s.Name)

	for i, mod := range s.Modules {
		if mod.IsDamaged {
			damagedCount++
			repairCost := engine.CalculateRepairCost(mod.Template, p, s, c.settings)
			fmt.Printf("[%d] %s | Repair Cost: %d C\n", i+1, mod.Template.Name, repairCost)
		}
	}

	if damagedCount == 0 {
		fmt.Println("All systems optimal. No repairs required.")
		return
	}

	fmt.Println("[0] Cancel")
	fmt.Print("Select module number to repair > ")

	c.scanner.Scan()
	modInput := strings.TrimSpace(c.scanner.Text())
	modIndex, err := strconv.Atoi(modInput)
	if err != nil || modIndex < 0 || modIndex > len(s.Modules) {
		fmt.Println("Invalid selection. Operation aborted.")
		return
	}
	if modIndex == 0 {
		return
	}

	targetMod := s.Modules[modIndex-1]
	if !targetMod.IsDamaged {
		fmt.Println("That module is already fully operational. Selection aborted.")
		return
	}

	err = engine.RepairModule(c.gameState, c.settings, s, p, modIndex-1)
	if err != nil {
		fmt.Printf("\nRepair Error: %v\n", err)
		return
	}

	fmt.Printf("\n>>> Repairs completed. %s is fully operational.\n", targetMod.Template.Name)
}
