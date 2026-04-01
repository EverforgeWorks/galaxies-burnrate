// Package player manages the entity data and progression for the user.
package player

import (
	"errors"
	"galaxies-burnrate/internal/ship"
)

// Player represents the logistics manager controlling the corporate entity.
type Player struct {
	Name        string
	CompanyName string
	Credits     int64
	Fleet       map[string]*ship.Ship
}

// New initializes and returns a new Player instance using provided starting values.
func New(name, companyName string, startingCredits int64) *Player {
	return &Player{
		Name:        name,
		CompanyName: companyName,
		Credits:     startingCredits,
		Fleet:       make(map[string]*ship.Ship),
	}
}

// AddShip integrates a new vessel into the player's managed fleet.
func (p *Player) AddShip(s *ship.Ship) error {
	if s == nil {
		return errors.New("cannot add a nil ship to the fleet")
	}

	if _, exists := p.Fleet[s.ID]; exists {
		return errors.New("ship with this ID already exists in the fleet")
	}

	p.Fleet[s.ID] = s
	return nil
}
