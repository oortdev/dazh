package main

type Item struct {
	ID    int // <- change from string to int
	Name  string
	URL   string
	Group string
	Color string
}

type DashboardData struct {
	Groups      map[string][]Item
	AllGroups   []string
	SelectedGrp string // for filtering
}

type ManageData struct {
	Items       []Item
	AllGroups   []string
	SelectedGrp string
	Theme       ThemeConfig
}

type ThemeConfig struct {
	BackgroundColor string `json:"background_color"`
	MenuColor       string `json:"menu_color"`
	BackgroundImage string `json:"background_image"`
}
