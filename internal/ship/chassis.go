// Package ship manages vessel architecture, module loadouts, and stat aggregation.
package ship

// ChassisTemplate defines the base hull characteristics and maximum capacities for a vessel.
type ChassisTemplate struct {
	ID              string   `yaml:"id"`
	Name            string   `yaml:"name"`
	Description     string   `yaml:"description"`
	BaseValue       int64    `yaml:"base_value"`
	MaxModules      int      `yaml:"max_modules"`
	BaseCargo       int      `yaml:"base_cargo"`
	BaseFuel        int      `yaml:"base_fuel"`
	BaseSpeed       int      `yaml:"base_speed"`
	BaseCabins      int      `yaml:"base_cabins"`
	BaseQuarters    int      `yaml:"base_quarters"`
	FuelConsumption int      `yaml:"fuel_consumption"`
	DefaultModules  []string `yaml:"default_modules"`
}
