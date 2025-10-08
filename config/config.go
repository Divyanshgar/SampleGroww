package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Email    EmailConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type EmailConfig struct {
	SMTPHost  string
	SMTPPort  string
	Username  string
	Password  string
	FromEmail string
	FromName  string
}

func LoadConfig() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8443"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "onboarding_app"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Email: EmailConfig{
			SMTPHost:  getEnv("SMTP_HOST", "smtp.gmail.com"),
			SMTPPort:  getEnv("SMTP_PORT", "587"),
			Username:  getEnv("SMTP_USERNAME", "bishtdaksh77@gmail.com"),
			Password:  getEnv("SMTP_PASSWORD", "irle kjkt soxg nxvt"),
			FromEmail: getEnv("SMTP_FROM_EMAIL", "bishtdaksh77@gmail.com"),
			FromName:  getEnv("SMTP_FROM_NAME", "GoPro"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
