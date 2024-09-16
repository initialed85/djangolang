package config

import (
	"fmt"
	_log "log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/initialed85/djangolang/pkg/helpers"
)

var log = helpers.GetLogger("config")

func ThisLogger() *_log.Logger {
	return log
}

// general
var debug bool
var queryDebug bool
var profile bool

// template
var modulePath string
var packageName string
var apiRoot string
var apiRootForOpenAPI string

// postgres
var postgresUser string
var postgresPassword string
var postgresHost string
var postgresPort uint16
var postgresDatabase string
var postgresSSLMode bool
var postgresSchema string

// stream
var nodeName string
var setReplicaIdentity string
var redisUrl string
var reloadChangeObjects bool

// server
var port uint16
var cors string

func init() {
	//
	// general
	//

	debug = GetEnvironmentVariableOrDefault("DJANGOLANG_DEBUG", "0") == "1"

	queryDebug = GetEnvironmentVariableOrDefault("DJANGOLANG_QUERY_DEBUG", "0") == "1"

	profile = GetEnvironmentVariableOrDefault("DJANGOLANG_PROFILE", "0") == "1"

	defaultModulePath := ""
	defaultPackageName := ""

	func() {
		wd, err := os.Getwd()
		if err != nil {
			return
		}

		b, err := os.ReadFile(path.Join(wd, "go.mod"))
		if err != nil {
			return
		}

		parts := strings.Split(strings.Split(strings.TrimSpace(string(b)), "\n")[0], "module ")

		if len(parts) < 2 {
			return
		}

		if defaultModulePath == "" {
			defaultModulePath = parts[1]
		}

		defaultModulePathParts := strings.Split(defaultModulePath, "/")

		if defaultPackageName == "" {
			defaultPackageName = defaultModulePathParts[len(defaultModulePathParts)-1]
		}
	}()

	if defaultModulePath != "" {
		modulePath = GetEnvironmentVariableOrDefault("DJANGOLANG_MODULE_PATH", defaultModulePath)
	} else {
		modulePath = GetEnvironmentVariableOrDefault("DJANGOLANG_MODULE_PATH", "github.com/initialed85/djangolang")
	}

	if defaultPackageName != "" {
		packageName = GetEnvironmentVariableOrDefault("DJANGOLANG_PACKAGE_NAME", defaultPackageName)
	} else {
		packageName = GetEnvironmentVariableOrDefault("DJANGOLANG_PACKAGE_NAME", "djangolang")
	}

	apiRoot = GetEnvironmentVariableOrDefault("DJANGOLANG_API_ROOT", "/")

	apiRootForOpenAPI = GetEnvironmentVariableOrDefault("DJANGOLANG_API_ROOT_FOR_OPENAPI", apiRoot)

	//
	// postgres
	//

	postgresUser = GetEnvironmentVariableOrDefault("POSTGRES_USER", "postgres")

	postgresPassword = GetEnvironmentVariableOrFatal("POSTGRES_PASSWORD")

	postgresHost = GetEnvironmentVariableOrDefault("POSTGRES_HOST", "localhost")

	rawPostgresPort := GetEnvironmentVariableOrDefault("POSTGRES_PORT", "5432")

	almostPostgresPort, err := strconv.ParseUint(rawPostgresPort, 10, 16)
	if err != nil {
		log.Fatalf("POSTGRES_PORT env var (%#+v) not parseable as uint16: %v", rawPostgresPort, err)
	}

	if almostPostgresPort > 65535 {
		log.Fatalf("POSTGRES_PORT env var (%#+v) not parseable as uint16: too large", rawPostgresPort)
	}

	postgresPort = uint16(almostPostgresPort)

	postgresDatabase = GetEnvironmentVariableOrFatal("POSTGRES_DB")

	postgresSSLMode = GetEnvironmentVariableOrDefault("POSTGRES_SSLMODE", "0") == "1"

	postgresSchema = GetEnvironmentVariableOrDefault("POSTGRES_SCHEMA", "public")

	//
	// stream
	//

	nodeName = GetEnvironmentVariableOrDefault("DJANGOLANG_NODE_NAME", "default")

	setReplicaIdentity = GetEnvironmentVariableOrDefault(
		"DJANGOLANG_SET_REPLICA_IDENTITY",
		"full",
		"warning: non-empty values for this flag will set the replica identity for all tables, which will persist across database restarts",
	)

	redisUrl = GetEnvironmentVariableOrFatal("REDIS_URL")

	reloadChangeObjects = GetEnvironmentVariableOrDefault(
		"DJANGOLANG_RELOAD_CHANGE_OBJECTS",
		"0",
		"warning: setting this flag can add load to the database (as each streamed change is followed up with a select for the object in question)",
	) == "1"

	//
	// server
	//

	rawPort := GetEnvironmentVariableOrDefault("PORT", "7070")

	almostPort, err := strconv.ParseUint(rawPort, 10, 16)
	if err != nil {
		log.Fatalf("PORT env var (%#+v) not parseable as uint16: %v", rawPort, err)
	}

	if almostPort > 65535 {
		log.Fatalf("PORT env var (%#+v) not parseable as uint16: too large", rawPort)
	}

	port = uint16(almostPort)

	cors = GetEnvironmentVariableOrDefault("DJANGOLANG_CORS", "*")
}

