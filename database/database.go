// database/database.go

package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./database/database.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Users table
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			userid INTEGER PRIMARY KEY,
			username TEXT
		);
	`)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	// Tournaments table
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS tournaments (
			tournament_id INTEGER PRIMARY KEY,
			tournament_name TEXT,
			format TEXT,
			winning_team TEXT
		);
	`)
	if err != nil {
		log.Fatalf("Failed to create tournaments table: %v", err)
	}

	// Team lists table
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS team_lists (
			id INTEGER PRIMARY KEY,
			tournament_id INTEGER REFERENCES tournaments(tournament_id),
			userid INTEGER REFERENCES users(userid),
			team_name TEXT
		);
	`)
	if err != nil {
		log.Fatalf("Failed to create team lists table: %v", err)
	}

	// Matches table
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS matches (
			id INTEGER PRIMARY KEY,
			tournament_id INTEGER REFERENCES tournaments(tournament_id),
			team1_id TEXT,
			team2_id TEXT,
			team1_score INTEGER,
			team2_score INTEGER,
			match_url TEXT
		);
	`)
	if err != nil {
		log.Fatalf("Failed to create matches table: %v", err)
	}

	// update team_lists table to have an extra column for placements
}

// Add a new user
func AddUser(userid int, username string) error {
	_, err := DB.Exec("INSERT INTO users(userid, username) VALUES (?, ?)", userid, username)
	return err
}

// Add a new tournament
func AddTournament(tournament_id int, tournament_name, format, winning_team string) error {
	_, err := DB.Exec("INSERT INTO tournaments(tournament_id, tournament_name, format, winning_team) VALUES (?, ?, ?, ?)", tournament_id, tournament_name, format, winning_team)
	return err
}
