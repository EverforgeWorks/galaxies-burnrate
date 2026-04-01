// Package crew manages the personnel, their roles, and their specific skills within the simulation.
package crew

// RoleTemplate defines the baseline responsibilities and compensation for a specific ship station.
type RoleTemplate struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	BaseSalary  int64  `yaml:"base_salary"`
}

// Origin defines the background of the crew member, altering their cost and physical traits.
type Origin struct {
	ID               string  `yaml:"id"`
	Name             string  `yaml:"name"`
	Description      string  `yaml:"description"`
	SalaryModifier   float64 `yaml:"salary_modifier"`
	HireMultiplier   float64 `yaml:"hire_multiplier"`
	ForceMaxSkill    bool    `yaml:"force_max_skill"`
	ImmuneToSickness bool    `yaml:"immune_to_sickness"`
}

// Specialty provides unique utility or secondary role coverage.
type Specialty struct {
	ID               string  `yaml:"id"`
	Name             string  `yaml:"name"`
	SecondaryRole    string  `yaml:"secondary_role"`
	ShipyardDiscount float64 `yaml:"shipyard_discount"`
}

// CrewMember represents an instanced, hirable entity with active states.
type CrewMember struct {
	ID             string
	Name           string
	Species        string // Tracks biological or artificial lineage for name mapping and flavor.
	Role           RoleTemplate
	Origin         Origin
	Specialty      Specialty
	SkillLevel     int
	DailySalary    int64
	HireCost       int64
	Morale         int    // Operates on a 1-10 scale influencing skill checks.
	ActiveIllness  string // The flavor string of the disease. Empty if healthy.
	IllnessDays    int    // Remaining days until the illness passes naturally.
	OnVacation     bool   // Flags if the member is left planetside for shore leave.
	RaiseRequested bool   // Flags if the crew member is currently demanding a raise.
	RaiseRefused   bool   // Flags if a previous raise request was denied, increasing desertion risk.
}
