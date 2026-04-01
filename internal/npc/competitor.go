// Package npc handles non-player entities that interact with the game world and economy.
package npc

import (
	crand "crypto/rand"
	"fmt"
	"math/rand"
	"time"

	"galaxies-burnrate/internal/market"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

var competitorPrefixes = []string{"Red Sun", "Void", "Stellar", "Apex", "Horizon", "Nebula", "Vector", "Core"}
var competitorSuffixes = []string{"Freight", "Logistics", "Runners", "Transit", "Heavy", "Express", "Hauling"}

// Competitor represents a primitive AI ship acting independently in the universe.
type Competitor struct {
	ID             string
	Name           string
	Status         string // "Idle" or "In Transit"
	LocationID     string
	DestinationID  string
	ArrivalDay     int
	ActiveContract *market.Contract
}

func generateID() string {
	b := make([]byte, 8)
	_, _ = crand.Read(b)
	return fmt.Sprintf("%x", b)
}

func generateCompetitorName() string {
	prefix := competitorPrefixes[rng.Intn(len(competitorPrefixes))]
	suffix := competitorSuffixes[rng.Intn(len(competitorSuffixes))]
	return fmt.Sprintf("%s %s", prefix, suffix)
}

// NewCompetitor instantiates a new independent rival ship starting at a specific location.
func NewCompetitor(locationID string) Competitor {
	return Competitor{
		ID:         generateID(),
		Name:       generateCompetitorName(),
		Status:     "Idle",
		LocationID: locationID,
	}
}
