package cmd

import (
	"strconv"

	"github.com/KommToby/NUOTbot/auth"
	"github.com/KommToby/NUOTbot/database"
	"github.com/bwmarrin/discordgo"
)

var ScanUsersCommand = &discordgo.ApplicationCommand{
	Name:        "scanusers",
	Description: "Scans and populates missing users from the team_members table into the users table",
}

func ScanUsersHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Send an initial "Processing..." message
	initialResponse, err := s.ChannelMessageSend(i.ChannelID, "Processing...")
	if err != nil {
		// Handle the error here
		return
	}

	// Get a list of all unique userIDs from the team_members table
	userIDs, err := database.GetAllUniqueUserIDsFromTeamMembers()
	if err != nil {
		// Handle the error here
		return
	}

	// Iterate through each userID and check if they exist in the users table
	for _, userID := range userIDs {
		userInDB, err := database.CheckUserInDatabase(userID)
		if err != nil {
			continue // You might choose to log the error instead of just continuing
		}

		if !userInDB {
			// User does not exist in users table, fetch their username from osu! API
			response, err := auth.GosuClient.GetUserData(strconv.Itoa(userID))
			if err != nil {
				continue // You might choose to log the error instead of just continuing
			}

			// Add the user to the users table
			database.AddUser(userID, response.Username)
		}
	}

	// Edit the initial "Processing..." message to indicate completion
	s.ChannelMessageEdit(i.ChannelID, initialResponse.ID, "Scan and population of missing users completed.")

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Done!",
		},
	})
}
