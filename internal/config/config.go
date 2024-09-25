package config

import (
	"log"
	"os"
)

// Config holds the application configuration
type Config struct {
	DbUser         string
	DbPassword     string
	DbHost         string
	DbPort         string
	DbName         string
	CACertPath     string
	ClientCertPath string
	ClientKeyPath  string
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		DbUser:         getEnvOrPanic("DB_USER"),
		DbPassword:     getEnvOrPanic("DB_PASSWORD"),
		DbHost:         getEnvOrPanic("DB_HOST"),
		DbPort:         getEnvOrPanic("DB_PORT"),
		DbName:         getEnvOrPanic("DB_NAME"),
		CACertPath:     getEnv("CACERT_PATH", ""),
		ClientCertPath: getEnv("CLIENT_CERT_PATH", ""),
		ClientKeyPath:  getEnv("CLIENT_KEY_PATH", ""),
	}
}

// getEnv retrieves an environment variable or returns a default value if not set
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

// getEnvOrPanic retrieves an environment variable or panics if none is found
func getEnvOrPanic(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Panicf("required environment variable %s is not set", key)
	}
	return value
}
