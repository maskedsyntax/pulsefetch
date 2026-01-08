package main

import (
	"fmt"
	"os"

	"pulsefetch/internal/config"
	"pulsefetch/internal/fetcher"
	"pulsefetch/internal/ui"
)

func main() {
	// Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Fetch System Info
	info, err := fetcher.Fetch(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching system info: %v\n", err)
		os.Exit(1)
	}

	// Get Logo
	logo := ui.GetLogo(cfg, info)

	// Render Output
	ui.Render(cfg, info, logo)
}
