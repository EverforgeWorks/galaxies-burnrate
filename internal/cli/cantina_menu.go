// Package cli provides the command-line interface and menu handling for the game.
package cli

import (
	"fmt"
	"galaxies-burnrate/internal/ship"
	"strconv"
	"strings"
)

// handleCantina routes the player to the local hiring hall if one exists.
func (c *CLI) handleCantina() {
	activeShip := c.getActiveShip()
	if activeShip == nil {
		fmt.Println("\nYou have no ships available.")
		return
	}

	if activeShip.Status != "Idle" {
		fmt.Printf("\nCantina services are unavailable. %s is currently in transit.\n", activeShip.Name)
		return
	}

	currentPlanet := c.gameState.Planets[activeShip.LocationID]
	if !currentPlanet.HasTag("cantina") {
		fmt.Printf("\nThere are no recognized hiring facilities or cantinas at %s.\n", currentPlanet.Name)
		return
	}

	c.runCantinaMenu(activeShip)
}

func (c *CLI) runCantinaMenu(s *ship.Ship) {
	planetID := s.LocationID
	planet := c.gameState.Planets[planetID]

	for {
		availableCrew := c.gameState.AvailableCrew[planetID]
		fmt.Printf("\n=== %s CANTINA ===\n", strings.ToUpper(planet.Name))
		fmt.Printf("Active Vessel: %s (Available Quarters: %d) | Corporate Funds: %d C\n", s.Name, s.AvailableCrewSpace(), c.gameState.Player.Credits)

		if len(availableCrew) == 0 {
			fmt.Println("\nThe hiring board is currently empty. Try back another day.")
			fmt.Println("\n[0] Return to Fleet Operations")
		} else {
			fmt.Println("\n--- Available Personnel ---")
			for i, member := range availableCrew {
				fmt.Printf("[%d] %s (%s) | Skill: %d/10\n", i+1, member.Name, member.Role.Name, member.SkillLevel)
				fmt.Printf("    Origin: %s | Specialty: %s\n", member.Origin.Name, member.Specialty.Name)
				fmt.Printf("    Hire Cost: %d C | Daily Salary: %d C\n", member.HireCost, member.DailySalary)
			}
			fmt.Println("[0] Cancel")
		}

		fmt.Print("\nSelect a candidate number to review and hire > ")

		if !c.scanner.Scan() {
			return
		}
		input := strings.TrimSpace(c.scanner.Text())

		index, err := strconv.Atoi(input)
		if err != nil || index < 0 || index > len(availableCrew) {
			fmt.Println("Invalid selection.")
			continue
		}
		if index == 0 {
			return
		}

		selected := availableCrew[index-1]

		if s.AvailableCrewSpace() < 1 {
			fmt.Println("\nWARNING: Insufficient active duty quarters on this vessel. Contract rejected.")
			continue
		}

		if c.gameState.Player.Credits < selected.HireCost {
			fmt.Println("\nWARNING: Insufficient corporate funds to cover the upfront hiring cost.")
			continue
		}

		fmt.Printf("\nAuthorize hiring contract for %s? Upfront cost: %d C (Y/N) > ", selected.Name, selected.HireCost)
		c.scanner.Scan()
		confirm := strings.TrimSpace(strings.ToLower(c.scanner.Text()))

		if confirm == "y" || confirm == "yes" {
			c.gameState.Player.Credits -= selected.HireCost
			s.Roster = append(s.Roster, selected)

			// Remove from the available pool
			c.gameState.AvailableCrew[planetID] = append(availableCrew[:index-1], availableCrew[index:]...)

			fmt.Printf("\n>>> Contract signed. %s has joined the crew of %s.\n", selected.Name, s.Name)
		} else {
			fmt.Println("\n>>> Contract negotiations terminated.")
		}
	}
}
