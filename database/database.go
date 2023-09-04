// database/database.go

package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

var DB *sql.DB

type PlayerStats struct {
	MatchesPlayed       int
	PointsScored        int
	PointsScoredAgainst int
}

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

func GetPlayerStats(username string) (PlayerStats, error) {
	var stats PlayerStats

	// Query for total matches played
	queryMatches := `
	SELECT count(m.id) 
	FROM matches m 
	JOIN team_members tm on tm.team_id = m.team1_id OR tm.team_id = m.team2_id 
	JOIN users u on u.userid = tm.userid 
	WHERE LOWER(u.username) = LOWER(?)
	`
	err := DB.QueryRow(queryMatches, username).Scan(&stats.MatchesPlayed)
	if err != nil {
		return stats, err
	}

	// Query for total points scored by the user's team
	queryPointsScored := `
	SELECT COALESCE(SUM(CASE WHEN tm.team_id = m.team1_id THEN m.team1_score ELSE m.team2_score END), 0) 
	FROM matches m 
	JOIN team_members tm on tm.team_id = m.team1_id OR tm.team_id = m.team2_id 
	JOIN users u on u.userid = tm.userid 
	WHERE LOWER(u.username) = LOWER(?)
	`
	err = DB.QueryRow(queryPointsScored, username).Scan(&stats.PointsScored)
	if err != nil {
		return stats, err
	}

	// Query for total points scored against the user's team
	queryPointsAgainst := `
	SELECT COALESCE(SUM(CASE WHEN tm.team_id = m.team1_id THEN m.team2_score ELSE m.team1_score END), 0) 
	FROM matches m 
	JOIN team_members tm on tm.team_id = m.team1_id OR tm.team_id = m.team2_id 
	JOIN users u on u.userid = tm.userid 
	WHERE LOWER(u.username) = LOWER(?)
	`
	err = DB.QueryRow(queryPointsAgainst, username).Scan(&stats.PointsScoredAgainst)
	if err != nil {
		return stats, err
	}

	return stats, nil
}

// CheckUserInDatabase checks if a user with the given ID exists in the database.
func CheckUserInDatabase(userID int) (bool, error) {
	query := `SELECT COUNT(userid) FROM users WHERE userid = ?`

	var count int
	err := DB.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetUserNameFromDatabase fetches the username associated with the given user ID.
func GetUserNameFromDatabase(userID int) (string, error) {
	query := `SELECT username FROM users WHERE userid = ?`

	var username string
	err := DB.QueryRow(query, userID).Scan(&username)
	if err != nil {
		return "", err
	}

	return username, nil
}

// UpdateUsernameInDatabase updates the username for a given user ID.
func UpdateUsernameInDatabase(userID int, newUsername string) error {
	query := `UPDATE users SET username = ? WHERE userid = ?`

	_, err := DB.Exec(query, newUsername, userID)
	return err
}
