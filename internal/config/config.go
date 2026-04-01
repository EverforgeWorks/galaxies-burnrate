// Package config handles the loading and parsing of game configuration files.
package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"galaxies-burnrate/internal/crew"
	"galaxies-burnrate/internal/event"
	"galaxies-burnrate/internal/market"
	"galaxies-burnrate/internal/passenger"
	"galaxies-burnrate/internal/planet"
	"galaxies-burnrate/internal/ship"
)

// GameSettings represents the master configuration struct for the simulation.
type GameSettings struct {
	Player            PlayerSettings                      `yaml:"player"`
	Contracts         ContractSettings                    `yaml:"contracts"`
	Crew              CrewSettings                        `yaml:"crew"`
	Cantina           CantinaSettings                     `yaml:"cantina"`
	CrewStates        CrewStateSettings                   `yaml:"crew_states"`
	Passenger         PassengerSettings                   `yaml:"passenger"`
	Events            EventSettings                       `yaml:"events"`
	CargoRisks        CargoRiskSettings                   `yaml:"cargo_risks"`
	Competitors       CompetitorSettings                  `yaml:"competitors"`
	Multipliers       Multipliers                         `yaml:"multipliers"`
	Planets           map[string]planet.Planet            `yaml:"-"`
	EconomicStatuses  map[string]planet.EconomicStatus    `yaml:"-"`
	SocialStatuses    map[string]planet.SocialStatus      `yaml:"-"`
	PoliticalStatuses map[string]planet.PoliticalStatus   `yaml:"-"`
	Chassis           map[string]ship.ChassisTemplate     `yaml:"-"`
	Modules           map[string]ship.ModuleTemplate      `yaml:"-"`
	Commodities       map[string]market.Commodity         `yaml:"-"`
	CrewRoles         map[string]crew.RoleTemplate        `yaml:"-"`
	Origins           map[string]crew.Origin              `yaml:"-"`
	Specialties       map[string]crew.Specialty           `yaml:"-"`
	PassengerClasses  map[string]passenger.PassengerClass `yaml:"-"`
	Names             map[string]SpeciesNames             `yaml:"-"`
	Encounters        map[string]event.EncounterTemplate  `yaml:"-"`
	IllnessNames      []string                            `yaml:"-"`
	MoraleStates      map[string]string                   `yaml:"-"`
}

// SpeciesNames holds the generated name banks for a given biological race or construct.
type SpeciesNames struct {
	FirstNames []string `yaml:"first_names"`
	LastNames  []string `yaml:"last_names"`
}

type PlayerSettings struct {
	StartingCredits int64  `yaml:"starting_credits"`
	StartingPlanet  string `yaml:"starting_planet"`
}

type ContractSettings struct {
	StartingMinPerPlanet int     `yaml:"starting_min_per_planet"`
	StartingMaxPerPlanet int     `yaml:"starting_max_per_planet"`
	MaxPerPlanet         int     `yaml:"max_per_planet"`
	QtyMin               int     `yaml:"qty_min"`
	QtyMax               int     `yaml:"qty_max"`
	DistanceModifier     float64 `yaml:"distance_modifier"`
	VarianceMin          float64 `yaml:"variance_min"`
	VarianceMax          float64 `yaml:"variance_max"`
	ExpirationBufferMin  int     `yaml:"expiration_buffer_min"`
	ExpirationBufferMax  int     `yaml:"expiration_buffer_max"`
}

type CrewSettings struct {
	PayPeriodDays int `yaml:"pay_period_days"`
}

