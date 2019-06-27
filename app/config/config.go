package config

import (
	"encoding/json"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServerPort    int    `envconfig:"SERVER_PORT"`
	ApiVersion    string `envconfig:"API_VERSION"`
	MongoHost     string `envconfig:"MONGO_HOST"`
	MongoPort     int    `envconfig:"MONGO_PORT"`
	MongoDatabase string `envconfig:"MONGO_DATABASE"`
	MongoUsername string `envconfig:"MONGO_USERNAME"`
	MongoPassword string `envconfig:"MONGO_PASSWORD"`
	LogLevel      string `envconfig:"LOG_LEVEL"`
}

func New() (*Config, error) {
	c := &Config{}
	if err := envconfig.Process("", c); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) Printable() string {
	printableConfig := *c
	printableConfig.MongoPassword = "***hidden***"
	marshaled, _ := json.Marshal(printableConfig)
	return string(marshaled)
}
