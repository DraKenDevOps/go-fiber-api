package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Cwd           string
	EnvMode       string
	ServiceName   string
	Host          string
	Port          string
	BasePath      string
	TZ            string
	AppVersion    string
	DBURI         string
	DBPass        string
	EncryptionKey string
	JWTPrivateKey string
	JWTPublicKey  string
	UploadLimit   int
	ImageCompress int
	RedisURI      string
	MQTTHost      string
	MQTPPort      string
	MQTTProto     string
	MQTTUser      string
	MQTPPassword  string
	MQTTPath      string
	MQTTTopic     string
	Feature       FeatureFlag
	LogLevel      string
}

type FeatureFlag struct {
	LimitMaxBalance bool `json:"limitMaxBalance"`
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("%s: %+v", err.Error(), err)
	}

	return &Config{
		Cwd:           cwd,
		EnvMode:       getEnv("ENV_MODE", "development"),
		ServiceName:   getEnv("SERVICE_NAME", "demo-rest-api"),
		Host:          getEnv("HOST", "0.0.0.0"),
		Port:          getEnv("PORT", "8000"),
		BasePath:      getEnv("BASE_PATH", "api"),
		TZ:            getEnv("TZ", "Asia/Bangkok"),
		AppVersion:    getEnv("APP_VERSION", "1.0.0"),
		DBURI:         getEnv("DB_URI", ""),
		DBPass:        getEnv("DB_PASS", ""),
		EncryptionKey: getEnv("ENCRYPTION_KEY", ""),
		JWTPrivateKey: getEnv("JWT_PRIVATE_KEY", ""),
		JWTPublicKey:  getEnv("JWT_PUBLIC_KEY", ""),
		UploadLimit:   getEnvAsInt("UPLOAD_LIMIT_SIZE", 10),
		ImageCompress: getEnvAsInt("IMAGE_COMPRESS_LEVEL", 70),
		RedisURI:      getEnv("REDIS_URI", ""),
		MQTTHost:      getEnv("MQTT_HOST", ""),
		MQTPPort:      getEnv("MQTT_PORT", ""),
		MQTTProto:     getEnv("MQTT_PROTOCOL", ""),
		MQTTUser:      getEnv("MQTT_USER", ""),
		MQTPPassword:  getEnv("MQTT_PASSWORD", ""),
		MQTTPath:      getEnv("MQTT_PATH", ""),
		MQTTTopic:     getEnv("MQTT_TOPIC", ""),
		Feature: FeatureFlag{
			LimitMaxBalance: getEnvAsBool("LIMIT_MAX_BALANCE", false),
		},
		LogLevel: getEnv("LOG_LEVEL", "debug"),
	}
}

func getEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		fmt.Sscanf(value, "%d", &result)
		return result
	}
	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	v := os.Getenv(key)
	val, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return val
}
