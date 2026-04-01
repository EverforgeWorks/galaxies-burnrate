// Package engine contains the core simulation mechanics and state mutation logic.
package engine

import (
	"fmt"

	"galaxies-burnrate/internal/config"
	"galaxies-burnrate/internal/game"
)

// processFleetPayroll calculates the total salary obligations for all active personnel across the fleet and deducts it from corporate funds.
func processFleetPayroll(state *game.State, settings *config.GameSettings) []string {
	var logs []string
	var totalPayroll int64 = 0

	// Calculate total payroll across all active ships.
	for _, s := range state.Player.Fleet {
		if s.Status == "Destroyed" {
			continue
		}
		for _, member := range s.Roster {
			if !member.OnVacation {
				totalPayroll += member.DailySalary * int64(settings.Crew.PayPeriodDays)
			}
		}
	}

	// Apply deductions and check for insolvency conditions.
	if totalPayroll > 0 {
		state.Player.Credits -= totalPayroll
		logs = append(logs, fmt.Sprintf(">>> PAYROLL DEDUCTION: %d C processed for fleet personnel salaries.", totalPayroll))

		if state.Player.Credits < 0 {
			logs = append(logs, "    - [WARNING] Corporate accounts are overdrawn. Missed payroll obligations have severely damaged fleet morale.")
			logs = append(logs, applyBankruptcyPenalties(state)...)
		}
	}

	return logs
}

// applyBankruptcyPenalties applies a universal morale degradation to all biological crew members when the corporation cannot meet its financial obligations.
func applyBankruptcyPenalties(state *game.State) []string {
	var logs []string
	penaltyApplied := false

	for _, s := range state.Player.Fleet {
		if s.Status == "Destroyed" {
			continue
		}

		// Iterate by index to allow direct struct mutation.
		for j := range s.Roster {
			member := &s.Roster[j]

			// Androids (daily salary of 0) and crew on vacation do not suffer morale hits from financial insolvency.
			if !member.OnVacation && member.DailySalary > 0 {
				member.Morale -= 2
				if member.Morale < 1 {
					member.Morale = 1
				}
				penaltyApplied = true
			}
		}
	}

	if penaltyApplied {
		logs = append(logs, "    - [FLEET STATUS] Biological crew morale has dropped across the fleet due to insolvency.")
	}

	return logs
}
