package config

import (
	"flag"
	"os"
	"crypto/rand"
)

const (
	envServerAddress     = "RUN_ADDRESS"
	flagServerAddress    = "a"
	defaultServerAddress = "localhost:8000"

	 
	envLogLevel     = "LOG_LEVEL"
	flagLogLevel    = "l"
	defaultLogLevel = "Info"

	envDatabaseDSN     = "DATABASE_URI"
	flagDatabaseDSN    = "d"
	defaultDatabaseDSN = ""

	envSecretKey     = "SECRET_KEY"
	flagSecretKey    = "k"
	

	envAccrualSystemAddress = "ACCRUAL_SYSTEM_ADDRESS"
	flagAccrualSystemAddress = "r"
	defaultAccrualSystemAddress = "http://localhost:8080"

)

type Config struct {
	ServerAddress   string
	BaseURL         string
	LogLevel        string
	DatabaseDSN     string
	SecretKey       string
	AccrualSystemAddress  string
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
	secretKey := flag.String(flagSecretKey, rand.Text(), "") // Если секретный ключ не указан, генерируем случайный 
	AccrualSystemAddress := flag.String(flagAccrualSystemAddress, defaultAccrualSystemAddress, "")

	flag.Parse()

	cfg := &Config{
		ServerAddress:   getParam(envServerAddress, *serverAddr),
		LogLevel:        getParam(envLogLevel, *logLevel),
		DatabaseDSN:     getParam(envDatabaseDSN, *databaseDSN),
		SecretKey:       getParam(envSecretKey, *secretKey),
		AccrualSystemAddress:   getParam(envAccrualSystemAddress, *AccrualSystemAddress),
	}

	return cfg
}
