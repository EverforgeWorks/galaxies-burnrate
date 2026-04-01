// Package market handles the economic simulation, including commodities and trading.
package market

import (
	"crypto/rand"
	"fmt"
)

// Contract represents a logistical delivery agreement available to the player.
type Contract struct {
	ID            string
	OriginID      string
	DestinationID string
	CommodityID   string
	Quantity      int
	Payout        int64
	ExpirationDay int
}

// NewContract instantiates a structured delivery agreement with a unique identifier.
func NewContract(origin, destination, commodity string, qty int, payout int64, expires int) Contract {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return Contract{
		ID:            fmt.Sprintf("%x", b),
		OriginID:      origin,
		DestinationID: destination,
		CommodityID:   commodity,
		Quantity:      qty,
		Payout:        payout,
		ExpirationDay: expires,
	}
}
