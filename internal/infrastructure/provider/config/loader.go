package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func NewConfig() (*Config, error) {
	var cfg Config

	viper.SetEnvPrefix("FUSION")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
