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
		Port:       getEnv("PORT", "8080"),
		DBHost:     getEnv("MYSQL_HOST", "mysql"),
		DBPort:     getEnv("MYSQL_PORT", "3306"),
		DBUser:     getEnv("MYSQL_USER", "gongchang"),
		DBPassword: getEnv("MYSQL_PASSWORD", "gongchang"),
		DBName:     getEnv("MYSQL_DATABASE", "gongchang"),
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