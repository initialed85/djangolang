package helpers

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func GetDBConn(ctx context.Context, dsn string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func GetDBConnFromEnvironment(ctx context.Context) (*pgx.Conn, error) {
	dsn, err := GetDSN()
	if err != nil {
		return nil, err
	}

	conn, err := GetDBConn(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
