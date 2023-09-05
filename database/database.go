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

type User struct {
	UserID   int    `json:"userid"`
	Username string `json:"username"`
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

// GetTopOpponent identifies the player who has scored the most points against the user.
func GetTopOpponent(username string) (string, error) {
	query := `
    WITH UserMatches AS (
        SELECT m.id, m.team1_id, m.team2_id, 
               CASE WHEN tm.team_id = m.team1_id THEN m.team1_score ELSE m.team2_score END as user_score,
               CASE WHEN tm.team_id = m.team1_id THEN m.team2_score ELSE m.team1_score END as opponent_score,
               CASE WHEN tm.team_id = m.team1_id THEN m.team2_id ELSE m.team1_id END as opponent_team_id
        FROM matches m
        JOIN team_members tm on tm.team_id = m.team1_id OR tm.team_id = m.team2_id
        JOIN users u on u.userid = tm.userid
        WHERE LOWER(u.username) = LOWER(?) AND (tm.team_id = m.team1_id OR tm.team_id = m.team2_id) AND tm.team_id != opponent_team_id
    ),
    OpponentScores AS (
        SELECT u.username, SUM(um.opponent_score) as total_score
        FROM UserMatches um
        JOIN team_members tm ON tm.team_id = um.opponent_team_id
        JOIN users u on u.userid = tm.userid
        WHERE u.username != ?
        GROUP BY u.username
    )
    SELECT username
    FROM OpponentScores
    WHERE total_score = (SELECT MAX(total_score) FROM OpponentScores)
    ORDER BY RANDOM()
    LIMIT 1;
    `

	var opponentName string
	err := DB.QueryRow(query, username, username).Scan(&opponentName)
	if err != nil {
		if err == sql.ErrNoRows {
			return "No opponents found", nil
		}
		return "", err
	}

	return opponentName, nil
}

// ScanAndAddMissingUsers scans the team_members table for any missing users in the users table and adds them.
func ScanAndAddMissingUsers() error {
	// Fetch all userids from the team_members table
	rows, err := DB.Query("SELECT DISTINCT userid FROM team_members")
	if err != nil {
		return err
	}
	defer rows.Close()

	// Iterate over each userid
	for rows.Next() {
		var userID int
		err := rows.Scan(&userID)
		if err != nil {
			return err
		}

		// Check if the userid exists in the users table
		userExists, err := CheckUserInDatabase(userID)
		if err != nil {
			return err
		}

		// If not exists, add them to the users table
		if !userExists {
			// Assuming that you would like to add the user to the `users` table with a blank username
			// as the actual username will be fetched and updated later from the osu API.
			err := AddUser(userID, "")
			if err != nil {
				return err
			}
		}
	}

	// Handle any error encountered during iteration
	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

// GetAllUniqueUserIDsFromTeamMembers fetches all unique user IDs from the team_members table.
func GetAllUniqueUserIDsFromTeamMembers() ([]int, error) {
	rows, err := DB.Query("SELECT DISTINCT userid FROM team_members")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []int
	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}

	// Handle any error encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userIDs, nil
}

// GetAllUsers retrieves all users from the users table.
func GetAllUsers() ([]User, error) {
	rows, err := DB.Query("SELECT userid, username FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.UserID, &user.Username); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	// Handle any error encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
