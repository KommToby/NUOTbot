// /cmd/leaderboard.go

package cmd

import (
	"fmt"
	"sort"

	"github.com/KommToby/NUOTbot/database"
	"github.com/bwmarrin/discordgo"
)

var LeaderboardCommand = &discordgo.ApplicationCommand{
	Name:        "leaderboard",
	Description: "Returns a leaderboard based on MatchesPlayed, PointsScored, and PointsScoredAgainst",
}

type UserLeaderboardEntry struct {
	Username            string
	MatchesPlayed       int
	PointsScored        int
	PointsScoredAgainst int
}

func BuildLeaderboard() ([]UserLeaderboardEntry, error) {
	users, err := database.GetAllUsers()
	if err != nil {
		return nil, err
	}

	var leaderboard []UserLeaderboardEntry
	for _, user := range users {
		stats, err := database.GetPlayerStats(user.Username)
		if err != nil {
			return nil, err
		}

		entry := UserLeaderboardEntry{
			Username:            user.Username,
			MatchesPlayed:       stats.MatchesPlayed,
			PointsScored:        stats.PointsScored,
			PointsScoredAgainst: stats.PointsScoredAgainst,
		}

		leaderboard = append(leaderboard, entry)
	}

	// Sorting the leaderboard based on MatchesPlayed
	// Adjust this part for different sorting criteria.
	sort.Slice(leaderboard, func(i, j int) bool {
		return leaderboard[i].MatchesPlayed > leaderboard[j].MatchesPlayed
	})

	return leaderboard, nil
}

func LeaderboardHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Defer the response
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		// Handle error if needed
		return
	}

	leaderboard, err := BuildLeaderboard()
	if err != nil {
		followupError(s, i, "Error building leaderboard.")
		return
	}

	// Create a message for the leaderboard.
	content := "Leaderboard:\n"
	maxEntries := 10
	if len(leaderboard) < maxEntries {
		maxEntries = len(leaderboard)
	}
	for index := 0; index < maxEntries; index++ {
		entry := leaderboard[index]
		content += fmt.Sprintf("Username: %s, Matches Played: %d, Points Scored: %d, Points Scored Against: %d\n",
			entry.Username, entry.MatchesPlayed, entry.PointsScored, entry.PointsScoredAgainst)
	}

	// Send follow-up with the leaderboard
	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: content,
	})
	if err != nil {
		// Handle error if needed
		return
	}
}

// Utility function to send follow-up error messages
func followupError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	_, _ = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: message,
	})
}
