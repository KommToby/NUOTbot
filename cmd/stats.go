// /cmd/stats.go

package cmd

import (
	"bytes"
	"fmt"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/KommToby/NUOTbot/auth"
	"github.com/KommToby/NUOTbot/database"
	"github.com/KommToby/NUOTbot/embed"
	"github.com/KommToby/NUOTbot/img"
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

	images, err := img.LoadImages()
	if err != nil {
		fmt.Println("Error loading images:", err)
		return
	}

	fmt.Println("Number of images:", len(images))
	for i, img := range images {
		if img == nil {
			fmt.Printf("Image at index %d is nil\n", i)
		} else {
			fmt.Printf("Image at index %d: width = %d, height = %d\n", i, img.Bounds().Dx(), img.Bounds().Dy())
		}
	}

	// Create the banner
	banner, err := img.CreateBanner(images)
	if err != nil {
		fmt.Println("Error creating banner:", err)
		return
	}

	// Save to an output file
	outFile, err := os.Create("./img/banner_from_stats.png")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outFile.Close()

	png.Encode(outFile, banner)

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
	stats, err := database.GetPlayerStats(response.Username)
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
	opponent, err := database.GetTopOpponent(response.Username)
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

	// Fetch the best teammate
	bestTeammate, err := database.GetTopTeammate(response.Username)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error fetching best teammate.",
			},
		})
		return
	}

	// Fetch the best tournament
	bestTournament, err := database.GetBestTournament(response.Username)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error fetching best tournament.",
			},
		})
		return
	}

	// Fetch the first tournament
	firstTournament, err := database.GetFirstTournament(response.Username)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error fetching first tournament.",
			},
		})
		return
	}

	imageBytes, err := readFileIntoByteSlice("./img/banner_from_stats.png")
	if err != nil {
		fmt.Println("Error reading image file:", err)
		return
	}

	msg, err := s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
		Content: "Uploading banner...",
		Files: []*discordgo.File{
			{
				Name:   "banner.png",
				Reader: bytes.NewReader(imageBytes),
			},
		},
	})
	if err != nil {
		fmt.Println("Error uploading image:", err)
		return
	}

	bannerURL := msg.Attachments[0].URL
	fmt.Println(bannerURL)
	statsEmbed := embed.CreateStatsEmbed(response.Username, stats.MatchesPlayed, winPercentage, opponent, response.AvatarURL, bestTeammate, bestTournament, firstTournament, bannerURL)

	// Respond to the user with the embed
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{statsEmbed},
		},
	})
	err = s.ChannelMessageDelete(i.ChannelID, msg.ID)
	if err != nil {
		fmt.Println("Error deleting message:", err)
		return
	}
}

func readFileIntoByteSlice(filepath string) ([]byte, error) {
	return ioutil.ReadFile(filepath)
}
