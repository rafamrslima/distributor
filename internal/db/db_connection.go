package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	dbPool   *pgxpool.Pool
	poolOnce sync.Once
	poolErr  error
)

func GetDB() (*pgxpool.Pool, error) {
	poolOnce.Do(func() {
		connString := os.Getenv("DATABASE_CONNECTION_STRING")
		if connString == "" {
			poolErr = errors.New("DATABASE_CONNECTION_STRING environment variable not set")
			return
		}

		config, err := pgxpool.ParseConfig(connString)
		if err != nil {
			poolErr = fmt.Errorf("failed to parse connection string: %w", err)
			return
		}

		// Configure connection pool
		config.MaxConns = 25
		config.MinConns = 5
		config.MaxConnLifetime = 5 * time.Minute
		config.MaxConnIdleTime = 1 * time.Minute
		config.HealthCheckPeriod = 30 * time.Second

		dbPool, err = pgxpool.NewWithConfig(context.Background(), config)
		if err != nil {
			poolErr = fmt.Errorf("failed to create connection pool: %w", err)
			return
		}

		// Verify connection
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := dbPool.Ping(ctx); err != nil {
			dbPool.Close()
			poolErr = fmt.Errorf("failed to ping database: %w", err)
			return
		}

		log.Println("Database connection pool initialized successfully")
	})

	if poolErr != nil {
		return nil, poolErr
	}

	return dbPool, nil
}

func CloseDB() {
	if dbPool != nil {
		dbPool.Close()
		log.Println("Database connection pool closed")
	}
}
