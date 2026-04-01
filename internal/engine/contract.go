// Package engine contains the core simulation mechanics and state mutation logic.
package engine

import (
	"math"
	"math/rand"
	"time"

	"galaxies-burnrate/internal/config"
	"galaxies-burnrate/internal/game"
	"galaxies-burnrate/internal/market"
)

// rng establishes a package-level pseudo-random number generator to avoid relying on deprecated global seeds.
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// GenerateInitialContracts populates the universe with a baseline set of delivery contracts, adjusted by planetary social traits.
func GenerateInitialContracts(state *game.State, settings *config.GameSettings) {
	state.Contracts = make(map[string][]market.Contract)

	for planetID, p := range state.Planets {
		socStatus := settings.SocialStatuses[p.SocialStatus]

		min := float64(settings.Contracts.StartingMinPerPlanet) * socStatus.ContractVolumeMultiplier
		max := float64(settings.Contracts.StartingMaxPerPlanet) * socStatus.ContractVolumeMultiplier

		count := int(math.Round(min))
		if int(math.Round(max)) > count {
			count = rng.Intn(int(math.Round(max))-count+1) + count
		}

		for i := 0; i < count; i++ {
			contract := createRandomContract(planetID, state, settings)
			if contract != nil {
				state.Contracts[planetID] = append(state.Contracts[planetID], *contract)
			}
		}
	}
}

// RefreshContracts clears expired contracts and generates new ones up to the trait-modified limit.
// A nil return from the generator breaks the loop, simulating economic stagnation if trade routes are heavily restricted.
func RefreshContracts(state *game.State, settings *config.GameSettings) {
	for planetID, contracts := range state.Contracts {
		p := state.Planets[planetID]
		socStatus := settings.SocialStatuses[p.SocialStatus]

		var validContracts []market.Contract

		for _, c := range contracts {
			if c.ExpirationDay > state.CurrentDay {
				validContracts = append(validContracts, c)
			}
		}

		limit := int(math.Round(float64(settings.Contracts.MaxPerPlanet) * socStatus.ContractVolumeMultiplier))

		for len(validContracts) < limit {
			newContract := createRandomContract(planetID, state, settings)
			if newContract != nil {
				validContracts = append(validContracts, *newContract)
			} else {
				break
			}
		}

		state.Contracts[planetID] = validContracts
	}
}

// createRandomContract generates a single logistical delivery agreement with scaled payouts based on distance, traits, and specialized cargo tags.
func createRandomContract(originID string, state *game.State, settings *config.GameSettings) *market.Contract {
	var destIDs []string
	for id := range state.Planets {
		if id != originID {
			destIDs = append(destIDs, id)
		}
	}
	if len(destIDs) == 0 {
		return nil
	}

	destID := destIDs[rng.Intn(len(destIDs))]
	origin := state.Planets[originID]
	dest := state.Planets[destID]

	var commIDs []string
	for id := range settings.Commodities {
		commIDs = append(commIDs, id)
	}
	if len(commIDs) == 0 {
		return nil
	}

	commID := commIDs[rng.Intn(len(commIDs))]
	commodity := settings.Commodities[commID]

	// Apply political constraints to determine if illicit goods can be successfully exported from the origin.
	if commodity.HasTag("illicit") {
		polStatus := settings.PoliticalStatuses[origin.PoliticalStatus]
		chance := int(math.Round(100.0 * polStatus.IllicitExportChanceMultiplier))
		if chance < 100 && rng.Intn(100) >= chance {
			return nil
		}
	}

	dx := float64(dest.X - origin.X)
	dy := float64(dest.Y - origin.Y)
	distance := int(math.Round(math.Sqrt(dx*dx + dy*dy)))
	if distance < 1 {
		distance = 1
	}

	qtyMin := settings.Contracts.QtyMin
	qtyMax := settings.Contracts.QtyMax
	qty := rng.Intn(qtyMax-qtyMin+1) + qtyMin

	baseValue := float64(commodity.BaseValue * int64(qty))
	distanceModifier := 1.0 + (float64(distance) * settings.Contracts.DistanceModifier)

	varianceMin := settings.Contracts.VarianceMin
	varianceMax := settings.Contracts.VarianceMax
	variance := varianceMin + (rng.Float64() * (varianceMax - varianceMin))

	tagMultiplier := 1.0
	if commodity.HasTag("hazardous") {
		tagMultiplier *= settings.Multipliers.HazardousPayout
	}
	if commodity.HasTag("fragile") {
		tagMultiplier *= settings.Multipliers.FragilePayout
	}
	if commodity.HasTag("illicit") {
		tagMultiplier *= settings.Multipliers.IllicitPayout
		destPolStatus := settings.PoliticalStatuses[dest.PoliticalStatus]
		tagMultiplier *= destPolStatus.IllicitImportPayoutMultiplier
	}
	if commodity.HasTag("perishable") {
		tagMultiplier *= settings.Multipliers.PerishablePayout
	}

	originEco := settings.EconomicStatuses[origin.EconomicStatus]
	destEco := settings.EconomicStatuses[dest.EconomicStatus]

	payout := int64(math.Round(baseValue * distanceModifier * variance * settings.Multipliers.ContractPayout * tagMultiplier * originEco.ExportContractMultiplier * destEco.ImportContractMultiplier))

	bufferMin := settings.Contracts.ExpirationBufferMin
	bufferMax := settings.Contracts.ExpirationBufferMax
	expires := state.CurrentDay + distance + rng.Intn(bufferMax-bufferMin+1) + bufferMin

	contract := market.NewContract(originID, destID, commID, qty, payout, expires)
	return &contract
}
