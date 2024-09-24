package config

import (
	"log"
	"os"
	"strings"
)

type AppConfig struct {
	GRPCPort string
	RESTPort string
	Mode     string
}

type DBConfig struct {
	Host                   string
	Port                   string
	User                   string
	Pass                   string
	Name                   string
	InstanceConnectionName string
	UseUnixSocket          bool
}

type JWTConfig struct {
	Secret string
}

type GRPCConfig struct {
	AuthAddress string
}

func GetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%s not provided on .env", key)
	}
	return value
}

func GetEnvAsBool(key string) bool {
	value := os.Getenv(key)
	if value == "" {
		return false
	}
	if strings.ToLower(value) == "true" {
		return true
	}
	return false
}
