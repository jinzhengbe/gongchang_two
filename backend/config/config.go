package config

import (
	"os"
)

type Config struct {
	Port       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	ServerHost string
}

func LoadConfig() (*Config, error) {
	return &Config{
		Port:       getEnv("PORT", "443"),
		DBHost:     getEnv("DB_HOST", "mysql"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "123456"),
		DBName:     getEnv("DB_NAME", "gongchang"),
		JWTSecret:  getEnv("JWT_SECRET", "your-secret-key"),
		ServerHost: getEnv("SERVER_HOST", "aneworder.com"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
} 