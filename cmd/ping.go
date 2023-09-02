// /cmd/ping.go

package cmd

import "github.com/bwmarrin/discordgo"

var PingCommand = &discordgo.ApplicationCommand{
	Name:        "ping",
	Description: "Returns pong",
}
