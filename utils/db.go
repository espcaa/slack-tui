package utils

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite" // Import the SQLite driver
)

// GetNameFromID retrieves the name associated with a given ID from the database.
func GetNameFromID(id string, realname bool) (string, error) {
	// Get the user's home directory
	home, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Error getting user home directory: %v", err)
		return "", err
	}

	// Construct the database path
	dbPath := filepath.Join(home, ".config", "slack-tui", "users.sqlite")
	db, err := sql.Open("sqlite", dbPath)

	if err != nil {
		log.Printf("Error opening database: %v", err)
		return "", err
	}
	defer db.Close()

	// Query the database for the name
	var name string
	query := "SELECT real_name FROM users WHERE id = ?"
	err = db.QueryRow(query, id).Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows found for the given ID
			return "", nil
		}
		log.Printf("Error querying database: %v", err)
		return "", err
	}

	return name, nil
}