type CrewStateSettings struct {
	MoraleStarting                  int `yaml:"morale_starting"`
	MoraleAbandonThreshold          int `yaml:"morale_abandon_threshold"`
	RaiseRequestChancePercent       int `yaml:"raise_request_chance_percent"`
	RaiseSalaryBumpPercent          int `yaml:"raise_salary_bump_percent"`
	IllnessContractionChancePercent int `yaml:"illness_contraction_chance_percent"`
	IllnessMinDays                  int `yaml:"illness_min_days"`
	IllnessMaxDays                  int `yaml:"illness_max_days"`
	IllnessBaseDC                   int `yaml:"illness_base_dc"`
	MoraleModifier1to2              int `yaml:"morale_modifier_1_to_2"`
	MoraleModifier3to4              int `yaml:"morale_modifier_3_to_4"`
	MoraleModifier5to6              int `yaml:"morale_modifier_5_to_6"`
	MoraleModifier7to8              int `yaml:"morale_modifier_7_to_8"`
	MoraleModifier9to10             int `yaml:"morale_modifier_9_to_10"`
}

type CantinaSettings struct {
	MinCrew      int   `yaml:"min_crew"`
	MaxCrew      int   `yaml:"max_crew"`
	SkillMin     int   `yaml:"skill_min"`
	SkillMax     int   `yaml:"skill_max"`
	SkillWeights []int `yaml:"skill_weights"`
}

type PassengerSettings struct {
	StartingMinPerPlanet int `yaml:"starting_min_per_planet"`
	StartingMaxPerPlanet int `yaml:"starting_max_per_planet"`
	MaxPerPlanet         int `yaml:"max_per_planet"`
	ExpirationBufferMin  int `yaml:"expiration_buffer_min"`
	ExpirationBufferMax  int `yaml:"expiration_buffer_max"`
}

type EventSettings struct {
	BaseChancePercent int `yaml:"base_chance_percent"`
}

type CargoRiskSettings struct {
	HazardousChainChancePercent   int     `yaml:"hazardous_chain_chance_percent"`
	HazardousShipDestroyPercent   int     `yaml:"hazardous_ship_destroy_percent"`
	IllicitDiscoveryChancePercent int     `yaml:"illicit_discovery_chance_percent"`
	IllicitFineMultiplier         float64 `yaml:"illicit_fine_multiplier"`
	FragileDamageChancePercent    int     `yaml:"fragile_damage_chance_percent"`
	FragileDamageQuantityPercent  int     `yaml:"fragile_damage_quantity_percent"`
	LateFeePercentPerDay          int     `yaml:"late_fee_percent_per_day"`
}

type CompetitorSettings struct {
	MinActive int `yaml:"min_active"`
	MaxActive int `yaml:"max_active"`
}

type Multipliers struct {
	ModuleRefundRate       float64 `yaml:"module_refund_rate"`
	ModuleRepairRate       float64 `yaml:"module_repair_rate"`
	CrewSalary             float64 `yaml:"crew_salary"`
	ContractPayout         float64 `yaml:"contract_payout"`
	FuelCost               float64 `yaml:"fuel_cost"`
	PassengerPayout        float64 `yaml:"passenger_payout"`
	VipPayout              float64 `yaml:"vip_payout"`
	FugitivePayout         float64 `yaml:"fugitive_payout"`
	VipDelayPenaltyPercent int     `yaml:"vip_delay_penalty_percent"`
	HazardousPayout        float64 `yaml:"hazardous_payout"`
	FragilePayout          float64 `yaml:"fragile_payout"`
	IllicitPayout          float64 `yaml:"illicit_payout"`
	PerishablePayout       float64 `yaml:"perishable_payout"`
}

type traitsYAML struct {
	EconomicStatuses  map[string]planet.EconomicStatus  `yaml:"economic_statuses"`
	SocialStatuses    map[string]planet.SocialStatus    `yaml:"social_statuses"`
	PoliticalStatuses map[string]planet.PoliticalStatus `yaml:"political_statuses"`
}

type shipsYAML struct {
	Chassis map[string]ship.ChassisTemplate `yaml:"chassis"`
}

type modulesYAML struct {
	Modules map[string]ship.ModuleTemplate `yaml:"modules"`
}

type planetsYAML struct {
	Planets map[string]planet.Planet `yaml:"planets"`
}

type commoditiesYAML struct {
	Commodities map[string]market.Commodity `yaml:"commodities"`
}

