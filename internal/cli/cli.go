// Package cli provides the command-line interface and menu handling for the game.
package cli

import (
	"bufio"
	"fmt"
	"os"
	"sort"

	"galaxies-burnrate/internal/config"
	"galaxies-burnrate/internal/game"
	"galaxies-burnrate/internal/ship"
)

// MenuState defines the current active interface layer of the application.
type MenuState int

const (
	// StateMainMenu represents the initial application landing screen.
	StateMainMenu MenuState = iota
	// StateGameMenu represents the active simulation loop and management interface.
	StateGameMenu
)

// CLI represents the command-line interface application state and dependencies.
type CLI struct {
	scanner      *bufio.Scanner
	gameState    *game.State
	settings     *config.GameSettings
	currentState MenuState
	activeShipID string
}

// New initializes and returns a new CLI instance ready for execution.
func New(settings *config.GameSettings) *CLI {
	return &CLI{
		scanner:      bufio.NewScanner(os.Stdin),
		gameState:    game.NewState(),
		settings:     settings,
		currentState: StateMainMenu,
	}
}

// Run starts the primary application loop, delegating to the active menu state.
func (c *CLI) Run() int {
	c.renderBanner()

	for {
		var keepRunning bool

		switch c.currentState {
		case StateMainMenu:
			keepRunning = c.runMainMenu()
		case StateGameMenu:
			keepRunning = c.runGameMenu()
		}

		if !keepRunning {
			return 0
		}
	}
}

// renderBanner outputs the ASCII art title for the application.
func (c *CLI) renderBanner() {
	banner := `
  _____       _           _            
 / ____|     | |         (_)           
| |  __  __ _| | __ ___  ___  ___  ___ 
| | |_ |/ _' | |/ _' \ \/ / |/ _ \/ __|
| |__| | (_| | | (_| |>  <| |  __/\__ \
 \_____|\__,_|_|\__,_/_/\_\_|\___||___/
                               
     B U R N   R A T E
`
	fmt.Println(banner)
}

// getActiveShip retrieves the currently managed vessel. If none is set, it defaults to the first alphabetical ship.
func (c *CLI) getActiveShip() *ship.Ship {
	if len(c.gameState.Player.Fleet) == 0 {
		c.activeShipID = ""
		return nil
	}

	if c.activeShipID != "" {
		if s, exists := c.gameState.Player.Fleet[c.activeShipID]; exists {
			return s
		}
	}

	var fleet []*ship.Ship
	for _, s := range c.gameState.Player.Fleet {
		fleet = append(fleet, s)
	}
	sort.Slice(fleet, func(i, j int) bool {
		return fleet[i].Name < fleet[j].Name
	})

	c.activeShipID = fleet[0].ID
	return fleet[0]
}
