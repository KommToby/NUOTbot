// /config/config.go

package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Discord struct {
		BotToken string `json:"BOT_TOKEN"`
	} `json:"DISCORD"`
}

func LoadConfig(path string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
