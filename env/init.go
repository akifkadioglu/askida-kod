package env

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

const (
	DB_CONNECTION   = "DB_CONNECTION"
	DB_HOST         = "DB_HOST"
	DB_PORT         = "DB_PORT"
	DB_DATABASE     = "DB_DATABASE"
	DB_USERNAME     = "DB_USERNAME"
	DB_PASSWORD     = "DB_PASSWORD"
	DB_EXTERNAL_URL = "DB_EXTERNAL_URL"

	APP_KEY = "APP_KEY"
)

const (
	LOCAL int = iota
	PROD
)

func InitEnv(t int) {
	var err error

	switch t {
	case LOCAL:
		err = godotenv.Load(filepath.Join("./env", ".env"))

	case PROD:
		err = godotenv.Load(filepath.Join("/etc/secrets", ".env"))
	}

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func Setenv(key string, value string) {
	os.Setenv(key, value)
}

func Getenv(key string) string {
	return os.Getenv(key)
}
