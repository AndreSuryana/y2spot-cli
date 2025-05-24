package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads the .env file and checks required variables
func LoadEnv() error {
	_ = godotenv.Load()

	if os.Getenv("SPOTIFY_CLIENT_ID") == "" || os.Getenv("SPOTIFY_CLIENT_SECRET") == "" {
		return fmt.Errorf("missing SPOTIFY_CLIENT_ID or SPOTIFY_CLIENT_SECRET environment variables")
	}
	return nil
}

// GetSpotifyCredentials returns (clientID, clientSecret)
func GetSpotifyCredentials() (string, string) {
	return os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET")
}
