// Package cli provides the command-line interface and menu handling for the game.
package cli

import (
	"fmt"
	"strings"
)

// handleCargo displays the current inventory, passenger manifests, active personnel roster, and cargo risk assessments.
func (c *CLI) handleCargo() {
	activeShip := c.getActiveShip()
	if activeShip == nil {
		fmt.Println("\nYou have no ships available.")
		return
	}

	stats := activeShip.CalculateStats()
	availCargo := activeShip.AvailableCargoSpace()
	availPass := activeShip.AvailablePassengerSpace()
	availCrew := activeShip.AvailableCrewSpace()

	fmt.Printf("\n--- Asset Manifest: %s [%s] ---\n", activeShip.Name, activeShip.Chassis.Name)

	fmt.Printf("\nCargo Capacity: %d / %d utilized\n", stats.CargoSpace-availCargo, stats.CargoSpace)
	if len(activeShip.Contracts) == 0 {
		fmt.Println("  No active contracts or cargo loaded.")
	} else {
		for _, contract := range activeShip.Contracts {
			dest := c.gameState.Planets[contract.DestinationID]
			comm := c.settings.Commodities[contract.CommodityID]

			tagStr := ""
			if len(comm.Tags) > 0 {
				tagStr = fmt.Sprintf(" %v", comm.Tags)
			}

			fmt.Printf("  - [CARGO] %d units of %s%s destined for %s (Payout: %d C)\n",
				contract.Quantity, comm.Name, strings.ToUpper(tagStr), dest.Name, contract.Payout)
		}

		fmt.Println("\n  [Specialized Storage Assessment]")
		tags := []string{"hazardous", "fragile", "illicit", "perishable"}
		hasSpecialized := false

		for _, tag := range tags {
			unprotected := activeShip.UnprotectedCargoAmount(tag, c.settings.Commodities)
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

			if unprotected > 0 || protection > 0 {
				hasSpecialized = true
				status := "SECURE"
				if unprotected > 0 {
					status = "AT RISK"
				}
				fmt.Printf("  > %s: Capacity %d | Unprotected Volume: %d [%s]\n", strings.ToUpper(tag), protection, unprotected, status)
			}
		}

		if !hasSpecialized {
			fmt.Println("  > No specialized or tagged cargo currently loaded.")
		}
	}

	fmt.Printf("\nCabin Capacity: %d / %d utilized\n", stats.PassengerSpace-availPass, stats.PassengerSpace)
	if len(activeShip.Passengers) == 0 {
		fmt.Println("  No passengers currently boarded.")
	} else {
		for _, p := range activeShip.Passengers {
			dest := c.gameState.Planets[p.DestinationID]
			fmt.Printf("  - [PASSENGER] %s [%s] destined for %s\n", p.Name, p.Class.Name, dest.Name)
		}
	}

	fmt.Printf("\nCrew Quarters: %d / %d utilized\n", stats.CrewCapacity-availCrew, stats.CrewCapacity)
	if len(activeShip.Roster) == 0 {
		fmt.Println("  No active duty personnel on board.")
	} else {
		for _, member := range activeShip.Roster {
			fmt.Printf("  - [CREW] %s | %s | Skill: %d/10 (Salary: %d C/day)\n", member.Name, member.Role.Name, member.SkillLevel, member.DailySalary)
		}
	}
}
