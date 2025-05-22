package clients

import (
	"backend/src/config"
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/lib/pq"
)

var (
	postgresDB   *sql.DB
	postgresOnce sync.Once
)

// NewPostgreSQLClient ensures thread-safe initialization of the PostgreSQL database connection.
func NewPostgreSQLClient() (*sql.DB, error) {
	var err error
	postgresOnce.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			config.GetEnv("POSTGRES_HOST", "localhost"),
			config.GetEnv("POSTGRES_PORT", "5432"),
			config.GetEnv("POSTGRES_USER", "postgres"),
			config.GetEnv("POSTGRES_PASSWORD", "password"),
			config.GetEnv("POSTGRES_DB", "postgres"),
		)

		postgresDB, err = sql.Open("postgres", dsn)

		if err == nil {
			// Test connection
			err = postgresDB.Ping()
		}
	})

	if err != nil {
		return nil, err
	}
	return postgresDB, nil
}
