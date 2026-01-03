package main

import (
	"encoding/json"
	"os"
)

func LoadTheme() ThemeConfig {
	theme := ThemeConfig{
		BackgroundColor: "#0f172a",
		MenuColor:       "#1e293b",
		BackgroundImage: "",
	}

	file, err := os.Open("config/theme.json")
	if err != nil {
		return theme
	}
	defer file.Close()

	json.NewDecoder(file).Decode(&theme)
	return theme
}
