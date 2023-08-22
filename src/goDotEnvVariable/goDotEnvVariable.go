package goDotEnvVariable

import (
	"log"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

func GoDotEnvVariable(key string) string {
	projectName := regexp.MustCompile(`^(.*` + "kafka-realtime-data-processing" + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))

	// load .env file
	err := godotenv.Load(string(rootPath) + `/.env`)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv(key)
}
