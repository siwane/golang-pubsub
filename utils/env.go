package utils

import (
    "os"
    "log"

    "github.com/joho/godotenv"
)

func Load() {
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatal(err)
    }
    log.Printf(".env file loaded.")
}

func MustGetenv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%s environment variable not set.", key)
	}
	return value
}

func Getenv(key string, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        value = defaultValue
    }
    return value
}