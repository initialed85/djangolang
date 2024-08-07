package stream

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/introspect"
	"github.com/jackc/pglogrepl"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgproto3"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/exp/maps"
)

var logger = helpers.GetLogger("djangolang/stream")

const (
	timeout = time.Second * 10
)

func Run(outerCtx context.Context, changes chan Change, tableByName map[string]*introspect.Table) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Printf("getting database from environment...")

	db, err := helpers.GetDBFromEnvironment(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = db.Close()
	}()

	logger.Printf("getting redis from environment...")

	redisURL := helpers.GetRedisURL()
	var redisConn redis.Conn
	if redisURL != "" {
		redisConn, err = redis.DialURLContext(ctx, redisURL)
		if err != nil {
			log.Fatalf("err: %v", err)
		}
		defer func() {
			_ = redisConn.Close()
		}()
	} else {
		logger.Printf("warning: redis disabled (no REDIS_URL set)")
	}

	logger.Printf("checking WAL level...")

	var walLevel string
	row := db.QueryRowContext(ctx, "SHOW wal_level;")

	err = row.Err()
	if err != nil {
		return fmt.Errorf("failed to check current wal level: %v", err)
	}

	err = row.Scan(&walLevel)
	if err != nil {
		return fmt.Errorf("failed to check current wal level: %v", err)
	}

	if walLevel != "logical" {
		return fmt.Errorf(
			"wal_level must be %#+v; got %#+v; hint: run \"ALTER SYSTEM SET wal_level = 'logical';\" and then restart Postgres",
			"logical",
			walLevel,
		)
	}

	setReplicaIdentity := helpers.GetEnvironmentVariableOrDefault("DJANGOLANG_SET_REPLICA_IDENTITY", "")
	if setReplicaIdentity != "" {
		logger.Printf("warning: DJANGOLANG_SET_REPLICA_IDENTITY=%v; changing replica identity on all tables (this persists at shutdown)...", setReplicaIdentity)
		for _, table := range tableByName {
			_, err = db.ExecContext(ctx, fmt.Sprintf("ALTER TABLE %v REPLICA IDENTITY %v", table.Name, setReplicaIdentity))
			if err != nil {
				return fmt.Errorf("failed to set replica identity to %v for %v: %v", setReplicaIdentity, table.Name, err)
			}
		}
	}

	nodeName := helpers.GetEnvironmentVariableOrDefault("DJANGOLANG_NODE_NAME", "default")

	logger.Printf("creating publication...")

	publicationName := fmt.Sprintf("%v_%v", "djangolang", nodeName)

	_, err = db.ExecContext(ctx, fmt.Sprintf("DROP PUBLICATION IF EXISTS %v;", publicationName))
	if err != nil {
		return fmt.Errorf("failed to ensure any pre-existing publication was removed: %v", err)
	}

	_, err = db.ExecContext(ctx, fmt.Sprintf("CREATE PUBLICATION %v FOR ALL TABLES;", publicationName))
	if err != nil {
		return fmt.Errorf("failed to ensure any pre-existing publication was removed: %v", err)
	}

	defer func() {
		_, _ = db.ExecContext(ctx, fmt.Sprintf("DROP PUBLICATION IF EXISTS %v;", publicationName))
	}()

	logger.Printf("getting connection from environment...")

	conn, err := helpers.GetConnFromEnvironment(ctx)
	if err != nil {
		return err
	}

	identifySystemResult, err := pglogrepl.IdentifySystem(ctx, conn)
	if err != nil {
		return fmt.Errorf("failed to identify system: %v", err)
	}

	logger.Printf("creating replication...")

	_, err = pglogrepl.CreateReplicationSlot(
		ctx,
		conn,
		publicationName,
		"pgoutput",
		pglogrepl.CreateReplicationSlotOptions{
			Temporary:      true,
			SnapshotAction: "",
			Mode:           pglogrepl.LogicalReplication,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create replication slot: %v", err)
	}

	defer func() {
		teardownCtx, teardownCancel := context.WithTimeout(context.Background(), time.Second*30)
		defer teardownCancel()

		_ = pglogrepl.DropReplicationSlot(
			teardownCtx,
			conn,
			publicationName,
			pglogrepl.DropReplicationSlotOptions{Wait: true},
		)
	}()

	logger.Printf("starting replication...")

	err = pglogrepl.StartReplication(
		ctx,
		conn,
		publicationName,
		identifySystemResult.XLogPos,
		pglogrepl.StartReplicationOptions{
			PluginArgs: []string{
				"proto_version '1'",
				fmt.Sprintf("publication_names '%v'", publicationName),
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to start replication: %v", err)
	}

	clientXLogPos := identifySystemResult.XLogPos
	nextStandbyMessageDeadline := time.Now().Add(timeout)
	relations := map[uint32]*pglogrepl.RelationMessage{}
	typeMap := pgtype.NewMap()

	var lastXid uint32

	for {
		select {
		case <-outerCtx.Done():
			return nil
		default:
		}

		if time.Now().After(nextStandbyMessageDeadline) {
			err = pglogrepl.SendStandbyStatusUpdate(
				ctx,
				conn,
				pglogrepl.StandbyStatusUpdate{WALWritePosition: clientXLogPos},
			)
			if err != nil {
				return err
			}

			nextStandbyMessageDeadline = time.Now().Add(timeout)
		}

		messageCtx, messageCancel := context.WithDeadline(ctx, nextStandbyMessageDeadline)
		backendMessage, err := conn.ReceiveMessage(messageCtx)
		messageCancel()

		if err != nil {
			if pgconn.Timeout(err) {
				continue
			}

			return fmt.Errorf("failed to receive backend message: %v", err)
		}

		errorResponse, ok := backendMessage.(*pgproto3.ErrorResponse)
		if ok {
			return fmt.Errorf("unexpectedly received %#+v", errorResponse)
		}

		copyData, ok := backendMessage.(*pgproto3.CopyData)
		if !ok {
			continue
		}

		var message pglogrepl.Message

		switch copyData.Data[0] {

		case pglogrepl.PrimaryKeepaliveMessageByteID:
			pkm, err := pglogrepl.ParsePrimaryKeepaliveMessage(copyData.Data[1:])
			if err != nil {
				return fmt.Errorf("failed to parse keepalive: %v", err)
			}

			if pkm.ReplyRequested {
				nextStandbyMessageDeadline = time.Time{}
			}

		case pglogrepl.XLogDataByteID:
			xld, err := pglogrepl.ParseXLogData(copyData.Data[1:])
			if err != nil {
				return fmt.Errorf("failed to parse xlog data: %v", err)
			}

			message, err = pglogrepl.Parse(xld.WALData)
			if err != nil {
				return fmt.Errorf("failed to parse logical replication message: %v", err)
			}

			switch message := message.(type) {

			case nil:
				continue

			case *pglogrepl.RelationMessage:
				relations[message.RelationID] = message

			default:
				var action Action
				var relationIDs []uint32
				var oldTupleColumns []*pglogrepl.TupleDataColumn
				var newTupleColumns []*pglogrepl.TupleDataColumn

				switch message := message.(type) {

				case *pglogrepl.InsertMessage:
					action = "INSERT"
					relationIDs = []uint32{message.RelationID}

					if message.Tuple != nil {
						newTupleColumns = message.Tuple.Columns
					}

				case *pglogrepl.UpdateMessage:
					action = "UPDATE"
					relationIDs = []uint32{message.RelationID}

					if message.OldTuple != nil {
						oldTupleColumns = message.OldTuple.Columns
					}

					if message.NewTuple != nil {
						newTupleColumns = message.NewTuple.Columns
					}

				case *pglogrepl.DeleteMessage:
					action = "DELETE"
					relationIDs = []uint32{message.RelationID}

					if message.OldTuple != nil {
						oldTupleColumns = message.OldTuple.Columns
					}

				case *pglogrepl.TruncateMessage:
					action = "TRUNCATE"
					relationIDs = message.RelationIDs

				case *pglogrepl.BeginMessage:
					lastXid = message.Xid

				case *pglogrepl.CommitMessage:
					lastXid = 0

				case *pglogrepl.TypeMessage, *pglogrepl.OriginMessage:
					// ignored

				default:
					return fmt.Errorf("unexpected message: %v", message)
				}

				if action == "" {
					continue
				}

				// should be a single item for everything except truncate
				for _, relationID := range relationIDs {
					relation, ok := relations[relationID]
					if !ok {
						logger.Printf("warning: failed to resolve relation ID %#+v to a relation", relationID)
						continue
					}

					table, ok := tableByName[relation.RelationName]
					if !ok {
						return fmt.Errorf(
							"assertion failed: relation name %#+v could not be resolved to an introspected table for %#+v",
							relation.RelationName,
							message,
						)
					}

					bothTupleColumns := [][]*pglogrepl.TupleDataColumn{
						oldTupleColumns,
						newTupleColumns,
					}

					bothItems := []map[string]any{
						nil,
						nil,
					}

					for i, thisTupleColumns := range bothTupleColumns {
						if thisTupleColumns == nil {
							continue
						}

						if action != TRUNCATE {
							if len(thisTupleColumns) != len(relation.Columns) {
								return fmt.Errorf(
									"assertion failed: tuple columns %#+v and relation columns %#+v should match for %#+v",
									len(thisTupleColumns),
									len(relation.Columns),
									message,
								)
							}
						}

						item := make(map[string]any)

						for j, tupleColumn := range thisTupleColumns {
							column := relation.Columns[j]

							switch tupleColumn.DataType {

							case 'n': // null
								item[column.Name] = nil

							case 't': // text
								var value any

								dt, ok := typeMap.TypeForOID(column.DataType)
								if !ok {
									value = string(tupleColumn.Data)
								} else {
									value, err = dt.Codec.DecodeValue(typeMap, column.DataType, pgtype.TextFormatCode, tupleColumn.Data)
									if err != nil {
										return fmt.Errorf("failed to decode %#+v: %v", column, err)
									}
								}

								item[column.Name] = value

							case 'u': // unchanged TOAST (The Oversized-Attribute Storage Technique)
								// TODO: figure out how to handle this- do we need to go fetch it?
								logger.Printf("warning: (not implemented) failed to handle unchanged toast: %#+v / %#+v", column, tupleColumn)

							default:
								return fmt.Errorf("failed to handle data type %v for %#+v / %#+v", column.DataType, column, tupleColumn)
							}
						}

						bothItems[i] = item
					}

					oldItem := bothItems[0]
					newItem := bothItems[1]

					relevantItem := newItem
					if action == DELETE {
						relevantItem = oldItem
					}

					// TODO: have a more configurable approach to soft-deletions
					if action == UPDATE {
						_, hasDeletedAtColumn := table.ColumnByName["deleted_at"]
						if hasDeletedAtColumn {
							if (oldItem == nil || oldItem["deleted_at"] == nil) && newItem["deleted_at"] != nil {
								action = SOFT_DELETE
							} else if (oldItem != nil && oldItem["deleted_at"] != nil) && newItem["deleted_at"] == nil {
								action = SOFT_RESTORE
							} else if (oldItem != nil && oldItem["deleted_at"] != nil) && newItem["deleted_at"] != nil {
								action = SOFT_UPDATE
							}
						}
					}

					if redisConn != nil {
						tableNames := make(map[string]struct{})
						tableNames[table.Name] = struct{}{}

						for _, otherTable := range table.ReferencedByTables {
							tableNames[otherTable.Name] = struct{}{}
						}

						keysToDelete := make([]string, 0)

						for _, tableName := range maps.Keys(tableNames) {
							scanResponses, err := redis.Scan(redis.Values(redisConn.Do("SCAN", 0, "MATCH", fmt.Sprintf("%v:*", tableName))))
							if err != nil {
								return fmt.Errorf("failed redis scan for %v:*: %v", tableName, err)
							}

							for _, scanResponse := range scanResponses {
								keys, err := redis.Strings(scanResponse, nil)
								if err != nil {
									return fmt.Errorf("failed redis scan for %v:*: %v", tableName, err)
								}

								keysToDelete = append(keysToDelete, keys...)
							}
						}

						for _, key := range keysToDelete {
							_, err = redisConn.Do("DEL", key)
							if err != nil {
								return fmt.Errorf("failed redis delete for %v: %v", key, err)
							}
						}
					}

					change := Change{
						Timestamp: time.Now().UTC(),
						ID:        uuid.Must(uuid.NewRandom()),
						Action:    Action(action),
						TableName: table.Name,
						Item:      relevantItem,
						Xid:       lastXid,
					}

					select {
					case <-time.After(timeout):
						logger.Printf(
							"warning: timed out after %v trying to send change %v; this change will now be thrown away",
							timeout, change,
						)
					case changes <- change:
					}
				}
			}

			clientXLogPos = xld.WALStart + pglogrepl.LSN(len(xld.WALData))
		}
	}
}
