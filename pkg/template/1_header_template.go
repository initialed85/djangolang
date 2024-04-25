package template

import "strings"

var headerTemplate = strings.TrimSpace(`
package templates

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/google/uuid"

	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/types"
	"github.com/initialed85/djangolang/pkg/introspect"
	"github.com/jmoiron/sqlx"
)
`) + "\n\n"
