package some_db

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/types"
	"github.com/jackc/pglogrepl"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgproto3"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	timeout = time.Second * 10
)

func runStream(outerCtx context.Context, changes chan types.Change) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := helpers.GetDBFromEnvironment(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = db.Close()
	}()

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

	nodeName := strings.TrimSpace(os.Getenv("DJANGOLANG_NODE_NAME"))
	if nodeName == "" {
		nodeName = "default"
	}

	publicationName := fmt.Sprintf("%v_%v", dbName, nodeName)

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

	conn, err := helpers.GetConnFromEnvironment(ctx)
	if err != nil {
		return err
	}

	identifySystemResult, err := pglogrepl.IdentifySystem(ctx, conn)
	if err != nil {
		return fmt.Errorf("failed to identify system: %v", err)
	}

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

	_ = typeMap

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
				var action types.Action
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

				case *pglogrepl.TypeMessage, *pglogrepl.OriginMessage, *pglogrepl.BeginMessage, *pglogrepl.CommitMessage:
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

					selectFunc, ok := selectFuncByTableName[relation.RelationName]
					if !ok {
						return fmt.Errorf(
							"assertion failed: relation name %#+v could not be resolved to a select function for %#+v",
							relation.RelationName,
							message,
						)
					}

					columnNames, ok := columnNamesByTableName[relation.RelationName]
					if !ok {
						return fmt.Errorf(
							"assertion failed: relation name %#+v could not be resolved to a set of column names for %#+v",
							relation.RelationName,
							message,
						)
					}

					var relevantTupleColumns []*pglogrepl.TupleDataColumn

					if action == types.INSERT || action == types.UPDATE {
						relevantTupleColumns = newTupleColumns
					} else if action == types.DELETE {
						relevantTupleColumns = oldTupleColumns
					}

					if action != types.TRUNCATE {
						if len(relevantTupleColumns) != len(relation.Columns) {
							return fmt.Errorf(
								"assertion failed: relevant tuple columns %#+v and relation columns %#+v should match for %#+v",
								len(relevantTupleColumns),
								len(relation.Columns),
								message,
							)
						}
					}

					values := make(map[string]any)

					for i, tupleColumn := range relevantTupleColumns {
						column := relation.Columns[i]

						switch tupleColumn.DataType {

						case 'n': // null
							values[column.Name] = nil

						case 't': // text
							var value any

							// TODO: PostGIS types

							dt, ok := typeMap.TypeForOID(column.DataType)
							if !ok {
								value = string(tupleColumn.Data)
							} else {
								value, err = dt.Codec.DecodeValue(typeMap, column.DataType, pgtype.TextFormatCode, tupleColumn.Data)
								if err != nil {
									return fmt.Errorf("failed to decode %#+v: %v", column, err)
								}
							}

							values[column.Name] = value

						case 'u': // unchanged TOAST (we'd have to go fetch this for ourselves if we needed it)
							return fmt.Errorf("not implemented: failed to handle unchanged toast: %#+v / %#+v", column, tupleColumn)

						default:
							return fmt.Errorf("failed to handle data type %v for %#+v / %#+v", column.DataType, column, tupleColumn)
						}
					}

					if actualDebug {
						logger.Printf("%v %v.%v %#+v", action, relation.Namespace, relation.RelationName, values)
					}

					if table.PrimaryKeyColumn == nil {
						logger.Printf("warning: table has no primary key for %v %v.%v; err: %v",
							action, relation.Namespace, relation.RelationName, err,
						)
						continue
					}

					primaryKeyColumn := table.PrimaryKeyColumn

					var primaryKeyValue any

					primaryKeyValue = values[primaryKeyColumn.Name]

					if primaryKeyColumn.DataType == "uuid" {
						rawValueAsBytes, ok := primaryKeyValue.([16]uint8)
						if !ok {
							logger.Printf("warning: failed to cast %v: %#+v to []uint8 for %v %v.%v; err: %v",
								primaryKeyColumn.Name, primaryKeyValue, action, relation.Namespace, relation.RelationName, err,
							)
							continue
						}

						rawValueAsUUID, err := uuid.FromBytes(rawValueAsBytes[:])
						if err != nil {
							logger.Printf("warning: failed to cast %v: %#+v to UUID for %v %v.%v; err: %v",
								primaryKeyColumn.Name, primaryKeyValue, action, relation.Namespace, relation.RelationName, err,
							)
							continue
						}

						primaryKeyValue = rawValueAsUUID
					}

					var object types.DjangolangObject

					if action == types.INSERT || action == types.UPDATE {
						objects, err := selectFunc(
							ctx,
							db,
							columnNames,
							nil,
							helpers.Ptr(1),
							nil,
							types.Clause(
								fmt.Sprintf(
									"%v.%v = $1",
									table.Name,
									primaryKeyColumn.Name),
								primaryKeyValue,
							),
						)
						if err != nil {
							logger.Printf("warning: failed to select object for %v %v.%v; err: %v",
								action, relation.Namespace, relation.RelationName, err,
							)
							continue
						}

						if len(objects) != 1 {
							logger.Printf("warning: failed to select object for %v %v.%v; err: %v",
								action, relation.Namespace, relation.RelationName,
								fmt.Errorf("wanted exactly 1 object, got %v objects", len(objects)),
							)
							continue
						}

						object = objects[0]
					}

					change := types.Change{
						ID:               uuid.Must(uuid.NewRandom()),
						Action:           types.Action(action),
						TableName:        table.Name,
						PrimaryKeyColumn: table.PrimaryKeyColumn.Name,
						PrimaryKeyValue:  primaryKeyValue,
						Object:           object,
					}

					if actualDebug {
						logger.Printf("change: %v", change.String())
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
