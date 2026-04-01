// Package cli provides the command-line interface and menu handling for the game.
package cli

import (
	"fmt"
	"strings"

	"galaxies-burnrate/internal/engine"
	"galaxies-burnrate/internal/planet"
	"galaxies-burnrate/internal/player"
	"galaxies-burnrate/internal/ship"
)

// runMainMenu renders the main menu options and processes the initial application flow.
func (c *CLI) runMainMenu() bool {
	fmt.Println("\nMain Menu:")
	fmt.Println("[1] (N)ew Game")
	fmt.Println("[2] (L)oad Game")
	fmt.Println("[3] (Q)uit")
	fmt.Print("\nSelect an option > ")

	if !c.scanner.Scan() {
		return false
	}

	input := strings.TrimSpace(strings.ToLower(c.scanner.Text()))

	switch input {
	case "n", "new", "1":
		c.handleNewGame()
		c.currentState = StateGameMenu
		return true
	case "l", "load", "2":
		c.handleLoadGame()
		return true
	case "q", "quit", "3":
		fmt.Println("\nExiting Galaxies: Burn Rate. Goodbye.")
		return false
	default:
		fmt.Println("\nInvalid selection. Please choose a valid option.")
		return true
	}
}

// handleNewGame processes the selection to start a new game session, prompting for initial data.
func (c *CLI) handleNewGame() {
	fmt.Println("\n--- Establishing New Corporation ---")

	fmt.Print("Enter your name, Manager > ")
	c.scanner.Scan()
	managerName := strings.TrimSpace(c.scanner.Text())
	if managerName == "" {
		managerName = "Unknown Manager"
	}

	fmt.Print("Enter your proposed Company Name > ")
	c.scanner.Scan()
	companyName := strings.TrimSpace(c.scanner.Text())
	if companyName == "" {
		companyName = "Independent Logistics"
	}

	fmt.Print("Enter a designation for your starting Scout vessel > ")
	c.scanner.Scan()
	shipName := strings.TrimSpace(c.scanner.Text())
	if shipName == "" {
		shipName = "Vanguard-1"
	}

	// Execute a deep copy of the static configuration planets into the dynamic game state.
	c.gameState.Planets = make(map[string]planet.Planet)
	for id, p := range c.settings.Planets {
		c.gameState.Planets[id] = p
	}

	startPlanetID := c.settings.Player.StartingPlanet
	startPlanet, planetExists := c.gameState.Planets[startPlanetID]
	if !planetExists {
		fmt.Printf("Error: Starting planet '%s' not found in configuration.\n", startPlanetID)
		return
	}

	c.gameState.Player = player.New(managerName, companyName, c.settings.Player.StartingCredits)

	scoutChassis, exists := c.settings.Chassis["scout"]
	if !exists {
		fmt.Println("Error: Default scout chassis not found in configuration.")
		return
	}

	starterShip := ship.New(shipName, scoutChassis, startPlanetID)

	for _, modID := range scoutChassis.DefaultModules {
		mod, modExists := c.settings.Modules[modID]
		if modExists {
			_ = starterShip.InstallModule(mod)
		} else {
			fmt.Printf("Warning: Default module '%s' not found in configuration.\n", modID)
		}
	}

	starterShip.Refuel()
	_ = c.gameState.Player.AddShip(starterShip)

	// Build the engine cache prior to executing entity generation.
	engine.Initialize(c.gameState, c.settings)

	engine.GenerateInitialContracts(c.gameState, c.settings)
	engine.GenerateInitialPassengers(c.gameState, c.settings)
	engine.GenerateInitialCrew(c.gameState, c.settings)
	engine.GenerateInitialCompetitors(c.gameState, c.settings)

	fmt.Printf("\n>>> Registration Complete. Welcome to %s, %s.\n", c.gameState.Player.CompanyName, c.gameState.Player.Name)
	fmt.Printf(">>> Operations established at %s (Sector %d, %d)\n", startPlanet.Name, startPlanet.X, startPlanet.Y)
}

// handleLoadGame processes the selection to load an existing game session.
func (c *CLI) handleLoadGame() {
	fmt.Println("\n>>> Registered selection: Load Game (Feature pending implementation)")
}
