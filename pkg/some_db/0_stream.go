package some_db

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/introspect"
	"github.com/jackc/pglogrepl"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgproto3"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	standbyMessageTimeout = time.Second * 10
)

type Action string

const (
	INSERT   Action = "INSERT"
	UPDATE   Action = "UPDATE"
	DELETE   Action = "DELETE"
	TRUNCATE Action = "TRUNCATE"
)

type Change struct {
	Action          Action
	Table           *introspect.Table
	PrimaryKeyValue any
	Object          DjangolangObject
}

func (c *Change) String() string {
	primaryKeyColumn := "(unknown)"
	if c.Table.PrimaryKeyColumn != nil {
		primaryKeyColumn = c.Table.PrimaryKeyColumn.Name
	}

	b, _ := json.Marshal(c.PrimaryKeyValue)
	primaryKeyValue := string(b)

	object, _ := json.Marshal(c.Object)

	return fmt.Sprintf(
		"%v %v (%v: %v): %v",
		c.Action, c.Table.Name, primaryKeyColumn, primaryKeyValue, string(object),
	)
}

func runStream(outerCtx context.Context) error {
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

	_, err = db.ExecContext(ctx, fmt.Sprintf("DROP PUBLICATION IF EXISTS %v;", dbName))
	if err != nil {
		return fmt.Errorf("failed to ensure any pre-existing publication was removed: %v", err)
	}

	_, err = db.ExecContext(ctx, fmt.Sprintf("CREATE PUBLICATION %v FOR ALL TABLES;", dbName))
	if err != nil {
		return fmt.Errorf("failed to ensure any pre-existing publication was removed: %v", err)
	}

	defer func() {
		_, _ = db.ExecContext(ctx, fmt.Sprintf("DROP PUBLICATION IF EXISTS %v;", dbName))
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
		dbName,
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
			dbName,
			pglogrepl.DropReplicationSlotOptions{Wait: true},
		)
	}()

	err = pglogrepl.StartReplication(
		ctx,
		conn,
		dbName,
		identifySystemResult.XLogPos,
		pglogrepl.StartReplicationOptions{
			PluginArgs: []string{
				"proto_version '1'",
				fmt.Sprintf("publication_names '%v'", dbName),
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to start replication: %v", err)
	}

	clientXLogPos := identifySystemResult.XLogPos
	nextStandbyMessageDeadline := time.Now().Add(standbyMessageTimeout)
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

			nextStandbyMessageDeadline = time.Now().Add(standbyMessageTimeout)
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
				var action string
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

				_ = oldTupleColumns
				_ = newTupleColumns

				// should be a single item for everything except truncate
				for _, relationID := range relationIDs {
					relation, ok := relations[relationID]
					if !ok {
						logger.Printf("warning: failed to resolve relation ID %#+v to a relation", relationID)
						continue
					}

					if action == "TRUNCATE" {
						// TODO
						logger.Printf("warning: not implemented: %v %v.%v", action, relation.Namespace, relation.RelationName)
						continue
					}

					var relevantTupleColumns []*pglogrepl.TupleDataColumn

					if action == "INSERT" || action == "UPDATE" {
						relevantTupleColumns = newTupleColumns
					} else if action == "DELETE" {
						relevantTupleColumns = oldTupleColumns
					}

					if relevantTupleColumns == nil {
						continue
					}

					if len(relevantTupleColumns) != len(relation.Columns) {
						return fmt.Errorf(
							"assertion failed: relevant tuple columns %#+v and relation columns %#+v should match for %#+v",
							len(relevantTupleColumns),
							len(relation.Columns),
							message,
						)
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

					if table.PrimaryKeyColumn != nil {
						primaryKeyColumn := table.PrimaryKeyColumn

						rawPrimaryKeyValue := values[primaryKeyColumn.Name]

						var primaryKeyValue any

						var primaryKeyValueAsString string

						if primaryKeyColumn.DataType == "uuid" {
							rawValueAsBytes, ok := rawPrimaryKeyValue.([16]uint8)
							if !ok {
								logger.Printf("warning: failed to cast %v: %#+v to []uint8 for %v %v.%v; err: %v",
									primaryKeyColumn.Name, rawPrimaryKeyValue, action, relation.Namespace, relation.RelationName, err,
								)
								continue
							}

							rawValueAsUUID, err := uuid.FromBytes(rawValueAsBytes[:])
							if err != nil {
								logger.Printf("warning: failed to cast %v: %#+v to UUID for %v %v.%v; err: %v",
									primaryKeyColumn.Name, rawPrimaryKeyValue, action, relation.Namespace, relation.RelationName, err,
								)
								continue
							}

							primaryKeyValue = rawValueAsUUID

							primaryKeyValueAsString = fmt.Sprintf("'%v'", rawValueAsUUID.String())
						} else {
							rawValueAsBytes, err := json.Marshal(rawPrimaryKeyValue)
							if err != nil {
								logger.Printf("warning: failed to marshal %v: %#+v to json for %v %v.%v; err: %v",
									primaryKeyColumn.Name, rawPrimaryKeyValue, action, relation.Namespace, relation.RelationName, err,
								)
								continue
							}

							err = json.Unmarshal(rawValueAsBytes, &primaryKeyValue)
							if err != nil {
								logger.Printf("warning: failed to unmarshal %v: %#+v from json for %v %v.%v; err: %v",
									primaryKeyColumn.Name, rawValueAsBytes, action, relation.Namespace, relation.RelationName, err,
								)
								continue
							}

							primaryKeyValueAsString = string(rawValueAsBytes)
							if strings.HasPrefix(primaryKeyValueAsString, "\"") && strings.HasSuffix(primaryKeyValueAsString, "\"") {
								primaryKeyValueAsString = fmt.Sprintf("'%v'", primaryKeyValueAsString[1:len(primaryKeyValueAsString)-1])
							}
						}

						var object DjangolangObject

						if action != "DELETE" {
							objects, err := selectFunc(
								ctx,
								db,
								columnNames,
								nil,
								helpers.Ptr(1),
								nil,
								fmt.Sprintf(
									"%v.%v = %v",
									table.Name,
									primaryKeyColumn.Name,
									primaryKeyValueAsString,
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

						change := Change{
							Action:          Action(action),
							Table:           table,
							PrimaryKeyValue: primaryKeyValue,
							Object:          object,
						}

						logger.Printf("change: %v", change.String())
					}
				}
			}

			clientXLogPos = xld.WALStart + pglogrepl.LSN(len(xld.WALData))
		}
	}
}
