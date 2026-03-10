package main

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"path/filepath"
	"sort"
)

func main() {
	driver := envOrDefault("OPS_API_DB_DRIVER", "postgres")
	dsn := os.Getenv("OPS_API_DB_DSN")
	if dsn == "" {
		log.Fatal("OPS_API_DB_DSN is required")
	}

	db, err := sql.Open(driver, dsn)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	files, err := findMigrationFiles()
	if err != nil {
		log.Fatalf("failed to find migrations: %v", err)
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("failed to read migration %s: %v", file, err)
		}

		log.Printf("applying migration: %s", file)
		if _, err := db.Exec(string(content)); err != nil {
			log.Fatalf("failed to execute migration %s: %v", file, err)
		}
	}

	log.Printf("applied %d migration(s)", len(files))
}

func findMigrationFiles() ([]string, error) {
	patterns := []string{"migrations/*.sql", "apps/backend/migrations/*.sql"}
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}
		if len(matches) > 0 {
			sort.Strings(matches)
			return matches, nil
		}
	}
	return nil, errors.New("no migration files found")
}

func envOrDefault(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
