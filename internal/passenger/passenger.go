// Package passenger manages the personnel, classes, and logic for transporting individuals.
package passenger

import (
	crand "crypto/rand"
	"fmt"
)

// PassengerClass defines the foundational blueprint and baseline economics for a traveler.
type PassengerClass struct {
	ID          string   `yaml:"id"`
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	BasePayout  int64    `yaml:"base_payout"`
	Tags        []string `yaml:"tags"`
}

// HasTag checks if the passenger class possesses a specific mechanical trait.
func (c PassengerClass) HasTag(tag string) bool {
	for _, t := range c.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// Passenger represents an instanced individual looking for transit.
type Passenger struct {
	ID            string
	Name          string
	Species       string
	Class         PassengerClass
	OriginID      string
	DestinationID string
	Payout        int64
	ExpirationDay int
}

// generateID creates a pseudo-random string for unique entity identification.
func generateID() string {
	b := make([]byte, 8)
	_, _ = crand.Read(b)
	return fmt.Sprintf("%x", b)
}

// New instantiates a new Passenger with a unique ID and dynamically assigned name.
func New(class PassengerClass, origin, destination string, payout int64, expires int, baseName string, species string) Passenger {
	finalName := "Commuter " + baseName

	// Apply flavor prefix based on specialized tags
	if class.HasTag("vip") {
		finalName = "Executive " + baseName
	} else if class.HasTag("fugitive") {
		finalName = "Alias " + baseName
	}

	return Passenger{
		ID:            generateID(),
		Name:          finalName,
		Species:       species,
		Class:         class,
		OriginID:      origin,
		DestinationID: destination,
		Payout:        payout,
		ExpirationDay: expires,
	}
}