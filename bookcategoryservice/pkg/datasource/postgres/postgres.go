package db

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"github.com/sir-shalahuddin/grpc-learn/bookcategoryservice/config"
)

func NewDB(config config.DBConfig) (*sql.DB, error) {

	port, err := strconv.ParseUint(config.Port, 10, 32)
	if err != nil {
		panic("failed to parse database port")
	}

	DB_URI := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		port,
		config.User,
		config.Pass,
		config.Name)

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
	log.Println("Database Connected ")

	return db, nil
}
