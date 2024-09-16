package helpers

import (
	_log "log"
	"os"
	"strings"
)

var log = GetLogger("helpers")

func GetLogger(prefix string) *_log.Logger {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		prefix = "(unknown)"
	}

	for len(prefix) > 24 {
		prefix = prefix[:len(prefix)-1]
	}

	for len(prefix) < 24 {
		prefix += " "
	}

	return _log.New(
		os.Stdout,
		prefix,
		// _log.Ldate|_log.Ltime|_log.Lmicroseconds|_log.Lshortfile|_log.Lmsgprefix,
		_log.Ldate|_log.Ltime|_log.Lmicroseconds|_log.Lmsgprefix,
	)
}
