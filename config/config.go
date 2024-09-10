package config

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

var (
	envServerPort string
	envDBString   string

	envJwtSecret          string
	envJwtIssuer          string
	envJwtExpiryInSeconds int

	envEmailHost      string
	envEmailPort      string
	envEmailUsername  string
	envEmailPassword  string
	envEmailOTPExpiry string

	envHelpCenterEmail   string
	envHelpCenterAddress string

	DevMode bool
)

func Init() {
	// Load the environment variables
	DevMode = os.Getenv("RASTA_DEV_MODE") == "true"
	if err := loadEnv(); DevMode && err != nil {
		panic(err)
	}
	log.Println("configs successfully loaded")
}

func loadEnv() error {
	if DevMode {
		err := godotenv.Load(".env")
		fmt.Println("Error load .env", err)
		if err != nil {
			return err
		}
	}

	var err error
	envServerPort, err = getEnv("SERVER_PORT", "3080")
	if err != nil {
		return err
	}
	envDBString, err = getEnv("DB_STRING", "")

	if err != nil {
		return err
	}

	envJwtSecret, err = getEnv("JWT_SECRET", "")
	if err != nil {
		return err
	}
	envJwtIssuer, err = getEnv("JWT_ISSUER", "")
	if err != nil {
		return err
	}
	EnvJwtExpiryStr, err := getEnv("JWT_EXPIRY", "")
	if err != nil {
		return err
	}
	envJwtExpiryInSeconds, err = strconv.Atoi(EnvJwtExpiryStr)
	if err != nil {
		return err
	}

	envEmailHost, err = getEnv("EMAIL_HOST", "")
	if err != nil {
		return err
	}
	envEmailPort, err = getEnv("EMAIL_PORT", "")
	if err != nil {
		return err
	}
	envEmailUsername, err = getEnv("EMAIL_USERNAME", "")
	if err != nil {
		return err
	}
	envEmailPassword, err = getEnv("EMAIL_PASSWORD", "")
	if err != nil {
		return err
	}
	envEmailOTPExpiry, err = getEnv("EMAIL_OTP_EXPIRY", "")
	if err != nil {
		return err
	}

	envHelpCenterEmail, err = getEnv("HELP_CENTER_EMAIL", "")
	if err != nil {
		return err
	}
	envHelpCenterAddress, err = getEnv("HELP_CENTER_ADDRESS", "")
	if err != nil {
		return err
	}

	return nil
}

func getEnv(key string, defaultVal string) (string, error) {
	val, ok := os.LookupEnv(key)
	if val != "" || ok {
		return val, nil
	}
	if defaultVal != "" && DevMode {
		return defaultVal, nil
	}
	return "", errors.New("failed to find environment variable: " + key)
}

func GetServerPort() string {
	return envServerPort
}

func GetDBString() string {
	return envDBString
}

func GetJwtSecret() string {
	return envJwtSecret
}

func GetJwtIssuer() string {
	return envJwtIssuer
}

func GetJwtExpiry() int {
	if envJwtExpiryInSeconds == 0 {
		envJwtExpiryInSeconds = 3600
	}

	return envJwtExpiryInSeconds
}

func GetEmailHost() string {
	return envEmailHost
}

func GetEmailPort() int {
	parseInt, err := strconv.ParseInt(envEmailPort, 10, 64)
	if err != nil {
		return 0
	}
	return int(parseInt)
}

func GetEmailUsername() string {
	return envEmailUsername
}

func GetEmailPassword() string {
	return envEmailPassword
}

func GetEnvEmailOTPExpiry() int {
	parseInt, err := strconv.ParseInt(envEmailOTPExpiry, 10, 64)
	if err != nil {
		return 0
	}
	return int(parseInt)
}

func GetHelpCenterEmail() string {
	return envHelpCenterEmail
}

func GetHelpCenterAddress() string {
	return envHelpCenterAddress
}
func GetEnvVars() map[string]any {
	return map[string]any{
		"SERVER_PORT":      envServerPort,
		"DB_STRING":        envDBString,
		"JWT_SECRET":       envJwtSecret,
		"JWT_ISSUER":       envJwtIssuer,
		"JWT_EXPIRY":       envJwtExpiryInSeconds,
		"EMAIL_HOST":       envEmailHost,
		"EMAIL_PORT":       envEmailPort,
		"EMAIL_USERNAME":   envEmailUsername,
		"EMAIL_PASSWORD":   envEmailPassword,
		"EMAIL_OTP_EXPIRY": envEmailOTPExpiry,
	}
}
