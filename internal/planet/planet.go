// Package planet manages the celestial bodies and locations within the simulation.
package planet

import "fmt"

// ShipyardInventory defines the specific chassis and modules available for purchase at a given location.
type ShipyardInventory struct {
	Chassis []string `yaml:"chassis"`
	Modules []string `yaml:"modules"`
}

// RefuelInfo defines the availability and cost of fuel at a specific location.
type RefuelInfo struct {
	Cost int `yaml:"cost"`
}

// Planet represents a fixed, navigable location within the game universe.
type Planet struct {
	ID              string            `yaml:"id"`
	Name            string            `yaml:"name"`
	EconomicStatus  string            `yaml:"economic_status"`
	SocialStatus    string            `yaml:"social_status"`
	PoliticalStatus string            `yaml:"political_status"`
	DockingFee      int               `yaml:"docking_fee"`
	X               int               `yaml:"x"`
	Y               int               `yaml:"y"`
	Tags            []string          `yaml:"tags"`
	Shipyard        ShipyardInventory `yaml:"shipyard"`
	Refuel          RefuelInfo        `yaml:"refuel"`
}

// HasTag checks if the planet contains a specific categorical identifier.
func (p *Planet) HasTag(tag string) bool {
	for _, t := range p.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// Descriptor generates the 3-word title representing the planet's combined traits.
func (p *Planet) Descriptor(eco map[string]EconomicStatus, soc map[string]SocialStatus, pol map[string]PoliticalStatus) string {
	ecoName := "Unknown"
	if e, ok := eco[p.EconomicStatus]; ok {
		ecoName = e.Name
	}

	socName := "Unknown"
	if s, ok := soc[p.SocialStatus]; ok {
		socName = s.Name
	}

	polName := "Unknown"
	if po, ok := pol[p.PoliticalStatus]; ok {
		polName = po.Name
	}

	return fmt.Sprintf("%s %s %s", ecoName, socName, polName)
}
