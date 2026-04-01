// Package engine contains the core simulation mechanics and state mutation logic.
package engine

import (
	"fmt"
	"math"
	mathrand "math/rand"

	"galaxies-burnrate/internal/config"
	"galaxies-burnrate/internal/crew"
	"galaxies-burnrate/internal/game"
	"galaxies-burnrate/internal/market"
	"galaxies-burnrate/internal/ship"
)

// calculateEffectiveSkill determines a crew member's operational capability based on their health, presence, and morale curve.
func calculateEffectiveSkill(member *crew.CrewMember, settings *config.GameSettings) (int, string) {
	if member.OnVacation {
		return 0, "On Vacation"
	}
	if member.ActiveIllness != "" {
		return 0, "Incapacitated (" + member.ActiveIllness + ")"
	}

	var moraleMod int
	var moraleState string

	switch {
	case member.Morale <= 2:
		moraleMod = settings.CrewStates.MoraleModifier1to2
		moraleState = settings.MoraleStates["state_1_2"]
	case member.Morale <= 4:
		moraleMod = settings.CrewStates.MoraleModifier3to4
		moraleState = settings.MoraleStates["state_3_4"]
	case member.Morale <= 6:
		moraleMod = settings.CrewStates.MoraleModifier5to6
		moraleState = settings.MoraleStates["state_5_6"]
	case member.Morale <= 8:
		moraleMod = settings.CrewStates.MoraleModifier7to8
		moraleState = settings.MoraleStates["state_7_8"]
	default:
		moraleMod = settings.CrewStates.MoraleModifier9to10
		moraleState = settings.MoraleStates["state_9_10"]
	}

	effectiveSkill := member.SkillLevel + moraleMod
	if effectiveSkill < 0 {
		effectiveSkill = 0
	}

	return effectiveSkill, moraleState
}

