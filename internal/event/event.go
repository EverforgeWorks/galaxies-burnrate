// Package event defines the structures and data models for randomized logistical encounters.
package event

// EncounterTemplate defines the blueprint for a specific hazard or scenario.
type EncounterTemplate struct {
	ID              string `yaml:"id"`
	Name            string `yaml:"name"`
	Description     string `yaml:"description"`
	TriggerPhase    string `yaml:"trigger_phase"`
	TargetRole      string `yaml:"target_role"`
	RequiredTag     string `yaml:"required_tag"`
	DifficultyClass int    `yaml:"difficulty_class"`
	SuccessLog      string `yaml:"success_log"`
	FailureLog      string `yaml:"failure_log"`
	FailurePenalty  string `yaml:"failure_penalty"`
	PenaltyAmount   int    `yaml:"penalty_amount"`
}
