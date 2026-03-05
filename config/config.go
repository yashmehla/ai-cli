package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GeminiAPIKey string
}

func Load() *Config {

	err := godotenv.Load()

	if err != nil {
		log.Println(".env not found, using system env")
	}

	return &Config{
		GeminiAPIKey: os.Getenv("GEMINI_API_KEY"),
	}
}
