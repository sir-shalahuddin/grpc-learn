package config

import (
	"log"
	"os"
)

type AppConfig struct {
	GRPCPort string
	RESTPort string
}

type DBConfig struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

type JWTConfig struct {
	Secret string
}

type GRPCConfig struct {
	AuthAddress     string
	CategoryAddress string
}

func GetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%s not provided on .env", key)
	}
	return value
}
