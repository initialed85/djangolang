package helpers

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func GetLogger(prefix string) *log.Logger {
	prefix = strings.TrimRight(strings.TrimSpace(prefix), ":")
	if prefix != "" {
		prefix = fmt.Sprintf("%v: ", prefix)
	}

	return log.New(
		os.Stdout,
		prefix,
		log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile|log.Lmsgprefix,
	)
}
