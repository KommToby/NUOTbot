// /embed/leaderboard.go

package embed

import (
	"fmt"

	"github.com/KommToby/NUOTbot/models"
	"github.com/bwmarrin/discordgo"
)

func CreateLeaderboardEmbed(entries []models.UserLeaderboardEntry, maxEntries int) *discordgo.MessageEmbed {
	var description string

	if len(entries) < maxEntries {
		maxEntries = len(entries)
	}

	for index := 0; index < maxEntries; index++ {
		entry := entries[index]
		description += fmt.Sprintf("**%s:** Matches Played: %d\n",
			entry.Username, entry.MatchesPlayed)
	}

	return &discordgo.MessageEmbed{
		Title:       "Leaderboard",
		Description: description,
		Color:       0x0099ff, // This is blue
	}
}
