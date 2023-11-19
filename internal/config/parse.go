package config

import (
	"fmt"

	"github.com/BurntSushi/toml"

	"github.com/lapitskyss/chat-service/internal/validator"
)

func ParseAndValidate(filename string) (Config, error) {
	var cfg Config
	if _, err := toml.DecodeFile(filename, &cfg); err != nil {
		return Config{}, fmt.Errorf("decode file: %v", err)
	}
	// Validate config
	err := validator.Validator.Struct(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("validate: %v", err)
	}
	return cfg, nil
}
