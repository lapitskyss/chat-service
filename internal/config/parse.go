package config

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"

	"github.com/lapitskyss/chat-service/internal/validator"
)

func ParseAndValidate(filename string) (Config, error) {
	// Read file
	file, err := os.ReadFile(filename)
	if err != nil {
		return Config{}, fmt.Errorf("read file: %v", err)
	}
	// Decode config
	var cfg Config
	err = toml.Unmarshal(file, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("decode file: %v", err)
	}
	// Validate config
	err = validator.Validator.Struct(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("validate: %v", err)
	}
	return cfg, nil
}
