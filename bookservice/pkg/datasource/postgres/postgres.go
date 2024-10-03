package db

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"github.com/sir-shalahuddin/grpc-learn/bookservice/config"
)

func NewDB(config config.DBConfig) (*sql.DB, error) {
	port, err := strconv.ParseUint(config.Port, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database port: %w", err)
	}

	var DB_URI string
	if config.UseUnixSocket {
		DB_URI = fmt.Sprintf(
			"user=%s password=%s dbname=%s sslmode=disable host=/cloudsql/%s",
			config.User,
			config.Pass,
			config.Name,
			config.InstanceConnectionName)
	} else {
		DB_URI = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			config.Host,
			port,
			config.User,
			config.Pass,
			config.Name)
	}

	db, err := sql.Open("postgres", DB_URI)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(20)
	db.SetConnMaxIdleTime(60 * time.Minute)
	db.SetConnMaxLifetime(10 * time.Minute)

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	connectionType := "SQL Connector"
	if config.UseUnixSocket {
		connectionType = "Unix Socket"
	}
	log.Printf("Database Connected using: %s", connectionType)

	return db, nil
}
