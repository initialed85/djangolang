package stream

import (
	"context"
	"fmt"
	_log "log"
	"maps"
	"slices"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/initialed85/djangolang/pkg/config"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/introspect"
	"github.com/jackc/pglogrepl"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgproto3"
	"github.com/jackc/pgx/v5/pgtype"
)

var nodeName = config.NodeName()

var log = helpers.GetLogger("stream")

func ThisLogger() *_log.Logger {
	return log
}

const (
	timeout = time.Second * 10
)

func Run(outerCtx context.Context, changes chan Change, tableByName introspect.TableByName) error {
	ctx, cancel := context.WithCancel(outerCtx)
	defer cancel()

	log.Printf("getting database from environment...")

	db, err := config.GetDBFromEnvironment(ctx)
	if err != nil {
		return err
	}

	// TODO: this hangs, probably something to do w/ a connection not being released back to the pool despite
	// the fact we're using the higher level API (i.e. implicit acquire / release)- who cares, we're shutting down
	// anyway
	// defer func() {
	// 	log.Printf("destroying db...")
	// 	defer log.Printf("destroyed db.")

	// 	db.Close()
	// }()

	log.Printf("getting redis from environment...")

	redisPool, err := config.GetRedisFromEnvironment()
	if err != nil {
		return err
	}
	defer func() {
		log.Printf("destroying redis...")
		defer log.Printf("destroyed redis.")

		_ = redisPool.Close()
	}()

	adjustedNodeName := strings.ReplaceAll(nodeName, "-", "_")

	publicationName := fmt.Sprintf("%v_%v", "djangolang", adjustedNodeName)

	log.Printf("publicationName: %s", publicationName)

	log.Printf("checking WAL level...")

	row := db.QueryRow(ctx, "SHOW wal_level;")

	var walLevel string
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

	for {
		err = func() error {
			tableNames := slices.Collect(maps.Keys(tableByName))

			setReplicaIdentity := config.SetReplicaIdentity()
			if setReplicaIdentity != "" {
				for _, tableName := range tableNames {
					table := tableByName[tableName]

					log.Printf("setting replica identity of table %s to %s", table.Name, setReplicaIdentity)

					if setReplicaIdentity == "full" {
						rows, err := db.Query(ctx, fmt.Sprintf("SELECT relreplident::text FROM pg_class WHERE oid = '%s'::regclass;", table.Name))
						if err != nil {
							return fmt.Errorf("failed to check current replica identity for %v; %v", table.Name, err)
						}

						if !rows.Next() {
							return fmt.Errorf("failed to check current replica identity for %v; %v", table.Name, fmt.Errorf("no rows returned"))
						}

						var currentReplicaIdentity string

						err = rows.Scan(&currentReplicaIdentity)
						if err != nil {
							return fmt.Errorf("failed to check current replica identity for %v; %v", table.Name, err)
						}

						if currentReplicaIdentity == "f" {
							log.Printf("replica identity is already %s for table %s", setReplicaIdentity, table.Name)
							continue
						}

						err = rows.Err()
						if err != nil {
							return fmt.Errorf("failed to check current replica identity for %v; %v", table.Name, err)
						}
					}

					for i := 0; i < 10; i++ {
						_, err = db.Exec(ctx, fmt.Sprintf("ALTER TABLE %v REPLICA IDENTITY %v", table.Name, setReplicaIdentity))
						if err != nil {
							time.Sleep(time.Millisecond * 100)
							continue
						}

						break
					}

					if err != nil {
						return fmt.Errorf("failed to set replica identity to %v for %v; %v", setReplicaIdentity, table.Name, err)
					}

					log.Printf("set replica identity to %s for table %s", setReplicaIdentity, table.Name)
				}
			}

			log.Printf("checking for conflicting publication / replication slot...")

			row = db.QueryRow(ctx, "SELECT count(*) FROM pg_publication WHERE pubname = $1;", publicationName)

			var count int
			err = row.Scan(&count)
			if err != nil {
				return fmt.Errorf("failed to check for conflicting publication: %v", err)
			}

			if count > 0 {
				row = db.QueryRow(ctx, "SELECT count(*) FROM pg_replication_slots WHERE slot_name = $1;", publicationName)

				var count int
				err = row.Scan(&count)
				if err != nil {
					return fmt.Errorf("failed to check for conflicting replication slot: %v", err)
				}

				if count > 0 {
					return fmt.Errorf(
						"cannot continue with DJANGOLANG_NODE_NAME=%v (publication name / replication slot name %v); there is already an active replication slot for that node name",
						adjustedNodeName,
						publicationName,
					)
				}
			}

			log.Printf("creating publication...")

			_, err = db.Exec(ctx, fmt.Sprintf("DROP PUBLICATION IF EXISTS %v;", publicationName))
			if err != nil {
				return fmt.Errorf("failed to ensure any pre-existing publication was removed: %v", err)
			}

			_, err = db.Exec(ctx, fmt.Sprintf("CREATE PUBLICATION %v FOR ALL TABLES;", publicationName))
			if err != nil {
				return fmt.Errorf("failed to ensure any pre-existing publication was removed: %v", err)
			}

			return nil
		}()
		if err != nil {
			log.Printf("warning: %s; retrying...", err.Error())
			time.Sleep(time.Millisecond * 100)
			continue
		}

		break
	}

	defer func() {
		teardownCtx, teardownCancel := context.WithTimeout(context.Background(), time.Second*10)
		defer teardownCancel()

		log.Printf("destroying publication...")
		defer log.Printf("destroyed publication.")

		_, _ = db.Exec(teardownCtx, fmt.Sprintf("DROP PUBLICATION IF EXISTS %v;", publicationName))
	}()

	log.Printf("getting connection from environment...")

	conn, err := config.GetConnFromEnvironment(ctx)
	if err != nil {
		return fmt.Errorf("failed to get database connection from pool: %v", err)
	}

	identifySystemResult, err := pglogrepl.IdentifySystem(ctx, conn)
	if err != nil {
		return fmt.Errorf("failed to identify system: %v", err)
	}

	log.Printf("creating replication slot...")

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
		teardownCtx, teardownCancel := context.WithTimeout(context.Background(), time.Second*10)
		defer teardownCancel()

		log.Printf("destroying replication slot...")
		defer log.Printf("destroyed replication slot.")

		_ = pglogrepl.DropReplicationSlot(
			teardownCtx,
			conn,
			publicationName,
			pglogrepl.DropReplicationSlotOptions{Wait: true},
		)
	}()

	log.Printf("starting replication...")

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

	log.Printf("handling replication messages...")

	clientXLogPos := identifySystemResult.XLogPos
	nextStandbyMessageDeadline := time.Now().Add(timeout)
	relations := map[uint32]*pglogrepl.RelationMessage{}
	typeMap := pgtype.NewMap()

	var lastXid uint32

