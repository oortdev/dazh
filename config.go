package main

import (
	"encoding/json"
	"os"
)

// LoadTheme loads theme configuration
func LoadTheme() ThemeConfig {
	// Default values
	theme := ThemeConfig{
		BackgroundColor: "#0f172a",
		MenuColor:       "#1e293b",
		BackgroundImage: "",
	}

	// Optional: read from JSON file
	file, err := os.Open("config/theme.json")
	if err != nil {
		return theme
	}
	defer file.Close()

	json.NewDecoder(file).Decode(&theme)
	return theme
}
