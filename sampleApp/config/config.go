package config

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Provider looks for configuration in environment variables
type Provider interface {
	ConfigFileUsed() string
	Get(key string) interface{}
	GetBool(key string) bool
	GetDuration(key string) time.Duration
	GetFloat64(key string) float64
	GetInt(key string) int
	GetInt64(key string) int64
	GetSizeInBytes(key string) uint
	GetString(key string) string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string
	GetStringSlice(key string) []string
	GetTime(key string) time.Time
	InConfig(key string) bool
	IsSet(key string) bool
}

// defaultConfig is a singleton to expose configuration
var defaultConfig *viper.Viper

// Config returns the singleton exposing local configuration
func Config() Provider {
	return defaultConfig
}

func init() {
	defaultConfig = readViperConfig()
}

func loadDotEnv() {
	_ = godotenv.Load(".env")
}
func readViperConfig() *viper.Viper {
	loadDotEnv()
	v := viper.New()
	v.AutomaticEnv()
	return v
}
