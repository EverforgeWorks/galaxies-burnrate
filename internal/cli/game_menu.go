// Package cli provides the command-line interface and menu handling for the game.
package cli

import (
	"fmt"
	"strconv"
	"strings"

	"galaxies-burnrate/internal/engine"
)

// runGameMenu acts as the primary operations loop for active gameplay.
func (c *CLI) runGameMenu() bool {
	for {
		c.renderDashboard()

		fmt.Print("\nSelect an action > ")
		if !c.scanner.Scan() {
			return false
		}
		input := strings.TrimSpace(strings.ToLower(c.scanner.Text()))

		switch input {
		case "1", "t", "travel":
			c.handleTravel()
		case "2", "w", "wait":
			c.handleWait()
		case "3", "m", "market":
			c.runMarketMenu()
		case "4", "s", "shipyard":
			c.handleShipyard()
		case "5", "c", "cantina":
			c.handleCantina()
		case "6", "i", "inventory", "roster":
			c.handleCargo()
		case "7", "p", "planet", "details":
			c.handlePlanetaryDetails()
		case "8", "q", "quit":
			c.currentState = StateMainMenu
			return true
		default:
			fmt.Println("\nInvalid selection. Please choose a valid operation.")
		}
	}
}

// handleWait processes the passage of time and triggers the simulation engine loop.
func (c *CLI) handleWait() {
	fmt.Print("\nEnter number of days to pass > ")
	c.scanner.Scan()
	input := strings.TrimSpace(c.scanner.Text())
	days, err := strconv.Atoi(input)
	if err != nil || days <= 0 {
		fmt.Println("Invalid input. Please enter a positive number of days.")
		return
	}

	fmt.Printf("\nProcessing %d days...\n", days)
	events := engine.AdvanceTime(c.gameState, c.settings, days)

	fmt.Printf(">>> Operations updated. It is now Day %d.\n", c.gameState.CurrentDay)
	if len(events) > 0 {
		fmt.Println("\n--- Operations Reports ---")
		for _, e := range events {
			fmt.Println(e)
		}
	}
}

// renderDashboard prints the high-level overview of corporate assets, active vessels, and cargo manifests.
func (c *CLI) renderDashboard() {
	fmt.Printf("\n=== DUTY CALLERS ===\n")
	fmt.Printf("Manager: %s | Credits: %d C | Day: %d\n", c.gameState.Player.Name, c.gameState.Player.Credits, c.gameState.CurrentDay)

	fmt.Println("\n--- Fleet Status ---")
	for _, s := range c.gameState.Player.Fleet {
		locDisplay := s.LocationID

		if s.Status == "Idle" {
			p := c.gameState.Planets[s.LocationID]
			desc := p.Descriptor(c.settings.EconomicStatuses, c.settings.SocialStatuses, c.settings.PoliticalStatuses)
			locDisplay = fmt.Sprintf("%s (%s)", p.Name, desc)
		} else {
			p := c.gameState.Planets[s.DestinationID]
			locDisplay = fmt.Sprintf("Dest: %s (ETA: Day %d)", p.Name, s.ArrivalDay)
		}

		stats := s.CalculateStats()
		fmt.Printf("- %s [%s] | Loc: %s | Fuel: %d/%d | Status: %s\n", s.Name, s.Chassis.Name, locDisplay, s.CurrentFuel, stats.FuelCapacity, s.Status)

		if len(s.Contracts) > 0 {
			for _, contract := range s.Contracts {
				destName := c.gameState.Planets[contract.DestinationID].Name
				commName := c.settings.Commodities[contract.CommodityID].Name
				fmt.Printf("    > Cargo: %d units of %s (To: %s)\n", contract.Quantity, commName, destName)
			}
		}

		if len(s.Passengers) > 0 {
			for _, pass := range s.Passengers {
				destName := c.gameState.Planets[pass.DestinationID].Name
				fmt.Printf("    > Passenger: %s [%s] (To: %s)\n", pass.Name, pass.Class.Name, destName)
			}
		}
	}

	activeShip := c.getActiveShip()
	if activeShip != nil {
		fmt.Printf("\n>>> Active Vessel: %s\n", activeShip.Name)
	}

	fmt.Println("\nActions:")
	fmt.Println("[1] (T)ravel Menu")
	fmt.Println("[2] (W)ait / Pass Time")
	fmt.Println("[3] (M)arket Services")
	fmt.Println("[4] (S)hipyard")
	fmt.Println("[5] (C)antina / Hiring Hall")
	fmt.Println("[6] (I)nventory / Roster")
	fmt.Println("[7] (P)lanetary Details")
	fmt.Println("[8] (Q)uit to Main Menu")
}

