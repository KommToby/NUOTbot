// database/database.go

package database

import (
	"database/sql"
	"fmt"
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
			winning_team_id INTEGER
		);
	`)
	if err != nil {
		log.Fatalf("Failed to create tournaments table: %v", err)
	}

	// Team lists table
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS teams (
			id INTEGER PRIMARY KEY,
			tournament_id INTEGER REFERENCES tournaments(tournament_id),
			team_name TEXT
		);
	`)
	if err != nil {
		log.Fatalf("Failed to create teams table: %v", err)
	}

	// Matches table
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS matches (
			id INTEGER PRIMARY KEY,
			tournament_id INTEGER REFERENCES tournaments(tournament_id),
			team1_id INTEGER,
			team2_id INTEGER,
			team1_score INTEGER,
			team2_score INTEGER,
			match_url TEXT
		);
	`)
	if err != nil {
		log.Fatalf("Failed to create matches table: %v", err)
	}

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS team_members (
			team_id INTEGER REFERENCES teams(id),
			userid INTEGER REFERENCES users(userid),
			PRIMARY KEY (team_id, userid)
		);
	`)
	if err != nil {
		log.Fatalf("Failed to create team members table: %v", err)
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

func AddTeamMember(team_id int, userid int) error {
	_, err := DB.Exec("INSERT INTO team_members(team_id, userid) VALUES (?, ?)", team_id, userid)
	return err
}

func GetPlayerStats(username string) (string, error) {
	query := `
	SELECT count(m.id) as total_matches 
	FROM matches m 
	JOIN team_members tm on tm.team_id = m.team1_id OR tm.team_id = m.team2_id 
	JOIN users u on u.userid = tm.userid 
	WHERE LOWER(u.username) = LOWER(?)
	`

	var totalMatches int
	err := DB.QueryRow(query, username).Scan(&totalMatches)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s has played in %d matches.", username, totalMatches), nil
}
