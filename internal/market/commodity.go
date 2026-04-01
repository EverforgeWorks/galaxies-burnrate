// Package market handles the economic simulation, including commodities and trading.
package market

// Commodity represents a distinct type of tradable good in the simulation.
type Commodity struct {
	ID          string   `yaml:"id"`
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	BaseValue   int64    `yaml:"base_value"`
	Tags        []string `yaml:"tags"`
}

// HasTag checks if the commodity possesses a specific risk or classification identifier.
func (c *Commodity) HasTag(tag string) bool {
	for _, t := range c.Tags {
		if t == tag {
			return true
		}
	}
	return false
}
