// /cmd/stats.go

package cmd

import (
	"github.com/KommToby/NUOTbot/database"
	"github.com/bwmarrin/discordgo"
)

var StatsCommand = &discordgo.ApplicationCommand{
	Name:        "stats",
	Description: "Returns tournament statistics for a specific player",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "username",
			Description: "Username of the player",
			Required:    true,
		},
	},
}

func StatsHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	username := i.ApplicationCommandData().Options[0].StringValue()

	// Fetch user's stats from the database
	stats, err := database.GetPlayerStats(username)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error fetching stats.",
			},
		})
		return
	}

	// Format the stats and send back to the user
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: stats, // assuming stats is a string at the moment
		},
	})

	// example
	// response, err := auth.GosuClient.GetUserBeatmapScore
}
