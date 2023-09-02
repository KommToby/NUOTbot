// config/config.go

package config

type Config struct {
	Discord struct {
		BotToken string `json:"BOT_TOKEN"`
	} `json:"DISCORD"`
}
