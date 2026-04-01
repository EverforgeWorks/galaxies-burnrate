// Package ship defines the physical vessels and modules in the simulation.
package ship

// ModuleTemplate defines a standardized installable ship component.
type ModuleTemplate struct {
	ID                     string `yaml:"id"`
	Name                   string `yaml:"name"`
	Description            string `yaml:"description"`
	Value                  int64  `yaml:"value"`
	SpeedBonus             int    `yaml:"speed_bonus"`
	FuelBonus              int    `yaml:"fuel_bonus"`
	FuelConsumption        int    `yaml:"fuel_consumption"`
	QuartersBonus          int    `yaml:"quarters_bonus"`
	PassengerSpace         int    `yaml:"passenger_space"`
	VipSpace               int    `yaml:"vip_space"`
	CargoBonus             int    `yaml:"cargo_bonus"`
	HazardousProtection    int    `yaml:"hazardous_protection"`
	FragileProtection      int    `yaml:"fragile_protection"`
	IllicitConcealment     int    `yaml:"illicit_concealment"`
	PerishablePreservation int    `yaml:"perishable_preservation"`
}

// InstalledModule tracks the active status of a module.
type InstalledModule struct {
	Template  ModuleTemplate
	IsDamaged bool
}
