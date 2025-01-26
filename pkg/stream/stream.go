package stream

import (
	"context"
	"fmt"
	_log "log"
	"maps"
	"slices"

	"github.com/gomodule/redigo/redis"
	"github.com/initialed85/djangolang/pkg/config"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/introspect"
	"github.com/initialed85/djangolang/pkg/raw_stream"
)

var log = helpers.GetLogger("stream")

func ThisLogger() *_log.Logger {
	return log
}

func Run(outerCtx context.Context, ready chan struct{}, changes chan *Change, tableByName introspect.TableByName, nodeNames ...string) error {
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

	handleChange := func(change *Change) error {
		if redisPool == nil {
			return nil
		}

		redisConn := redisPool.Get()
		defer func() {
			redisConn.Close()
		}()

		table := tableByName[change.TableName]
		if table == nil {
			return fmt.Errorf("table %s unknown in change %#+v", change.TableName, change)
		}

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
	}

	return raw_stream.Run(outerCtx, ready, changes, tableByName, handleChange, nodeNames...)
}
