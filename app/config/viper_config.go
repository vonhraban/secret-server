package config

import (
	"github.com/spf13/viper"

	"fmt"
)

type ViperConfig struct {
	path   string
	format string
	name   string
}

func NewViperConfig(name string, path string, format string) (*ViperConfig, error) {
	viper.SetConfigName(name)
	viper.AddConfigPath(path)
	viper.SetConfigType(format)
	err := viper.ReadInConfig()
	if err != nil {
		// TODO! Return error
		return nil, fmt.Errorf("Error reading config: %s", err)
	}
	return &ViperConfig{
		name:   name,
		format: format,
		path:   path,
	}, nil
}

func (cfg *ViperConfig) GetString(key string) (string, error) {
	return viper.GetString(key), nil
}

func (cfg *ViperConfig) GetBool(key string) (bool, error) {
	return viper.GetBool(key), nil
}

func (cfg *ViperConfig) GetInt(key string) (int, error) {
	return viper.GetInt(key), nil
}
