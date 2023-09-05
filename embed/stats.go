// /embed/stats.go

package embed

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func CreateStatsEmbed(username string, matchesPlayed int, winPercentage float64, topOpponent string, avatarURL string, bestTeammate string, bestTournament string, firstTournament string) *discordgo.MessageEmbed {
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
			{
				Name:   "Best Teammate",
				Value:  bestTeammate,
				Inline: true,
			},
			{
				Name:   "Best Tournament",
				Value:  bestTournament,
				Inline: true,
			},
			{
				Name:   "First Tournament",
				Value:  firstTournament,
				Inline: true,
			},
		},
		Color: 0x0099ff, // This is blue
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: avatarURL,
		},
	}
}
