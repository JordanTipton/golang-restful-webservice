package app

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config stores global application settings
var Config appConfig

type appConfig struct {
	DSN                string `json:"dsn"`
	JWTSigningMethod   string `json:"jwtSigningMethod"`
	JWTSigningKey      string `json:"jwtSigningKey"`
	JWTVerificationKey string `json:"jwtVerificationKey"`
	LogPath            string `json:"logPath"`
	ServerPort         int    `json:"serverPort"`
}

// LoadConfig loads configuration from a specified json configuration file
func LoadConfig(configPath string) error {
	configFile, err := os.Open(configPath)
	defer configFile.Close()
	if err != nil {
		return fmt.Errorf("Failed to read the configuration file: %s", err)
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&Config)
	return nil
}
