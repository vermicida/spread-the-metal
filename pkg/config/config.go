package config

import (
	"os"
)

// TwitterConfig struct represents the Twitter configuration
type TwitterConfig struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
	Hashtags       string
}

// DynamoDBConfig struct represents the DynamoDB configuration
type DynamoDBConfig struct {
	DefaultRegion           string
	SongsTableName          string
	SongsTableKeyDateFormat string
	SongsTableKeyHourFormat string
}

// Config struct represents the application configuration
type Config struct {
	Twitter  TwitterConfig
	DynamoDB DynamoDBConfig
}

// New creates a new application configuration struct
func New() *Config {
	return &Config{
		Twitter: TwitterConfig{
			ConsumerKey:    getEnv("CONSUMER_KEY", ""),
			ConsumerSecret: getEnv("CONSUMER_SECRET", ""),
			AccessToken:    getEnv("ACCESS_TOKEN", ""),
			AccessSecret:   getEnv("ACCESS_SECRET", ""),
			Hashtags:       getEnv("STATUS_HASHTAGS", ""),
		},
		DynamoDB: DynamoDBConfig{
			DefaultRegion:           getEnv("DEFAULT_REGION", "eu-west-1"),
			SongsTableName:          getEnv("SONGS_TABLE_NAME", "stm-songs"),
			SongsTableKeyDateFormat: getEnv("DATE_FORMAT", "20060102"),
			SongsTableKeyHourFormat: getEnv("TIME_FORMAT", "15"),
		},
	}
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