type crewYAML struct {
	Roles        map[string]crew.RoleTemplate `yaml:"roles"`
	Origins      map[string]crew.Origin       `yaml:"origins"`
	Specialties  map[string]crew.Specialty    `yaml:"specialties"`
	IllnessNames []string                     `yaml:"illness_names"`
	MoraleStates map[string]string            `yaml:"morale_states"`
}

type passengersYAML struct {
	Classes map[string]passenger.PassengerClass `yaml:"classes"`
}

type namesYAML struct {
	Species map[string]SpeciesNames `yaml:"species"`
}

type eventsYAML struct {
	Events map[string]event.EncounterTemplate `yaml:"events"`
}

// Load parses all YAML configuration files from the specified directory and returns the aggregated GameSettings.
func Load(configDir string) (*GameSettings, error) {
	settingsData, err := os.ReadFile(filepath.Join(configDir, "settings.yaml"))
	if err != nil {
		return nil, err
	}

	var settings GameSettings
	if err := yaml.Unmarshal(settingsData, &settings); err != nil {
		return nil, err
	}

	traitsData, err := os.ReadFile(filepath.Join(configDir, "planet_traits.yaml"))
	if err != nil {
		return nil, err
	}

	var ty traitsYAML
	if err := yaml.Unmarshal(traitsData, &ty); err != nil {
		return nil, err
	}
	settings.EconomicStatuses = ty.EconomicStatuses
	settings.SocialStatuses = ty.SocialStatuses
	settings.PoliticalStatuses = ty.PoliticalStatuses

	shipsData, err := os.ReadFile(filepath.Join(configDir, "ships.yaml"))
	if err != nil {
		return nil, err
	}

	var sy shipsYAML
	if err := yaml.Unmarshal(shipsData, &sy); err != nil {
		return nil, err
	}
	settings.Chassis = sy.Chassis

	modulesData, err := os.ReadFile(filepath.Join(configDir, "modules.yaml"))
	if err != nil {
		return nil, err
	}

	var my modulesYAML
	if err := yaml.Unmarshal(modulesData, &my); err != nil {
		return nil, err
	}
	settings.Modules = my.Modules

	planetsData, err := os.ReadFile(filepath.Join(configDir, "planets.yaml"))
	if err != nil {
		return nil, err
	}

	var py planetsYAML
	if err := yaml.Unmarshal(planetsData, &py); err != nil {
		return nil, err
	}
	settings.Planets = py.Planets

	commoditiesData, err := os.ReadFile(filepath.Join(configDir, "commodities.yaml"))
	if err != nil {
		return nil, err
	}

	var cy commoditiesYAML
	if err := yaml.Unmarshal(commoditiesData, &cy); err != nil {
		return nil, err
	}
	settings.Commodities = cy.Commodities

	crewData, err := os.ReadFile(filepath.Join(configDir, "crew.yaml"))
	if err != nil {
		return nil, err
	}

	var crY crewYAML
	if err := yaml.Unmarshal(crewData, &crY); err != nil {
		return nil, err
	}
	settings.CrewRoles = crY.Roles
	settings.Origins = crY.Origins
	settings.Specialties = crY.Specialties
	settings.IllnessNames = crY.IllnessNames
	settings.MoraleStates = crY.MoraleStates

	passengersData, err := os.ReadFile(filepath.Join(configDir, "passengers.yaml"))
	if err != nil {
		return nil, err
	}

	var pY passengersYAML
	if err := yaml.Unmarshal(passengersData, &pY); err != nil {
		return nil, err
	}
	settings.PassengerClasses = pY.Classes

	namesData, err := os.ReadFile(filepath.Join(configDir, "names.yaml"))
	if err != nil {
		return nil, err
	}

	var nY namesYAML
	if err := yaml.Unmarshal(namesData, &nY); err != nil {
		return nil, err
	}
	settings.Names = nY.Species

	eventsData, err := os.ReadFile(filepath.Join(configDir, "events.yaml"))
	if err != nil {
		return nil, err
	}

	var eV eventsYAML
	if err := yaml.Unmarshal(eventsData, &eV); err != nil {
		return nil, err
	}
	settings.Encounters = eV.Events

	return &settings, nil
}
