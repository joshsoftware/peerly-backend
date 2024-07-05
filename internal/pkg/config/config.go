package config

import (
	"fmt"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	"github.com/spf13/viper"
)

var (
	appName string
	appPort int
)

// Load - loads all the environment variables and/or params in application.yml
func Load() {

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		panic(fmt.Sprintf("Error loading .env file: %v", err))
	}

	viper.ReadInConfig()
	viper.AutomaticEnv()

	viper.SetDefault(constants.AppName, "app")
	viper.SetDefault(constants.AppPort, "8002")

	// Check for the presence of JWT_KEY and JWT_EXPIRY_DURATION_HOURS
	JWTKey()
	JWTExpiryDurationHours()
}

// AppName - returns the app name
func AppName() string {
	if appName == "" {
		appName = ReadEnvString(constants.AppName)
	}
	return appName
}

// AppPort - returns application http port
func AppPort() int {
	if appPort == 0 {
		appPort = ReadEnvInt(constants.AppPort)
	}
	return appPort
}

// JWTKey - returns the JSON Web Token key
func JWTKey() []byte {
	return []byte(ReadEnvString(constants.JWTSecret))
}

// JWTExpiryDurationHours - returns duration for jwt expiry in int
func JWTExpiryDurationHours() int {
	return int(ReadEnvInt(constants.JWTExpiryDurationHours))
}

// ReadEnvInt - reads an environment variable as an integer
func ReadEnvInt(key string) int {
	checkIfSet(key)
	v, err := strconv.Atoi(viper.GetString(key))
	if err != nil {
		panic(fmt.Sprintf("key %s is not a valid integer", key))
	}
	return v
}

// ReadEnvString - reads an environment variable as a string
func ReadEnvString(key string) string {
	checkIfSet(key)
	return viper.GetString(key)
}

// ReadEnvBool - reads environment variable as a boolean
func ReadEnvBool(key string) bool {
	checkIfSet(key)
	return viper.GetBool(key)
}

func checkIfSet(key string) {
	if !viper.IsSet(key) {
		panic(apperrors.ErrKeyNotSet(key))
	}
}
