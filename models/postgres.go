package models

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/dmcclung/pixelparade/migrations"
	"github.com/pressly/goose/v3"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

var DefaultPostgresConfig = PostgresConfig{
	Host:     "localhost",
	Port:     "5432",
	User:     "admin",
	Password: "admin",
	Database: "pixelparade",
	SSLMode:  "disable",
}

func Open(pg PostgresConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		pg.Host, pg.Port, pg.User, pg.Password, pg.Database, pg.SSLMode)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("opening connection to db: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("connecting to db: %w", err)
	}

	log.Printf("Connected to postgres host %s and port %s", pg.Host, pg.Port)
	return db, nil
}

func Migrate(db *sql.DB) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	goose.SetBaseFS(migrations.FS)
	err = goose.Up(db, ".")
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	return nil
}
