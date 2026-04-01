// Package cli provides the command-line interface and menu handling for the game.
package cli

import (
	"fmt"
	"strconv"
	"strings"

	"galaxies-burnrate/internal/engine"
	"galaxies-burnrate/internal/ship"
)

// runMarketMenu traps the user in the planetary market loop for the active ship.
func (c *CLI) runMarketMenu() {
	activeShip := c.getActiveShip()
	if activeShip == nil {
		fmt.Println("\nYou have no ships available.")
		return
	}

	if activeShip.Status != "Idle" {
		fmt.Printf("\nMarket services are unavailable. %s is currently in transit.\n", activeShip.Name)
		return
	}

	currentPlanet := c.gameState.Planets[activeShip.LocationID]

	for {
		fmt.Printf("\n=== %s MARKET ===\n", strings.ToUpper(currentPlanet.Name))
		fmt.Printf("Active Vessel: %s | Corporate Funds: %d C\n", activeShip.Name, c.gameState.Player.Credits)
		fmt.Println("\nServices:")
		fmt.Println("[1] Refuel Vessel")
		fmt.Println("[2] Contract Board")
		fmt.Println("[3] Passenger Terminal")
		fmt.Println("[0] Return to Fleet Operations")
		fmt.Print("\nSelect a service > ")

		if !c.scanner.Scan() {
			return
		}
		input := strings.TrimSpace(c.scanner.Text())

		switch input {
		case "1":
			c.handleRefuel(activeShip)
		case "2":
			c.handleContracts(activeShip)
		case "3":
			c.handlePassengers(activeShip)
		case "0":
			return
		default:
			fmt.Println("\nInvalid selection.")
		}
	}
}

func (c *CLI) handleRefuel(s *ship.Ship) {
	currentPlanet := c.gameState.Planets[s.LocationID]

	if !currentPlanet.HasTag("refuel") {
		fmt.Printf("\nThere are no recognized refueling facilities at %s.\n", currentPlanet.Name)
		return
	}

	invoice, err := engine.CalculateRefuel(s, currentPlanet, c.gameState.Player.Credits, c.settings)
	if err != nil {
		fmt.Printf("\nRefuel Error: %v\n", err)
		return
	}

	fmt.Println("\n--- Refuel Invoice Preview ---")
	fmt.Printf("Fuel to Add: %d units\n", invoice.FuelToAdd)
	fmt.Printf("Total Cost: %d C (Rate: %d C/unit)\n", invoice.TotalCost, currentPlanet.Refuel.Cost)
	fmt.Printf("Remaining Corporate Credits: %d C\n", c.gameState.Player.Credits-invoice.TotalCost)

	fmt.Print("\nAuthorize transaction? (Y/N) > ")
	c.scanner.Scan()
	confirm := strings.TrimSpace(strings.ToLower(c.scanner.Text()))

	if confirm == "y" || confirm == "yes" {
		err := engine.ExecuteRefuel(s, c.gameState.Player, currentPlanet, c.settings)
		if err != nil {
			fmt.Printf("Transaction Error: %v\n", err)
			return
		}
		fmt.Printf("\n>>> Transaction successful. %s refueled.\n", s.Name)
	} else {
		fmt.Println("\n>>> Transaction authorization revoked.")
	}
}

