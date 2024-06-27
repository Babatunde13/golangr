package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func GetConfig() map[string]string {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	config := make(map[string]string)
	for _, line := range os.Environ() {
		kv := strings.Split(line, "=")
		config[kv[0]] = kv[1]
	}
	return config
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
