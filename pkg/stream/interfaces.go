package stream

import (
	"github.com/initialed85/djangolang/pkg/raw_stream"
)

type Action = raw_stream.Action

const (
	INSERT       = raw_stream.INSERT
	UPDATE       = raw_stream.UPDATE
	DELETE       = raw_stream.DELETE
	TRUNCATE     = raw_stream.TRUNCATE
	SOFT_DELETE  = raw_stream.SOFT_DELETE
	SOFT_UPDATE  = raw_stream.SOFT_UPDATE
	SOFT_RESTORE = raw_stream.SOFT_RESTORE
)

type Change = raw_stream.Change
