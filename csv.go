package main

import (
	"encoding/csv"
	"os"
	"strconv"
)

func LoadItems() []Item {
	file, err := os.Open("config/url.csv")
	if err != nil {
		return []Item{}
	}
	defer file.Close()

	r := csv.NewReader(file)
	rows, err := r.ReadAll()
	if err != nil {
		return []Item{}
	}

	items := []Item{}
	for _, row := range rows {
		if len(row) < 6 {
			continue
		}
		id, _ := strconv.Atoi(row[0])
		items = append(items, Item{
			ID:    id,
			Name:  row[1],
			URL:   row[2],
			Group: row[3],
			Color: row[4],
			Image: row[5],
		})
	}
	return items
}

func SaveItems(items []Item) error {
	file, err := os.Create("config/url.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	for _, item := range items {
		row := []string{
			strconv.Itoa(item.ID),
			item.Name,
			item.URL,
			item.Group,
			item.Color,
			item.Image,
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}
	return nil
}

func getNextID(items []Item) int {
	maxID := 0
	for _, i := range items {
		if i.ID > maxID {
			maxID = i.ID
		}
	}
	return maxID + 1
}
