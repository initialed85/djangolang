package helpers

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func GetDSN() (string, error) {
	postgresUser := strings.TrimSpace(os.Getenv("POSTGRES_USER"))
	if postgresUser == "" {
		postgresUser = "postgres"
	}

	postgresPassword := strings.TrimSpace(os.Getenv("POSTGRES_PASSWORD"))
	if postgresPassword == "" {
		return "", fmt.Errorf("POSTGRES_PASSWORD env var empty or unset")
	}

	postgresHost := strings.TrimSpace(os.Getenv("POSTGRES_HOST"))
	if postgresHost == "" {
		postgresHost = "localhost"
	}

	postgresPort := strings.TrimSpace(os.Getenv("POSTGRES_PORT"))
	if postgresPort == "" {
		postgresPort = "5432"
	}

	postgresDatabase := strings.TrimSpace(os.Getenv("POSTGRES_DB"))
	if postgresDatabase == "" {
		return "", fmt.Errorf("POSTGRES_DB env var empty or unset")
	}

	postgresSSLModeString := "?sslmode=disable"
	if strings.TrimSpace(os.Getenv("POSTGRES_SSLMODE")) == "1" {
		postgresSSLModeString = "?sslmode=enable"
	}

	return fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v%v",
		postgresUser,
		postgresPassword,
		postgresHost,
		postgresPort,
		postgresDatabase,
		postgresSSLModeString,
	), nil
}

func GetSchema() string {
	postgresSchema := strings.TrimSpace(os.Getenv("POSTGRES_SCHEMA"))
	if postgresSchema == "" {
		postgresSchema = "public"
	}

	return postgresSchema
}

func GetDB(ctx context.Context, dsn string, maxIdleConns int, maxOpenConns int, connMaxIdleTime time.Duration, connMaxLifetime time.Duration) (*sqlx.DB, error) {
	db, err := sqlx.ConnectContext(ctx, "postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxIdleTime(connMaxIdleTime)
	db.SetConnMaxLifetime(connMaxLifetime)

	return db, nil
}

func GetDBFromEnvironment(ctx context.Context) (*sqlx.DB, error) {
	dsn, err := GetDSN()
	if err != nil {
		return nil, err
	}

	maxIdleConns := 2
	maxOpenConns := 100
	connMaxIdleTime := time.Second * 300
	connMaxLifetime := time.Second * 86400

	conn, err := GetDB(ctx, dsn, maxIdleConns, maxOpenConns, connMaxIdleTime, connMaxLifetime)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