loop:
	for {
		select {
		case <-ctx.Done():
			log.Printf("stream ctx done at top of loop")
			break loop
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
						log.Printf("warning: failed to resolve relation ID %#+v to a relation", relationID)
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
										return fmt.Errorf("failed to decode %#+v; %v", column, err)
									}
								}

								item[column.Name] = value

							case 'u': // unchanged TOAST (The Oversized-Attribute Storage Technique)
								// TODO: figure out how to handle this- do we need to go fetch it?
								log.Printf("warning: (not implemented) failed to handle unchanged toast: %#+v / %#+v", column, tupleColumn)

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

					err = func() error {
						if redisPool == nil {
							return nil
						}

						redisConn := redisPool.Get()
						defer func() {
							redisConn.Close()
						}()

						tableNames := make(map[string]struct{})
						tableNames[table.Name] = struct{}{}

						for _, column := range table.Columns {
							if column.ForeignColumn == nil {
								continue
							}

							tableNames[column.ForeignColumn.TableName] = struct{}{}
						}

						// if the change was for Camera, it should propagate to Video and Detection (using referenced-by)
						for _, referencedByColumn := range table.ReferencedByColumns {
							tableNames[referencedByColumn.TableName] = struct{}{}
						}

						keysToDelete := make([]string, 0)

						for _, tableName := range slices.Collect(maps.Keys(tableNames)) {
							scanResponses, err := redis.Scan(redis.Values(redisConn.Do("SCAN", 0, "MATCH", fmt.Sprintf("%v:*", tableName))))
							if err != nil {
								return fmt.Errorf("failed redis scan for %v:*; %v", tableName, err)
							}

							for _, scanResponse := range scanResponses {
								keys, err := redis.Strings(scanResponse, nil)
								if err != nil {
									return fmt.Errorf("failed redis scan for %v:*; %v", tableName, err)
								}

								keysToDelete = append(keysToDelete, keys...)
							}
						}

						for _, key := range keysToDelete {
							_, err = redisConn.Do("DEL", key)
							if err != nil {
								return fmt.Errorf("failed redis delete for %v; %v", key, err)
							}
						}

						return nil
					}()
					if err != nil {
						return err
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
					case <-ctx.Done():
						log.Printf("ctx context done while trying to push change")
						break loop
					case <-time.After(timeout):
						log.Printf(
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

	return fmt.Errorf("context canceled or something like that")
}
