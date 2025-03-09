package config

import "github.com/caarlos0/env/v11"

type Config struct {
	Token        string `env:"TOKEN,notEmpty"`
	AppID        string `env:"APP_ID,notEmpty"`
	GuildID      string `env:"GUILD_ID,notEmpty"`
	LogChannelID string `env:"LOG_CHANNEL_ID,notEmpty"`
}

func Load() (Config, error) {
	var config Config
	err := env.Parse(&config)

	return config, err
}
