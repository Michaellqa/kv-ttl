package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

const migrationsDir = "./repository/postgres/migrations"

type Connection struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func NewPostgresDb(c Connection) (*sql.DB, error) {
	if err := ensureDbExists(c); err != nil {
		return nil, err
	}
	connection := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.User, c.Password, c.Database)
	db, err := sql.Open("postgres", connection)
	if err != nil {
		return nil, err
	}
	if err = goose.Up(db, migrationsDir); err != nil {
		return nil, err
	}
	return db, nil
}

func ensureDbExists(c Connection) error {
	connection := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password)
	db, err := sql.Open("postgres", connection)
	if err != nil {
		return err
	}
	var dbExists bool
	err = db.QueryRow(`select exists(select * from pg_database where datname = $1)`, c.Database).Scan(&dbExists)
	if !dbExists {
		_, err = db.Exec(`create database cache_app;`)
		return err
	}
	return nil
}
