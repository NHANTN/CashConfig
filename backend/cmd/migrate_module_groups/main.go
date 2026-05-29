package main

import (
	"encoding/json"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/cashier-config/server/internal/model"
)

type stepEntry struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type moduleGroup struct {
	Name  string      `json:"name"`
	Steps []stepEntry `json:"steps"`
}

func main() {
	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=cashier_config sslmode=disable"
	if dsn == "" {
		log.Fatal("DATABASE_DSN not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	var modules []model.Module
	if err := db.Find(&modules).Error; err != nil {
		log.Fatalf("failed to query modules: %v", err)
	}
	fmt.Printf("Found %d modules\n", len(modules))

	updated := 0
	skipped := 0

	for _, m := range modules {
		// Skip if already numeric IDs
		var ids []int64
		if json.Unmarshal([]byte(m.Modules), &ids) == nil {
			skipped++
			continue
		}

		// Parse as inline group objects [{name, steps}, ...]
		var inline []moduleGroup
		if err := json.Unmarshal([]byte(m.Modules), &inline); err != nil {
			fmt.Printf("  [SKIP] module %d (%s): cannot parse Modules: %v\n", m.ID, m.Name, err)
			skipped++
			continue
		}
		if len(inline) == 0 {
			skipped++
			continue
		}

		// Resolve each inline group name -> group ID
		var resolvedIDs []int64
		for _, g := range inline {
			var group model.Group
			if err := db.Where("name = ?", g.Name).First(&group).Error; err != nil {
				fmt.Printf("  [WARN] module %d (%s): group %q not found, skipping\n", m.ID, m.Name, g.Name)
				continue
			}
			resolvedIDs = append(resolvedIDs, group.ID)
		}
		if len(resolvedIDs) == 0 {
			fmt.Printf("  [SKIP] module %d (%s): no groups could be resolved\n", m.ID, m.Name)
			skipped++
			continue
		}

		idJSON, _ := json.Marshal(resolvedIDs)
		if err := db.Model(&m).Update("modules", string(idJSON)).Error; err != nil {
			fmt.Printf("  [ERR] module %d (%s): %v\n", m.ID, m.Name, err)
			continue
		}
		fmt.Printf("  [OK] module %d (%s): resolved to %s\n", m.ID, m.Name, string(idJSON))
		updated++
	}

	fmt.Printf("\nDone. Updated: %d, Skipped: %d\n", updated, skipped)
}
