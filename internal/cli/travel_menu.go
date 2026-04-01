// Package cli provides the command-line interface and menu handling for the game.
package cli

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"galaxies-burnrate/internal/engine"
	"galaxies-burnrate/internal/planet"
)

// handleTravel walks the user through dispatching a jump for the active vessel.
func (c *CLI) handleTravel() {
	activeShip := c.getActiveShip()
	if activeShip == nil {
		fmt.Println("\nYou have no ships available.")
		return
	}

	if activeShip.Status == "In Transit" {
		fmt.Printf("\n%s is currently in transit and cannot receive navigation orders.\n", activeShip.Name)
		return
	}

	currentPlanet := c.gameState.Planets[activeShip.LocationID]

	var destinations []planet.Planet
	for _, p := range c.gameState.Planets {
		if p.ID != activeShip.LocationID {
			destinations = append(destinations, p)
		}
	}
	sort.Slice(destinations, func(i, j int) bool {
		return destinations[i].Name < destinations[j].Name
	})

	fmt.Printf("\n--- Navigation: %s ---\n", activeShip.Name)
	for i, p := range destinations {
		fmt.Printf("[%d] %s (Sector %d, %d)\n", i+1, p.Name, p.X, p.Y)
	}
	fmt.Println("[0] Cancel")
	fmt.Print("Select destination number > ")

	c.scanner.Scan()
	destInput := strings.TrimSpace(c.scanner.Text())
	destIndex, err := strconv.Atoi(destInput)
	if err != nil || destIndex < 0 || destIndex > len(destinations) {
		fmt.Println("Invalid selection. Canceling travel.")
		return
	}
	if destIndex == 0 {
		return
	}

	targetPlanet := destinations[destIndex-1]

	plan, err := engine.CalculateTravel(activeShip, currentPlanet, targetPlanet)
	if err != nil {
		fmt.Printf("Calculation Error: %v\n", err)
		return
	}

	fmt.Println("\n--- Jump Manifest Preview ---")
	fmt.Printf("Route: %s -> %s\n", currentPlanet.Name, targetPlanet.Name)
	fmt.Printf("Distance: %d units\n", plan.Distance)
	fmt.Printf("Fuel Cost: %d units (Current: %d)\n", plan.FuelCost, activeShip.CurrentFuel)
	fmt.Printf("Time Required: %d days\n", plan.TimeCost)

	if activeShip.CurrentFuel < plan.FuelCost {
		fmt.Println("\nWARNING: Insufficient fuel for this journey. Jump aborted.")
		return
	}

	fmt.Print("\nAuthorize jump? (Y/N) > ")
	c.scanner.Scan()
	confirm := strings.TrimSpace(strings.ToLower(c.scanner.Text()))

	if confirm == "y" || confirm == "yes" {
		err := engine.ExecuteTravel(activeShip, currentPlanet, targetPlanet, c.gameState.CurrentDay)
		if err != nil {
			fmt.Printf("Navigation Error: %v\n", err)
			return
		}
		fmt.Printf("\n>>> Orders received. %s is en route to %s.\n", activeShip.Name, targetPlanet.Name)
	} else {
		fmt.Println("\n>>> Jump authorization revoked.")
	}
}
