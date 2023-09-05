// /cmd/leaderboard.go

package cmd

import (
	"sort"

	"github.com/KommToby/NUOTbot/database"
	"github.com/KommToby/NUOTbot/embed"
	"github.com/KommToby/NUOTbot/models"
	"github.com/bwmarrin/discordgo"
)

var LeaderboardCommand = &discordgo.ApplicationCommand{
	Name:        "leaderboard",
	Description: "Returns a leaderboard based on MatchesPlayed, PointsScored, and PointsScoredAgainst",
}

const EntriesPerPage = 10

func BuildLeaderboard() ([]models.UserLeaderboardEntry, error) {
	users, err := database.GetAllUsers()
	if err != nil {
		return nil, err
	}

	var leaderboard []models.UserLeaderboardEntry
	for _, user := range users {
		stats, err := database.GetPlayerStats(user.Username)
		if err != nil {
			return nil, err
		}

		entry := models.UserLeaderboardEntry{
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

	// Create leaderboard embed
	embed := embed.CreateLeaderboardEmbed(leaderboard, 10) // assuming top 10

	// Send follow-up with the leaderboard embed
	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{embed},
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

func splitIntoPages(leaderboard []models.UserLeaderboardEntry) [][]models.UserLeaderboardEntry {
	var pages [][]models.UserLeaderboardEntry
	for i := 0; i < len(leaderboard); i += EntriesPerPage {
		end := i + EntriesPerPage
		if end > len(leaderboard) {
			end = len(leaderboard)
		}
		pages = append(pages, leaderboard[i:end])
	}
	return pages
}
