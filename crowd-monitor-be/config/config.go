package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	TwilioAccountSid  string
	TwilioAuthToken   string
	TwilioServiceSid  string
	FirebaseCredPath  string
	FirebaseDatabaseURL string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	config := &Config{
		Port:              getEnv("PORT", "8080"),
		TwilioAccountSid:  getEnv("TWILIO_ACCOUNT_SID", ""),
		TwilioAuthToken:   getEnv("TWILIO_AUTH_TOKEN", ""),
		TwilioServiceSid:  getEnv("TWILIO_SERVICE_SID", ""),
		FirebaseCredPath:  getEnv("FIREBASE_CRED_PATH", ""),
		FirebaseDatabaseURL: getEnv("FIREBASE_DATABASE_URL", ""),
	}

	return config
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}