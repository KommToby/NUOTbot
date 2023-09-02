package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/kommtoby/NUOTbot/config"
)

func main() {
	// Open and read the config file
	configFile, err := os.Open("config/config.json")
	if err != nil {
		fmt.Println("Error opening config file:", err)
		return
	}
	defer configFile.Close()

	// Decode JSON content into the Config struct
	var conf config.Config
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&conf)
	if err != nil {
		fmt.Println("Error decoding config:", err)
		return
	}

	// Create a new Discord session
	dg, err := discordgo.New("Bot " + conf.Discord.BotToken)
	if err != nil {
		fmt.Println("Error creating Discord session:", err)
		return
	}

	// Add message create handler
	dg.AddHandler(messageCreate)

	// Open Discord session
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection:", err)
		return
	}

	fmt.Println("Bot is now running. Press Ctrl+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Close Discord session
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
}
