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
	err = conn.Ping(ctx)
	if err != nil {
		return nil, err
	}

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

	// TODO: some versions of Postgres (or something) return SYNTAX ERROR for this... idk
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
	connMaxIdleTime := time.Second * 30
	connMaxLifetime := time.Second * 600

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

func DoQuery(ctx context.Context, db *pgxpool.Pool, query string, params ...any) ([][]any, error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	rows, err := tx.Query(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	defer func() {
		rows.Close()
	}()

	allValues := make([][]any, 0)

	for rows.Next() {
		err = rows.Err()
		if err != nil {
			return nil, err
		}

		values, err := rows.Values()
		if err != nil {
			return nil, err
		}

		allValues = append(allValues, values)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return allValues, nil
}

func DoExec(ctx context.Context, db *pgxpool.Pool, query string, params ...any) (int, error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		return 0, err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	result, err := tx.Exec(ctx, query, params...)
	if err != nil {
		return 0, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, err
	}

	return int(result.RowsAffected()), nil
}
