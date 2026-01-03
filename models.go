package main

type Item struct {
	ID    int
	Name  string
	URL   string
	Group string
	Color string
	Image string // Optional tile image
}

type ThemeConfig struct {
	BackgroundColor string `json:"background_color"`
	MenuColor       string `json:"menu_color"`
	BackgroundImage string `json:"background_image"`
}

type DashboardData struct {
	Groups      map[string][]Item
	AllGroups   []string
	SelectedGrp string
	Theme       ThemeConfig
}

type ManageData struct {
	Items       []Item
	AllGroups   []string
	SelectedGrp string
	Theme       ThemeConfig
}
