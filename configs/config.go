package configs

import (
	"log"
	"os"
	"strconv"

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
	Port         string
	Host         string
	ReadTimeout  int
	WriteTimeout int
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Type     string
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// GameConfig holds game-specific configuration
type GameConfig struct {
	WordsPerGame       int
	BlueTeamStartsProb float64
	AssassinCount      int
	NeutralCount       int
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Try to load .env file if it exists
	_ = godotenv.Load()

	return &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			ReadTimeout:  getEnvAsInt("SERVER_READ_TIMEOUT", 10),
			WriteTimeout: getEnvAsInt("SERVER_WRITE_TIMEOUT", 10),
		},
		Database: DatabaseConfig{
			Type:     getEnv("DB_TYPE", "postgres"), // memory, postgres, mysql
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "11052004ARAd."),
			Name:     getEnv("DB_NAME", "codenames"),
		},
		Game: GameConfig{
			WordsPerGame:       getEnvAsInt("GAME_WORDS_PER_GAME", 25),
			BlueTeamStartsProb: getEnvAsFloat("GAME_BLUE_TEAM_STARTS_PROB", 0.5),
			AssassinCount:      getEnvAsInt("GAME_ASSASSIN_COUNT", 1),
			NeutralCount:       getEnvAsInt("GAME_NEUTRAL_COUNT", 7),
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
