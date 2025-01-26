package config

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetDSN() string {
	postgresUser := PostgresUser()

	postgresPassword := PostgresPassword()

	postgresHost := PostgresHost()

	postgresPort := PostgresPort()

	postgresDatabase := PostgresDatabase()

	postgresSSLModeString := "?sslmode=disable"
	if PostgresSSLMode() {
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
	)
}

func GetSchema() string {
	return PostgresSchema()
}

func GetConn(ctx context.Context, dsn string) (*pgconn.PgConn, error) {
	conn, err := pgconn.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}

	// TODO: some versions of Postgres (or something) return SYNTAX ERROR for this... idk
	// err = conn.Ping(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	return conn, nil
}

func GetDB(ctx context.Context, dsn string, maxIdleConns int, maxOpenConns int, connMaxIdleTime time.Duration, connMaxLifetime time.Duration) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	config.AfterConnect = func(ctx context.Context, c *pgx.Conn) error {
		// TODO: figure out what to do with c.TypeMap().RegisterType() as required
		return nil
	}

	config.MinConns = int32(maxIdleConns)
	config.MaxConns = int32(maxOpenConns)
	config.MaxConnIdleTime = connMaxIdleTime
	config.MaxConnLifetime = connMaxLifetime
	config.HealthCheckPeriod = time.Second * 10

	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetDBFromEnvironment(ctx context.Context) (*pgxpool.Pool, error) {
	dsn := GetDSN()

	maxIdleConns := 5
	maxOpenConns := 500
	connMaxIdleTime := time.Second * 300
	connMaxLifetime := time.Second * 86400

	db, err := GetDB(ctx, dsn, maxIdleConns, maxOpenConns, connMaxIdleTime, connMaxLifetime)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetConnFromEnvironment(ctx context.Context) (*pgconn.PgConn, error) {
	dsn := GetDSN()

	conn, err := GetConn(ctx, fmt.Sprintf("%v&replication=database", dsn))
	if err != nil {
		return nil, err
	}

	return conn, nil
}
