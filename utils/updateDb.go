package utils

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"slacktui/config"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func UpdateDB() (success bool) {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return false
	}

	// Open the database file
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	dbPath := filepath.Join(home, ".config", "slack-tui", "users.sqlite")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return false
	}
	defer db.Close()

	createTableQuery := `
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			name TEXT,
			real_name TEXT,
			tz TEXT
		);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return false
	}

	// Pagination variables
	baseURL := "https://slack.com/api/users.list?limit=1000"
	cursor := ""
	client := &http.Client{Timeout: 5 * time.Second}

	for {
		// Build the request URL with the cursor
		url := baseURL
		if cursor != "" {
			url += "&cursor=" + cursor
		}

		// Make the API request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return false
		}

		req.Header.Set("Authorization", "Bearer "+cfg.SlackToken)
		req.Header.Set("Cookie", cfg.Cookies)

		resp, err := client.Do(req)
		if err != nil {
			return false
		}
		defer resp.Body.Close()

		var result struct {
			Ok      bool   `json:"ok"`
			Error   string `json:"error,omitempty"`
			Members []struct {
				ID       string `json:"id"`
				Name     string `json:"name"`
				RealName string `json:"real_name"`
				Tz       string `json:"tz"`
			} `json:"members"`
			ResponseMetadata struct {
				NextCursor string `json:"next_cursor"`
			} `json:"response_metadata"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return false
		}

		if !result.Ok {
			return false
		}

		// Insert members into the database
		for _, member := range result.Members {

			if member.RealName == "" { // Check if RealName is an empty string
				member.RealName = member.Name // Fallback to name if real_name is empty
			}
			if member.Tz == "" {
				member.Tz = "UTC" // Default timezone if not set
			}

			insertQuery := `
			INSERT INTO users (id, name, real_name, tz)
			VALUES (?, ?, ?, ?)
			ON CONFLICT(id) DO UPDATE SET
				name = excluded.name,
				real_name = excluded.real_name,
				tz = excluded.tz;`
			_, err = db.Exec(insertQuery, member.ID, member.Name, member.RealName, member.Tz)
			if err != nil {
				return false
			}
		}

		// Check if there is a next cursor
		cursor = result.ResponseMetadata.NextCursor
		if cursor == "" {
			break
		}
	}
	return true
}