// handlePlanetaryDetails renders a comprehensive dossier of the local planet's configured traits and multipliers.
func (c *CLI) handlePlanetaryDetails() {
	activeShip := c.getActiveShip()
	if activeShip == nil {
		fmt.Println("\nYou have no ships available.")
		return
	}

	if activeShip.Status != "Idle" {
		fmt.Printf("\n%s is currently in transit. Planetary scanners are offline.\n", activeShip.Name)
		return
	}

	p := c.gameState.Planets[activeShip.LocationID]
	eco := c.settings.EconomicStatuses[p.EconomicStatus]
	soc := c.settings.SocialStatuses[p.SocialStatus]
	pol := c.settings.PoliticalStatuses[p.PoliticalStatus]

	desc := p.Descriptor(c.settings.EconomicStatuses, c.settings.SocialStatuses, c.settings.PoliticalStatuses)

	fmt.Printf("\n=== %s: PLANETARY DOSSIER ===\n", strings.ToUpper(p.Name))
	fmt.Printf("Classification: %s\n", desc)

	fmt.Printf("\n[Economic Status: %s]\n", eco.Name)
	fmt.Printf("  - Export Contract Value: %.2fx\n", eco.ExportContractMultiplier)
	fmt.Printf("  - Import Contract Value: %.2fx\n", eco.ImportContractMultiplier)
	fmt.Printf("  - Export Passenger Fare: %.2fx\n", eco.ExportPassengerMultiplier)
	fmt.Printf("  - Import Passenger Fare: %.2fx\n", eco.ImportPassengerMultiplier)
	fmt.Printf("  - Local Fuel Cost: %.2fx\n", eco.FuelCostMultiplier)
	fmt.Printf("  - Shipyard Pricing: %.2fx\n", eco.ShipyardCostMultiplier)

	fmt.Printf("\n[Social Status: %s]\n", soc.Name)
	fmt.Printf("  - Freight Volume Generation: %.2fx\n", soc.ContractVolumeMultiplier)
	fmt.Printf("  - Emigration (Departing Passengers): %.2fx\n", soc.PassengerExportVolumeMultiplier)
	fmt.Printf("  - Immigration (Arriving Passengers): %.2fx\n", soc.PassengerImportVolumeMultiplier)

	fmt.Printf("\n[Political Status: %s]\n", pol.Name)
	currentDockingFee := int64(float64(p.DockingFee) * pol.DockingFeeMultiplier)
	fmt.Printf("  - Base Docking Fee: %d C (Modifier: %.2fx = %d C)\n", p.DockingFee, pol.DockingFeeMultiplier, currentDockingFee)
	fmt.Printf("  - Customs Inspection Modifier: %+d%%\n", pol.InspectionChanceModifier)
	fmt.Printf("  - Illicit Export Availability: %.2fx\n", pol.IllicitExportChanceMultiplier)
	fmt.Printf("  - Illicit Import Payout: %.2fx\n", pol.IllicitImportPayoutMultiplier)

	fmt.Println("\nPress Enter to return to operations...")
	c.scanner.Scan()
}
