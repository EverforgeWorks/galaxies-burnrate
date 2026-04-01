// Package engine contains the core simulation mechanics and state mutation logic.
package engine

import (
	"fmt"

	"galaxies-burnrate/internal/config"
	"galaxies-burnrate/internal/ship"
)

// processDailyCrewHealth handles the contraction of illnesses, natural recovery, and active Medic healing checks.
func processDailyCrewHealth(s *ship.Ship, settings *config.GameSettings) []string {
	var logs []string

	// Identify the most capable, healthy Medic on board.
	bestMedicSkill := 0
	bestMedicName := ""
	for _, c := range s.Roster {
		if c.ActiveIllness == "" && !c.OnVacation {
			// A crew member can act as a Medic if it is their primary role or secondary specialty.
			if c.Role.ID == "medic" || c.Specialty.SecondaryRole == "medic" {
				if c.SkillLevel > bestMedicSkill {
					bestMedicSkill = c.SkillLevel
					bestMedicName = c.Name
				}
			}
		}
	}

	// Iterate through the roster by index to allow direct mutation of the structs.
	for i := range s.Roster {
		member := &s.Roster[i]

		if member.OnVacation {
			continue
		}

		// 1. Illness Contraction
		if member.ActiveIllness == "" && !member.Origin.ImmuneToSickness {
			if rng.Intn(100) < settings.CrewStates.IllnessContractionChancePercent {
				if len(settings.IllnessNames) > 0 {
					illness := settings.IllnessNames[rng.Intn(len(settings.IllnessNames))]
					member.ActiveIllness = illness

					// Calculate a random duration within the configured bounds.
					durationRange := settings.CrewStates.IllnessMaxDays - settings.CrewStates.IllnessMinDays + 1
					member.IllnessDays = rng.Intn(durationRange) + settings.CrewStates.IllnessMinDays

					logs = append(logs, fmt.Sprintf("    - [MEDICAL ALERT] %s has contracted %s. Estimated natural recovery: %d days.", member.Name, member.ActiveIllness, member.IllnessDays))
				}
			}
		}

		// 2. Illness Progression & Healing
		if member.ActiveIllness != "" {
			// Natural recovery progresses by 1 day automatically.
			member.IllnessDays--

			// Active Medic Intervention
			if bestMedicName != "" && bestMedicName != member.Name {
				roll := rng.Intn(20) + 1
				total := roll + bestMedicSkill
				if total >= settings.CrewStates.IllnessBaseDC {
					member.IllnessDays--
					logs = append(logs, fmt.Sprintf("      >> %s successfully treated %s's symptoms, accelerating recovery.", bestMedicName, member.Name))
				}
			}

			// Clear the illness if the duration has fully elapsed.
			if member.IllnessDays <= 0 {
				logs = append(logs, fmt.Sprintf("    - [MEDICAL UPDATE] %s has fully recovered from %s and returned to active duty.", member.Name, member.ActiveIllness))
				member.ActiveIllness = ""
				member.IllnessDays = 0
			}
		}
	}

	return logs
}
