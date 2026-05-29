package main

import (
	"encoding/json"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/cashier-config/server/internal/model"
)

type moduleGroup struct {
	Name  string       `json:"name"`
	Steps []stepEntry  `json:"steps"`
}

type stepEntry struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func main() {
	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=cashier_config sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	var modules []model.Module
	if err := db.Find(&modules).Error; err != nil {
		log.Fatalf("failed to query modules: %v", err)
	}
	fmt.Printf("Found %d modules\n", len(modules))

	seen := map[string]bool{}
	created := 0
	updated := 0

	for _, m := range modules {
		var groups []moduleGroup
		if err := json.Unmarshal([]byte(m.Modules), &groups); err != nil {
			fmt.Printf("  skip module %s: invalid json: %v\n", m.Name, err)
			continue
		}
		for _, g := range groups {
			stepsJSON, _ := json.Marshal(g.Steps)
			record := model.Group{
				Name:  g.Name,
				Steps: string(stepsJSON),
			}

			var existing model.Group
			err := db.Where("name = ?", g.Name).First(&existing).Error
			if err == nil {
				if existing.Steps != record.Steps {
					db.Model(&existing).Update("steps", record.Steps)
					updated++
					fmt.Printf("  updated group: %s\n", g.Name)
				}
			} else {
				db.Create(&record)
				created++
				fmt.Printf("  created group: %s\n", g.Name)
			}
			seen[g.Name] = true
		}
	}

	fmt.Printf("\nDone: %d created, %d updated, %d total unique groups\n", created, updated, len(seen))
}