// resolveTransitEvent checks if a random encounter occurs and processes the crew's attempt to mitigate it.
func resolveTransitEvent(state *game.State, s *ship.Ship, settings *config.GameSettings) []string {
	var logs []string

	// 1. Roll to see if an event triggers based on the global base chance percent.
	if mathrand.Intn(100) >= settings.Events.BaseChancePercent {
		return logs
	}

	// 2. Select a valid encounter from the cache based on the vessel's currently exposed tags.
	var validEncounterIDs []string
	validEncounterIDs = append(validEncounterIDs, cache.EncounterIDsByTag[""]...)

	for tag, ids := range cache.EncounterIDsByTag {
		if tag != "" && s.HasUnprotectedCargo(tag, settings.Commodities) {
			validEncounterIDs = append(validEncounterIDs, ids...)
		}
	}

	if len(validEncounterIDs) == 0 {
		return logs
	}

	encounterID := validEncounterIDs[mathrand.Intn(len(validEncounterIDs))]
	encounter := settings.Encounters[encounterID]

	logs = append(logs, fmt.Sprintf("    - [ALERT] Encounter: %s", encounter.Name))
	logs = append(logs, fmt.Sprintf("      >> %s", encounter.Description))

	// 3. Identify the most capable, healthy crew member to handle the crisis.
	bestSkill := -1
	bestName := "No Assigned Crew"
	bestState := "Unmanned"

	for i := range s.Roster {
		member := &s.Roster[i]
		if member.Role.ID == encounter.TargetRole || member.Specialty.SecondaryRole == encounter.TargetRole {
			effSkill, moraleState := calculateEffectiveSkill(member, settings)
			if effSkill > bestSkill {
				bestSkill = effSkill
				bestName = member.Name
				bestState = moraleState
			}
		}
	}

	// 4. Resolve the 1d20 check.
	roll := mathrand.Intn(20) + 1
	var total int

	if bestSkill == -1 {
		total = roll
		logs = append(logs, fmt.Sprintf("      >> No capable personnel available to respond! Roll: %d. DC: %d.", roll, encounter.DifficultyClass))
	} else {
		total = roll + bestSkill
		logs = append(logs, fmt.Sprintf("      >> %s (%s) attempts to resolve. Roll: %d + Skill %d = %d. DC: %d.", bestName, bestState, roll, bestSkill, total, encounter.DifficultyClass))
	}

	if total >= encounter.DifficultyClass {
		logs = append(logs, "      >> SUCCESS: The crisis was expertly averted.")
		return logs
	}

	logs = append(logs, "      >> FAILURE: The crew failed to mitigate the hazard.")

	// 5. Apply standard and specialized failure penalties.
	switch encounter.FailurePenalty {
	case "credit_loss":
		state.Player.Credits -= int64(encounter.PenaltyAmount)
		logs = append(logs, fmt.Sprintf("      >> Penalty: Corporate account debited %d C for damages and extortions.", encounter.PenaltyAmount))
	case "time_delay":
		s.ArrivalDay += encounter.PenaltyAmount
		logs = append(logs, fmt.Sprintf("      >> Penalty: Navigational delay extends arrival by %d day(s).", encounter.PenaltyAmount))
	case "module_damage":
		modName, success := s.DamageRandomModule()
		if success {
			logs = append(logs, fmt.Sprintf("      >> CRITICAL DAMAGE: %s was disabled and requires shipyard repair.", modName))
		} else {
			logs = append(logs, "      >> Hull took heavy impacts, but all remaining modules held together.")
		}
	case "inspection":
		logs = append(logs, resolveCustomsInspection(state, s, settings)...)
	case "hazardous_spill":
		logs = append(logs, "      >> CONTAINMENT BREACH: Unprotected volatile cargo has ruptured!")

		if mathrand.Intn(100) < settings.CargoRisks.HazardousShipDestroyPercent {
			s.Status = "Destroyed"
			logs = append(logs, "      >> CATASTROPHIC FAILURE: The vessel was completely vaporized in the resulting explosion.")
			return logs
		}

		if len(s.Contracts) > 0 {
			targetIdx := mathrand.Intn(len(s.Contracts))
			lostComm := settings.Commodities[s.Contracts[targetIdx].CommodityID]
			logs = append(logs, fmt.Sprintf("      >> CHAIN REACTION: The spill spread, destroying %d units of %s.", s.Contracts[targetIdx].Quantity, lostComm.Name))
			s.Contracts = append(s.Contracts[:targetIdx], s.Contracts[targetIdx+1:]...)

			modName, success := s.DamageRandomModule()
			if success {
				logs = append(logs, fmt.Sprintf("      >> STRUCTURAL FAILURE: The explosion disabled the %s.", modName))
			}
		}
	}

	// 6. Check for collateral cargo damage if Fragile goods are unprotected.
	if s.HasUnprotectedCargo("fragile", settings.Commodities) {
		if mathrand.Intn(100) < settings.CargoRisks.FragileDamageChancePercent {
			for i := range s.Contracts {
				comm := settings.Commodities[s.Contracts[i].CommodityID]
				if comm.HasTag("fragile") {
					loss := int(math.Ceil(float64(s.Contracts[i].Quantity) * (float64(settings.CargoRisks.FragileDamageQuantityPercent) / 100.0)))
					if loss < 1 {
						loss = 1
					}
					s.Contracts[i].Quantity -= loss
					logs = append(logs, fmt.Sprintf("      >> COLLATERAL: %d units of %s were shattered during the incident.", loss, comm.Name))
					break
				}
			}
		}
	}

	return logs
}

// resolveCustomsInspection processes the port authority check for illicit goods.
func resolveCustomsInspection(state *game.State, s *ship.Ship, settings *config.GameSettings) []string {
	var logs []string

	if !s.HasUnprotectedCargo("illicit", settings.Commodities) {
		logs = append(logs, "      >> Inspection passed. No unshielded contraband detected.")
		return logs
	}

	var remainingContracts []market.Contract
	totalFine := int64(0)

	for _, c := range s.Contracts {
		comm := settings.Commodities[c.CommodityID]
		if comm.HasTag("illicit") {
			if mathrand.Intn(100) < settings.CargoRisks.IllicitDiscoveryChancePercent {
				fine := int64(math.Round(float64(comm.BaseValue*int64(c.Quantity)) * settings.CargoRisks.IllicitFineMultiplier))
				totalFine += fine
				logs = append(logs, fmt.Sprintf("      >> CONTRABAND SEIZED: Authorities discovered %d units of %s. Fined %d C.", c.Quantity, comm.Name, fine))
			} else {
				remainingContracts = append(remainingContracts, c)
			}
		} else {
			remainingContracts = append(remainingContracts, c)
		}
	}

	s.Contracts = remainingContracts

	if totalFine > 0 {
		state.Player.Credits -= totalFine
		logs = append(logs, fmt.Sprintf("      >> Total Port Authority fines levied: %d C.", totalFine))
	} else {
		logs = append(logs, "      >> The inspector overlooked the contraband. You got lucky.")
	}

	return logs
}
