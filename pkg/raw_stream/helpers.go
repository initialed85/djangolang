package raw_stream

import (
	"context"
	"fmt"

	"github.com/initialed85/djangolang/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CheckWALLevel(ctx context.Context, dbPool *pgxpool.Pool) error {
	allValues, err := config.DoQuery(ctx, dbPool, "SHOW wal_level;")
	if err != nil {
		return fmt.Errorf("check wal level failed: %s", err)
	}

	if len(allValues) != 1 {
		return fmt.Errorf("check wal level failed to get exactly 1 row; %#+v", allValues)
	}

	walLevel := allValues[0][0].(string)

	if walLevel != "logical" {
		return fmt.Errorf(
			"wal_level must be %#+v; got %#+v; hint: run \"ALTER SYSTEM SET wal_level = 'logical';\" and then restart Postgres",
			"logical",
			walLevel,
		)
	}

	return nil
}

func GetReplicaIdentity(ctx context.Context, dbPool *pgxpool.Pool, tableName string) (string, error) {
	query := fmt.Sprintf("SELECT relreplident::text FROM pg_class WHERE oid = '%s'::regclass;", tableName)

	allValues, err := config.DoQuery(ctx, dbPool, query)
	if err != nil {
		return "", fmt.Errorf("get replica identify failed: %s", err)
	}

	if len(allValues) != 1 {
		return "", fmt.Errorf("get replica identity failed to get exactly 1 row; %#+v", allValues)
	}

	replicaIdentity := allValues[0][0].(string)

	return replicaIdentity, nil
}

func CheckForConflictingPublication(ctx context.Context, dbPool *pgxpool.Pool, adjustedNodeName string, publicationName string) error {
	allValues, err := config.DoQuery(ctx, dbPool, "SELECT count(*) FROM pg_publication WHERE pubname = $1;", publicationName)
	if err != nil {
		return err
	}

	if len(allValues) != 1 {
		return fmt.Errorf("check for conflicting replication slot failed to get exactly 1 row; %#+v", allValues)
	}

	count := allValues[0][0].(int64)

	if count > 0 {
		return fmt.Errorf(
			"cannot continue with DJANGOLANG_NODE_NAME=%v (publication name / replication slot name %v); there is already %d publication(s) for that node name",
			adjustedNodeName,
			publicationName,
			count,
		)
	}

	return nil
}

func CheckForConflictingReplicationSlot(ctx context.Context, dbPool *pgxpool.Pool, adjustedNodeName string, publicationName string) error {
	allValues, err := config.DoQuery(ctx, dbPool, "SELECT count(*) FROM pg_replication_slots WHERE slot_name = $1;", publicationName)
	if err != nil {
		return err
	}

	if len(allValues) != 1 {
		return fmt.Errorf("check for conflicting replication slot failed to get exactly 1 row; %#+v", allValues)
	}

	count := allValues[0][0].(int64)

	if count > 0 {
		return fmt.Errorf(
			"cannot continue with DJANGOLANG_NODE_NAME=%v (publication name / replication slot name %v); there is already %d active replication slot(s) for that node name",
			adjustedNodeName,
			publicationName,
			count,
		)
	}

	return nil
}