func (c *CLI) handleContracts(s *ship.Ship) {
	planetID := s.LocationID
	contracts := c.gameState.Contracts[planetID]

	if len(contracts) == 0 {
		fmt.Println("\nNo contracts currently available at this location.")
		return
	}

	fmt.Printf("\n--- Contract Board: %s ---\n", c.gameState.Planets[planetID].Name)
	fmt.Printf("Vessel: %s (Available Cargo: %d)\n\n", s.Name, s.AvailableCargoSpace())

	for i, contract := range contracts {
		dest := c.gameState.Planets[contract.DestinationID]
		comm := c.settings.Commodities[contract.CommodityID]

		tagStr := ""
		if len(comm.Tags) > 0 {
			tagStr = fmt.Sprintf(" %v", comm.Tags)
		}

		fmt.Printf("[%d] Deliver %d units of %s%s to %s | Payout: %d C | Expires: Day %d\n",
			i+1, contract.Quantity, comm.Name, strings.ToUpper(tagStr), dest.Name, contract.Payout, contract.ExpirationDay)
	}
	fmt.Println("[0] Cancel")
	fmt.Print("Select contract number to accept > ")

	c.scanner.Scan()
	contractInput := strings.TrimSpace(c.scanner.Text())
	contractIndex, err := strconv.Atoi(contractInput)
	if err != nil || contractIndex < 0 || contractIndex > len(contracts) {
		fmt.Println("Invalid selection. Canceling assignment.")
		return
	}
	if contractIndex == 0 {
		return
	}

	selectedContract := contracts[contractIndex-1]
	selectedComm := c.settings.Commodities[selectedContract.CommodityID]

	if s.AvailableCargoSpace() < selectedContract.Quantity {
		fmt.Println("\nWARNING: Insufficient general cargo space on this vessel. Contract rejected.")
		return
	}

	// Calculate and warn about unprotected cargo loads.
	requiresOverride := false
	stats := s.CalculateStats()

	for _, tag := range selectedComm.Tags {
		currentTotal := 0
		for _, existing := range s.Contracts {
			if exComm := c.settings.Commodities[existing.CommodityID]; exComm.HasTag(tag) {
				currentTotal += existing.Quantity
			}
		}

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

		if currentTotal+selectedContract.Quantity > protection {
			fmt.Printf("\n[WARNING] Accepting this contract will result in unprotected %s cargo, exposing the vessel to significant transit risks and penalties.", strings.ToUpper(tag))
			requiresOverride = true
		}
	}

	if requiresOverride {
		fmt.Print("\nDo you wish to authorize loading this cargo despite the risks? (Y/N) > ")
		c.scanner.Scan()
		confirm := strings.TrimSpace(strings.ToLower(c.scanner.Text()))
		if confirm != "y" && confirm != "yes" {
			fmt.Println("\n>>> Load authorization revoked. Contract rejected.")
			return
		}
	}

	s.Contracts = append(s.Contracts, selectedContract)
	c.gameState.Contracts[planetID] = append(contracts[:contractIndex-1], contracts[contractIndex:]...)

	fmt.Printf("\n>>> Contract Accepted. %s loaded onto %s.\n", selectedComm.Name, s.Name)
}

func (c *CLI) handlePassengers(s *ship.Ship) {
	planetID := s.LocationID
	passengers := c.gameState.Passengers[planetID]

	if len(passengers) == 0 {
		fmt.Println("\nNo passengers currently waiting at this terminal.")
		return
	}

	fmt.Printf("\n--- Passenger Terminal: %s ---\n", c.gameState.Planets[planetID].Name)
	fmt.Printf("Vessel: %s (Available Cabins: %d)\n\n", s.Name, s.AvailablePassengerSpace())

	for i, p := range passengers {
		dest := c.gameState.Planets[p.DestinationID]
		fmt.Printf("[%d] %s [%s] traveling to %s | Fare: %d C | Departs By: Day %d\n",
			i+1, p.Name, p.Class.Name, dest.Name, p.Payout, p.ExpirationDay)
	}
	fmt.Println("[0] Cancel")
	fmt.Print("Select passenger number to board > ")

	c.scanner.Scan()
	passengerInput := strings.TrimSpace(c.scanner.Text())
	passengerIndex, err := strconv.Atoi(passengerInput)
	if err != nil || passengerIndex < 0 || passengerIndex > len(passengers) {
		fmt.Println("Invalid selection. Canceling boarding.")
		return
	}
	if passengerIndex == 0 {
		return
	}

	selectedPassenger := passengers[passengerIndex-1]

	if s.AvailablePassengerSpace() < 1 {
		fmt.Println("\nWARNING: Insufficient cabin space on this vessel. Passenger rejected.")
		return
	}

	s.Passengers = append(s.Passengers, selectedPassenger)
	c.gameState.Passengers[planetID] = append(passengers[:passengerIndex-1], passengers[passengerIndex:]...)

	fmt.Printf("\n>>> Boarding Complete. %s has boarded %s.\n", selectedPassenger.Name, s.Name)
}
