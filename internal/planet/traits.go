// Package planet manages the celestial bodies and locations within the simulation.
package planet

// EconomicStatus represents the financial health of a planet and its effect on trade values.
type EconomicStatus struct {
	ID                        string  `yaml:"id"`
	Name                      string  `yaml:"name"`
	ExportContractMultiplier  float64 `yaml:"export_contract_multiplier"`
	ImportContractMultiplier  float64 `yaml:"import_contract_multiplier"`
	ExportPassengerMultiplier float64 `yaml:"export_passenger_multiplier"`
	ImportPassengerMultiplier float64 `yaml:"import_passenger_multiplier"`
	FuelCostMultiplier        float64 `yaml:"fuel_cost_multiplier"`
	ShipyardCostMultiplier    float64 `yaml:"shipyard_cost_multiplier"`
	RepairCostMultiplier      float64 `yaml:"repair_cost_multiplier"`
}

// SocialStatus represents the societal conditions of a planet and its effect on population movement and industry volume.
type SocialStatus struct {
	ID                              string  `yaml:"id"`
	Name                            string  `yaml:"name"`
	ContractVolumeMultiplier        float64 `yaml:"contract_volume_multiplier"`
	PassengerExportVolumeMultiplier float64 `yaml:"passenger_export_volume_multiplier"`
	PassengerImportVolumeMultiplier float64 `yaml:"passenger_import_volume_multiplier"`
}

// PoliticalStatus represents the governmental structure of a planet and its effect on law enforcement and illicit trade.
type PoliticalStatus struct {
	ID                            string  `yaml:"id"`
	Name                          string  `yaml:"name"`
	DockingFeeMultiplier          float64 `yaml:"docking_fee_multiplier"`
	InspectionChanceModifier      int     `yaml:"inspection_chance_modifier"`
	IllicitExportChanceMultiplier float64 `yaml:"illicit_export_chance_multiplier"`
	IllicitImportPayoutMultiplier float64 `yaml:"illicit_import_payout_multiplier"`
}
