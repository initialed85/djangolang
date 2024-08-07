package helpers

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func GetDSN() (string, error) {
	postgresUser := GetEnvironmentVariableOrDefault("POSTGRES_USER", "postgres")

	postgresPassword := GetEnvironmentVariable("POSTGRES_PASSWORD")
	if postgresPassword == "" {
		return "", fmt.Errorf("POSTGRES_PASSWORD env var empty or unset")
	}

	postgresHost := GetEnvironmentVariableOrDefault("POSTGRES_HOST", "localhost")

	postgresPort := GetEnvironmentVariableOrDefault("POSTGRES_PORT", "5432")

	postgresDatabase := GetEnvironmentVariable("POSTGRES_DB")
	if postgresDatabase == "" {
		return "", fmt.Errorf("POSTGRES_DB env var empty or unset")
	}

	postgresSSLModeString := "?sslmode=disable"
	if GetEnvironmentVariableOrDefault("POSTGRES_SSLMODE", "0") == "1" {
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
		log.Printf("POSTGRES_SCHEMA empty or unset; defaulted to %v", postgresSchema)
	}

	return postgresSchema
}

func GetConn(ctx context.Context, dsn string) (*pgconn.PgConn, error) {
	conn, err := pgconn.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return conn, nil
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

func GetConnFromEnvironment(ctx context.Context) (*pgconn.PgConn, error) {
	dsn, err := GetDSN()
	if err != nil {
		return nil, err
	}

	conn, err := GetConn(ctx, fmt.Sprintf("%v&replication=database", dsn))
	if err != nil {
		return nil, err
	}

	return conn, nil
}
