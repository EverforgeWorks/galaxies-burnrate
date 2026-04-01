// Package main serves as the primary entry point for the Galaxies: Burn Rate CLI.
package main

import (
	"fmt"
	"os"

	"galaxies-burnrate/internal/cli"
	"galaxies-burnrate/internal/config"
)

// main initializes and hands over control to the CLI interface.
func main() {
	// Load the foundational game settings from the configs directory prior to initializing the application.
	settings, err := config.Load("configs")
	if err != nil {
		fmt.Printf("Fatal error loading configuration: %v\n", err)
		os.Exit(1)
	}

	app := cli.New(settings)

	// Execute the application loop and exit with the returned status code.
	os.Exit(app.Run())
}
