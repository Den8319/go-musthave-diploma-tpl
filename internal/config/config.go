package config

import (
	"flag"
	"os"
)

const (
	envServerAddress     = "RUN_ADDRESS"
	flagServerAddress    = "a"
	defaultServerAddress = ":8080"

	 
	envLogLevel     = "LOG_LEVEL"
	flagLogLevel    = "l"
	defaultLogLevel = "Info"

	envDatabaseDSN     = "DATABASE_URI"
	flagDatabaseDSN    = "d"
	defaultDatabaseDSN = ""

	envSecretKey     = "SECRET_KEY"
	flagSecretKey    = "k"
	defaultSecretKey = "123"
)

type Config struct {
	ServerAddress   string
	BaseURL         string
	LogLevel        string
	DatabaseDSN     string
	SecretKey       string
}

func getParam(envName, flagValue string) string {
	if envValue, exists := os.LookupEnv(envName); exists && envValue != "" {
		return envValue
	}

	return flagValue
}

func New() *Config {
	serverAddr := flag.String(flagServerAddress, defaultServerAddress, "")
	logLevel := flag.String(flagLogLevel, defaultLogLevel, "")
	databaseDSN := flag.String(flagDatabaseDSN, defaultDatabaseDSN, "")
	secretKey := flag.String(flagSecretKey, defaultSecretKey, "")

	flag.Parse()

	cfg := &Config{
		ServerAddress:   getParam(envServerAddress, *serverAddr),
		LogLevel:        getParam(envLogLevel, *logLevel),
		DatabaseDSN:     getParam(envDatabaseDSN, *databaseDSN),
		SecretKey:       getParam(envSecretKey, *secretKey),
	}

	return cfg
}
