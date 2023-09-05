// /cmd/stats.go

package cmd

import (
	"github.com/KommToby/NUOTbot/auth"
	"github.com/KommToby/NUOTbot/database"
	"github.com/KommToby/NUOTbot/embed"
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

	// api call the username
	response, err := auth.GosuClient.GetUserData(username)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Username does not exist on osu! servers", // usually typos, unlikely to be api issues
			},
		})
		return
	}

	// if username exists, check their user id
	userID := response.UserCompact.ID
	if userID == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "User not found", // Secondary check
			},
		})
		return
	}

	// Check if the userid is in the database
	// Check if userID is in the database
	userInDB, err := database.CheckUserInDatabase(userID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error checking user in database.",
			},
		})
		return
	}

	if !userInDB {
		// Add user to database
		database.AddUser(userID, response.Username)
	} else {
		// Check if response.username is different and update
		currentUsername, err := database.GetUserNameFromDatabase(userID)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Error fetching username from database.",
				},
			})
			return
		}
		if currentUsername != response.Username {
			database.UpdateUsernameInDatabase(userID, response.Username)
		}
	}

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

	// Calculate win percentage based on points scored vs points scored against
	var winPercentage float64
	if stats.PointsScored+stats.PointsScoredAgainst != 0 {
		winPercentage = float64(stats.PointsScored) / float64(stats.PointsScored+stats.PointsScoredAgainst) * 100
	}

	// Fetch the opponent against whom the user has lost the most.
	opponent, err := database.GetTopOpponent(username)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error fetching top opponent.",
			},
		})
		return
	}

	// If the opponent's username doesn't exist in the database, similar check and add/update.
	// Note: This can be optimized by avoiding repeated checks or making a bulk check at the start.
	opponentExistsInDB, err := auth.GosuClient.GetUserData(opponent)
	if err == nil {
		userInDB, err := database.CheckUserInDatabase(opponentExistsInDB.UserCompact.ID)
		if err == nil && !userInDB {
			// Add opponent to database
			database.AddUser(opponentExistsInDB.UserCompact.ID, opponentExistsInDB.Username)
		}
	}

	// Create the embed
	statsEmbed := embed.CreateStatsEmbed(username, stats.MatchesPlayed, winPercentage, opponent)

	// Respond to the user with the embed
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{statsEmbed},
		},
	})

}
