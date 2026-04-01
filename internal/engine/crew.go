// Package engine contains the core simulation mechanics and state mutation logic.
package engine

import (
	"crypto/rand"
	"fmt"
	"math"
	mathrand "math/rand"

	"galaxies-burnrate/internal/config"
	"galaxies-burnrate/internal/crew"
	"galaxies-burnrate/internal/game"
)

// generateCrewID creates a secure, unique hexadecimal identifier.
func generateCrewID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// generateCrewNameAndSpecies dynamically selects a species and generates a name utilizing the configuration banks.
func generateCrewNameAndSpecies(originID string, names map[string]config.SpeciesNames) (string, string) {
	var validSpecies []string
	for sp := range names {
		if sp != "android" {
			validSpecies = append(validSpecies, sp)
		}
	}

	selectedSpecies := "human"
	if originID == "android" {
		selectedSpecies = "android"
	} else if len(validSpecies) > 0 {
		selectedSpecies = validSpecies[mathrand.Intn(len(validSpecies))]
	}

	spData := names[selectedSpecies]

	first := "Unknown"
	if len(spData.FirstNames) > 0 {
		first = spData.FirstNames[mathrand.Intn(len(spData.FirstNames))]
	}

	last := "Entity"
	if len(spData.LastNames) > 0 {
		last = spData.LastNames[mathrand.Intn(len(spData.LastNames))]
	}

	if selectedSpecies == "android" {
		return fmt.Sprintf("%s-%s", first, last), selectedSpecies
	}
	return fmt.Sprintf("%s %s", first, last), selectedSpecies
}

// GenerateInitialCrew populates the hiring boards at cantinas across the galaxy.
func GenerateInitialCrew(state *game.State, settings *config.GameSettings) {
	state.AvailableCrew = make(map[string][]crew.CrewMember)

	for planetID, p := range state.Planets {
		if !p.HasTag("cantina") {
			continue
		}

		min := settings.Cantina.MinCrew
		max := settings.Cantina.MaxCrew
		count := min
		if max > min {
			count = mathrand.Intn(max-min+1) + min
		}

		for i := 0; i < count; i++ {
			member := createRandomCrew(settings)
			if member != nil {
				state.AvailableCrew[planetID] = append(state.AvailableCrew[planetID], *member)
			}
		}
	}
}

// RefreshCrew cycles the hiring board to simulate a living galaxy where independent ships hire crew.
func RefreshCrew(state *game.State, settings *config.GameSettings) {
	for planetID, p := range state.Planets {
		if !p.HasTag("cantina") {
			continue
		}

		current := state.AvailableCrew[planetID]
		if len(current) > 0 {
			keep := mathrand.Intn(len(current))
			current = current[:keep]
		}

		limit := settings.Cantina.MaxCrew
		for len(current) < limit {
			member := createRandomCrew(settings)
			if member != nil {
				current = append(current, *member)
			}
		}

		state.AvailableCrew[planetID] = current
	}
}

// createRandomCrew builds a hirable entity utilizing cached templates to map traits and hiring fees.
func createRandomCrew(settings *config.GameSettings) *crew.CrewMember {
	if len(cache.CrewRoles) == 0 {
		return nil
	}
	role := cache.CrewRoles[mathrand.Intn(len(cache.CrewRoles))]
	origin := cache.Origins[mathrand.Intn(len(cache.Origins))]
	specialty := cache.Specialties[mathrand.Intn(len(cache.Specialties))]

	skill := 1
	if origin.ForceMaxSkill {
		skill = settings.Cantina.SkillMax
	} else {
		totalWeight := 0
		for _, w := range settings.Cantina.SkillWeights {
			totalWeight += w
		}
		roll := mathrand.Intn(totalWeight)
		currentWeight := 0
		for i, w := range settings.Cantina.SkillWeights {
			currentWeight += w
			if roll < currentWeight {
				skill = i + 1
				break
			}
		}
	}

	baseDaily := float64(role.BaseSalary) * (1.0 + (float64(skill-1) * 0.5))
	finalDaily := int64(math.Round(baseDaily * settings.Multipliers.CrewSalary * origin.SalaryModifier))

	hireCost := int64(math.Round(float64(role.BaseSalary) * float64(skill) * origin.HireMultiplier))

	generatedName, generatedSpecies := generateCrewNameAndSpecies(origin.ID, settings.Names)

	member := crew.CrewMember{
		ID:            generateCrewID(),
		Name:          generatedName,
		Species:       generatedSpecies,
		Role:          role,
		Origin:        origin,
		Specialty:     specialty,
		SkillLevel:    skill,
		DailySalary:   finalDaily,
		HireCost:      hireCost,
		Morale:        settings.CrewStates.MoraleStarting,
		ActiveIllness: "",
		IllnessDays:   0,
		OnVacation:    false,
	}

	return &member
}
