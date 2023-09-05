// models/user.go
package models

type UserLeaderboardEntry struct {
	Username            string
	MatchesPlayed       int
	PointsScored        int
	PointsScoredAgainst int
}
