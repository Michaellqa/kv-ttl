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

// NewPostgresDb connects to given database and executes initial schema migrations
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

// ensureDbExists connects to given data source and creates the application database
// if it doesn't exist int the system table of databases
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
		// cannot use sql parameters here
		_, err = db.Exec(fmt.Sprintf("create database %s", c.Database))
		return err
	}
	return nil
}
