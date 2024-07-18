package config

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
)

func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	return nil
}

func GetJWTKey() string {
	return os.Getenv("JWT_SECRET")
}

func GetPort() (string, error) {
	if port := os.Getenv("HTTP_PORT"); port == "" {
		return "", errors.New("port not set")
	} else {
		return port, nil
	}
}
