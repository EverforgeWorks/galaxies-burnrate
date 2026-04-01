// Package engine contains the core simulation mechanics and state mutation logic.
package engine

import (
	"fmt"
	"math"
	mathrand "math/rand"

	"galaxies-burnrate/internal/config"
	"galaxies-burnrate/internal/game"
	"galaxies-burnrate/internal/passenger"
)

// generatePassengerNameAndSpecies randomly selects a species from the configuration and constructs a valid name.
func generatePassengerNameAndSpecies(names map[string]config.SpeciesNames) (string, string) {
	var speciesList []string
	for sp := range names {
		speciesList = append(speciesList, sp)
	}

	selectedSpecies := "human"
	if len(speciesList) > 0 {
		selectedSpecies = speciesList[mathrand.Intn(len(speciesList))]
	}

	spData := names[selectedSpecies]

	first := "Unknown"
	if len(spData.FirstNames) > 0 {
		first = spData.FirstNames[mathrand.Intn(len(spData.FirstNames))]
	}

	last := "Traveler"
	if len(spData.LastNames) > 0 {
		last = spData.LastNames[mathrand.Intn(len(spData.LastNames))]
	}

	if selectedSpecies == "android" {
		return fmt.Sprintf("%s-%s", first, last), selectedSpecies
	}
	return fmt.Sprintf("%s %s", first, last), selectedSpecies
}

// GenerateInitialPassengers populates the universe with a baseline set of travelers, restricted by local social export traits.
func GenerateInitialPassengers(state *game.State, settings *config.GameSettings) {
	state.Passengers = make(map[string][]passenger.Passenger)

	for planetID, p := range state.Planets {
		socStatus := settings.SocialStatuses[p.SocialStatus]

		min := float64(settings.Passenger.StartingMinPerPlanet) * socStatus.PassengerExportVolumeMultiplier
		max := float64(settings.Passenger.StartingMaxPerPlanet) * socStatus.PassengerExportVolumeMultiplier

		count := int(math.Round(min))
		if int(math.Round(max)) > count {
			count = mathrand.Intn(int(math.Round(max))-count+1) + count
		}

		for i := 0; i < count; i++ {
			pass := createRandomPassenger(planetID, state, settings)
			if pass != nil {
				state.Passengers[planetID] = append(state.Passengers[planetID], *pass)
			}
		}
	}
}

// RefreshPassengers clears expired travelers and generates new ones up to the trait-modified terminal capacity limit.
func RefreshPassengers(state *game.State, settings *config.GameSettings) {
	for planetID, passengers := range state.Passengers {
		p := state.Planets[planetID]
		socStatus := settings.SocialStatuses[p.SocialStatus]

		var validPassengers []passenger.Passenger

		for _, pass := range passengers {
			if pass.ExpirationDay > state.CurrentDay {
				validPassengers = append(validPassengers, pass)
			}
		}

		limit := int(math.Round(float64(settings.Passenger.MaxPerPlanet) * socStatus.PassengerExportVolumeMultiplier))

		attempts := 0
		maxAttempts := limit * 3

		for len(validPassengers) < limit && attempts < maxAttempts {
			attempts++
			newPassenger := createRandomPassenger(planetID, state, settings)
			if newPassenger != nil {
				validPassengers = append(validPassengers, *newPassenger)
			}
		}

		state.Passengers[planetID] = validPassengers
	}
}

// createRandomPassenger generates a single traveler with scaled ticket payouts utilizing the engine cache.
func createRandomPassenger(originID string, state *game.State, settings *config.GameSettings) *passenger.Passenger {
	if len(cache.PlanetIDs) <= 1 || len(cache.PassengerClassIDs) == 0 {
		return nil
	}

	destID := cache.PlanetIDs[mathrand.Intn(len(cache.PlanetIDs))]
	if destID == originID {
		return nil
	}

	origin := state.Planets[originID]
	dest := state.Planets[destID]

	destSocStatus := settings.SocialStatuses[dest.SocialStatus]
	importChance := int(math.Round(100.0 * destSocStatus.PassengerImportVolumeMultiplier))
	if importChance < 100 && mathrand.Intn(100) >= importChance {
		return nil
	}

	classID := cache.PassengerClassIDs[mathrand.Intn(len(cache.PassengerClassIDs))]
	pClass := settings.PassengerClasses[classID]

	dx := float64(dest.X - origin.X)
	dy := float64(dest.Y - origin.Y)
	distance := int(math.Round(math.Sqrt(dx*dx + dy*dy)))
	if distance < 1 {
		distance = 1
	}

	baseFare := float64(pClass.BasePayout)

	originEco := settings.EconomicStatuses[origin.EconomicStatus]
	destEco := settings.EconomicStatuses[dest.EconomicStatus]

	// Apply base route multiplier.
	basePayout := baseFare * float64(distance) * settings.Multipliers.PassengerPayout * originEco.ExportPassengerMultiplier * destEco.ImportPassengerMultiplier

	// Apply specialized passenger risk/luxury multipliers.
	if pClass.HasTag("vip") {
		basePayout *= settings.Multipliers.VipPayout
	}
	if pClass.HasTag("fugitive") || pClass.HasTag("illicit") {
		basePayout *= settings.Multipliers.FugitivePayout
	}

	payout := int64(math.Round(basePayout))

	bufferMin := settings.Passenger.ExpirationBufferMin
	bufferMax := settings.Passenger.ExpirationBufferMax
	expires := state.CurrentDay + distance + mathrand.Intn(bufferMax-bufferMin+1) + bufferMin

	// Generate the name and species in the engine to break the import cycle
	baseName, generatedSpecies := generatePassengerNameAndSpecies(settings.Names)

	p := passenger.New(pClass, originID, destID, payout, expires, baseName, generatedSpecies)
	return &p
}
