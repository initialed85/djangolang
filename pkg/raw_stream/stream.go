package raw_stream

import (
	"context"
	"fmt"
	_log "log"
	"maps"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/initialed85/djangolang/pkg/config"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/introspect"
	"github.com/jackc/pglogrepl"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgproto3"
	"github.com/jackc/pgx/v5/pgtype"
)

var defaultNodeName = config.NodeName()

var log = helpers.GetLogger("raw_stream")

func ThisLogger() *_log.Logger {
	return log
}

const (
	timeout = time.Second * 10
)

func Run(outerCtx context.Context, ready chan struct{}, changes chan *Change, tableByName introspect.TableByName, handleChange func(*Change) error, nodeNames ...string) error {
	nodeName := defaultNodeName
	if len(nodeNames) > 0 {
		nodeName = nodeNames[0]
	}

	ctx, cancel := context.WithCancel(outerCtx)
	defer cancel()

	log.Printf("getting database from environment...")

	dbPool, err := config.GetDBFromEnvironment(ctx)
	if err != nil {
		return err
	}

	defer func() {
		log.Printf("destroying db...")
		defer log.Printf("destroyed db.")

		dbPool.Close()
	}()

	if tableByName == nil {
		err = func() error {
			schema := config.GetSchema()

			tx, err := dbPool.Begin(ctx)
			if err != nil {
				return err
			}

			defer func() {
				_ = tx.Rollback(ctx)
			}()

			tableByName, err = introspect.Introspect(ctx, tx, schema)
			if err != nil {
				return err
			}
			return nil
		}()
		if err != nil {
			return err
		}
	}

	adjustedNodeName := strings.ReplaceAll(nodeName, "-", "_")

	publicationName := fmt.Sprintf("%v_%v", "djangolang", adjustedNodeName)

	log.Printf("publicationName: %s", publicationName)

	log.Printf("checking WAL level...")
	err = CheckWALLevel(ctx, dbPool)
	if err != nil {
		return err
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
						currentReplicaIdentity, err := GetReplicaIdentity(ctx, dbPool, tableName)
						if err != nil {
							return fmt.Errorf("failed to check current replica identity for %v; %v", table.Name, err)
						}

						if currentReplicaIdentity == "f" {
							log.Printf("replica identity is already %s for table %s", setReplicaIdentity, table.Name)
							continue
						}
					}

					_, err = config.DoExec(ctx, dbPool, fmt.Sprintf("ALTER TABLE %v REPLICA IDENTITY %v", table.Name, setReplicaIdentity))
					if err != nil {
						return fmt.Errorf("failed to set replica identity to %v for %v; %v", setReplicaIdentity, table.Name, err)
					}

					log.Printf("set replica identity to %s for table %s", setReplicaIdentity, table.Name)
				}
			}

			// log.Printf("checking for conflicting publication / replication slot...")

			// err = CheckForConflictingPublication(ctx, dbPool, adjustedNodeName, publicationName)
			// if err != nil {
			// 	return err
			// }

			// err = CheckForConflictingReplicationSlot(ctx, dbPool, adjustedNodeName, publicationName)
			// if err != nil {
			// 	return err
			// }

			log.Printf("creating publication...")

			_, err = config.DoExec(ctx, dbPool, fmt.Sprintf("DROP PUBLICATION IF EXISTS %v;", publicationName))
			if err != nil {
				return fmt.Errorf("failed to ensure any pre-existing publication was removed: %v", err)
			}

			_, err = config.DoExec(ctx, dbPool, fmt.Sprintf("CREATE PUBLICATION %v FOR ALL TABLES;", publicationName))
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

		_, _ = config.DoExec(teardownCtx, dbPool, fmt.Sprintf("DROP PUBLICATION IF EXISTS %v;", publicationName))
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

	_, _ = config.DoExec(ctx, dbPool, fmt.Sprintf("SELECT pg_drop_replication_slot('%s')", publicationName))

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

	select {
	case ready <- struct{}{}:
	default:
	}

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

					bothTupleColumns := [][]*pglogrepl.TupleDataColumn{
						oldTupleColumns,
						newTupleColumns,
					}

					bothItems := []map[string]any{
						nil,
						nil,
					}

					hasDeletedAtColumn := false

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

							if column.Name == "deleted_at" {
								hasDeletedAtColumn = true
							}

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

					change := &Change{
						Timestamp: time.Now().UTC(),
						ID:        uuid.Must(uuid.NewRandom()),
						Action:    Action(action),
						TableName: relation.RelationName,
						Item:      relevantItem,
						Xid:       lastXid,
					}

					if handleChange != nil {
						err = handleChange(change)
						if err != nil {
							return err
						}
					}

					if changes != nil {
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
			}

			clientXLogPos = xld.WALStart + pglogrepl.LSN(len(xld.WALData))
		}
	}

	return fmt.Errorf("context canceled or something like that")
}
