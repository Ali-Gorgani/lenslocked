package models

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func DefultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "root",
		Password: "110963",
		Database: "lenslocked",
		SSLMode:  "disable",
	}
}

func (cfg PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode,
	)
}

// Open opens a database connection using the provided PostgresConfig.
// It attempts to establish a connection and verify it with a ping.
// Caller must ensure that the connection is closed via db.Close() method.
func Open(cfg PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.String())
	if err != nil {
		return nil, fmt.Errorf("Open: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Open: ping failed: %w", err)
	}
	fmt.Println("Database connected!")
	return db, nil
}
