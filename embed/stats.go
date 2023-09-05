// /embed/stats.go

package embed

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func CreateStatsEmbed(username string, matchesPlayed int, winPercentage float64, topOpponent string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Player Statistics",
		Description: fmt.Sprintf("%s's Stats", username),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Matches Played",
				Value:  fmt.Sprintf("%d", matchesPlayed),
				Inline: true,
			},
			{
				Name:   "Win Percentage",
				Value:  fmt.Sprintf("%.2f%%", winPercentage),
				Inline: true,
			},
			{
				Name:   "Lost Most Against",
				Value:  topOpponent,
				Inline: true,
			},
		},
		Color: 0x0099ff, // This is blue
	}
}
