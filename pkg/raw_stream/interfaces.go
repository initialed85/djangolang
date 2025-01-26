package raw_stream

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Action string

const (
	INSERT       Action = "INSERT"
	UPDATE       Action = "UPDATE"
	DELETE       Action = "DELETE"
	TRUNCATE     Action = "TRUNCATE"
	SOFT_DELETE  Action = "SOFT_DELETE"
	SOFT_UPDATE  Action = "SOFT_UPDATE"
	SOFT_RESTORE Action = "SOFT_RESTORE"
)

type Change struct {
	Timestamp time.Time      `json:"timestamp"`
	ID        uuid.UUID      `json:"id"`
	Action    Action         `json:"action"`
	TableName string         `json:"table_name"`
	Item      map[string]any `json:"item"`
	Xid       uint32         `json:"xid"`
}

func (c *Change) String() string {
	return fmt.Sprintf(
		"(%s / %d) @ %s; %s %s",
		c.ID, c.Xid, c.Timestamp, c.Action, c.TableName,
	)
}
