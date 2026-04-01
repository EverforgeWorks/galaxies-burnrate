// Package game houses the core simulation state and loop management.
package game

import (
	"galaxies-burnrate/internal/crew"
	"galaxies-burnrate/internal/market"
	"galaxies-burnrate/internal/npc"
	"galaxies-burnrate/internal/passenger"
	"galaxies-burnrate/internal/planet"
	"galaxies-burnrate/internal/player"
)

// State acts as the central repository for all active game data in memory.
type State struct {
	Player        *player.Player
	Planets       map[string]planet.Planet
	Contracts     map[string][]market.Contract
	Passengers    map[string][]passenger.Passenger
	AvailableCrew map[string][]crew.CrewMember
	Competitors   []*npc.Competitor
	CurrentDay    int
}

// NewState initializes a clean, empty game state starting at Day 1.
func NewState() *State {
	return &State{
		Planets:       make(map[string]planet.Planet),
		Contracts:     make(map[string][]market.Contract),
		Passengers:    make(map[string][]passenger.Passenger),
		AvailableCrew: make(map[string][]crew.CrewMember),
		Competitors:   make([]*npc.Competitor, 0),
		CurrentDay:    1,
	}
}
