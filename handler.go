package main

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// DashboardHandler renders the main dashboard
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	items := LoadItems()
	selectedGrp := r.URL.Query().Get("group")

	// Group items by Group
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

	// Load theme (background image, colors, etc.)
	theme := LoadTheme() // implement this to read from config.json or CSV

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

// ManageHandler
func ManageHandler(w http.ResponseWriter, r *http.Request) {
	items := LoadItems()
	selectedGrp := r.URL.Query().Get("group")
	theme := LoadTheme()

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

// SaveHandler
func SaveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	r.ParseMultipartForm(10 << 20) // 10MB max

	items := LoadItems()
	nextID := getNextID(items)

	// Update existing items
	for i, item := range items {
		if r.FormValue("id_"+strconv.Itoa(item.ID)+"_delete") == "on" {
			items[i] = Item{}
			continue
		}

		item.Name = r.FormValue("id_" + strconv.Itoa(item.ID) + "_name")
		item.URL = r.FormValue("id_" + strconv.Itoa(item.ID) + "_url")
		item.Group = r.FormValue("id_" + strconv.Itoa(item.ID) + "_group")
		item.Color = r.FormValue("id_" + strconv.Itoa(item.ID) + "_color")
		item.Image = r.FormValue("id_" + strconv.Itoa(item.ID) + "_image")

		// Handle uploaded file
		file, header, err := r.FormFile("id_" + strconv.Itoa(item.ID) + "_file")
		if err == nil && file != nil {
			defer file.Close()
			dst := filepath.Join("static/uploads", header.Filename)
			os.MkdirAll("static/uploads", os.ModePerm)
			out, _ := os.Create(dst)
			defer out.Close()
			io.Copy(out, file)
			item.Image = "/static/uploads/" + header.Filename
		}

		items[i] = item
	}

	// Remove deleted items
	filtered := []Item{}
	for _, i := range items {
		if i.Name != "" || i.URL != "" {
			filtered = append(filtered, i)
		}
	}
	items = filtered

	// Add new items
	names := r.Form["new_name[]"]
	urls := r.Form["new_url[]"]
	groups := r.Form["new_group[]"]
	colors := r.Form["new_color[]"]
	images := r.Form["new_image[]"]

	for idx := range names {
		if names[idx] == "" && urls[idx] == "" {
			continue
		}
		newItem := Item{
			ID:    nextID,
			Name:  names[idx],
			URL:   urls[idx],
			Group: groups[idx],
			Color: colors[idx],
			Image: images[idx],
		}

		// Handle uploaded file for new row
		file, header, err := r.FormFile("new_file[]")
		if err == nil && file != nil {
			defer file.Close()
			dst := filepath.Join("static/uploads", header.Filename)
			os.MkdirAll("static/uploads", os.ModePerm)
			out, _ := os.Create(dst)
			defer out.Close()
			io.Copy(out, file)
			newItem.Image = "/static/uploads/" + header.Filename
		}

		items = append(items, newItem)
		nextID++
	}

	// Theme settings
	theme := ThemeConfig{
		BackgroundColor: r.FormValue("background_color"),
		MenuColor:       r.FormValue("menu_color"),
		BackgroundImage: r.FormValue("background_image"),
	}

	file, header, err := r.FormFile("background_file")
	if err == nil && file != nil {
		defer file.Close()
		dst := filepath.Join("static/uploads", header.Filename)
		os.MkdirAll("static/uploads", os.ModePerm)
		out, _ := os.Create(dst)
		defer out.Close()
		io.Copy(out, file)
		theme.BackgroundImage = "/static/uploads/" + header.Filename
	}

	SaveItems(items)

	f, _ := os.Create("config/theme.json")
	defer f.Close()
	json.NewEncoder(f).Encode(theme)

	http.Redirect(w, r, "/manage", http.StatusSeeOther)
}