//
// general
//

func Debug() bool {
	return debug
}

func QueryDebug() bool {
	return queryDebug
}

func Profile() bool {
	return profile
}

//
// template
//

func ModulePath() string {
	return modulePath
}

func PackageName() string {
	return packageName
}

func APIRoot() string {
	return apiRoot
}

func APIRootForOpenAPI() string {
	return apiRootForOpenAPI
}

//
// postgres
//

func PostgresUser() string {
	return postgresUser
}

func PostgresPassword() string {
	return postgresPassword
}

func PostgresHost() string {
	return postgresHost
}

func PostgresPort() uint16 {
	return postgresPort
}

func PostgresDatabase() string {
	return postgresDatabase
}

func PostgresSSLMode() bool {
	return postgresSSLMode
}

func PostgresSchema() string {
	return postgresSchema
}

//
// stream
//

func NodeName() string {
	return nodeName
}

func SetReplicaIdentity() string {
	return setReplicaIdentity
}

func RedisURL() string {
	return redisUrl
}

func ReloadChangeObjects() bool {
	return reloadChangeObjects
}

//
// port
//

func Port() uint16 {
	return port
}

func CORS() string {
	return cors
}

//
// other
//

func DumpConfig() {
	fmt.Printf("DJANGOLANG_DEBUG=%v\n", Debug())
	fmt.Printf("DJANGOLANG_QUERY_DEBUG=%v\n", QueryDebug())
	fmt.Printf("DJANGOLANG_PROFILE=%v\n", Profile())
	fmt.Printf("DJANGOLANG_MODULE_PATH=%v\n", ModulePath())
	fmt.Printf("DJANGOLANG_PACKAGE_NAME=%v\n", PackageName())
	fmt.Printf("DJANGOLANG_API_ROOT=%v\n", APIRoot())
	fmt.Printf("DJANGOLANG_API_ROOT_FOR_OPENAPI=%v\n", APIRootForOpenAPI())
	fmt.Printf("POSTGRES_USER=%v\n", PostgresUser())
	fmt.Printf("POSTGRES_PASSWORD=%v\n", PostgresPassword())
	fmt.Printf("POSTGRES_HOST=%v\n", PostgresHost())
	fmt.Printf("POSTGRES_PORT=%v\n", PostgresPort())
	fmt.Printf("POSTGRES_DB=%v\n", PostgresDatabase())
	fmt.Printf("POSTGRES_SSLMODE=%v\n", PostgresSSLMode())
	fmt.Printf("POSTGRES_SCHEMA=%v\n", PostgresSchema())
	fmt.Printf("DJANGOLANG_NODE_NAME=%v\n", NodeName())
	fmt.Printf("DJANGOLANG_SET_REPLICA_IDENTITY=%v\n", SetReplicaIdentity())
	fmt.Printf("REDIS_URL=%v\n", RedisURL())
	fmt.Printf("DJANGOLANG_RELOAD_CHANGE_OBJECTS=%v\n", ReloadChangeObjects())
	fmt.Printf("PORT=%v\n", Port())
	fmt.Printf("DJANGOLANG_CORS=%v\n", CORS())
}
