package config

import "os"

type Config struct {
	Token string
}

func NewConfigWithEnv() *Config {
	c := &Config{}

	if token, ok := os.LookupEnv("BOT_TOKEN"); ok {
		c.Token = token
	}

	return c
}
