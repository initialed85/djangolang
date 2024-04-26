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

	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/pg_types"
	"github.com/initialed85/djangolang/pkg/types"
	"github.com/initialed85/djangolang/pkg/introspect"
	"github.com/jmoiron/sqlx"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/lib/pq/hstore"
	geom "github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
)
`) + "\n\n"
