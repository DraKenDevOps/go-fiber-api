package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	EnvMode        string
	ServiceName    string
	Host           string
	Port           string
	BasePath       string
	Timezone       string
	DBURI          string
	DBPass         string
	EncryptionKey  string
	JWTPrivateKey  string
	JWTPublicKey   string
	UploadLimit    int
	ImageCompress  int
	RedisURI       string
	MQTTHost       string
	MQTPPort       string
	MQTTProto      string
	MQTTUser       string
	MQTPPassword   string
	MQTTPath       string
	MQTTTopic      string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get all configuration values from environment variables
	return &Config{
		EnvMode:      getEnv("ENV_MODE", "development"),
		ServiceName:  getEnv("SERVICE_NAME", "demo-rest-api"),
		Host:         getEnv("HOST", "0.0.0.0"),
		Port:         getEnv("PORT", "8000"),
		BasePath:     getEnv("BASE_PATH", "api"),
		Timezone:     getEnv("TZ", "Asia/Bangkok"),
		DBURI:        getEnv("DB_URI", ""),
		DBPass:       getEnv("DB_PASS", ""),
		EncryptionKey:getEnv("ENCRYPTION_KEY", ""),
		JWTPrivateKey:getEnv("JWT_PRIVATE_KEY", ""),
		JWTPublicKey: getEnv("JWT_PUBLIC_KEY", ""),
		UploadLimit:  getEnvAsInt("UPLOAD_LIMIT_SIZE", 10),
		ImageCompress:getEnvAsInt("IMAGE_COMPRESS_LEVEL", 70),
		RedisURI:     getEnv("REDIS_URI", ""),
		MQTTHost:     getEnv("MQTT_HOST", ""),
		MQTPPort:     getEnv("MQTT_PORT", ""),
		MQTTProto:    getEnv("MQTT_PROTOCOL", ""),
		MQTTUser:     getEnv("MQTT_USER", ""),
		MQTPPassword: getEnv("MQTT_PASSWORD", ""),
		MQTTPath:     getEnv("MQTT_PATH", ""),
		MQTTTopic:    getEnv("MQTT_TOPIC", ""),
	}
}

// Helper function to get environment variable with fallback
func getEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// Helper function to get environment variable as integer with fallback
func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		fmt.Sscanf(value, "%d", &result)
		return result
	}
	return fallback
}