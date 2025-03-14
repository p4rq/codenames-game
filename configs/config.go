package configs

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Game     GameConfig
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Host         string
	Port         string
	ReadTimeout  int
	WriteTimeout int
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	URI      string
	Database string
}

// GameConfig holds game-specific configuration
type GameConfig struct {
	DefaultTeamSize int
	MaxPlayers      int
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Try to load .env file if it exists
	_ = godotenv.Load()

	return &Config{
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "localhost"),
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getEnvAsInt("SERVER_READ_TIMEOUT", 10),
			WriteTimeout: getEnvAsInt("SERVER_WRITE_TIMEOUT", 10),
		},
		Database: DatabaseConfig{
			URI:      getEnv("DB_URI", ""),
			Database: getEnv("DB_NAME", "codenames"),
		},
		Game: GameConfig{
			DefaultTeamSize: getEnvAsInt("GAME_DEFAULT_TEAM_SIZE", 4),
			MaxPlayers:      getEnvAsInt("GAME_MAX_PLAYERS", 10),
		},
	}
}

// Helper functions to get environment variables
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		intVal, err := strconv.Atoi(value)
		if err != nil {
			log.Printf("Error parsing %s as int: %v. Using default: %d", key, err, defaultValue)
			return defaultValue
		}
		return intVal
	}
	return defaultValue
}

func getEnvAsFloat(key string, defaultValue float64) float64 {
	if value, exists := os.LookupEnv(key); exists {
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			log.Printf("Error parsing %s as float: %v. Using default: %f", key, err, defaultValue)
			return defaultValue
		}
		return floatVal
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		durationVal, err := time.ParseDuration(value)
		if err != nil {
			log.Printf("Error parsing %s as duration: %v. Using default: %v", key, err, defaultValue)
			return defaultValue
		}
		return durationVal
	}
	return defaultValue
}
