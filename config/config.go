package config

import "os"

type Config struct {
	APIKey string
}

func NewConfigWithEnv() *Config {
	c := &Config{}

	if apiKey, ok := os.LookupEnv("API_KEY"); ok {
		c.APIKey = apiKey
	}

	return c
}
