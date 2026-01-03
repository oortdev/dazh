package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"strconv"
)

// DashboardHandler renders the main dashboard
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	items := LoadItems()
	selectedGrp := r.URL.Query().Get("group")

	groups := make(map[string][]Item)
	allGroups := make(map[string]struct{})
	for _, item := range items {
		groups[item.Group] = append(groups[item.Group], item)
		allGroups[item.Group] = struct{}{}
	}

	allGroupsList := []string{}
	for g := range allGroups {
		allGroupsList = append(allGroupsList, g)
	}

	tmpl, err := template.ParseFiles("templates/dashboard.html")
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	theme := LoadTheme() // load theme

	data := struct {
		Groups      map[string][]Item
		AllGroups   []string
		SelectedGrp string
		Theme       ThemeConfig
	}{
		Groups:      groups,
		AllGroups:   allGroupsList,
		SelectedGrp: selectedGrp,
		Theme:       theme,
	}

	tmpl.Execute(w, data)
}

// ManageHandler renders the manage page
func ManageHandler(w http.ResponseWriter, r *http.Request) {
	items := LoadItems()
	selectedGrp := r.URL.Query().Get("group")
	theme := LoadTheme() // load current theme

	// Collect all group names for filter
	allGroups := make(map[string]struct{})
	for _, item := range items {
		allGroups[item.Group] = struct{}{}
	}
	allGroupsList := []string{}
	for g := range allGroups {
		allGroupsList = append(allGroupsList, g)
	}

	tmpl, err := template.ParseFiles("templates/manage.html")
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := ManageData{
		Items:       items,
		AllGroups:   allGroupsList,
		SelectedGrp: selectedGrp,
		Theme:       theme,
	}

	tmpl.Execute(w, data)
}

// SaveHandler processes form submission from manage page
func SaveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	r.ParseForm()
	items := LoadItems()
	newItems := []Item{}
	nextID := getNextID(items) // you can implement this in csv.go

	// Update existing items & handle deletions
	for i, item := range items {
		if r.FormValue("id_"+strconv.Itoa(item.ID)+"_delete") == "on" {
			continue
		}
		item.Name = r.FormValue("id_" + strconv.Itoa(item.ID) + "_name")
		item.URL = r.FormValue("id_" + strconv.Itoa(item.ID) + "_url")
		item.Group = r.FormValue("id_" + strconv.Itoa(item.ID) + "_group")
		item.Color = r.FormValue("id_" + strconv.Itoa(item.ID) + "_color")
		items[i] = item
	}

	// Add new items
	names := r.Form["new_name[]"]
	urls := r.Form["new_url[]"]
	groups := r.Form["new_group[]"]
	colors := r.Form["new_color[]"]

	for i := range names {
		if names[i] == "" && urls[i] == "" && groups[i] == "" {
			continue
		}
		newItem := Item{
			ID:    nextID,
			Name:  names[i],
			URL:   urls[i],
			Group: groups[i],
			Color: colors[i],
		}
		newItems = append(newItems, newItem)
		nextID++
	}

	// Merge
	items = append(items, newItems...)

	if err := SaveItems(items); err != nil {
		http.Error(w, "Failed to save: "+err.Error(), http.StatusInternalServerError)
		return
	}

	theme := ThemeConfig{
		BackgroundColor: r.FormValue("background_color"),
		MenuColor:       r.FormValue("menu_color"),
		BackgroundImage: r.FormValue("background_image"),
	}

	// Save theme to JSON
	file, _ := os.Create("config/theme.json")
	defer file.Close()
	json.NewEncoder(file).Encode(theme)

	http.Redirect(w, r, "/manage", http.StatusSeeOther)
}
