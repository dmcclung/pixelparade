package models

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/dmcclung/pixelparade/migrations"
	"github.com/pressly/goose/v3"
)

type PostgresConfig struct {
	host     string
	port     string
	user     string
	password string
	dbname   string
	sslmode  string
}

var DefaultPostgresConfig = PostgresConfig{
	host:     "localhost",
	port:     "5432",
	user:     "admin",
	password: "admin",
	dbname:   "pixelparade",
	sslmode:  "disable",
}

func (pg PostgresConfig) Open() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		pg.host, pg.port, pg.user, pg.password, pg.dbname, pg.sslmode)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("opening connection to db: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("connecting to db: %w", err)
	}

	log.Printf("Connected to postgres host %s and port %s", pg.host, pg.port)
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
