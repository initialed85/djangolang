package djangolang

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chanced/caps"
	"github.com/google/uuid"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/jmoiron/sqlx"
)

type SelectFunc = func(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error)
type Handler = func(w http.ResponseWriter, r *http.Request)
type SelectHandler = Handler

var (
	logger                = helpers.GetLogger("djangolang")
	mu                    = new(sync.RWMutex)
	actualDebug           = false
	selectFuncByTableName = make(map[string]SelectFunc)
	columnsByTableName    = make(map[string][]string)
)

func init() {
	mu.Lock()
	defer mu.Unlock()
	rawDesiredDebug := os.Getenv("DJANGOLANG_DEBUG")
	actualDebug = rawDesiredDebug == "1"
	logger.Printf("DJANGOLANG_DEBUG=%v, debugging enabled: %v", rawDesiredDebug, actualDebug)

	selectFuncByTableName["amenities"] = genericSelectAmenities
	columnsByTableName["amenities"] = AmenitiesColumns
	selectFuncByTableName["charging_station_log_error_codes"] = genericSelectChargingStationLogErrorCodes
	columnsByTableName["charging_station_log_error_codes"] = ChargingStationLogErrorCodesColumns
	selectFuncByTableName["charging_stations"] = genericSelectChargingStations
	columnsByTableName["charging_stations"] = ChargingStationsColumns
	selectFuncByTableName["configurations"] = genericSelectConfigurations
	columnsByTableName["configurations"] = ConfigurationsColumns
	selectFuncByTableName["connectors"] = genericSelectConnectors
	columnsByTableName["connectors"] = ConnectorsColumns
	selectFuncByTableName["cpo_api_keys"] = genericSelectCPOAPIKeys
	columnsByTableName["cpo_api_keys"] = CPOAPIKeysColumns
	selectFuncByTableName["discounts"] = genericSelectDiscounts
	columnsByTableName["discounts"] = DiscountsColumns
	selectFuncByTableName["files"] = genericSelectFiles
	columnsByTableName["files"] = FilesColumns
	selectFuncByTableName["groups"] = genericSelectGroups
	columnsByTableName["groups"] = GroupsColumns
	selectFuncByTableName["ihomer_station_logs"] = genericSelectIhomerStationLogs
	columnsByTableName["ihomer_station_logs"] = IhomerStationLogsColumns
	selectFuncByTableName["ihomer_transactions"] = genericSelectIhomerTransactions
	columnsByTableName["ihomer_transactions"] = IhomerTransactionsColumns
	selectFuncByTableName["local_token_charging_stations"] = genericSelectLocalTokenChargingStations
	columnsByTableName["local_token_charging_stations"] = LocalTokenChargingStationsColumns
	selectFuncByTableName["location_amenities"] = genericSelectLocationAmenities
	columnsByTableName["location_amenities"] = LocationAmenitiesColumns
	selectFuncByTableName["location_groups"] = genericSelectLocationGroups
	columnsByTableName["location_groups"] = LocationGroupsColumns
	selectFuncByTableName["location_images"] = genericSelectLocationImages
	columnsByTableName["location_images"] = LocationImagesColumns
	selectFuncByTableName["location_opening_times"] = genericSelectLocationOpeningTimes
	columnsByTableName["location_opening_times"] = LocationOpeningTimesColumns
	selectFuncByTableName["locations"] = genericSelectLocations
	columnsByTableName["locations"] = LocationsColumns
	selectFuncByTableName["networks"] = genericSelectNetworks
	columnsByTableName["networks"] = NetworksColumns
	selectFuncByTableName["payment_intents"] = genericSelectPaymentIntents
	columnsByTableName["payment_intents"] = PaymentIntentsColumns
	selectFuncByTableName["raw_cdrs"] = genericSelectRawCdrs
	columnsByTableName["raw_cdrs"] = RawCdrsColumns
	selectFuncByTableName["schema_migrations"] = genericSelectSchemaMigrations
	columnsByTableName["schema_migrations"] = SchemaMigrationsColumns
	selectFuncByTableName["stats"] = genericSelectStats
	columnsByTableName["stats"] = StatsColumns
	selectFuncByTableName["stripe_invoices"] = genericSelectStripeInvoices
	columnsByTableName["stripe_invoices"] = StripeInvoicesColumns
	selectFuncByTableName["stripe_subscriptions"] = genericSelectStripeSubscriptions
	columnsByTableName["stripe_subscriptions"] = StripeSubscriptionsColumns
	selectFuncByTableName["subscription_types"] = genericSelectSubscriptionTypes
	columnsByTableName["subscription_types"] = SubscriptionTypesColumns
	selectFuncByTableName["token_groups"] = genericSelectTokenGroups
	columnsByTableName["token_groups"] = TokenGroupsColumns
	selectFuncByTableName["token_readonly_users"] = genericSelectTokenReadonlyUsers
	columnsByTableName["token_readonly_users"] = TokenReadonlyUsersColumns
	selectFuncByTableName["tokens"] = genericSelectTokens
	columnsByTableName["tokens"] = TokensColumns
	selectFuncByTableName["user_evs"] = genericSelectUserEvs
	columnsByTableName["user_evs"] = UserEvsColumns
	selectFuncByTableName["user_groups"] = genericSelectUserGroups
	columnsByTableName["user_groups"] = UserGroupsColumns
	selectFuncByTableName["user_readonly_tokens"] = genericSelectUserReadonlyTokens
	columnsByTableName["user_readonly_tokens"] = UserReadonlyTokensColumns
	selectFuncByTableName["user_roles"] = genericSelectUserRoles
	columnsByTableName["user_roles"] = UserRolesColumns
	selectFuncByTableName["user_sessions"] = genericSelectUserSessions
	columnsByTableName["user_sessions"] = UserSessionsColumns
	selectFuncByTableName["user_subscriptions"] = genericSelectUserSubscriptions
	columnsByTableName["user_subscriptions"] = UserSubscriptionsColumns
	selectFuncByTableName["users"] = genericSelectUsers
	columnsByTableName["users"] = UsersColumns
	selectFuncByTableName["v_charging_stations_private"] = genericSelectVChargingStationsPrivate
	columnsByTableName["v_charging_stations_private"] = VChargingStationsPrivateColumns
	selectFuncByTableName["v_charging_stations_private_over_time"] = genericSelectVChargingStationsPrivateOverTime
	columnsByTableName["v_charging_stations_private_over_time"] = VChargingStationsPrivateOverTimeColumns
	selectFuncByTableName["v_charging_stations_public"] = genericSelectVChargingStationsPublic
	columnsByTableName["v_charging_stations_public"] = VChargingStationsPublicColumns
	selectFuncByTableName["v_charging_stations_public_over_time"] = genericSelectVChargingStationsPublicOverTime
	columnsByTableName["v_charging_stations_public_over_time"] = VChargingStationsPublicOverTimeColumns
	selectFuncByTableName["v_cpos"] = genericSelectVCpos
	columnsByTableName["v_cpos"] = VCposColumns
	selectFuncByTableName["v_locations_private"] = genericSelectVLocationsPrivate
	columnsByTableName["v_locations_private"] = VLocationsPrivateColumns
	selectFuncByTableName["v_locations_public"] = genericSelectVLocationsPublic
	columnsByTableName["v_locations_public"] = VLocationsPublicColumns
	selectFuncByTableName["v_user_sessions_private"] = genericSelectVUserSessionsPrivate
	columnsByTableName["v_user_sessions_private"] = VUserSessionsPrivateColumns
	selectFuncByTableName["v_user_sessions_public"] = genericSelectVUserSessionsPublic
	columnsByTableName["v_user_sessions_public"] = VUserSessionsPublicColumns
	selectFuncByTableName["v_users"] = genericSelectVUsers
	columnsByTableName["v_users"] = VUsersColumns
	selectFuncByTableName["v_users_registered_over_time"] = genericSelectVUsersRegisteredOverTime
	columnsByTableName["v_users_registered_over_time"] = VUsersRegisteredOverTimeColumns
	selectFuncByTableName["voltguard_events"] = genericSelectVoltguardEvents
	columnsByTableName["voltguard_events"] = VoltguardEventsColumns
}

func SetDebug(desiredDebug bool) {
	mu.Lock()
	defer mu.Unlock()
	actualDebug = desiredDebug
	logger.Printf("runtime SetDebug() called, debugging enabled: %v", actualDebug)
}

func Descending(columns ...string) *string {
	return helpers.Ptr(
		fmt.Sprintf(
			"(%v) DESC",
			strings.Join(columns, ", "),
		),
	)
}

func Ascending(columns ...string) *string {
	return helpers.Ptr(
		fmt.Sprintf(
			"(%v) ASC",
			strings.Join(columns, ", "),
		),
	)
}

func Columns(includeColumns []string, excludeColumns ...string) []string {
	excludeColumnLookup := make(map[string]bool)
	for _, column := range excludeColumns {
		excludeColumnLookup[column] = true
	}

	columns := make([]string, 0)
	for _, column := range includeColumns {
		_, ok := excludeColumnLookup[column]
		if ok {
			continue
		}

		columns = append(columns, column)
	}

	return columns
}

func GetSelectFuncByTableName() map[string]SelectFunc {
	thisSelectFuncByTableName := make(map[string]SelectFunc)

	for tableName, selectFunc := range selectFuncByTableName {
		thisSelectFuncByTableName[tableName] = selectFunc
	}

	return thisSelectFuncByTableName
}

func GetSelectHandlerForTableName(tableName string, db *sqlx.DB) (SelectHandler, error) {
	selectFunc, ok := selectFuncByTableName[tableName]
	if !ok {
		return nil, fmt.Errorf("no selectFuncByTableName entry for tableName %v", tableName)
	}

	columns, ok := columnsByTableName[tableName]
	if !ok {
		return nil, fmt.Errorf("no columnsByTableName entry for tableName %v", tableName)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK
		var body []byte
		var err error

		defer func() {
			w.Header().Add("Content-Type", "application/json")

			w.WriteHeader(status)

			if err != nil {
				_, _ = w.Write([]byte(fmt.Sprintf("{\"error\": %#+v}", err.Error())))
				return
			}

			_, _ = w.Write(body)
		}()

		if r.Method != http.MethodGet {
			status = http.StatusMethodNotAllowed
			err = fmt.Errorf("%v; wanted %v, got %v", "StatusMethodNotAllowed", http.MethodGet, r.Method)
			return
		}

		limit := helpers.Ptr(50)
		rawLimit := strings.TrimSpace(r.URL.Query().Get("limit"))
		if len(rawLimit) > 0 {
			var parsedLimit int64
			parsedLimit, err = strconv.ParseInt(rawLimit, 10, 64)
			if err != nil {
				status = http.StatusBadRequest
				err = fmt.Errorf("failed to parse limit %#+v as int; err: %v", limit, err.Error())
				return
			}
			limit = helpers.Ptr(int(parsedLimit))
		}

		offset := helpers.Ptr(50)
		rawOffset := strings.TrimSpace(r.URL.Query().Get("offset"))
		if len(rawOffset) > 0 {
			var parsedOffset int64
			parsedOffset, err = strconv.ParseInt(rawOffset, 10, 64)
			if err != nil {
				status = http.StatusBadRequest
				err = fmt.Errorf("failed to parse offset %#+v as int; err: %v", offset, err.Error())
				return
			}
			offset = helpers.Ptr(int(parsedOffset))
		}

		var order *string

		wheres := make([]string, 0)

		var items []any
		items, err = selectFunc(r.Context(), db, columns, order, limit, offset, wheres...)
		if err != nil {
			status = http.StatusInternalServerError
			err = fmt.Errorf("query failed; err: %v", err.Error())
			return
		}

		body, err = json.Marshal(items)
		if err != nil {
			status = http.StatusInternalServerError
			err = fmt.Errorf("serialization failed; err: %v", err.Error())
		}
	}, nil
}

func GetSelectHandlerByEndpointName(db *sqlx.DB) (map[string]SelectHandler, error) {
	thisSelectHandlerByEndpointName := make(map[string]SelectHandler)

	for _, tableName := range TableNames {
		endpointName := caps.ToKebab[string](tableName)

		selectHandler, err := GetSelectHandlerForTableName(tableName, db)
		if err != nil {
			return nil, err
		}

		thisSelectHandlerByEndpointName[endpointName] = selectHandler
	}

	return thisSelectHandlerByEndpointName, nil
}

var TableNames = []string{"amenities", "charging_station_log_error_codes", "charging_stations", "configurations", "connectors", "cpo_api_keys", "discounts", "files", "groups", "ihomer_station_logs", "ihomer_transactions", "local_token_charging_stations", "location_amenities", "location_groups", "location_images", "location_opening_times", "locations", "networks", "payment_intents", "raw_cdrs", "schema_migrations", "stats", "stripe_invoices", "stripe_subscriptions", "subscription_types", "token_groups", "token_readonly_users", "tokens", "user_evs", "user_groups", "user_readonly_tokens", "user_roles", "user_sessions", "user_subscriptions", "users", "v_charging_stations_private", "v_charging_stations_private_over_time", "v_charging_stations_public", "v_charging_stations_public_over_time", "v_cpos", "v_locations_private", "v_locations_public", "v_user_sessions_private", "v_user_sessions_public", "v_users", "v_users_registered_over_time", "voltguard_events"}

var AmenitiesTable = "amenities"
var AmenitiesIDColumn = "id"
var AmenitiesNameColumn = "name"
var AmenitiesCreatedAtColumn = "created_at"
var AmenitiesUpdatedAtColumn = "updated_at"
var AmenitiesDeletedAtColumn = "deleted_at"
var AmenitiesColumns = []string{"id", "name", "created_at", "updated_at", "deleted_at"}

type Amenities struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

func SelectAmenities(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*Amenities, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM amenities%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*Amenities, 0)
	for rows.Next() {
		rowCount++

		var item Amenities
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectAmenities(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectAmenities(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var ChargingStationLogErrorCodesTable = "charging_station_log_error_codes"
var ChargingStationLogErrorCodesIDColumn = "id"
var ChargingStationLogErrorCodesErrorCodeColumn = "error_code"
var ChargingStationLogErrorCodesDescriptionColumn = "description"
var ChargingStationLogErrorCodesLabelColumn = "label"
var ChargingStationLogErrorCodesIsActiveColumn = "is_active"
var ChargingStationLogErrorCodesIsOutsideStationLogColumn = "is_outside_station_log"
var ChargingStationLogErrorCodesColumns = []string{"id", "error_code", "description", "label", "is_active", "is_outside_station_log"}

type ChargingStationLogErrorCodes struct {
	ID                  uuid.UUID `json:"id" db:"id"`
	ErrorCode           string    `json:"error_code" db:"error_code"`
	Description         string    `json:"description" db:"description"`
	Label               string    `json:"label" db:"label"`
	IsActive            *bool     `json:"is_active" db:"is_active"`
	IsOutsideStationLog *bool     `json:"is_outside_station_log" db:"is_outside_station_log"`
}

func SelectChargingStationLogErrorCodes(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*ChargingStationLogErrorCodes, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM charging_station_log_error_codes%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*ChargingStationLogErrorCodes, 0)
	for rows.Next() {
		rowCount++

		var item ChargingStationLogErrorCodes
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectChargingStationLogErrorCodes(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectChargingStationLogErrorCodes(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var ChargingStationsTable = "charging_stations"
var ChargingStationsIDColumn = "id"
var ChargingStationsLocationIDColumn = "location_id"
var ChargingStationsUserIDColumn = "user_id"
var ChargingStationsProtocolColumn = "protocol"
var ChargingStationsAcceptedColumn = "accepted"
var ChargingStationsReservableColumn = "reservable"
var ChargingStationsConfiguredColumn = "configured"
var ChargingStationsAccessibilityColumn = "accessibility"
var ChargingStationsAvailabilityColumn = "availability"
var ChargingStationsStatusColumn = "status"
var ChargingStationsSubStatusColumn = "sub_status"
var ChargingStationsLocalAuthorizationListVersionColumn = "local_authorization_list_version"
var ChargingStationsIsActiveColumn = "is_active"
var ChargingStationsIsPublicColumn = "is_public"
var ChargingStationsCreatedAtColumn = "created_at"
var ChargingStationsUpdatedAtColumn = "updated_at"
var ChargingStationsDeletedAtColumn = "deleted_at"
var ChargingStationsAliasColumn = "alias"
var ChargingStationsPlugshareJSONDataColumn = "plugshare_json_data"
var ChargingStationsIsConfiguredColumn = "is_configured"
var ChargingStationsPlugshareCostTypeColumn = "plugshare_cost_type"
var ChargingStationsNetworkIDColumn = "network_id"
var ChargingStationsCredentialsGeneratedColumn = "credentials_generated"
var ChargingStationsRanGenerateCredentialsColumn = "ran_generate_credentials"
var ChargingStationsNeedsUnlockFirstColumn = "needs_unlock_first"
var ChargingStationsCurrentColumn = "current"
var ChargingStationsIsV2ChargerColumn = "is_v2_charger"
var ChargingStationsLastSeenColumn = "last_seen"
var ChargingStationsOcppServerIDColumn = "ocpp_server_id"
var ChargingStationsColumns = []string{"id", "location_id", "user_id", "protocol", "accepted", "reservable", "configured", "accessibility", "availability", "status", "sub_status", "local_authorization_list_version", "is_active", "is_public", "created_at", "updated_at", "deleted_at", "alias", "plugshare_json_data", "is_configured", "plugshare_cost_type", "network_id", "credentials_generated", "ran_generate_credentials", "needs_unlock_first", "current", "is_v2_charger", "last_seen", "ocpp_server_id"}

type ChargingStations struct {
	ID                            string     `json:"id" db:"id"`
	LocationID                    uuid.UUID  `json:"location_id" db:"location_id"`
	UserID                        uuid.UUID  `json:"user_id" db:"user_id"`
	Protocol                      string     `json:"protocol" db:"protocol"`
	Accepted                      bool       `json:"accepted" db:"accepted"`
	Reservable                    bool       `json:"reservable" db:"reservable"`
	Configured                    bool       `json:"configured" db:"configured"`
	Accessibility                 string     `json:"accessibility" db:"accessibility"`
	Availability                  string     `json:"availability" db:"availability"`
	Status                        string     `json:"status" db:"status"`
	SubStatus                     string     `json:"sub_status" db:"sub_status"`
	LocalAuthorizationListVersion int64      `json:"local_authorization_list_version" db:"local_authorization_list_version"`
	IsActive                      bool       `json:"is_active" db:"is_active"`
	IsPublic                      bool       `json:"is_public" db:"is_public"`
	CreatedAt                     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt                     time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt                     *time.Time `json:"deleted_at" db:"deleted_at"`
	Alias                         string     `json:"alias" db:"alias"`
	PlugshareJSONData             any        `json:"plugshare_json_data" db:"plugshare_json_data"`
	IsConfigured                  bool       `json:"is_configured" db:"is_configured"`
	PlugshareCostType             int64      `json:"plugshare_cost_type" db:"plugshare_cost_type"`
	NetworkID                     *int64     `json:"network_id" db:"network_id"`
	CredentialsGenerated          *bool      `json:"credentials_generated" db:"credentials_generated"`
	RanGenerateCredentials        *bool      `json:"ran_generate_credentials" db:"ran_generate_credentials"`
	NeedsUnlockFirst              bool       `json:"needs_unlock_first" db:"needs_unlock_first"`
	Current                       string     `json:"current" db:"current"`
	IsV2Charger                   *bool      `json:"is_v2_charger" db:"is_v2_charger"`
	LastSeen                      *time.Time `json:"last_seen" db:"last_seen"`
	OcppServerID                  *uuid.UUID `json:"ocpp_server_id" db:"ocpp_server_id"`
}

func SelectChargingStations(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*ChargingStations, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM charging_stations%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*ChargingStations, 0)
	for rows.Next() {
		rowCount++

		var item ChargingStations
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		var temp1PlugshareJSONData any
		var temp2PlugshareJSONData []any

		err = json.Unmarshal(item.PlugshareJSONData.([]byte), &temp1PlugshareJSONData)
		if err != nil {
			err = json.Unmarshal(item.PlugshareJSONData.([]byte), &temp2PlugshareJSONData)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal %#+v as json: %v", string(item.PlugshareJSONData.([]byte)), err)
			} else {
				item.PlugshareJSONData = temp2PlugshareJSONData
			}
		} else {
			item.PlugshareJSONData = temp1PlugshareJSONData
		}

		item.PlugshareJSONData = temp1PlugshareJSONData

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectChargingStations(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectChargingStations(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var ConfigurationsTable = "configurations"
var ConfigurationsIDColumn = "id"
var ConfigurationsParkingGraceMinsColumn = "parking_grace_mins"
var ConfigurationsMinChargePriceCentsColumn = "min_charge_price_cents"
var ConfigurationsCreatedByColumn = "created_by"
var ConfigurationsUpdatedByColumn = "updated_by"
var ConfigurationsCreatedAtColumn = "created_at"
var ConfigurationsUpdatedAtColumn = "updated_at"
var ConfigurationsDeletedAtColumn = "deleted_at"
var ConfigurationsEnableChargeFlowV2Column = "enable_charge_flow_v2"
var ConfigurationsAuthorizationHoldAmountCentsColumn = "authorization_hold_amount_cents"
var ConfigurationsEnablePrecalculateStripeFeesColumn = "enable_precalculate_stripe_fees"
var ConfigurationsColumns = []string{"id", "parking_grace_mins", "min_charge_price_cents", "created_by", "updated_by", "created_at", "updated_at", "deleted_at", "enable_charge_flow_v2", "authorization_hold_amount_cents", "enable_precalculate_stripe_fees"}

type Configurations struct {
	ID                           uuid.UUID  `json:"id" db:"id"`
	ParkingGraceMins             *int64     `json:"parking_grace_mins" db:"parking_grace_mins"`
	MinChargePriceCents          *int64     `json:"min_charge_price_cents" db:"min_charge_price_cents"`
	CreatedBy                    *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy                    *uuid.UUID `json:"updated_by" db:"updated_by"`
	CreatedAt                    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt                    time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt                    *time.Time `json:"deleted_at" db:"deleted_at"`
	EnableChargeFlowV2           bool       `json:"enable_charge_flow_v2" db:"enable_charge_flow_v2"`
	AuthorizationHoldAmountCents int64      `json:"authorization_hold_amount_cents" db:"authorization_hold_amount_cents"`
	EnablePrecalculateStripeFees bool       `json:"enable_precalculate_stripe_fees" db:"enable_precalculate_stripe_fees"`
}

func SelectConfigurations(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*Configurations, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM configurations%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*Configurations, 0)
	for rows.Next() {
		rowCount++

		var item Configurations
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectConfigurations(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectConfigurations(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var ConnectorsTable = "connectors"
var ConnectorsIDColumn = "id"
var ConnectorsPhaseColumn = "phase"
var ConnectorsMaxAmpColumn = "max_amp"
var ConnectorsVoltageColumn = "voltage"
var ConnectorsMaxElectricPowerColumn = "max_electric_power"
var ConnectorsChargingProtocolColumn = "charging_protocol"
var ConnectorsStatusColumn = "status"
var ConnectorsSubStatusColumn = "sub_status"
var ConnectorsCurrentColumn = "current"
var ConnectorsChargingStationIDColumn = "charging_station_id"
var ConnectorsConnectorTypeColumn = "connector_type"
var ConnectorsCreatedAtColumn = "created_at"
var ConnectorsUpdatedAtColumn = "updated_at"
var ConnectorsDeletedAtColumn = "deleted_at"
var ConnectorsConnectorIDColumn = "connector_id"
var ConnectorsIsPlugshareColumn = "is_plugshare"
var ConnectorsLastSeenColumn = "last_seen"
var ConnectorsAvailabilityColumn = "availability"
var ConnectorsColumns = []string{"id", "phase", "max_amp", "voltage", "max_electric_power", "charging_protocol", "status", "sub_status", "current", "charging_station_id", "connector_type", "created_at", "updated_at", "deleted_at", "connector_id", "is_plugshare", "last_seen", "availability"}

type Connectors struct {
	ID                uuid.UUID  `json:"id" db:"id"`
	Phase             float64    `json:"phase" db:"phase"`
	MaxAmp            float64    `json:"max_amp" db:"max_amp"`
	Voltage           float64    `json:"voltage" db:"voltage"`
	MaxElectricPower  float64    `json:"max_electric_power" db:"max_electric_power"`
	ChargingProtocol  string     `json:"charging_protocol" db:"charging_protocol"`
	Status            string     `json:"status" db:"status"`
	SubStatus         string     `json:"sub_status" db:"sub_status"`
	Current           string     `json:"current" db:"current"`
	ChargingStationID string     `json:"charging_station_id" db:"charging_station_id"`
	ConnectorType     string     `json:"connector_type" db:"connector_type"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at" db:"deleted_at"`
	ConnectorID       int64      `json:"connector_id" db:"connector_id"`
	IsPlugshare       bool       `json:"is_plugshare" db:"is_plugshare"`
	LastSeen          *time.Time `json:"last_seen" db:"last_seen"`
	Availability      string     `json:"availability" db:"availability"`
}

func SelectConnectors(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*Connectors, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM connectors%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*Connectors, 0)
	for rows.Next() {
		rowCount++

		var item Connectors
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectConnectors(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectConnectors(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var CPOAPIKeysTable = "cpo_api_keys"
var CPOAPIKeysIDColumn = "id"
var CPOAPIKeysCPOIDColumn = "cpo_id"
var CPOAPIKeysAPIKeyColumn = "api_key"
var CPOAPIKeysLabelColumn = "label"
var CPOAPIKeysEnabledColumn = "enabled"
var CPOAPIKeysCreatedAtColumn = "created_at"
var CPOAPIKeysUpdatedAtColumn = "updated_at"
var CPOAPIKeysDeletedAtColumn = "deleted_at"
var CPOAPIKeysColumns = []string{"id", "cpo_id", "api_key", "label", "enabled", "created_at", "updated_at", "deleted_at"}

type CPOAPIKeys struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	CPOID     uuid.UUID  `json:"cpo_id" db:"cpo_id"`
	APIKey    string     `json:"api_key" db:"api_key"`
	Label     string     `json:"label" db:"label"`
	Enabled   bool       `json:"enabled" db:"enabled"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

func SelectCPOAPIKeys(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*CPOAPIKeys, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM cpo_api_keys%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*CPOAPIKeys, 0)
	for rows.Next() {
		rowCount++

		var item CPOAPIKeys
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectCPOAPIKeys(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectCPOAPIKeys(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var DiscountsTable = "discounts"
var DiscountsIDColumn = "id"
var DiscountsSubscriptionTypeIDColumn = "subscription_type_id"
var DiscountsNameColumn = "name"
var DiscountsDiscountPerChargePointPercentColumn = "discount_per_charge_point_percent"
var DiscountsDiscountPerTransactionPercentColumn = "discount_per_transaction_percent"
var DiscountsStartDateColumn = "start_date"
var DiscountsEndDateColumn = "end_date"
var DiscountsIsActiveColumn = "is_active"
var DiscountsStripeCouponIDColumn = "stripe_coupon_id"
var DiscountsCreatedByColumn = "created_by"
var DiscountsCreatedAtColumn = "created_at"
var DiscountsUpdatedAtColumn = "updated_at"
var DiscountsDeletedAtColumn = "deleted_at"
var DiscountsStripePromotionCodeColumn = "stripe_promotion_code"
var DiscountsStripePromotionCodeIDColumn = "stripe_promotion_code_id"
var DiscountsTwoYearStripeCouponIDColumn = "two_year_stripe_coupon_id"
var DiscountsTwoYearStripePromotionCodeColumn = "two_year_stripe_promotion_code"
var DiscountsTwoYearStripePromotionCodeIDColumn = "two_year_stripe_promotion_code_id"
var DiscountsColumns = []string{"id", "subscription_type_id", "name", "discount_per_charge_point_percent", "discount_per_transaction_percent", "start_date", "end_date", "is_active", "stripe_coupon_id", "created_by", "created_at", "updated_at", "deleted_at", "stripe_promotion_code", "stripe_promotion_code_id", "two_year_stripe_coupon_id", "two_year_stripe_promotion_code", "two_year_stripe_promotion_code_id"}

type Discounts struct {
	ID                            uuid.UUID  `json:"id" db:"id"`
	SubscriptionTypeID            uuid.UUID  `json:"subscription_type_id" db:"subscription_type_id"`
	Name                          string     `json:"name" db:"name"`
	DiscountPerChargePointPercent *float64   `json:"discount_per_charge_point_percent" db:"discount_per_charge_point_percent"`
	DiscountPerTransactionPercent *float64   `json:"discount_per_transaction_percent" db:"discount_per_transaction_percent"`
	StartDate                     time.Time  `json:"start_date" db:"start_date"`
	EndDate                       time.Time  `json:"end_date" db:"end_date"`
	IsActive                      bool       `json:"is_active" db:"is_active"`
	StripeCouponID                *string    `json:"stripe_coupon_id" db:"stripe_coupon_id"`
	CreatedBy                     *uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt                     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt                     time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt                     *time.Time `json:"deleted_at" db:"deleted_at"`
	StripePromotionCode           *string    `json:"stripe_promotion_code" db:"stripe_promotion_code"`
	StripePromotionCodeID         *string    `json:"stripe_promotion_code_id" db:"stripe_promotion_code_id"`
	TwoYearStripeCouponID         *string    `json:"two_year_stripe_coupon_id" db:"two_year_stripe_coupon_id"`
	TwoYearStripePromotionCode    *string    `json:"two_year_stripe_promotion_code" db:"two_year_stripe_promotion_code"`
	TwoYearStripePromotionCodeID  *string    `json:"two_year_stripe_promotion_code_id" db:"two_year_stripe_promotion_code_id"`
}

func SelectDiscounts(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*Discounts, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM discounts%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*Discounts, 0)
	for rows.Next() {
		rowCount++

		var item Discounts
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectDiscounts(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectDiscounts(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var FilesTable = "files"
var FilesIDColumn = "id"
var FilesURLColumn = "url"
var FilesNameColumn = "name"
var FilesHashColumn = "hash"
var FilesMimeTypeColumn = "mime_type"
var FilesFileSizeBytesColumn = "file_size_bytes"
var FilesExtensionColumn = "extension"
var FilesCreatedAtColumn = "created_at"
var FilesUpdatedAtColumn = "updated_at"
var FilesDeletedAtColumn = "deleted_at"
var FilesPlugshareIDColumn = "plugshare_id"
var FilesColumns = []string{"id", "url", "name", "hash", "mime_type", "file_size_bytes", "extension", "created_at", "updated_at", "deleted_at", "plugshare_id"}

type Files struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	URL           string     `json:"url" db:"url"`
	Name          string     `json:"name" db:"name"`
	Hash          string     `json:"hash" db:"hash"`
	MimeType      string     `json:"mime_type" db:"mime_type"`
	FileSizeBytes int64      `json:"file_size_bytes" db:"file_size_bytes"`
	Extension     string     `json:"extension" db:"extension"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at" db:"deleted_at"`
	PlugshareID   *int64     `json:"plugshare_id" db:"plugshare_id"`
}

func SelectFiles(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*Files, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM files%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*Files, 0)
	for rows.Next() {
		rowCount++

		var item Files
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectFiles(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectFiles(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var GroupsTable = "groups"
var GroupsIDColumn = "id"
var GroupsNameColumn = "name"
var GroupsTypeColumn = "type"
var GroupsDescriptionColumn = "description"
var GroupsActiveColumn = "active"
var GroupsCreatedByColumn = "created_by"
var GroupsCreatedAtColumn = "created_at"
var GroupsUpdatedAtColumn = "updated_at"
var GroupsDeletedAtColumn = "deleted_at"
var GroupsColumns = []string{"id", "name", "type", "description", "active", "created_by", "created_at", "updated_at", "deleted_at"}

type Groups struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Type        string     `json:"type" db:"type"`
	Description string     `json:"description" db:"description"`
	Active      bool       `json:"active" db:"active"`
	CreatedBy   *uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at" db:"deleted_at"`
}

func SelectGroups(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*Groups, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM groups%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*Groups, 0)
	for rows.Next() {
		rowCount++

		var item Groups
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectGroups(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectGroups(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var IhomerStationLogsTable = "ihomer_station_logs"
var IhomerStationLogsIDColumn = "id"
var IhomerStationLogsCreatedAtColumn = "created_at"
var IhomerStationLogsUpdatedAtColumn = "updated_at"
var IhomerStationLogsDeletedAtColumn = "deleted_at"
var IhomerStationLogsChargingStationIDColumn = "charging_station_id"
var IhomerStationLogsIhomerLogIDColumn = "ihomer_log_id"
var IhomerStationLogsCallIDColumn = "call_id"
var IhomerStationLogsMessageTypeColumn = "message_type"
var IhomerStationLogsActionTypeColumn = "action_type"
var IhomerStationLogsLogTimestampColumn = "log_timestamp"
var IhomerStationLogsPayloadColumn = "payload"
var IhomerStationLogsColumns = []string{"id", "created_at", "updated_at", "deleted_at", "charging_station_id", "ihomer_log_id", "call_id", "message_type", "action_type", "log_timestamp", "payload"}

type IhomerStationLogs struct {
	ID                uuid.UUID  `json:"id" db:"id"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at" db:"deleted_at"`
	ChargingStationID *string    `json:"charging_station_id" db:"charging_station_id"`
	IhomerLogID       int64      `json:"ihomer_log_id" db:"ihomer_log_id"`
	CallID            string     `json:"call_id" db:"call_id"`
	MessageType       string     `json:"message_type" db:"message_type"`
	ActionType        string     `json:"action_type" db:"action_type"`
	LogTimestamp      *string    `json:"log_timestamp" db:"log_timestamp"`
	Payload           *any       `json:"payload" db:"payload"`
}

func SelectIhomerStationLogs(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*IhomerStationLogs, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM ihomer_station_logs%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*IhomerStationLogs, 0)
	for rows.Next() {
		rowCount++

		var item IhomerStationLogs
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectIhomerStationLogs(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectIhomerStationLogs(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var IhomerTransactionsTable = "ihomer_transactions"
var IhomerTransactionsIDColumn = "id"
var IhomerTransactionsCreatedAtColumn = "created_at"
var IhomerTransactionsUpdatedAtColumn = "updated_at"
var IhomerTransactionsDeletedAtColumn = "deleted_at"
var IhomerTransactionsChargingStationIDColumn = "charging_station_id"
var IhomerTransactionsTransactionIDColumn = "transaction_id"
var IhomerTransactionsIhomerTransactionIDColumn = "ihomer_transaction_id"
var IhomerTransactionsTokenUIDColumn = "token_uid"
var IhomerTransactionsEvseIDColumn = "evse_id"
var IhomerTransactionsStartedAtColumn = "started_at"
var IhomerTransactionsEndedAtColumn = "ended_at"
var IhomerTransactionsMeterStartColumn = "meter_start"
var IhomerTransactionsMeterStopColumn = "meter_stop"
var IhomerTransactionsMeterValuesColumn = "meter_values"
var IhomerTransactionsStopReasonColumn = "stop_reason"
var IhomerTransactionsForceStoppedColumn = "force_stopped"
var IhomerTransactionsPoolIDColumn = "pool_id"
var IhomerTransactionsColumns = []string{"id", "created_at", "updated_at", "deleted_at", "charging_station_id", "transaction_id", "ihomer_transaction_id", "token_uid", "evse_id", "started_at", "ended_at", "meter_start", "meter_stop", "meter_values", "stop_reason", "force_stopped", "pool_id"}

type IhomerTransactions struct {
	ID                  uuid.UUID  `json:"id" db:"id"`
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at" db:"deleted_at"`
	ChargingStationID   *string    `json:"charging_station_id" db:"charging_station_id"`
	TransactionID       int64      `json:"transaction_id" db:"transaction_id"`
	IhomerTransactionID string     `json:"ihomer_transaction_id" db:"ihomer_transaction_id"`
	TokenUID            string     `json:"token_uid" db:"token_uid"`
	EvseID              int64      `json:"evse_id" db:"evse_id"`
	StartedAt           time.Time  `json:"started_at" db:"started_at"`
	EndedAt             time.Time  `json:"ended_at" db:"ended_at"`
	MeterStart          int64      `json:"meter_start" db:"meter_start"`
	MeterStop           int64      `json:"meter_stop" db:"meter_stop"`
	MeterValues         *any       `json:"meter_values" db:"meter_values"`
	StopReason          *string    `json:"stop_reason" db:"stop_reason"`
	ForceStopped        *bool      `json:"force_stopped" db:"force_stopped"`
	PoolID              *string    `json:"pool_id" db:"pool_id"`
}

func SelectIhomerTransactions(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*IhomerTransactions, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM ihomer_transactions%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*IhomerTransactions, 0)
	for rows.Next() {
		rowCount++

		var item IhomerTransactions
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectIhomerTransactions(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectIhomerTransactions(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var LocalTokenChargingStationsTable = "local_token_charging_stations"
var LocalTokenChargingStationsTokenIDColumn = "token_id"
var LocalTokenChargingStationsChargingStationIDColumn = "charging_station_id"
var LocalTokenChargingStationsColumns = []string{"token_id", "charging_station_id"}

type LocalTokenChargingStations struct {
	TokenID           uuid.UUID `json:"token_id" db:"token_id"`
	ChargingStationID string    `json:"charging_station_id" db:"charging_station_id"`
}

func SelectLocalTokenChargingStations(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*LocalTokenChargingStations, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM local_token_charging_stations%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*LocalTokenChargingStations, 0)
	for rows.Next() {
		rowCount++

		var item LocalTokenChargingStations
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectLocalTokenChargingStations(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectLocalTokenChargingStations(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var LocationAmenitiesTable = "location_amenities"
var LocationAmenitiesAmenityIDColumn = "amenity_id"
var LocationAmenitiesLocationIDColumn = "location_id"
var LocationAmenitiesColumns = []string{"amenity_id", "location_id"}

type LocationAmenities struct {
	AmenityID  uuid.UUID `json:"amenity_id" db:"amenity_id"`
	LocationID uuid.UUID `json:"location_id" db:"location_id"`
}

func SelectLocationAmenities(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*LocationAmenities, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM location_amenities%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*LocationAmenities, 0)
	for rows.Next() {
		rowCount++

		var item LocationAmenities
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectLocationAmenities(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectLocationAmenities(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var LocationGroupsTable = "location_groups"
var LocationGroupsLocationIDColumn = "location_id"
var LocationGroupsGroupIDColumn = "group_id"
var LocationGroupsColumns = []string{"location_id", "group_id"}

type LocationGroups struct {
	LocationID uuid.UUID `json:"location_id" db:"location_id"`
	GroupID    uuid.UUID `json:"group_id" db:"group_id"`
}

func SelectLocationGroups(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*LocationGroups, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM location_groups%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*LocationGroups, 0)
	for rows.Next() {
		rowCount++

		var item LocationGroups
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectLocationGroups(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectLocationGroups(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var LocationImagesTable = "location_images"
var LocationImagesFileIDColumn = "file_id"
var LocationImagesLocationIDColumn = "location_id"
var LocationImagesOrderNoColumn = "order_no"
var LocationImagesColumns = []string{"file_id", "location_id", "order_no"}

type LocationImages struct {
	FileID     uuid.UUID `json:"file_id" db:"file_id"`
	LocationID uuid.UUID `json:"location_id" db:"location_id"`
	OrderNo    int64     `json:"order_no" db:"order_no"`
}

func SelectLocationImages(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*LocationImages, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM location_images%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*LocationImages, 0)
	for rows.Next() {
		rowCount++

		var item LocationImages
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectLocationImages(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectLocationImages(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var LocationOpeningTimesTable = "location_opening_times"
var LocationOpeningTimesIDColumn = "id"
var LocationOpeningTimesLocationIDColumn = "location_id"
var LocationOpeningTimesDayNumberColumn = "day_number"
var LocationOpeningTimesOpenTimeColumn = "open_time"
var LocationOpeningTimesCloseTimeColumn = "close_time"
var LocationOpeningTimesCreatedAtColumn = "created_at"
var LocationOpeningTimesUpdatedAtColumn = "updated_at"
var LocationOpeningTimesDeletedAtColumn = "deleted_at"
var LocationOpeningTimesEnergyPricePerKWHColumn = "energy_price_per_kwh"
var LocationOpeningTimesChargingTimePricePerHourColumn = "charging_time_price_per_hour"
var LocationOpeningTimesParkingTimePricePerHourColumn = "parking_time_price_per_hour"
var LocationOpeningTimesFreePeriodMinutesColumn = "free_period_minutes"
var LocationOpeningTimesMaxChargeTimeMinutesColumn = "max_charge_time_minutes"
var LocationOpeningTimesColumns = []string{"id", "location_id", "day_number", "open_time", "close_time", "created_at", "updated_at", "deleted_at", "energy_price_per_kwh", "charging_time_price_per_hour", "parking_time_price_per_hour", "free_period_minutes", "max_charge_time_minutes"}

type LocationOpeningTimes struct {
	ID                       uuid.UUID  `json:"id" db:"id"`
	LocationID               uuid.UUID  `json:"location_id" db:"location_id"`
	DayNumber                int64      `json:"day_number" db:"day_number"`
	OpenTime                 time.Time  `json:"open_time" db:"open_time"`
	CloseTime                time.Time  `json:"close_time" db:"close_time"`
	CreatedAt                time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt                *time.Time `json:"deleted_at" db:"deleted_at"`
	EnergyPricePerKWH        float64    `json:"energy_price_per_kwh" db:"energy_price_per_kwh"`
	ChargingTimePricePerHour float64    `json:"charging_time_price_per_hour" db:"charging_time_price_per_hour"`
	ParkingTimePricePerHour  float64    `json:"parking_time_price_per_hour" db:"parking_time_price_per_hour"`
	FreePeriodMinutes        int64      `json:"free_period_minutes" db:"free_period_minutes"`
	MaxChargeTimeMinutes     int64      `json:"max_charge_time_minutes" db:"max_charge_time_minutes"`
}

func SelectLocationOpeningTimes(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*LocationOpeningTimes, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM location_opening_times%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*LocationOpeningTimes, 0)
	for rows.Next() {
		rowCount++

		var item LocationOpeningTimes
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectLocationOpeningTimes(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectLocationOpeningTimes(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var LocationsTable = "locations"
var LocationsIDColumn = "id"
var LocationsCountryCodeColumn = "country_code"
var LocationsPartyIDColumn = "party_id"
var LocationsIsPublishedColumn = "is_published"
var LocationsPublishOnlyToColumn = "publish_only_to"
var LocationsNameColumn = "name"
var LocationsAddressColumn = "address"
var LocationsCityColumn = "city"
var LocationsPostalcodeColumn = "postalcode"
var LocationsStateColumn = "state"
var LocationsCountryColumn = "country"
var LocationsLatitudeColumn = "latitude"
var LocationsLongitudeColumn = "longitude"
var LocationsParkingTypeColumn = "parking_type"
var LocationsTimezoneColumn = "timezone"
var LocationsLastUpdatedColumn = "last_updated"
var LocationsEnabledColumn = "enabled"
var LocationsCreatedAtColumn = "created_at"
var LocationsUpdatedAtColumn = "updated_at"
var LocationsDeletedAtColumn = "deleted_at"
var LocationsIhomerIDColumn = "ihomer_id"
var LocationsChargeStationCountColumn = "charge_station_count"
var LocationsPlugshareIDColumn = "plugshare_id"
var LocationsOpen247Column = "open_247"
var LocationsIsActiveColumn = "is_active"
var LocationsCreatedByColumn = "created_by"
var LocationsIsPublicColumn = "is_public"
var LocationsPlugshareLocationsJSONDataColumn = "plugshare_locations_json_data"
var LocationsDefaultEnergyPricePerKWHColumn = "default_energy_price_per_kwh"
var LocationsDefaultChargingTimePricePerHourColumn = "default_charging_time_price_per_hour"
var LocationsDefaultParkingTimePricePerHourColumn = "default_parking_time_price_per_hour"
var LocationsAccessColumn = "access"
var LocationsPlugshareScoreColumn = "plugshare_score"
var LocationsDefaultFreePeriodMinutesColumn = "default_free_period_minutes"
var LocationsDefaultMaxChargeTimeMinutesColumn = "default_max_charge_time_minutes"
var LocationsNotesColumn = "notes"
var LocationsLinkedPlugshareIDColumn = "linked_plugshare_id"
var LocationsCalcOverrideKeyColumn = "calc_override_key"
var LocationsOcppLocationIDColumn = "ocpp_location_id"
var LocationsColumns = []string{"id", "country_code", "party_id", "is_published", "publish_only_to", "name", "address", "city", "postalcode", "state", "country", "latitude", "longitude", "parking_type", "timezone", "last_updated", "enabled", "created_at", "updated_at", "deleted_at", "ihomer_id", "charge_station_count", "plugshare_id", "open_247", "is_active", "created_by", "is_public", "plugshare_locations_json_data", "default_energy_price_per_kwh", "default_charging_time_price_per_hour", "default_parking_time_price_per_hour", "access", "plugshare_score", "default_free_period_minutes", "default_max_charge_time_minutes", "notes", "linked_plugshare_id", "calc_override_key", "ocpp_location_id"}

type Locations struct {
	ID                              uuid.UUID  `json:"id" db:"id"`
	CountryCode                     string     `json:"country_code" db:"country_code"`
	PartyID                         string     `json:"party_id" db:"party_id"`
	IsPublished                     bool       `json:"is_published" db:"is_published"`
	PublishOnlyTo                   *[]string  `json:"publish_only_to" db:"publish_only_to"`
	Name                            string     `json:"name" db:"name"`
	Address                         string     `json:"address" db:"address"`
	City                            string     `json:"city" db:"city"`
	Postalcode                      string     `json:"postalcode" db:"postalcode"`
	State                           string     `json:"state" db:"state"`
	Country                         string     `json:"country" db:"country"`
	Latitude                        float64    `json:"latitude" db:"latitude"`
	Longitude                       float64    `json:"longitude" db:"longitude"`
	ParkingType                     string     `json:"parking_type" db:"parking_type"`
	Timezone                        string     `json:"timezone" db:"timezone"`
	LastUpdated                     time.Time  `json:"last_updated" db:"last_updated"`
	Enabled                         bool       `json:"enabled" db:"enabled"`
	CreatedAt                       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt                       time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt                       *time.Time `json:"deleted_at" db:"deleted_at"`
	IhomerID                        string     `json:"ihomer_id" db:"ihomer_id"`
	ChargeStationCount              int64      `json:"charge_station_count" db:"charge_station_count"`
	PlugshareID                     *int64     `json:"plugshare_id" db:"plugshare_id"`
	Open247                         bool       `json:"open_247" db:"open_247"`
	IsActive                        bool       `json:"is_active" db:"is_active"`
	CreatedBy                       *uuid.UUID `json:"created_by" db:"created_by"`
	IsPublic                        bool       `json:"is_public" db:"is_public"`
	PlugshareLocationsJSONData      any        `json:"plugshare_locations_json_data" db:"plugshare_locations_json_data"`
	DefaultEnergyPricePerKWH        float64    `json:"default_energy_price_per_kwh" db:"default_energy_price_per_kwh"`
	DefaultChargingTimePricePerHour float64    `json:"default_charging_time_price_per_hour" db:"default_charging_time_price_per_hour"`
	DefaultParkingTimePricePerHour  float64    `json:"default_parking_time_price_per_hour" db:"default_parking_time_price_per_hour"`
	Access                          string     `json:"access" db:"access"`
	PlugshareScore                  *float64   `json:"plugshare_score" db:"plugshare_score"`
	DefaultFreePeriodMinutes        int64      `json:"default_free_period_minutes" db:"default_free_period_minutes"`
	DefaultMaxChargeTimeMinutes     int64      `json:"default_max_charge_time_minutes" db:"default_max_charge_time_minutes"`
	Notes                           *string    `json:"notes" db:"notes"`
	LinkedPlugshareID               *int64     `json:"linked_plugshare_id" db:"linked_plugshare_id"`
	CalcOverrideKey                 string     `json:"calc_override_key" db:"calc_override_key"`
	OcppLocationID                  *uuid.UUID `json:"ocpp_location_id" db:"ocpp_location_id"`
}

func SelectLocations(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*Locations, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM locations%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*Locations, 0)
	for rows.Next() {
		rowCount++

		var item Locations
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		var temp1PlugshareLocationsJSONData any
		var temp2PlugshareLocationsJSONData []any

		err = json.Unmarshal(item.PlugshareLocationsJSONData.([]byte), &temp1PlugshareLocationsJSONData)
		if err != nil {
			err = json.Unmarshal(item.PlugshareLocationsJSONData.([]byte), &temp2PlugshareLocationsJSONData)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal %#+v as json: %v", string(item.PlugshareLocationsJSONData.([]byte)), err)
			} else {
				item.PlugshareLocationsJSONData = temp2PlugshareLocationsJSONData
			}
		} else {
			item.PlugshareLocationsJSONData = temp1PlugshareLocationsJSONData
		}

		item.PlugshareLocationsJSONData = temp1PlugshareLocationsJSONData

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectLocations(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectLocations(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var NetworksTable = "networks"
var NetworksIDColumn = "id"
var NetworksDescriptionColumn = "description"
var NetworksURLColumn = "url"
var NetworksImageColumn = "image"
var NetworksActionNameColumn = "action_name"
var NetworksPhoneColumn = "phone"
var NetworksE164PhoneNumberColumn = "e164_phone_number"
var NetworksFormattedPhoneNumberColumn = "formatted_phone_number"
var NetworksActionURLColumn = "action_url"
var NetworksNameColumn = "name"
var NetworksColumns = []string{"id", "description", "url", "image", "action_name", "phone", "e164_phone_number", "formatted_phone_number", "action_url", "name"}

type Networks struct {
	ID                   int64  `json:"id" db:"id"`
	Description          string `json:"description" db:"description"`
	URL                  string `json:"url" db:"url"`
	Image                string `json:"image" db:"image"`
	ActionName           string `json:"action_name" db:"action_name"`
	Phone                string `json:"phone" db:"phone"`
	E164PhoneNumber      string `json:"e164_phone_number" db:"e164_phone_number"`
	FormattedPhoneNumber string `json:"formatted_phone_number" db:"formatted_phone_number"`
	ActionURL            string `json:"action_url" db:"action_url"`
	Name                 string `json:"name" db:"name"`
}

func SelectNetworks(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*Networks, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM networks%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*Networks, 0)
	for rows.Next() {
		rowCount++

		var item Networks
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectNetworks(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectNetworks(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var PaymentIntentsTable = "payment_intents"
var PaymentIntentsIDColumn = "id"
var PaymentIntentsUserIDColumn = "user_id"
var PaymentIntentsDescriptionColumn = "description"
var PaymentIntentsStatusColumn = "status"
var PaymentIntentsCreatedAtColumn = "created_at"
var PaymentIntentsUpdatedAtColumn = "updated_at"
var PaymentIntentsCanceledAtColumn = "canceled_at"
var PaymentIntentsDeletedAtColumn = "deleted_at"
var PaymentIntentsUserSessionIDColumn = "user_session_id"
var PaymentIntentsColumns = []string{"id", "user_id", "description", "status", "created_at", "updated_at", "canceled_at", "deleted_at", "user_session_id"}

type PaymentIntents struct {
	ID            string     `json:"id" db:"id"`
	UserID        uuid.UUID  `json:"user_id" db:"user_id"`
	Description   string     `json:"description" db:"description"`
	Status        string     `json:"status" db:"status"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
	CanceledAt    *time.Time `json:"canceled_at" db:"canceled_at"`
	DeletedAt     *time.Time `json:"deleted_at" db:"deleted_at"`
	UserSessionID *uuid.UUID `json:"user_session_id" db:"user_session_id"`
}

func SelectPaymentIntents(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*PaymentIntents, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM payment_intents%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*PaymentIntents, 0)
	for rows.Next() {
		rowCount++

		var item PaymentIntents
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectPaymentIntents(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectPaymentIntents(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var RawCdrsTable = "raw_cdrs"
var RawCdrsIDColumn = "id"
var RawCdrsCountryCodeColumn = "country_code"
var RawCdrsPartyIDColumn = "party_id"
var RawCdrsStartDateTimeColumn = "start_date_time"
var RawCdrsEndDateTimeColumn = "end_date_time"
var RawCdrsContractIDColumn = "contract_id"
var RawCdrsAuthMethodColumn = "auth_method"
var RawCdrsCdrLocationIDColumn = "cdr_location_id"
var RawCdrsCdrLocationNameColumn = "cdr_location_name"
var RawCdrsCdrLocationAddressColumn = "cdr_location_address"
var RawCdrsCdrLocationCityColumn = "cdr_location_city"
var RawCdrsCdrLocationPostalCodeColumn = "cdr_location_postal_code"
var RawCdrsCdrLocationCountryColumn = "cdr_location_country"
var RawCdrsCdrLocationLatitudeColumn = "cdr_location_latitude"
var RawCdrsCdrLocationLongitudeColumn = "cdr_location_longitude"
var RawCdrsCdrLocationEvseUIDColumn = "cdr_location_evse_uid"
var RawCdrsCdrLocationEvseIDColumn = "cdr_location_evse_id"
var RawCdrsCdrLocationConnectorIDColumn = "cdr_location_connector_id"
var RawCdrsCdrLocationConnectorStandardColumn = "cdr_location_connector_standard"
var RawCdrsCdrLocationConnectorFormatColumn = "cdr_location_connector_format"
var RawCdrsCdrLocationConnectorPowerTypeColumn = "cdr_location_connector_power_type"
var RawCdrsChargingPeriodsColumn = "charging_periods"
var RawCdrsTotalCostExclVatColumn = "total_cost_excl_vat"
var RawCdrsTotalEnergyColumn = "total_energy"
var RawCdrsTotalTimeColumn = "total_time"
var RawCdrsLastUpdatedColumn = "last_updated"
var RawCdrsUserIDColumn = "user_id"
var RawCdrsCreatedAtColumn = "created_at"
var RawCdrsColumns = []string{"id", "country_code", "party_id", "start_date_time", "end_date_time", "contract_id", "auth_method", "cdr_location_id", "cdr_location_name", "cdr_location_address", "cdr_location_city", "cdr_location_postal_code", "cdr_location_country", "cdr_location_latitude", "cdr_location_longitude", "cdr_location_evse_uid", "cdr_location_evse_id", "cdr_location_connector_id", "cdr_location_connector_standard", "cdr_location_connector_format", "cdr_location_connector_power_type", "charging_periods", "total_cost_excl_vat", "total_energy", "total_time", "last_updated", "user_id", "created_at"}

type RawCdrs struct {
	ID                            string     `json:"id" db:"id"`
	CountryCode                   *string    `json:"country_code" db:"country_code"`
	PartyID                       *string    `json:"party_id" db:"party_id"`
	StartDateTime                 *time.Time `json:"start_date_time" db:"start_date_time"`
	EndDateTime                   *time.Time `json:"end_date_time" db:"end_date_time"`
	ContractID                    *string    `json:"contract_id" db:"contract_id"`
	AuthMethod                    *string    `json:"auth_method" db:"auth_method"`
	CdrLocationID                 *string    `json:"cdr_location_id" db:"cdr_location_id"`
	CdrLocationName               *string    `json:"cdr_location_name" db:"cdr_location_name"`
	CdrLocationAddress            *string    `json:"cdr_location_address" db:"cdr_location_address"`
	CdrLocationCity               *string    `json:"cdr_location_city" db:"cdr_location_city"`
	CdrLocationPostalCode         *string    `json:"cdr_location_postal_code" db:"cdr_location_postal_code"`
	CdrLocationCountry            *string    `json:"cdr_location_country" db:"cdr_location_country"`
	CdrLocationLatitude           *string    `json:"cdr_location_latitude" db:"cdr_location_latitude"`
	CdrLocationLongitude          *string    `json:"cdr_location_longitude" db:"cdr_location_longitude"`
	CdrLocationEvseUID            *string    `json:"cdr_location_evse_uid" db:"cdr_location_evse_uid"`
	CdrLocationEvseID             *string    `json:"cdr_location_evse_id" db:"cdr_location_evse_id"`
	CdrLocationConnectorID        *string    `json:"cdr_location_connector_id" db:"cdr_location_connector_id"`
	CdrLocationConnectorStandard  *string    `json:"cdr_location_connector_standard" db:"cdr_location_connector_standard"`
	CdrLocationConnectorFormat    *string    `json:"cdr_location_connector_format" db:"cdr_location_connector_format"`
	CdrLocationConnectorPowerType *string    `json:"cdr_location_connector_power_type" db:"cdr_location_connector_power_type"`
	ChargingPeriods               *any       `json:"charging_periods" db:"charging_periods"`
	TotalCostExclVat              *float64   `json:"total_cost_excl_vat" db:"total_cost_excl_vat"`
	TotalEnergy                   *float64   `json:"total_energy" db:"total_energy"`
	TotalTime                     *float64   `json:"total_time" db:"total_time"`
	LastUpdated                   *time.Time `json:"last_updated" db:"last_updated"`
	UserID                        *uuid.UUID `json:"user_id" db:"user_id"`
	CreatedAt                     time.Time  `json:"created_at" db:"created_at"`
}

func SelectRawCdrs(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*RawCdrs, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM raw_cdrs%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*RawCdrs, 0)
	for rows.Next() {
		rowCount++

		var item RawCdrs
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectRawCdrs(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectRawCdrs(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var SchemaMigrationsTable = "schema_migrations"
var SchemaMigrationsVersionColumn = "version"
var SchemaMigrationsDirtyColumn = "dirty"
var SchemaMigrationsColumns = []string{"version", "dirty"}

type SchemaMigrations struct {
	Version int64 `json:"version" db:"version"`
	Dirty   bool  `json:"dirty" db:"dirty"`
}

func SelectSchemaMigrations(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*SchemaMigrations, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM schema_migrations%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*SchemaMigrations, 0)
	for rows.Next() {
		rowCount++

		var item SchemaMigrations
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectSchemaMigrations(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectSchemaMigrations(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var StatsTable = "stats"
var StatsNameColumn = "name"
var StatsValueColumn = "value"
var StatsCreatedAtColumn = "created_at"
var StatsUpdatedAtColumn = "updated_at"
var StatsColumns = []string{"name", "value", "created_at", "updated_at"}

type Stats struct {
	Name      string    `json:"name" db:"name"`
	Value     int64     `json:"value" db:"value"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func SelectStats(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*Stats, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM stats%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*Stats, 0)
	for rows.Next() {
		rowCount++

		var item Stats
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectStats(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectStats(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var StripeInvoicesTable = "stripe_invoices"
var StripeInvoicesIDColumn = "id"
var StripeInvoicesCustomerIDColumn = "customer_id"
var StripeInvoicesSubscriptionTypeIDColumn = "subscription_type_id"
var StripeInvoicesQuantityColumn = "quantity"
var StripeInvoicesCurrencyColumn = "currency"
var StripeInvoicesStatusColumn = "status"
var StripeInvoicesInvoicePdfColumn = "invoice_pdf"
var StripeInvoicesHostedInvoiceURLColumn = "hosted_invoice_url"
var StripeInvoicesCollectionMethodColumn = "collection_method"
var StripeInvoicesInvoiceAmountDueColumn = "invoice_amount_due"
var StripeInvoicesInvoiceAmountPaidColumn = "invoice_amount_paid"
var StripeInvoicesReceiptNumberColumn = "receipt_number"
var StripeInvoicesRawJSONColumn = "raw_json"
var StripeInvoicesCreatedAtColumn = "created_at"
var StripeInvoicesUpdatedAtColumn = "updated_at"
var StripeInvoicesDeletedAtColumn = "deleted_at"
var StripeInvoicesColumns = []string{"id", "customer_id", "subscription_type_id", "quantity", "currency", "status", "invoice_pdf", "hosted_invoice_url", "collection_method", "invoice_amount_due", "invoice_amount_paid", "receipt_number", "raw_json", "created_at", "updated_at", "deleted_at"}

type StripeInvoices struct {
	ID                 string     `json:"id" db:"id"`
	CustomerID         string     `json:"customer_id" db:"customer_id"`
	SubscriptionTypeID uuid.UUID  `json:"subscription_type_id" db:"subscription_type_id"`
	Quantity           *int64     `json:"quantity" db:"quantity"`
	Currency           *string    `json:"currency" db:"currency"`
	Status             *string    `json:"status" db:"status"`
	InvoicePdf         *string    `json:"invoice_pdf" db:"invoice_pdf"`
	HostedInvoiceURL   *string    `json:"hosted_invoice_url" db:"hosted_invoice_url"`
	CollectionMethod   *string    `json:"collection_method" db:"collection_method"`
	InvoiceAmountDue   *int64     `json:"invoice_amount_due" db:"invoice_amount_due"`
	InvoiceAmountPaid  *int64     `json:"invoice_amount_paid" db:"invoice_amount_paid"`
	ReceiptNumber      *string    `json:"receipt_number" db:"receipt_number"`
	RawJSON            *any       `json:"raw_json" db:"raw_json"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at" db:"deleted_at"`
}

func SelectStripeInvoices(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*StripeInvoices, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM stripe_invoices%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*StripeInvoices, 0)
	for rows.Next() {
		rowCount++

		var item StripeInvoices
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectStripeInvoices(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectStripeInvoices(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var StripeSubscriptionsTable = "stripe_subscriptions"
var StripeSubscriptionsIDColumn = "id"
var StripeSubscriptionsCustomerIDColumn = "customer_id"
var StripeSubscriptionsSubscriptionTypeIDColumn = "subscription_type_id"
var StripeSubscriptionsBillingCycleAnchorColumn = "billing_cycle_anchor"
var StripeSubscriptionsQuantityColumn = "quantity"
var StripeSubscriptionsCurrencyColumn = "currency"
var StripeSubscriptionsStatusColumn = "status"
var StripeSubscriptionsInvoicePdfColumn = "invoice_pdf"
var StripeSubscriptionsHostedInvoiceURLColumn = "hosted_invoice_url"
var StripeSubscriptionsCanceledAtColumn = "canceled_at"
var StripeSubscriptionsRawJSONColumn = "raw_json"
var StripeSubscriptionsCreatedAtColumn = "created_at"
var StripeSubscriptionsUpdatedAtColumn = "updated_at"
var StripeSubscriptionsDeletedAtColumn = "deleted_at"
var StripeSubscriptionsCollectionMethodColumn = "collection_method"
var StripeSubscriptionsDefaultPaymentMethodColumn = "default_payment_method"
var StripeSubscriptionsLatestInvoiceIDColumn = "latest_invoice_id"
var StripeSubscriptionsIntervalColumn = "interval"
var StripeSubscriptionsInvoiceAmountPaidColumn = "invoice_amount_paid"
var StripeSubscriptionsColumns = []string{"id", "customer_id", "subscription_type_id", "billing_cycle_anchor", "quantity", "currency", "status", "invoice_pdf", "hosted_invoice_url", "canceled_at", "raw_json", "created_at", "updated_at", "deleted_at", "collection_method", "default_payment_method", "latest_invoice_id", "interval", "invoice_amount_paid"}

type StripeSubscriptions struct {
	ID                   string     `json:"id" db:"id"`
	CustomerID           string     `json:"customer_id" db:"customer_id"`
	SubscriptionTypeID   uuid.UUID  `json:"subscription_type_id" db:"subscription_type_id"`
	BillingCycleAnchor   *int64     `json:"billing_cycle_anchor" db:"billing_cycle_anchor"`
	Quantity             *int64     `json:"quantity" db:"quantity"`
	Currency             *string    `json:"currency" db:"currency"`
	Status               *string    `json:"status" db:"status"`
	InvoicePdf           *string    `json:"invoice_pdf" db:"invoice_pdf"`
	HostedInvoiceURL     *string    `json:"hosted_invoice_url" db:"hosted_invoice_url"`
	CanceledAt           *time.Time `json:"canceled_at" db:"canceled_at"`
	RawJSON              *any       `json:"raw_json" db:"raw_json"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt            *time.Time `json:"deleted_at" db:"deleted_at"`
	CollectionMethod     *string    `json:"collection_method" db:"collection_method"`
	DefaultPaymentMethod *string    `json:"default_payment_method" db:"default_payment_method"`
	LatestInvoiceID      *string    `json:"latest_invoice_id" db:"latest_invoice_id"`
	Interval             *string    `json:"interval" db:"interval"`
	InvoiceAmountPaid    *int64     `json:"invoice_amount_paid" db:"invoice_amount_paid"`
}

func SelectStripeSubscriptions(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*StripeSubscriptions, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM stripe_subscriptions%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*StripeSubscriptions, 0)
	for rows.Next() {
		rowCount++

		var item StripeSubscriptions
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectStripeSubscriptions(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectStripeSubscriptions(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var SubscriptionTypesTable = "subscription_types"
var SubscriptionTypesIDColumn = "id"
var SubscriptionTypesTitleColumn = "title"
var SubscriptionTypesDescriptionColumn = "description"
var SubscriptionTypesPublicChargingStationsColumn = "public_charging_stations"
var SubscriptionTypesTransactionFeeThresholdCentsColumn = "transaction_fee_threshold_cents"
var SubscriptionTypesTransactionFeeBelowThresholdCentsColumn = "transaction_fee_below_threshold_cents"
var SubscriptionTypesTransactionFeeAboveThresholdPercentColumn = "transaction_fee_above_threshold_percent"
var SubscriptionTypesMaxChargingStationsColumn = "max_charging_stations"
var SubscriptionTypesCreatedAtColumn = "created_at"
var SubscriptionTypesUpdatedAtColumn = "updated_at"
var SubscriptionTypesDeletedAtColumn = "deleted_at"
var SubscriptionTypesStripeProductIDColumn = "stripe_product_id"
var SubscriptionTypesPriceIsPerColumn = "price_is_per"
var SubscriptionTypesMonthlyPriceCentsColumn = "monthly_price_cents"
var SubscriptionTypesYearlyPriceCentsColumn = "yearly_price_cents"
var SubscriptionTypesTwoYearlyPriceCentsColumn = "two_yearly_price_cents"
var SubscriptionTypesMaxChargingStationsAcColumn = "max_charging_stations_ac"
var SubscriptionTypesMaxChargingStationsDcColumn = "max_charging_stations_dc"
var SubscriptionTypesMaxLocationsColumn = "max_locations"
var SubscriptionTypesStripeMonthlyPriceIDColumn = "stripe_monthly_price_id"
var SubscriptionTypesStripeYearlyPriceIDColumn = "stripe_yearly_price_id"
var SubscriptionTypesStripeTwoYearlyPriceIDColumn = "stripe_two_yearly_price_id"
var SubscriptionTypesActiveColumn = "active"
var SubscriptionTypesVisibleColumn = "visible"
var SubscriptionTypesColumns = []string{"id", "title", "description", "public_charging_stations", "transaction_fee_threshold_cents", "transaction_fee_below_threshold_cents", "transaction_fee_above_threshold_percent", "max_charging_stations", "created_at", "updated_at", "deleted_at", "stripe_product_id", "price_is_per", "monthly_price_cents", "yearly_price_cents", "two_yearly_price_cents", "max_charging_stations_ac", "max_charging_stations_dc", "max_locations", "stripe_monthly_price_id", "stripe_yearly_price_id", "stripe_two_yearly_price_id", "active", "visible"}

type SubscriptionTypes struct {
	ID                                  uuid.UUID  `json:"id" db:"id"`
	Title                               string     `json:"title" db:"title"`
	Description                         string     `json:"description" db:"description"`
	PublicChargingStations              bool       `json:"public_charging_stations" db:"public_charging_stations"`
	TransactionFeeThresholdCents        int64      `json:"transaction_fee_threshold_cents" db:"transaction_fee_threshold_cents"`
	TransactionFeeBelowThresholdCents   int64      `json:"transaction_fee_below_threshold_cents" db:"transaction_fee_below_threshold_cents"`
	TransactionFeeAboveThresholdPercent int64      `json:"transaction_fee_above_threshold_percent" db:"transaction_fee_above_threshold_percent"`
	MaxChargingStations                 *int64     `json:"max_charging_stations" db:"max_charging_stations"`
	CreatedAt                           time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt                           time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt                           *time.Time `json:"deleted_at" db:"deleted_at"`
	StripeProductID                     *string    `json:"stripe_product_id" db:"stripe_product_id"`
	PriceIsPer                          *string    `json:"price_is_per" db:"price_is_per"`
	MonthlyPriceCents                   int64      `json:"monthly_price_cents" db:"monthly_price_cents"`
	YearlyPriceCents                    int64      `json:"yearly_price_cents" db:"yearly_price_cents"`
	TwoYearlyPriceCents                 *int64     `json:"two_yearly_price_cents" db:"two_yearly_price_cents"`
	MaxChargingStationsAc               *int64     `json:"max_charging_stations_ac" db:"max_charging_stations_ac"`
	MaxChargingStationsDc               *int64     `json:"max_charging_stations_dc" db:"max_charging_stations_dc"`
	MaxLocations                        *int64     `json:"max_locations" db:"max_locations"`
	StripeMonthlyPriceID                *string    `json:"stripe_monthly_price_id" db:"stripe_monthly_price_id"`
	StripeYearlyPriceID                 *string    `json:"stripe_yearly_price_id" db:"stripe_yearly_price_id"`
	StripeTwoYearlyPriceID              *string    `json:"stripe_two_yearly_price_id" db:"stripe_two_yearly_price_id"`
	Active                              bool       `json:"active" db:"active"`
	Visible                             bool       `json:"visible" db:"visible"`
}

func SelectSubscriptionTypes(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*SubscriptionTypes, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM subscription_types%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*SubscriptionTypes, 0)
	for rows.Next() {
		rowCount++

		var item SubscriptionTypes
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectSubscriptionTypes(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectSubscriptionTypes(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var TokenGroupsTable = "token_groups"
var TokenGroupsTokenIDColumn = "token_id"
var TokenGroupsGroupIDColumn = "group_id"
var TokenGroupsColumns = []string{"token_id", "group_id"}

type TokenGroups struct {
	TokenID uuid.UUID `json:"token_id" db:"token_id"`
	GroupID uuid.UUID `json:"group_id" db:"group_id"`
}

func SelectTokenGroups(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*TokenGroups, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM token_groups%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*TokenGroups, 0)
	for rows.Next() {
		rowCount++

		var item TokenGroups
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectTokenGroups(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectTokenGroups(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var TokenReadonlyUsersTable = "token_readonly_users"
var TokenReadonlyUsersTokenIDColumn = "token_id"
var TokenReadonlyUsersReadonlyIDColumn = "readonly_id"
var TokenReadonlyUsersColumns = []string{"token_id", "readonly_id"}

type TokenReadonlyUsers struct {
	TokenID    uuid.UUID `json:"token_id" db:"token_id"`
	ReadonlyID uuid.UUID `json:"readonly_id" db:"readonly_id"`
}

func SelectTokenReadonlyUsers(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*TokenReadonlyUsers, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM token_readonly_users%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*TokenReadonlyUsers, 0)
	for rows.Next() {
		rowCount++

		var item TokenReadonlyUsers
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectTokenReadonlyUsers(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectTokenReadonlyUsers(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var TokensTable = "tokens"
var TokensIDColumn = "id"
var TokensCountryCodeColumn = "country_code"
var TokensPartyIDColumn = "party_id"
var TokensUIDColumn = "uid"
var TokensTypeColumn = "type"
var TokensContractIDColumn = "contract_id"
var TokensContractSupplierNameColumn = "contract_supplier_name"
var TokensVisualNumberColumn = "visual_number"
var TokensIssuerColumn = "issuer"
var TokensGroupIDColumn = "group_id"
var TokensValidColumn = "valid"
var TokensWhitelistColumn = "whitelist"
var TokensLanguageColumn = "language"
var TokensDefaultProfileTypeColumn = "default_profile_type"
var TokensCreatedByColumn = "created_by"
var TokensCreatedAtColumn = "created_at"
var TokensUpdatedAtColumn = "updated_at"
var TokensDeletedAtColumn = "deleted_at"
var TokensAssignedToColumn = "assigned_to"
var TokensIsInitActivatedColumn = "is_init_activated"
var TokensIhomerIDColumn = "ihomer_id"
var TokensDefaultPaymentMethodIDColumn = "default_payment_method_id"
var TokensIsLocalTokenColumn = "is_local_token"
var TokensOcppServerIDColumn = "ocpp_server_id"
var TokensColumns = []string{"id", "country_code", "party_id", "uid", "type", "contract_id", "contract_supplier_name", "visual_number", "issuer", "group_id", "valid", "whitelist", "language", "default_profile_type", "created_by", "created_at", "updated_at", "deleted_at", "assigned_to", "is_init_activated", "ihomer_id", "default_payment_method_id", "is_local_token", "ocpp_server_id"}

type Tokens struct {
	ID                     uuid.UUID  `json:"id" db:"id"`
	CountryCode            string     `json:"country_code" db:"country_code"`
	PartyID                *string    `json:"party_id" db:"party_id"`
	UID                    string     `json:"uid" db:"uid"`
	Type                   string     `json:"type" db:"type"`
	ContractID             *string    `json:"contract_id" db:"contract_id"`
	ContractSupplierName   *string    `json:"contract_supplier_name" db:"contract_supplier_name"`
	VisualNumber           string     `json:"visual_number" db:"visual_number"`
	Issuer                 string     `json:"issuer" db:"issuer"`
	GroupID                *string    `json:"group_id" db:"group_id"`
	Valid                  bool       `json:"valid" db:"valid"`
	Whitelist              string     `json:"whitelist" db:"whitelist"`
	Language               string     `json:"language" db:"language"`
	DefaultProfileType     string     `json:"default_profile_type" db:"default_profile_type"`
	CreatedBy              *uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt              time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt              *time.Time `json:"deleted_at" db:"deleted_at"`
	AssignedTo             *uuid.UUID `json:"assigned_to" db:"assigned_to"`
	IsInitActivated        bool       `json:"is_init_activated" db:"is_init_activated"`
	IhomerID               string     `json:"ihomer_id" db:"ihomer_id"`
	DefaultPaymentMethodID *string    `json:"default_payment_method_id" db:"default_payment_method_id"`
	IsLocalToken           bool       `json:"is_local_token" db:"is_local_token"`
	OcppServerID           *uuid.UUID `json:"ocpp_server_id" db:"ocpp_server_id"`
}

func SelectTokens(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*Tokens, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM tokens%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*Tokens, 0)
	for rows.Next() {
		rowCount++

		var item Tokens
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectTokens(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectTokens(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var UserEvsTable = "user_evs"
var UserEvsIDColumn = "id"
var UserEvsEvIDColumn = "ev_id"
var UserEvsUserIDColumn = "user_id"
var UserEvsManufacturerColumn = "manufacturer"
var UserEvsModelColumn = "model"
var UserEvsYearColumn = "year"
var UserEvsLicensePlateNumberColumn = "license_plate_number"
var UserEvsBatteryCapacityColumn = "battery_capacity"
var UserEvsPlugTypeColumn = "plug_type"
var UserEvsChargingTypeColumn = "charging_type"
var UserEvsImageColumn = "image"
var UserEvsNicknameColumn = "nickname"
var UserEvsIsPrimaryColumn = "is_primary"
var UserEvsCreatedAtColumn = "created_at"
var UserEvsUpdatedAtColumn = "updated_at"
var UserEvsDeletedAtColumn = "deleted_at"
var UserEvsRangeColumn = "range"
var UserEvsColumns = []string{"id", "ev_id", "user_id", "manufacturer", "model", "year", "license_plate_number", "battery_capacity", "plug_type", "charging_type", "image", "nickname", "is_primary", "created_at", "updated_at", "deleted_at", "range"}

type UserEvs struct {
	ID                 uuid.UUID  `json:"id" db:"id"`
	EvID               int64      `json:"ev_id" db:"ev_id"`
	UserID             *uuid.UUID `json:"user_id" db:"user_id"`
	Manufacturer       string     `json:"manufacturer" db:"manufacturer"`
	Model              string     `json:"model" db:"model"`
	Year               *int64     `json:"year" db:"year"`
	LicensePlateNumber *string    `json:"license_plate_number" db:"license_plate_number"`
	BatteryCapacity    *float64   `json:"battery_capacity" db:"battery_capacity"`
	PlugType           *string    `json:"plug_type" db:"plug_type"`
	ChargingType       *string    `json:"charging_type" db:"charging_type"`
	Image              *string    `json:"image" db:"image"`
	Nickname           *string    `json:"nickname" db:"nickname"`
	IsPrimary          bool       `json:"is_primary" db:"is_primary"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at" db:"deleted_at"`
	Range              *float64   `json:"range" db:"range"`
}

func SelectUserEvs(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*UserEvs, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM user_evs%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*UserEvs, 0)
	for rows.Next() {
		rowCount++

		var item UserEvs
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectUserEvs(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectUserEvs(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var UserGroupsTable = "user_groups"
var UserGroupsUserIDColumn = "user_id"
var UserGroupsGroupIDColumn = "group_id"
var UserGroupsColumns = []string{"user_id", "group_id"}

type UserGroups struct {
	UserID  uuid.UUID `json:"user_id" db:"user_id"`
	GroupID uuid.UUID `json:"group_id" db:"group_id"`
}

func SelectUserGroups(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*UserGroups, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM user_groups%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*UserGroups, 0)
	for rows.Next() {
		rowCount++

		var item UserGroups
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectUserGroups(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectUserGroups(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var UserReadonlyTokensTable = "user_readonly_tokens"
var UserReadonlyTokensUserIDColumn = "user_id"
var UserReadonlyTokensReadonlyIDColumn = "readonly_id"
var UserReadonlyTokensColumns = []string{"user_id", "readonly_id"}

type UserReadonlyTokens struct {
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	ReadonlyID uuid.UUID `json:"readonly_id" db:"readonly_id"`
}

func SelectUserReadonlyTokens(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*UserReadonlyTokens, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM user_readonly_tokens%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*UserReadonlyTokens, 0)
	for rows.Next() {
		rowCount++

		var item UserReadonlyTokens
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectUserReadonlyTokens(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectUserReadonlyTokens(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var UserRolesTable = "user_roles"
var UserRolesIDColumn = "id"
var UserRolesNameColumn = "name"
var UserRolesCreatedAtColumn = "created_at"
var UserRolesUpdatedAtColumn = "updated_at"
var UserRolesDeletedAtColumn = "deleted_at"
var UserRolesColumns = []string{"id", "name", "created_at", "updated_at", "deleted_at"}

type UserRoles struct {
	ID        int64      `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

func SelectUserRoles(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*UserRoles, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM user_roles%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*UserRoles, 0)
	for rows.Next() {
		rowCount++

		var item UserRoles
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectUserRoles(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectUserRoles(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var UserSessionsTable = "user_sessions"
var UserSessionsIDColumn = "id"
var UserSessionsUserIDColumn = "user_id"
var UserSessionsStartedAtColumn = "started_at"
var UserSessionsEndedAtColumn = "ended_at"
var UserSessionsCreatedAtColumn = "created_at"
var UserSessionsUpdatedAtColumn = "updated_at"
var UserSessionsDeletedAtColumn = "deleted_at"
var UserSessionsMspTransactionIDColumn = "msp_transaction_id"
var UserSessionsKWHColumn = "kwh"
var UserSessionsEvseIDColumn = "evse_id"
var UserSessionsMaxChargeMinutesColumn = "max_charge_minutes"
var UserSessionsMaxChargePercentColumn = "max_charge_percent"
var UserSessionsMaxChargePriceCentsColumn = "max_charge_price_cents"
var UserSessionsChargingStationIDColumn = "charging_station_id"
var UserSessionsIhomerTransactionIDColumn = "ihomer_transaction_id"
var UserSessionsTokenIDColumn = "token_id"
var UserSessionsRawCdrIDColumn = "raw_cdr_id"
var UserSessionsTotalCostCentsColumn = "total_cost_cents"
var UserSessionsUserEvIDColumn = "user_ev_id"
var UserSessionsStripePaymentIntentIDColumn = "stripe_payment_intent_id"
var UserSessionsPaidColumn = "paid"
var UserSessionsStripePaymentMethodIDColumn = "stripe_payment_method_id"
var UserSessionsStripeReceiptURLColumn = "stripe_receipt_url"
var UserSessionsStripeChargeIDColumn = "stripe_charge_id"
var UserSessionsEnergyPricePerKWHColumn = "energy_price_per_kwh"
var UserSessionsChargingTimePricePerHourColumn = "charging_time_price_per_hour"
var UserSessionsParkingTimePricePerHourColumn = "parking_time_price_per_hour"
var UserSessionsTransactionFeeThresholdCentsColumn = "transaction_fee_threshold_cents"
var UserSessionsTransactionFeeBelowThresholdCentsColumn = "transaction_fee_below_threshold_cents"
var UserSessionsTransactionFeeAboveThresholdPercentColumn = "transaction_fee_above_threshold_percent"
var UserSessionsStripeReceiptNumberColumn = "stripe_receipt_number"
var UserSessionsTotalCPOSplitCentsColumn = "total_cpo_split_cents"
var UserSessionsTotalPlatformFeeCentsColumn = "total_platform_fee_cents"
var UserSessionsTotalStripeFeeCentsColumn = "total_stripe_fee_cents"
var UserSessionsTotalTaxCentsColumn = "total_tax_cents"
var UserSessionsTotalCPONetCentsColumn = "total_cpo_net_cents"
var UserSessionsCardBrandColumn = "card_brand"
var UserSessionsCardLast4Column = "card_last4"
var UserSessionsChargerDisconnectedAtColumn = "charger_disconnected_at"
var UserSessionsPayableParkingChargeCentsColumn = "payable_parking_charge_cents"
var UserSessionsStripeTaxCentsColumn = "stripe_tax_cents"
var UserSessionsStatusColumn = "status"
var UserSessionsTotalAmountCentsColumn = "total_amount_cents"
var UserSessionsFreePeriodMinutesColumn = "free_period_minutes"
var UserSessionsIsLocalColumn = "is_local"
var UserSessionsTimeSuspendedEvColumn = "time_suspended_ev"
var UserSessionsTimeSuspendedEvseColumn = "time_suspended_evse"
var UserSessionsReasonColumn = "reason"
var UserSessionsHistoryColumn = "history"
var UserSessionsSessionTypeColumn = "session_type"
var UserSessionsChargingStationNeedsUnlockFirstColumn = "charging_station_needs_unlock_first"
var UserSessionsIsV2SessionColumn = "is_v2_session"
var UserSessionsOcppServerIDColumn = "ocpp_server_id"
var UserSessionsOcppServerTransactionIDColumn = "ocpp_server_transaction_id"
var UserSessionsLimitSessionSetupByUserColumn = "limit_session_setup_by_user"
var UserSessionsCalculatedStripeFeesCentsColumn = "calculated_stripe_fees_cents"
var UserSessionsTotalAmountMinusStripeFeesCentsColumn = "total_amount_minus_stripe_fees_cents"
var UserSessionsHasPrecalculatedStripeFeesColumn = "has_precalculated_stripe_fees"
var UserSessionsColumns = []string{"id", "user_id", "started_at", "ended_at", "created_at", "updated_at", "deleted_at", "msp_transaction_id", "kwh", "evse_id", "max_charge_minutes", "max_charge_percent", "max_charge_price_cents", "charging_station_id", "ihomer_transaction_id", "token_id", "raw_cdr_id", "total_cost_cents", "user_ev_id", "stripe_payment_intent_id", "paid", "stripe_payment_method_id", "stripe_receipt_url", "stripe_charge_id", "energy_price_per_kwh", "charging_time_price_per_hour", "parking_time_price_per_hour", "transaction_fee_threshold_cents", "transaction_fee_below_threshold_cents", "transaction_fee_above_threshold_percent", "stripe_receipt_number", "total_cpo_split_cents", "total_platform_fee_cents", "total_stripe_fee_cents", "total_tax_cents", "total_cpo_net_cents", "card_brand", "card_last4", "charger_disconnected_at", "payable_parking_charge_cents", "stripe_tax_cents", "status", "total_amount_cents", "free_period_minutes", "is_local", "time_suspended_ev", "time_suspended_evse", "reason", "history", "session_type", "charging_station_needs_unlock_first", "is_v2_session", "ocpp_server_id", "ocpp_server_transaction_id", "limit_session_setup_by_user", "calculated_stripe_fees_cents", "total_amount_minus_stripe_fees_cents", "has_precalculated_stripe_fees"}

type UserSessions struct {
	ID                                  uuid.UUID  `json:"id" db:"id"`
	UserID                              *uuid.UUID `json:"user_id" db:"user_id"`
	StartedAt                           time.Time  `json:"started_at" db:"started_at"`
	EndedAt                             *time.Time `json:"ended_at" db:"ended_at"`
	CreatedAt                           time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt                           time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt                           *time.Time `json:"deleted_at" db:"deleted_at"`
	MspTransactionID                    *string    `json:"msp_transaction_id" db:"msp_transaction_id"`
	KWH                                 *float64   `json:"kwh" db:"kwh"`
	EvseID                              int64      `json:"evse_id" db:"evse_id"`
	MaxChargeMinutes                    *int64     `json:"max_charge_minutes" db:"max_charge_minutes"`
	MaxChargePercent                    *int64     `json:"max_charge_percent" db:"max_charge_percent"`
	MaxChargePriceCents                 *int64     `json:"max_charge_price_cents" db:"max_charge_price_cents"`
	ChargingStationID                   string     `json:"charging_station_id" db:"charging_station_id"`
	IhomerTransactionID                 *string    `json:"ihomer_transaction_id" db:"ihomer_transaction_id"`
	TokenID                             string     `json:"token_id" db:"token_id"`
	RawCdrID                            *string    `json:"raw_cdr_id" db:"raw_cdr_id"`
	TotalCostCents                      *int64     `json:"total_cost_cents" db:"total_cost_cents"`
	UserEvID                            *uuid.UUID `json:"user_ev_id" db:"user_ev_id"`
	StripePaymentIntentID               *string    `json:"stripe_payment_intent_id" db:"stripe_payment_intent_id"`
	Paid                                bool       `json:"paid" db:"paid"`
	StripePaymentMethodID               *string    `json:"stripe_payment_method_id" db:"stripe_payment_method_id"`
	StripeReceiptURL                    *string    `json:"stripe_receipt_url" db:"stripe_receipt_url"`
	StripeChargeID                      *string    `json:"stripe_charge_id" db:"stripe_charge_id"`
	EnergyPricePerKWH                   float64    `json:"energy_price_per_kwh" db:"energy_price_per_kwh"`
	ChargingTimePricePerHour            float64    `json:"charging_time_price_per_hour" db:"charging_time_price_per_hour"`
	ParkingTimePricePerHour             float64    `json:"parking_time_price_per_hour" db:"parking_time_price_per_hour"`
	TransactionFeeThresholdCents        int64      `json:"transaction_fee_threshold_cents" db:"transaction_fee_threshold_cents"`
	TransactionFeeBelowThresholdCents   int64      `json:"transaction_fee_below_threshold_cents" db:"transaction_fee_below_threshold_cents"`
	TransactionFeeAboveThresholdPercent int64      `json:"transaction_fee_above_threshold_percent" db:"transaction_fee_above_threshold_percent"`
	StripeReceiptNumber                 *string    `json:"stripe_receipt_number" db:"stripe_receipt_number"`
	TotalCPOSplitCents                  *int64     `json:"total_cpo_split_cents" db:"total_cpo_split_cents"`
	TotalPlatformFeeCents               *int64     `json:"total_platform_fee_cents" db:"total_platform_fee_cents"`
	TotalStripeFeeCents                 *int64     `json:"total_stripe_fee_cents" db:"total_stripe_fee_cents"`
	TotalTaxCents                       *int64     `json:"total_tax_cents" db:"total_tax_cents"`
	TotalCPONetCents                    *int64     `json:"total_cpo_net_cents" db:"total_cpo_net_cents"`
	CardBrand                           string     `json:"card_brand" db:"card_brand"`
	CardLast4                           string     `json:"card_last4" db:"card_last4"`
	ChargerDisconnectedAt               *time.Time `json:"charger_disconnected_at" db:"charger_disconnected_at"`
	PayableParkingChargeCents           int64      `json:"payable_parking_charge_cents" db:"payable_parking_charge_cents"`
	StripeTaxCents                      int64      `json:"stripe_tax_cents" db:"stripe_tax_cents"`
	Status                              string     `json:"status" db:"status"`
	TotalAmountCents                    int64      `json:"total_amount_cents" db:"total_amount_cents"`
	FreePeriodMinutes                   int64      `json:"free_period_minutes" db:"free_period_minutes"`
	IsLocal                             bool       `json:"is_local" db:"is_local"`
	TimeSuspendedEv                     *int64     `json:"time_suspended_ev" db:"time_suspended_ev"`
	TimeSuspendedEvse                   *int64     `json:"time_suspended_evse" db:"time_suspended_evse"`
	Reason                              string     `json:"reason" db:"reason"`
	History                             any        `json:"history" db:"history"`
	SessionType                         string     `json:"session_type" db:"session_type"`
	ChargingStationNeedsUnlockFirst     bool       `json:"charging_station_needs_unlock_first" db:"charging_station_needs_unlock_first"`
	IsV2Session                         *bool      `json:"is_v2_session" db:"is_v2_session"`
	OcppServerID                        *uuid.UUID `json:"ocpp_server_id" db:"ocpp_server_id"`
	OcppServerTransactionID             *int64     `json:"ocpp_server_transaction_id" db:"ocpp_server_transaction_id"`
	LimitSessionSetupByUser             *bool      `json:"limit_session_setup_by_user" db:"limit_session_setup_by_user"`
	CalculatedStripeFeesCents           *int64     `json:"calculated_stripe_fees_cents" db:"calculated_stripe_fees_cents"`
	TotalAmountMinusStripeFeesCents     *int64     `json:"total_amount_minus_stripe_fees_cents" db:"total_amount_minus_stripe_fees_cents"`
	HasPrecalculatedStripeFees          *bool      `json:"has_precalculated_stripe_fees" db:"has_precalculated_stripe_fees"`
}

func SelectUserSessions(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*UserSessions, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM user_sessions%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*UserSessions, 0)
	for rows.Next() {
		rowCount++

		var item UserSessions
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		var temp1History any
		var temp2History []any

		err = json.Unmarshal(item.History.([]byte), &temp1History)
		if err != nil {
			err = json.Unmarshal(item.History.([]byte), &temp2History)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal %#+v as json: %v", string(item.History.([]byte)), err)
			} else {
				item.History = temp2History
			}
		} else {
			item.History = temp1History
		}

		item.History = temp1History

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectUserSessions(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectUserSessions(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var UserSubscriptionsTable = "user_subscriptions"
var UserSubscriptionsIDColumn = "id"
var UserSubscriptionsUserIDColumn = "user_id"
var UserSubscriptionsSubscriptionTypeIDColumn = "subscription_type_id"
var UserSubscriptionsActiveColumn = "active"
var UserSubscriptionsStartedAtColumn = "started_at"
var UserSubscriptionsPaidUntilColumn = "paid_until"
var UserSubscriptionsCreatedAtColumn = "created_at"
var UserSubscriptionsUpdatedAtColumn = "updated_at"
var UserSubscriptionsDeletedAtColumn = "deleted_at"
var UserSubscriptionsMaxChargingStationsColumn = "max_charging_stations"
var UserSubscriptionsMaxLocationsColumn = "max_locations"
var UserSubscriptionsStripeSubscriptionIDColumn = "stripe_subscription_id"
var UserSubscriptionsSubscriptionExpiredHandledColumn = "subscription_expired_handled"
var UserSubscriptionsPublicChargingStationsColumn = "public_charging_stations"
var UserSubscriptionsTransactionFeeThresholdCentsColumn = "transaction_fee_threshold_cents"
var UserSubscriptionsTransactionFeeBelowThresholdCentsColumn = "transaction_fee_below_threshold_cents"
var UserSubscriptionsTransactionFeeAboveThresholdPercentColumn = "transaction_fee_above_threshold_percent"
var UserSubscriptionsDiscountIDColumn = "discount_id"
var UserSubscriptionsMaxChargingStationsAcColumn = "max_charging_stations_ac"
var UserSubscriptionsMaxChargingStationsDcColumn = "max_charging_stations_dc"
var UserSubscriptionsPendingColumn = "pending"
var UserSubscriptionsStripeInvoiceIDColumn = "stripe_invoice_id"
var UserSubscriptionsColumns = []string{"id", "user_id", "subscription_type_id", "active", "started_at", "paid_until", "created_at", "updated_at", "deleted_at", "max_charging_stations", "max_locations", "stripe_subscription_id", "subscription_expired_handled", "public_charging_stations", "transaction_fee_threshold_cents", "transaction_fee_below_threshold_cents", "transaction_fee_above_threshold_percent", "discount_id", "max_charging_stations_ac", "max_charging_stations_dc", "pending", "stripe_invoice_id"}

type UserSubscriptions struct {
	ID                                  uuid.UUID  `json:"id" db:"id"`
	UserID                              uuid.UUID  `json:"user_id" db:"user_id"`
	SubscriptionTypeID                  uuid.UUID  `json:"subscription_type_id" db:"subscription_type_id"`
	Active                              bool       `json:"active" db:"active"`
	StartedAt                           time.Time  `json:"started_at" db:"started_at"`
	PaidUntil                           *time.Time `json:"paid_until" db:"paid_until"`
	CreatedAt                           time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt                           time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt                           *time.Time `json:"deleted_at" db:"deleted_at"`
	MaxChargingStations                 *int64     `json:"max_charging_stations" db:"max_charging_stations"`
	MaxLocations                        *int64     `json:"max_locations" db:"max_locations"`
	StripeSubscriptionID                *string    `json:"stripe_subscription_id" db:"stripe_subscription_id"`
	SubscriptionExpiredHandled          bool       `json:"subscription_expired_handled" db:"subscription_expired_handled"`
	PublicChargingStations              bool       `json:"public_charging_stations" db:"public_charging_stations"`
	TransactionFeeThresholdCents        int64      `json:"transaction_fee_threshold_cents" db:"transaction_fee_threshold_cents"`
	TransactionFeeBelowThresholdCents   int64      `json:"transaction_fee_below_threshold_cents" db:"transaction_fee_below_threshold_cents"`
	TransactionFeeAboveThresholdPercent int64      `json:"transaction_fee_above_threshold_percent" db:"transaction_fee_above_threshold_percent"`
	DiscountID                          *uuid.UUID `json:"discount_id" db:"discount_id"`
	MaxChargingStationsAc               *int64     `json:"max_charging_stations_ac" db:"max_charging_stations_ac"`
	MaxChargingStationsDc               *int64     `json:"max_charging_stations_dc" db:"max_charging_stations_dc"`
	Pending                             bool       `json:"pending" db:"pending"`
	StripeInvoiceID                     *string    `json:"stripe_invoice_id" db:"stripe_invoice_id"`
}

func SelectUserSubscriptions(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*UserSubscriptions, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM user_subscriptions%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*UserSubscriptions, 0)
	for rows.Next() {
		rowCount++

		var item UserSubscriptions
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectUserSubscriptions(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectUserSubscriptions(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var UsersTable = "users"
var UsersIDColumn = "id"
var UsersEmailColumn = "email"
var UsersEnabledColumn = "enabled"
var UsersDefaultTokenIDColumn = "default_token_id"
var UsersCreatedAtColumn = "created_at"
var UsersUpdatedAtColumn = "updated_at"
var UsersDeletedAtColumn = "deleted_at"
var UsersActiveUserSubscriptionIDColumn = "active_user_subscription_id"
var UsersNameColumn = "name"
var UsersDOBColumn = "dob"
var UsersNationalityColumn = "nationality"
var UsersCountryColumn = "country"
var UsersAddressColumn = "address"
var UsersCompanyNameColumn = "company_name"
var UsersCompanyWebsiteColumn = "company_website"
var UsersCompanyCountryColumn = "company_country"
var UsersCompanyAddressColumn = "company_address"
var UsersCompanyLicenseNumberColumn = "company_license_number"
var UsersCPOTypeColumn = "cpo_type"
var UsersCPOReasonColumn = "cpo_reason"
var UsersStripeCustomerIDColumn = "stripe_customer_id"
var UsersFirstNameColumn = "first_name"
var UsersLastNameColumn = "last_name"
var UsersPhoneColumn = "phone"
var UsersCountryDialCodeColumn = "country_dial_code"
var UsersImageColumn = "image"
var UsersStripeAccountIDColumn = "stripe_account_id"
var UsersStripeAccountTypeColumn = "stripe_account_type"
var UsersRoleIDColumn = "role_id"
var UsersNotificationTokenColumn = "notification_token"
var UsersEmailValidateCodeColumn = "email_validate_code"
var UsersCustomerIDColumn = "customer_id"
var UsersDiscountIDColumn = "discount_id"
var UsersIsEmailVerifiedColumn = "is_email_verified"
var UsersManagerUserIDColumn = "manager_user_id"
var UsersLastLoggedInAtColumn = "last_logged_in_at"
var UsersVoltguardEnabledColumn = "voltguard_enabled"
var UsersVoltguardNotificationsColumn = "voltguard_notifications"
var UsersVoltguardTriggerEventsColumn = "voltguard_trigger_events"
var UsersColumns = []string{"id", "email", "enabled", "default_token_id", "created_at", "updated_at", "deleted_at", "active_user_subscription_id", "name", "dob", "nationality", "country", "address", "company_name", "company_website", "company_country", "company_address", "company_license_number", "cpo_type", "cpo_reason", "stripe_customer_id", "first_name", "last_name", "phone", "country_dial_code", "image", "stripe_account_id", "stripe_account_type", "role_id", "notification_token", "email_validate_code", "customer_id", "discount_id", "is_email_verified", "manager_user_id", "last_logged_in_at", "voltguard_enabled", "voltguard_notifications", "voltguard_trigger_events"}

type Users struct {
	ID                       uuid.UUID  `json:"id" db:"id"`
	Email                    *string    `json:"email" db:"email"`
	Enabled                  bool       `json:"enabled" db:"enabled"`
	DefaultTokenID           string     `json:"default_token_id" db:"default_token_id"`
	CreatedAt                time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt                *time.Time `json:"deleted_at" db:"deleted_at"`
	ActiveUserSubscriptionID *uuid.UUID `json:"active_user_subscription_id" db:"active_user_subscription_id"`
	Name                     *string    `json:"name" db:"name"`
	DOB                      *time.Time `json:"dob" db:"dob"`
	Nationality              *string    `json:"nationality" db:"nationality"`
	Country                  *string    `json:"country" db:"country"`
	Address                  *string    `json:"address" db:"address"`
	CompanyName              *string    `json:"company_name" db:"company_name"`
	CompanyWebsite           *string    `json:"company_website" db:"company_website"`
	CompanyCountry           *string    `json:"company_country" db:"company_country"`
	CompanyAddress           *string    `json:"company_address" db:"company_address"`
	CompanyLicenseNumber     *string    `json:"company_license_number" db:"company_license_number"`
	CPOType                  *string    `json:"cpo_type" db:"cpo_type"`
	CPOReason                *string    `json:"cpo_reason" db:"cpo_reason"`
	StripeCustomerID         *string    `json:"stripe_customer_id" db:"stripe_customer_id"`
	FirstName                *string    `json:"first_name" db:"first_name"`
	LastName                 *string    `json:"last_name" db:"last_name"`
	Phone                    *string    `json:"phone" db:"phone"`
	CountryDialCode          *string    `json:"country_dial_code" db:"country_dial_code"`
	Image                    *string    `json:"image" db:"image"`
	StripeAccountID          *string    `json:"stripe_account_id" db:"stripe_account_id"`
	StripeAccountType        *string    `json:"stripe_account_type" db:"stripe_account_type"`
	RoleID                   int64      `json:"role_id" db:"role_id"`
	NotificationToken        *string    `json:"notification_token" db:"notification_token"`
	EmailValidateCode        *string    `json:"email_validate_code" db:"email_validate_code"`
	CustomerID               int64      `json:"customer_id" db:"customer_id"`
	DiscountID               *uuid.UUID `json:"discount_id" db:"discount_id"`
	IsEmailVerified          bool       `json:"is_email_verified" db:"is_email_verified"`
	ManagerUserID            *uuid.UUID `json:"manager_user_id" db:"manager_user_id"`
	LastLoggedInAt           *time.Time `json:"last_logged_in_at" db:"last_logged_in_at"`
	VoltguardEnabled         bool       `json:"voltguard_enabled" db:"voltguard_enabled"`
	VoltguardNotifications   any        `json:"voltguard_notifications" db:"voltguard_notifications"`
	VoltguardTriggerEvents   any        `json:"voltguard_trigger_events" db:"voltguard_trigger_events"`
}

func SelectUsers(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*Users, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM users%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*Users, 0)
	for rows.Next() {
		rowCount++

		var item Users
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		var temp1VoltguardNotifications any
		var temp2VoltguardNotifications []any

		err = json.Unmarshal(item.VoltguardNotifications.([]byte), &temp1VoltguardNotifications)
		if err != nil {
			err = json.Unmarshal(item.VoltguardNotifications.([]byte), &temp2VoltguardNotifications)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal %#+v as json: %v", string(item.VoltguardNotifications.([]byte)), err)
			} else {
				item.VoltguardNotifications = temp2VoltguardNotifications
			}
		} else {
			item.VoltguardNotifications = temp1VoltguardNotifications
		}

		item.VoltguardNotifications = temp1VoltguardNotifications

		var temp1VoltguardTriggerEvents any
		var temp2VoltguardTriggerEvents []any

		err = json.Unmarshal(item.VoltguardTriggerEvents.([]byte), &temp1VoltguardTriggerEvents)
		if err != nil {
			err = json.Unmarshal(item.VoltguardTriggerEvents.([]byte), &temp2VoltguardTriggerEvents)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal %#+v as json: %v", string(item.VoltguardTriggerEvents.([]byte)), err)
			} else {
				item.VoltguardTriggerEvents = temp2VoltguardTriggerEvents
			}
		} else {
			item.VoltguardTriggerEvents = temp1VoltguardTriggerEvents
		}

		item.VoltguardTriggerEvents = temp1VoltguardTriggerEvents

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectUsers(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectUsers(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var VChargingStationsPrivateTable = "v_charging_stations_private"
var VChargingStationsPrivateIDColumn = "id"
var VChargingStationsPrivateLocationIDColumn = "location_id"
var VChargingStationsPrivateUserIDColumn = "user_id"
var VChargingStationsPrivateProtocolColumn = "protocol"
var VChargingStationsPrivateAcceptedColumn = "accepted"
var VChargingStationsPrivateReservableColumn = "reservable"
var VChargingStationsPrivateConfiguredColumn = "configured"
var VChargingStationsPrivateAccessibilityColumn = "accessibility"
var VChargingStationsPrivateAvailabilityColumn = "availability"
var VChargingStationsPrivateStatusColumn = "status"
var VChargingStationsPrivateSubStatusColumn = "sub_status"
var VChargingStationsPrivateLocalAuthorizationListVersionColumn = "local_authorization_list_version"
var VChargingStationsPrivateIsActiveColumn = "is_active"
var VChargingStationsPrivateIsPublicColumn = "is_public"
var VChargingStationsPrivateCreatedAtColumn = "created_at"
var VChargingStationsPrivateUpdatedAtColumn = "updated_at"
var VChargingStationsPrivateDeletedAtColumn = "deleted_at"
var VChargingStationsPrivateAliasColumn = "alias"
var VChargingStationsPrivatePlugshareJSONDataColumn = "plugshare_json_data"
var VChargingStationsPrivateIsConfiguredColumn = "is_configured"
var VChargingStationsPrivatePlugshareCostTypeColumn = "plugshare_cost_type"
var VChargingStationsPrivateNetworkIDColumn = "network_id"
var VChargingStationsPrivateCredentialsGeneratedColumn = "credentials_generated"
var VChargingStationsPrivateRanGenerateCredentialsColumn = "ran_generate_credentials"
var VChargingStationsPrivateNeedsUnlockFirstColumn = "needs_unlock_first"
var VChargingStationsPrivateColumns = []string{"id", "location_id", "user_id", "protocol", "accepted", "reservable", "configured", "accessibility", "availability", "status", "sub_status", "local_authorization_list_version", "is_active", "is_public", "created_at", "updated_at", "deleted_at", "alias", "plugshare_json_data", "is_configured", "plugshare_cost_type", "network_id", "credentials_generated", "ran_generate_credentials", "needs_unlock_first"}

type VChargingStationsPrivate struct {
	ID                            *string    `json:"id" db:"id"`
	LocationID                    *uuid.UUID `json:"location_id" db:"location_id"`
	UserID                        *uuid.UUID `json:"user_id" db:"user_id"`
	Protocol                      *string    `json:"protocol" db:"protocol"`
	Accepted                      *bool      `json:"accepted" db:"accepted"`
	Reservable                    *bool      `json:"reservable" db:"reservable"`
	Configured                    *bool      `json:"configured" db:"configured"`
	Accessibility                 *string    `json:"accessibility" db:"accessibility"`
	Availability                  *string    `json:"availability" db:"availability"`
	Status                        *string    `json:"status" db:"status"`
	SubStatus                     *string    `json:"sub_status" db:"sub_status"`
	LocalAuthorizationListVersion *int64     `json:"local_authorization_list_version" db:"local_authorization_list_version"`
	IsActive                      *bool      `json:"is_active" db:"is_active"`
	IsPublic                      *bool      `json:"is_public" db:"is_public"`
	CreatedAt                     *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                     *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt                     *time.Time `json:"deleted_at" db:"deleted_at"`
	Alias                         *string    `json:"alias" db:"alias"`
	PlugshareJSONData             *any       `json:"plugshare_json_data" db:"plugshare_json_data"`
	IsConfigured                  *bool      `json:"is_configured" db:"is_configured"`
	PlugshareCostType             *int64     `json:"plugshare_cost_type" db:"plugshare_cost_type"`
	NetworkID                     *int64     `json:"network_id" db:"network_id"`
	CredentialsGenerated          *bool      `json:"credentials_generated" db:"credentials_generated"`
	RanGenerateCredentials        *bool      `json:"ran_generate_credentials" db:"ran_generate_credentials"`
	NeedsUnlockFirst              *bool      `json:"needs_unlock_first" db:"needs_unlock_first"`
}

func SelectVChargingStationsPrivate(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*VChargingStationsPrivate, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM v_charging_stations_private%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*VChargingStationsPrivate, 0)
	for rows.Next() {
		rowCount++

		var item VChargingStationsPrivate
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectVChargingStationsPrivate(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectVChargingStationsPrivate(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var VChargingStationsPrivateOverTimeTable = "v_charging_stations_private_over_time"
var VChargingStationsPrivateOverTimeTimeColumn = "time"
var VChargingStationsPrivateOverTimeCountColumn = "count"
var VChargingStationsPrivateOverTimeColumns = []string{"time", "count"}

type VChargingStationsPrivateOverTime struct {
	Time  *time.Time `json:"time" db:"time"`
	Count *float64   `json:"count" db:"count"`
}

func SelectVChargingStationsPrivateOverTime(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*VChargingStationsPrivateOverTime, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM v_charging_stations_private_over_time%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*VChargingStationsPrivateOverTime, 0)
	for rows.Next() {
		rowCount++

		var item VChargingStationsPrivateOverTime
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectVChargingStationsPrivateOverTime(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectVChargingStationsPrivateOverTime(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var VChargingStationsPublicTable = "v_charging_stations_public"
var VChargingStationsPublicIDColumn = "id"
var VChargingStationsPublicLocationIDColumn = "location_id"
var VChargingStationsPublicUserIDColumn = "user_id"
var VChargingStationsPublicProtocolColumn = "protocol"
var VChargingStationsPublicAcceptedColumn = "accepted"
var VChargingStationsPublicReservableColumn = "reservable"
var VChargingStationsPublicConfiguredColumn = "configured"
var VChargingStationsPublicAccessibilityColumn = "accessibility"
var VChargingStationsPublicAvailabilityColumn = "availability"
var VChargingStationsPublicStatusColumn = "status"
var VChargingStationsPublicSubStatusColumn = "sub_status"
var VChargingStationsPublicLocalAuthorizationListVersionColumn = "local_authorization_list_version"
var VChargingStationsPublicIsActiveColumn = "is_active"
var VChargingStationsPublicIsPublicColumn = "is_public"
var VChargingStationsPublicCreatedAtColumn = "created_at"
var VChargingStationsPublicUpdatedAtColumn = "updated_at"
var VChargingStationsPublicDeletedAtColumn = "deleted_at"
var VChargingStationsPublicAliasColumn = "alias"
var VChargingStationsPublicPlugshareJSONDataColumn = "plugshare_json_data"
var VChargingStationsPublicIsConfiguredColumn = "is_configured"
var VChargingStationsPublicPlugshareCostTypeColumn = "plugshare_cost_type"
var VChargingStationsPublicNetworkIDColumn = "network_id"
var VChargingStationsPublicCredentialsGeneratedColumn = "credentials_generated"
var VChargingStationsPublicRanGenerateCredentialsColumn = "ran_generate_credentials"
var VChargingStationsPublicNeedsUnlockFirstColumn = "needs_unlock_first"
var VChargingStationsPublicColumns = []string{"id", "location_id", "user_id", "protocol", "accepted", "reservable", "configured", "accessibility", "availability", "status", "sub_status", "local_authorization_list_version", "is_active", "is_public", "created_at", "updated_at", "deleted_at", "alias", "plugshare_json_data", "is_configured", "plugshare_cost_type", "network_id", "credentials_generated", "ran_generate_credentials", "needs_unlock_first"}

type VChargingStationsPublic struct {
	ID                            *string    `json:"id" db:"id"`
	LocationID                    *uuid.UUID `json:"location_id" db:"location_id"`
	UserID                        *uuid.UUID `json:"user_id" db:"user_id"`
	Protocol                      *string    `json:"protocol" db:"protocol"`
	Accepted                      *bool      `json:"accepted" db:"accepted"`
	Reservable                    *bool      `json:"reservable" db:"reservable"`
	Configured                    *bool      `json:"configured" db:"configured"`
	Accessibility                 *string    `json:"accessibility" db:"accessibility"`
	Availability                  *string    `json:"availability" db:"availability"`
	Status                        *string    `json:"status" db:"status"`
	SubStatus                     *string    `json:"sub_status" db:"sub_status"`
	LocalAuthorizationListVersion *int64     `json:"local_authorization_list_version" db:"local_authorization_list_version"`
	IsActive                      *bool      `json:"is_active" db:"is_active"`
	IsPublic                      *bool      `json:"is_public" db:"is_public"`
	CreatedAt                     *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                     *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt                     *time.Time `json:"deleted_at" db:"deleted_at"`
	Alias                         *string    `json:"alias" db:"alias"`
	PlugshareJSONData             *any       `json:"plugshare_json_data" db:"plugshare_json_data"`
	IsConfigured                  *bool      `json:"is_configured" db:"is_configured"`
	PlugshareCostType             *int64     `json:"plugshare_cost_type" db:"plugshare_cost_type"`
	NetworkID                     *int64     `json:"network_id" db:"network_id"`
	CredentialsGenerated          *bool      `json:"credentials_generated" db:"credentials_generated"`
	RanGenerateCredentials        *bool      `json:"ran_generate_credentials" db:"ran_generate_credentials"`
	NeedsUnlockFirst              *bool      `json:"needs_unlock_first" db:"needs_unlock_first"`
}

func SelectVChargingStationsPublic(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*VChargingStationsPublic, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM v_charging_stations_public%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*VChargingStationsPublic, 0)
	for rows.Next() {
		rowCount++

		var item VChargingStationsPublic
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectVChargingStationsPublic(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectVChargingStationsPublic(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var VChargingStationsPublicOverTimeTable = "v_charging_stations_public_over_time"
var VChargingStationsPublicOverTimeTimeColumn = "time"
var VChargingStationsPublicOverTimeCountColumn = "count"
var VChargingStationsPublicOverTimeColumns = []string{"time", "count"}

type VChargingStationsPublicOverTime struct {
	Time  *time.Time `json:"time" db:"time"`
	Count *float64   `json:"count" db:"count"`
}

func SelectVChargingStationsPublicOverTime(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*VChargingStationsPublicOverTime, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM v_charging_stations_public_over_time%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*VChargingStationsPublicOverTime, 0)
	for rows.Next() {
		rowCount++

		var item VChargingStationsPublicOverTime
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectVChargingStationsPublicOverTime(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectVChargingStationsPublicOverTime(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var VCposTable = "v_cpos"
var VCposIDColumn = "id"
var VCposEmailColumn = "email"
var VCposEnabledColumn = "enabled"
var VCposDefaultTokenIDColumn = "default_token_id"
var VCposCreatedAtColumn = "created_at"
var VCposUpdatedAtColumn = "updated_at"
var VCposDeletedAtColumn = "deleted_at"
var VCposActiveUserSubscriptionIDColumn = "active_user_subscription_id"
var VCposNameColumn = "name"
var VCposDOBColumn = "dob"
var VCposNationalityColumn = "nationality"
var VCposCountryColumn = "country"
var VCposAddressColumn = "address"
var VCposCompanyNameColumn = "company_name"
var VCposCompanyWebsiteColumn = "company_website"
var VCposCompanyCountryColumn = "company_country"
var VCposCompanyAddressColumn = "company_address"
var VCposCompanyLicenseNumberColumn = "company_license_number"
var VCposCPOTypeColumn = "cpo_type"
var VCposCPOReasonColumn = "cpo_reason"
var VCposStripeCustomerIDColumn = "stripe_customer_id"
var VCposFirstNameColumn = "first_name"
var VCposLastNameColumn = "last_name"
var VCposPhoneColumn = "phone"
var VCposCountryDialCodeColumn = "country_dial_code"
var VCposImageColumn = "image"
var VCposStripeAccountIDColumn = "stripe_account_id"
var VCposStripeAccountTypeColumn = "stripe_account_type"
var VCposRoleIDColumn = "role_id"
var VCposNotificationTokenColumn = "notification_token"
var VCposEmailValidateCodeColumn = "email_validate_code"
var VCposCustomerIDColumn = "customer_id"
var VCposDiscountIDColumn = "discount_id"
var VCposIsEmailVerifiedColumn = "is_email_verified"
var VCposManagerUserIDColumn = "manager_user_id"
var VCposColumns = []string{"id", "email", "enabled", "default_token_id", "created_at", "updated_at", "deleted_at", "active_user_subscription_id", "name", "dob", "nationality", "country", "address", "company_name", "company_website", "company_country", "company_address", "company_license_number", "cpo_type", "cpo_reason", "stripe_customer_id", "first_name", "last_name", "phone", "country_dial_code", "image", "stripe_account_id", "stripe_account_type", "role_id", "notification_token", "email_validate_code", "customer_id", "discount_id", "is_email_verified", "manager_user_id"}

type VCpos struct {
	ID                       *uuid.UUID `json:"id" db:"id"`
	Email                    *string    `json:"email" db:"email"`
	Enabled                  *bool      `json:"enabled" db:"enabled"`
	DefaultTokenID           *string    `json:"default_token_id" db:"default_token_id"`
	CreatedAt                *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt                *time.Time `json:"deleted_at" db:"deleted_at"`
	ActiveUserSubscriptionID *uuid.UUID `json:"active_user_subscription_id" db:"active_user_subscription_id"`
	Name                     *string    `json:"name" db:"name"`
	DOB                      *time.Time `json:"dob" db:"dob"`
	Nationality              *string    `json:"nationality" db:"nationality"`
	Country                  *string    `json:"country" db:"country"`
	Address                  *string    `json:"address" db:"address"`
	CompanyName              *string    `json:"company_name" db:"company_name"`
	CompanyWebsite           *string    `json:"company_website" db:"company_website"`
	CompanyCountry           *string    `json:"company_country" db:"company_country"`
	CompanyAddress           *string    `json:"company_address" db:"company_address"`
	CompanyLicenseNumber     *string    `json:"company_license_number" db:"company_license_number"`
	CPOType                  *string    `json:"cpo_type" db:"cpo_type"`
	CPOReason                *string    `json:"cpo_reason" db:"cpo_reason"`
	StripeCustomerID         *string    `json:"stripe_customer_id" db:"stripe_customer_id"`
	FirstName                *string    `json:"first_name" db:"first_name"`
	LastName                 *string    `json:"last_name" db:"last_name"`
	Phone                    *string    `json:"phone" db:"phone"`
	CountryDialCode          *string    `json:"country_dial_code" db:"country_dial_code"`
	Image                    *string    `json:"image" db:"image"`
	StripeAccountID          *string    `json:"stripe_account_id" db:"stripe_account_id"`
	StripeAccountType        *string    `json:"stripe_account_type" db:"stripe_account_type"`
	RoleID                   *int64     `json:"role_id" db:"role_id"`
	NotificationToken        *string    `json:"notification_token" db:"notification_token"`
	EmailValidateCode        *string    `json:"email_validate_code" db:"email_validate_code"`
	CustomerID               *int64     `json:"customer_id" db:"customer_id"`
	DiscountID               *uuid.UUID `json:"discount_id" db:"discount_id"`
	IsEmailVerified          *bool      `json:"is_email_verified" db:"is_email_verified"`
	ManagerUserID            *uuid.UUID `json:"manager_user_id" db:"manager_user_id"`
}

func SelectVCpos(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*VCpos, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM v_cpos%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*VCpos, 0)
	for rows.Next() {
		rowCount++

		var item VCpos
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectVCpos(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectVCpos(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var VLocationsPrivateTable = "v_locations_private"
var VLocationsPrivateIDColumn = "id"
var VLocationsPrivateCountryCodeColumn = "country_code"
var VLocationsPrivatePartyIDColumn = "party_id"
var VLocationsPrivateIsPublishedColumn = "is_published"
var VLocationsPrivatePublishOnlyToColumn = "publish_only_to"
var VLocationsPrivateNameColumn = "name"
var VLocationsPrivateAddressColumn = "address"
var VLocationsPrivateCityColumn = "city"
var VLocationsPrivatePostalcodeColumn = "postalcode"
var VLocationsPrivateStateColumn = "state"
var VLocationsPrivateCountryColumn = "country"
var VLocationsPrivateLatitudeColumn = "latitude"
var VLocationsPrivateLongitudeColumn = "longitude"
var VLocationsPrivateParkingTypeColumn = "parking_type"
var VLocationsPrivateTimezoneColumn = "timezone"
var VLocationsPrivateLastUpdatedColumn = "last_updated"
var VLocationsPrivateEnabledColumn = "enabled"
var VLocationsPrivateCreatedAtColumn = "created_at"
var VLocationsPrivateUpdatedAtColumn = "updated_at"
var VLocationsPrivateDeletedAtColumn = "deleted_at"
var VLocationsPrivateIhomerIDColumn = "ihomer_id"
var VLocationsPrivateChargeStationCountColumn = "charge_station_count"
var VLocationsPrivatePlugshareIDColumn = "plugshare_id"
var VLocationsPrivateOpen247Column = "open_247"
var VLocationsPrivateIsActiveColumn = "is_active"
var VLocationsPrivateCreatedByColumn = "created_by"
var VLocationsPrivateIsPublicColumn = "is_public"
var VLocationsPrivatePlugshareLocationsJSONDataColumn = "plugshare_locations_json_data"
var VLocationsPrivateDefaultEnergyPricePerKWHColumn = "default_energy_price_per_kwh"
var VLocationsPrivateDefaultChargingTimePricePerHourColumn = "default_charging_time_price_per_hour"
var VLocationsPrivateDefaultParkingTimePricePerHourColumn = "default_parking_time_price_per_hour"
var VLocationsPrivateAccessColumn = "access"
var VLocationsPrivatePlugshareScoreColumn = "plugshare_score"
var VLocationsPrivateDefaultFreePeriodMinutesColumn = "default_free_period_minutes"
var VLocationsPrivateDefaultMaxChargeTimeMinutesColumn = "default_max_charge_time_minutes"
var VLocationsPrivateNotesColumn = "notes"
var VLocationsPrivateLinkedPlugshareIDColumn = "linked_plugshare_id"
var VLocationsPrivateCalcOverrideKeyColumn = "calc_override_key"
var VLocationsPrivateColumns = []string{"id", "country_code", "party_id", "is_published", "publish_only_to", "name", "address", "city", "postalcode", "state", "country", "latitude", "longitude", "parking_type", "timezone", "last_updated", "enabled", "created_at", "updated_at", "deleted_at", "ihomer_id", "charge_station_count", "plugshare_id", "open_247", "is_active", "created_by", "is_public", "plugshare_locations_json_data", "default_energy_price_per_kwh", "default_charging_time_price_per_hour", "default_parking_time_price_per_hour", "access", "plugshare_score", "default_free_period_minutes", "default_max_charge_time_minutes", "notes", "linked_plugshare_id", "calc_override_key"}

type VLocationsPrivate struct {
	ID                              *uuid.UUID `json:"id" db:"id"`
	CountryCode                     *string    `json:"country_code" db:"country_code"`
	PartyID                         *string    `json:"party_id" db:"party_id"`
	IsPublished                     *bool      `json:"is_published" db:"is_published"`
	PublishOnlyTo                   *[]string  `json:"publish_only_to" db:"publish_only_to"`
	Name                            *string    `json:"name" db:"name"`
	Address                         *string    `json:"address" db:"address"`
	City                            *string    `json:"city" db:"city"`
	Postalcode                      *string    `json:"postalcode" db:"postalcode"`
	State                           *string    `json:"state" db:"state"`
	Country                         *string    `json:"country" db:"country"`
	Latitude                        *float64   `json:"latitude" db:"latitude"`
	Longitude                       *float64   `json:"longitude" db:"longitude"`
	ParkingType                     *string    `json:"parking_type" db:"parking_type"`
	Timezone                        *string    `json:"timezone" db:"timezone"`
	LastUpdated                     *time.Time `json:"last_updated" db:"last_updated"`
	Enabled                         *bool      `json:"enabled" db:"enabled"`
	CreatedAt                       *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                       *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt                       *time.Time `json:"deleted_at" db:"deleted_at"`
	IhomerID                        *string    `json:"ihomer_id" db:"ihomer_id"`
	ChargeStationCount              *int64     `json:"charge_station_count" db:"charge_station_count"`
	PlugshareID                     *int64     `json:"plugshare_id" db:"plugshare_id"`
	Open247                         *bool      `json:"open_247" db:"open_247"`
	IsActive                        *bool      `json:"is_active" db:"is_active"`
	CreatedBy                       *uuid.UUID `json:"created_by" db:"created_by"`
	IsPublic                        *bool      `json:"is_public" db:"is_public"`
	PlugshareLocationsJSONData      *any       `json:"plugshare_locations_json_data" db:"plugshare_locations_json_data"`
	DefaultEnergyPricePerKWH        *float64   `json:"default_energy_price_per_kwh" db:"default_energy_price_per_kwh"`
	DefaultChargingTimePricePerHour *float64   `json:"default_charging_time_price_per_hour" db:"default_charging_time_price_per_hour"`
	DefaultParkingTimePricePerHour  *float64   `json:"default_parking_time_price_per_hour" db:"default_parking_time_price_per_hour"`
	Access                          *string    `json:"access" db:"access"`
	PlugshareScore                  *float64   `json:"plugshare_score" db:"plugshare_score"`
	DefaultFreePeriodMinutes        *int64     `json:"default_free_period_minutes" db:"default_free_period_minutes"`
	DefaultMaxChargeTimeMinutes     *int64     `json:"default_max_charge_time_minutes" db:"default_max_charge_time_minutes"`
	Notes                           *string    `json:"notes" db:"notes"`
	LinkedPlugshareID               *int64     `json:"linked_plugshare_id" db:"linked_plugshare_id"`
	CalcOverrideKey                 *string    `json:"calc_override_key" db:"calc_override_key"`
}

func SelectVLocationsPrivate(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*VLocationsPrivate, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM v_locations_private%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*VLocationsPrivate, 0)
	for rows.Next() {
		rowCount++

		var item VLocationsPrivate
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectVLocationsPrivate(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectVLocationsPrivate(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var VLocationsPublicTable = "v_locations_public"
var VLocationsPublicIDColumn = "id"
var VLocationsPublicCountryCodeColumn = "country_code"
var VLocationsPublicPartyIDColumn = "party_id"
var VLocationsPublicIsPublishedColumn = "is_published"
var VLocationsPublicPublishOnlyToColumn = "publish_only_to"
var VLocationsPublicNameColumn = "name"
var VLocationsPublicAddressColumn = "address"
var VLocationsPublicCityColumn = "city"
var VLocationsPublicPostalcodeColumn = "postalcode"
var VLocationsPublicStateColumn = "state"
var VLocationsPublicCountryColumn = "country"
var VLocationsPublicLatitudeColumn = "latitude"
var VLocationsPublicLongitudeColumn = "longitude"
var VLocationsPublicParkingTypeColumn = "parking_type"
var VLocationsPublicTimezoneColumn = "timezone"
var VLocationsPublicLastUpdatedColumn = "last_updated"
var VLocationsPublicEnabledColumn = "enabled"
var VLocationsPublicCreatedAtColumn = "created_at"
var VLocationsPublicUpdatedAtColumn = "updated_at"
var VLocationsPublicDeletedAtColumn = "deleted_at"
var VLocationsPublicIhomerIDColumn = "ihomer_id"
var VLocationsPublicChargeStationCountColumn = "charge_station_count"
var VLocationsPublicPlugshareIDColumn = "plugshare_id"
var VLocationsPublicOpen247Column = "open_247"
var VLocationsPublicIsActiveColumn = "is_active"
var VLocationsPublicCreatedByColumn = "created_by"
var VLocationsPublicIsPublicColumn = "is_public"
var VLocationsPublicPlugshareLocationsJSONDataColumn = "plugshare_locations_json_data"
var VLocationsPublicDefaultEnergyPricePerKWHColumn = "default_energy_price_per_kwh"
var VLocationsPublicDefaultChargingTimePricePerHourColumn = "default_charging_time_price_per_hour"
var VLocationsPublicDefaultParkingTimePricePerHourColumn = "default_parking_time_price_per_hour"
var VLocationsPublicAccessColumn = "access"
var VLocationsPublicPlugshareScoreColumn = "plugshare_score"
var VLocationsPublicDefaultFreePeriodMinutesColumn = "default_free_period_minutes"
var VLocationsPublicDefaultMaxChargeTimeMinutesColumn = "default_max_charge_time_minutes"
var VLocationsPublicNotesColumn = "notes"
var VLocationsPublicLinkedPlugshareIDColumn = "linked_plugshare_id"
var VLocationsPublicCalcOverrideKeyColumn = "calc_override_key"
var VLocationsPublicColumns = []string{"id", "country_code", "party_id", "is_published", "publish_only_to", "name", "address", "city", "postalcode", "state", "country", "latitude", "longitude", "parking_type", "timezone", "last_updated", "enabled", "created_at", "updated_at", "deleted_at", "ihomer_id", "charge_station_count", "plugshare_id", "open_247", "is_active", "created_by", "is_public", "plugshare_locations_json_data", "default_energy_price_per_kwh", "default_charging_time_price_per_hour", "default_parking_time_price_per_hour", "access", "plugshare_score", "default_free_period_minutes", "default_max_charge_time_minutes", "notes", "linked_plugshare_id", "calc_override_key"}

type VLocationsPublic struct {
	ID                              *uuid.UUID `json:"id" db:"id"`
	CountryCode                     *string    `json:"country_code" db:"country_code"`
	PartyID                         *string    `json:"party_id" db:"party_id"`
	IsPublished                     *bool      `json:"is_published" db:"is_published"`
	PublishOnlyTo                   *[]string  `json:"publish_only_to" db:"publish_only_to"`
	Name                            *string    `json:"name" db:"name"`
	Address                         *string    `json:"address" db:"address"`
	City                            *string    `json:"city" db:"city"`
	Postalcode                      *string    `json:"postalcode" db:"postalcode"`
	State                           *string    `json:"state" db:"state"`
	Country                         *string    `json:"country" db:"country"`
	Latitude                        *float64   `json:"latitude" db:"latitude"`
	Longitude                       *float64   `json:"longitude" db:"longitude"`
	ParkingType                     *string    `json:"parking_type" db:"parking_type"`
	Timezone                        *string    `json:"timezone" db:"timezone"`
	LastUpdated                     *time.Time `json:"last_updated" db:"last_updated"`
	Enabled                         *bool      `json:"enabled" db:"enabled"`
	CreatedAt                       *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                       *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt                       *time.Time `json:"deleted_at" db:"deleted_at"`
	IhomerID                        *string    `json:"ihomer_id" db:"ihomer_id"`
	ChargeStationCount              *int64     `json:"charge_station_count" db:"charge_station_count"`
	PlugshareID                     *int64     `json:"plugshare_id" db:"plugshare_id"`
	Open247                         *bool      `json:"open_247" db:"open_247"`
	IsActive                        *bool      `json:"is_active" db:"is_active"`
	CreatedBy                       *uuid.UUID `json:"created_by" db:"created_by"`
	IsPublic                        *bool      `json:"is_public" db:"is_public"`
	PlugshareLocationsJSONData      *any       `json:"plugshare_locations_json_data" db:"plugshare_locations_json_data"`
	DefaultEnergyPricePerKWH        *float64   `json:"default_energy_price_per_kwh" db:"default_energy_price_per_kwh"`
	DefaultChargingTimePricePerHour *float64   `json:"default_charging_time_price_per_hour" db:"default_charging_time_price_per_hour"`
	DefaultParkingTimePricePerHour  *float64   `json:"default_parking_time_price_per_hour" db:"default_parking_time_price_per_hour"`
	Access                          *string    `json:"access" db:"access"`
	PlugshareScore                  *float64   `json:"plugshare_score" db:"plugshare_score"`
	DefaultFreePeriodMinutes        *int64     `json:"default_free_period_minutes" db:"default_free_period_minutes"`
	DefaultMaxChargeTimeMinutes     *int64     `json:"default_max_charge_time_minutes" db:"default_max_charge_time_minutes"`
	Notes                           *string    `json:"notes" db:"notes"`
	LinkedPlugshareID               *int64     `json:"linked_plugshare_id" db:"linked_plugshare_id"`
	CalcOverrideKey                 *string    `json:"calc_override_key" db:"calc_override_key"`
}

func SelectVLocationsPublic(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*VLocationsPublic, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM v_locations_public%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*VLocationsPublic, 0)
	for rows.Next() {
		rowCount++

		var item VLocationsPublic
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectVLocationsPublic(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectVLocationsPublic(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var VUserSessionsPrivateTable = "v_user_sessions_private"
var VUserSessionsPrivateIDColumn = "id"
var VUserSessionsPrivateUserIDColumn = "user_id"
var VUserSessionsPrivateStartedAtColumn = "started_at"
var VUserSessionsPrivateEndedAtColumn = "ended_at"
var VUserSessionsPrivateCreatedAtColumn = "created_at"
var VUserSessionsPrivateUpdatedAtColumn = "updated_at"
var VUserSessionsPrivateDeletedAtColumn = "deleted_at"
var VUserSessionsPrivateMspTransactionIDColumn = "msp_transaction_id"
var VUserSessionsPrivateKWHColumn = "kwh"
var VUserSessionsPrivateEvseIDColumn = "evse_id"
var VUserSessionsPrivateMaxChargeMinutesColumn = "max_charge_minutes"
var VUserSessionsPrivateMaxChargePercentColumn = "max_charge_percent"
var VUserSessionsPrivateMaxChargePriceCentsColumn = "max_charge_price_cents"
var VUserSessionsPrivateChargingStationIDColumn = "charging_station_id"
var VUserSessionsPrivateIhomerTransactionIDColumn = "ihomer_transaction_id"
var VUserSessionsPrivateTokenIDColumn = "token_id"
var VUserSessionsPrivateRawCdrIDColumn = "raw_cdr_id"
var VUserSessionsPrivateTotalCostCentsColumn = "total_cost_cents"
var VUserSessionsPrivateUserEvIDColumn = "user_ev_id"
var VUserSessionsPrivateStripePaymentIntentIDColumn = "stripe_payment_intent_id"
var VUserSessionsPrivatePaidColumn = "paid"
var VUserSessionsPrivateStripePaymentMethodIDColumn = "stripe_payment_method_id"
var VUserSessionsPrivateStripeReceiptURLColumn = "stripe_receipt_url"
var VUserSessionsPrivateStripeChargeIDColumn = "stripe_charge_id"
var VUserSessionsPrivateEnergyPricePerKWHColumn = "energy_price_per_kwh"
var VUserSessionsPrivateChargingTimePricePerHourColumn = "charging_time_price_per_hour"
var VUserSessionsPrivateParkingTimePricePerHourColumn = "parking_time_price_per_hour"
var VUserSessionsPrivateTransactionFeeThresholdCentsColumn = "transaction_fee_threshold_cents"
var VUserSessionsPrivateTransactionFeeBelowThresholdCentsColumn = "transaction_fee_below_threshold_cents"
var VUserSessionsPrivateTransactionFeeAboveThresholdPercentColumn = "transaction_fee_above_threshold_percent"
var VUserSessionsPrivateStripeReceiptNumberColumn = "stripe_receipt_number"
var VUserSessionsPrivateTotalCPOSplitCentsColumn = "total_cpo_split_cents"
var VUserSessionsPrivateTotalPlatformFeeCentsColumn = "total_platform_fee_cents"
var VUserSessionsPrivateTotalStripeFeeCentsColumn = "total_stripe_fee_cents"
var VUserSessionsPrivateTotalTaxCentsColumn = "total_tax_cents"
var VUserSessionsPrivateTotalCPONetCentsColumn = "total_cpo_net_cents"
var VUserSessionsPrivateCardBrandColumn = "card_brand"
var VUserSessionsPrivateCardLast4Column = "card_last4"
var VUserSessionsPrivateChargerDisconnectedAtColumn = "charger_disconnected_at"
var VUserSessionsPrivatePayableParkingChargeCentsColumn = "payable_parking_charge_cents"
var VUserSessionsPrivateStripeTaxCentsColumn = "stripe_tax_cents"
var VUserSessionsPrivateStatusColumn = "status"
var VUserSessionsPrivateTotalAmountCentsColumn = "total_amount_cents"
var VUserSessionsPrivateFreePeriodMinutesColumn = "free_period_minutes"
var VUserSessionsPrivateIsLocalColumn = "is_local"
var VUserSessionsPrivateTimeSuspendedEvColumn = "time_suspended_ev"
var VUserSessionsPrivateTimeSuspendedEvseColumn = "time_suspended_evse"
var VUserSessionsPrivateReasonColumn = "reason"
var VUserSessionsPrivateHistoryColumn = "history"
var VUserSessionsPrivateSessionTypeColumn = "session_type"
var VUserSessionsPrivateChargingStationNeedsUnlockFirstColumn = "charging_station_needs_unlock_first"
var VUserSessionsPrivateColumns = []string{"id", "user_id", "started_at", "ended_at", "created_at", "updated_at", "deleted_at", "msp_transaction_id", "kwh", "evse_id", "max_charge_minutes", "max_charge_percent", "max_charge_price_cents", "charging_station_id", "ihomer_transaction_id", "token_id", "raw_cdr_id", "total_cost_cents", "user_ev_id", "stripe_payment_intent_id", "paid", "stripe_payment_method_id", "stripe_receipt_url", "stripe_charge_id", "energy_price_per_kwh", "charging_time_price_per_hour", "parking_time_price_per_hour", "transaction_fee_threshold_cents", "transaction_fee_below_threshold_cents", "transaction_fee_above_threshold_percent", "stripe_receipt_number", "total_cpo_split_cents", "total_platform_fee_cents", "total_stripe_fee_cents", "total_tax_cents", "total_cpo_net_cents", "card_brand", "card_last4", "charger_disconnected_at", "payable_parking_charge_cents", "stripe_tax_cents", "status", "total_amount_cents", "free_period_minutes", "is_local", "time_suspended_ev", "time_suspended_evse", "reason", "history", "session_type", "charging_station_needs_unlock_first"}

type VUserSessionsPrivate struct {
	ID                                  *uuid.UUID `json:"id" db:"id"`
	UserID                              *uuid.UUID `json:"user_id" db:"user_id"`
	StartedAt                           *time.Time `json:"started_at" db:"started_at"`
	EndedAt                             *time.Time `json:"ended_at" db:"ended_at"`
	CreatedAt                           *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                           *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt                           *time.Time `json:"deleted_at" db:"deleted_at"`
	MspTransactionID                    *string    `json:"msp_transaction_id" db:"msp_transaction_id"`
	KWH                                 *float64   `json:"kwh" db:"kwh"`
	EvseID                              *int64     `json:"evse_id" db:"evse_id"`
	MaxChargeMinutes                    *int64     `json:"max_charge_minutes" db:"max_charge_minutes"`
	MaxChargePercent                    *int64     `json:"max_charge_percent" db:"max_charge_percent"`
	MaxChargePriceCents                 *int64     `json:"max_charge_price_cents" db:"max_charge_price_cents"`
	ChargingStationID                   *string    `json:"charging_station_id" db:"charging_station_id"`
	IhomerTransactionID                 *string    `json:"ihomer_transaction_id" db:"ihomer_transaction_id"`
	TokenID                             *string    `json:"token_id" db:"token_id"`
	RawCdrID                            *string    `json:"raw_cdr_id" db:"raw_cdr_id"`
	TotalCostCents                      *int64     `json:"total_cost_cents" db:"total_cost_cents"`
	UserEvID                            *uuid.UUID `json:"user_ev_id" db:"user_ev_id"`
	StripePaymentIntentID               *string    `json:"stripe_payment_intent_id" db:"stripe_payment_intent_id"`
	Paid                                *bool      `json:"paid" db:"paid"`
	StripePaymentMethodID               *string    `json:"stripe_payment_method_id" db:"stripe_payment_method_id"`
	StripeReceiptURL                    *string    `json:"stripe_receipt_url" db:"stripe_receipt_url"`
	StripeChargeID                      *string    `json:"stripe_charge_id" db:"stripe_charge_id"`
	EnergyPricePerKWH                   *float64   `json:"energy_price_per_kwh" db:"energy_price_per_kwh"`
	ChargingTimePricePerHour            *float64   `json:"charging_time_price_per_hour" db:"charging_time_price_per_hour"`
	ParkingTimePricePerHour             *float64   `json:"parking_time_price_per_hour" db:"parking_time_price_per_hour"`
	TransactionFeeThresholdCents        *int64     `json:"transaction_fee_threshold_cents" db:"transaction_fee_threshold_cents"`
	TransactionFeeBelowThresholdCents   *int64     `json:"transaction_fee_below_threshold_cents" db:"transaction_fee_below_threshold_cents"`
	TransactionFeeAboveThresholdPercent *int64     `json:"transaction_fee_above_threshold_percent" db:"transaction_fee_above_threshold_percent"`
	StripeReceiptNumber                 *string    `json:"stripe_receipt_number" db:"stripe_receipt_number"`
	TotalCPOSplitCents                  *int64     `json:"total_cpo_split_cents" db:"total_cpo_split_cents"`
	TotalPlatformFeeCents               *int64     `json:"total_platform_fee_cents" db:"total_platform_fee_cents"`
	TotalStripeFeeCents                 *int64     `json:"total_stripe_fee_cents" db:"total_stripe_fee_cents"`
	TotalTaxCents                       *int64     `json:"total_tax_cents" db:"total_tax_cents"`
	TotalCPONetCents                    *int64     `json:"total_cpo_net_cents" db:"total_cpo_net_cents"`
	CardBrand                           *string    `json:"card_brand" db:"card_brand"`
	CardLast4                           *string    `json:"card_last4" db:"card_last4"`
	ChargerDisconnectedAt               *time.Time `json:"charger_disconnected_at" db:"charger_disconnected_at"`
	PayableParkingChargeCents           *int64     `json:"payable_parking_charge_cents" db:"payable_parking_charge_cents"`
	StripeTaxCents                      *int64     `json:"stripe_tax_cents" db:"stripe_tax_cents"`
	Status                              *string    `json:"status" db:"status"`
	TotalAmountCents                    *int64     `json:"total_amount_cents" db:"total_amount_cents"`
	FreePeriodMinutes                   *int64     `json:"free_period_minutes" db:"free_period_minutes"`
	IsLocal                             *bool      `json:"is_local" db:"is_local"`
	TimeSuspendedEv                     *int64     `json:"time_suspended_ev" db:"time_suspended_ev"`
	TimeSuspendedEvse                   *int64     `json:"time_suspended_evse" db:"time_suspended_evse"`
	Reason                              *string    `json:"reason" db:"reason"`
	History                             *any       `json:"history" db:"history"`
	SessionType                         *string    `json:"session_type" db:"session_type"`
	ChargingStationNeedsUnlockFirst     *bool      `json:"charging_station_needs_unlock_first" db:"charging_station_needs_unlock_first"`
}

func SelectVUserSessionsPrivate(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*VUserSessionsPrivate, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM v_user_sessions_private%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*VUserSessionsPrivate, 0)
	for rows.Next() {
		rowCount++

		var item VUserSessionsPrivate
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectVUserSessionsPrivate(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectVUserSessionsPrivate(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var VUserSessionsPublicTable = "v_user_sessions_public"
var VUserSessionsPublicIDColumn = "id"
var VUserSessionsPublicUserIDColumn = "user_id"
var VUserSessionsPublicStartedAtColumn = "started_at"
var VUserSessionsPublicEndedAtColumn = "ended_at"
var VUserSessionsPublicCreatedAtColumn = "created_at"
var VUserSessionsPublicUpdatedAtColumn = "updated_at"
var VUserSessionsPublicDeletedAtColumn = "deleted_at"
var VUserSessionsPublicMspTransactionIDColumn = "msp_transaction_id"
var VUserSessionsPublicKWHColumn = "kwh"
var VUserSessionsPublicEvseIDColumn = "evse_id"
var VUserSessionsPublicMaxChargeMinutesColumn = "max_charge_minutes"
var VUserSessionsPublicMaxChargePercentColumn = "max_charge_percent"
var VUserSessionsPublicMaxChargePriceCentsColumn = "max_charge_price_cents"
var VUserSessionsPublicChargingStationIDColumn = "charging_station_id"
var VUserSessionsPublicIhomerTransactionIDColumn = "ihomer_transaction_id"
var VUserSessionsPublicTokenIDColumn = "token_id"
var VUserSessionsPublicRawCdrIDColumn = "raw_cdr_id"
var VUserSessionsPublicTotalCostCentsColumn = "total_cost_cents"
var VUserSessionsPublicUserEvIDColumn = "user_ev_id"
var VUserSessionsPublicStripePaymentIntentIDColumn = "stripe_payment_intent_id"
var VUserSessionsPublicPaidColumn = "paid"
var VUserSessionsPublicStripePaymentMethodIDColumn = "stripe_payment_method_id"
var VUserSessionsPublicStripeReceiptURLColumn = "stripe_receipt_url"
var VUserSessionsPublicStripeChargeIDColumn = "stripe_charge_id"
var VUserSessionsPublicEnergyPricePerKWHColumn = "energy_price_per_kwh"
var VUserSessionsPublicChargingTimePricePerHourColumn = "charging_time_price_per_hour"
var VUserSessionsPublicParkingTimePricePerHourColumn = "parking_time_price_per_hour"
var VUserSessionsPublicTransactionFeeThresholdCentsColumn = "transaction_fee_threshold_cents"
var VUserSessionsPublicTransactionFeeBelowThresholdCentsColumn = "transaction_fee_below_threshold_cents"
var VUserSessionsPublicTransactionFeeAboveThresholdPercentColumn = "transaction_fee_above_threshold_percent"
var VUserSessionsPublicStripeReceiptNumberColumn = "stripe_receipt_number"
var VUserSessionsPublicTotalCPOSplitCentsColumn = "total_cpo_split_cents"
var VUserSessionsPublicTotalPlatformFeeCentsColumn = "total_platform_fee_cents"
var VUserSessionsPublicTotalStripeFeeCentsColumn = "total_stripe_fee_cents"
var VUserSessionsPublicTotalTaxCentsColumn = "total_tax_cents"
var VUserSessionsPublicTotalCPONetCentsColumn = "total_cpo_net_cents"
var VUserSessionsPublicCardBrandColumn = "card_brand"
var VUserSessionsPublicCardLast4Column = "card_last4"
var VUserSessionsPublicChargerDisconnectedAtColumn = "charger_disconnected_at"
var VUserSessionsPublicPayableParkingChargeCentsColumn = "payable_parking_charge_cents"
var VUserSessionsPublicStripeTaxCentsColumn = "stripe_tax_cents"
var VUserSessionsPublicStatusColumn = "status"
var VUserSessionsPublicTotalAmountCentsColumn = "total_amount_cents"
var VUserSessionsPublicFreePeriodMinutesColumn = "free_period_minutes"
var VUserSessionsPublicIsLocalColumn = "is_local"
var VUserSessionsPublicTimeSuspendedEvColumn = "time_suspended_ev"
var VUserSessionsPublicTimeSuspendedEvseColumn = "time_suspended_evse"
var VUserSessionsPublicReasonColumn = "reason"
var VUserSessionsPublicHistoryColumn = "history"
var VUserSessionsPublicSessionTypeColumn = "session_type"
var VUserSessionsPublicChargingStationNeedsUnlockFirstColumn = "charging_station_needs_unlock_first"
var VUserSessionsPublicColumns = []string{"id", "user_id", "started_at", "ended_at", "created_at", "updated_at", "deleted_at", "msp_transaction_id", "kwh", "evse_id", "max_charge_minutes", "max_charge_percent", "max_charge_price_cents", "charging_station_id", "ihomer_transaction_id", "token_id", "raw_cdr_id", "total_cost_cents", "user_ev_id", "stripe_payment_intent_id", "paid", "stripe_payment_method_id", "stripe_receipt_url", "stripe_charge_id", "energy_price_per_kwh", "charging_time_price_per_hour", "parking_time_price_per_hour", "transaction_fee_threshold_cents", "transaction_fee_below_threshold_cents", "transaction_fee_above_threshold_percent", "stripe_receipt_number", "total_cpo_split_cents", "total_platform_fee_cents", "total_stripe_fee_cents", "total_tax_cents", "total_cpo_net_cents", "card_brand", "card_last4", "charger_disconnected_at", "payable_parking_charge_cents", "stripe_tax_cents", "status", "total_amount_cents", "free_period_minutes", "is_local", "time_suspended_ev", "time_suspended_evse", "reason", "history", "session_type", "charging_station_needs_unlock_first"}

type VUserSessionsPublic struct {
	ID                                  *uuid.UUID `json:"id" db:"id"`
	UserID                              *uuid.UUID `json:"user_id" db:"user_id"`
	StartedAt                           *time.Time `json:"started_at" db:"started_at"`
	EndedAt                             *time.Time `json:"ended_at" db:"ended_at"`
	CreatedAt                           *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                           *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt                           *time.Time `json:"deleted_at" db:"deleted_at"`
	MspTransactionID                    *string    `json:"msp_transaction_id" db:"msp_transaction_id"`
	KWH                                 *float64   `json:"kwh" db:"kwh"`
	EvseID                              *int64     `json:"evse_id" db:"evse_id"`
	MaxChargeMinutes                    *int64     `json:"max_charge_minutes" db:"max_charge_minutes"`
	MaxChargePercent                    *int64     `json:"max_charge_percent" db:"max_charge_percent"`
	MaxChargePriceCents                 *int64     `json:"max_charge_price_cents" db:"max_charge_price_cents"`
	ChargingStationID                   *string    `json:"charging_station_id" db:"charging_station_id"`
	IhomerTransactionID                 *string    `json:"ihomer_transaction_id" db:"ihomer_transaction_id"`
	TokenID                             *string    `json:"token_id" db:"token_id"`
	RawCdrID                            *string    `json:"raw_cdr_id" db:"raw_cdr_id"`
	TotalCostCents                      *int64     `json:"total_cost_cents" db:"total_cost_cents"`
	UserEvID                            *uuid.UUID `json:"user_ev_id" db:"user_ev_id"`
	StripePaymentIntentID               *string    `json:"stripe_payment_intent_id" db:"stripe_payment_intent_id"`
	Paid                                *bool      `json:"paid" db:"paid"`
	StripePaymentMethodID               *string    `json:"stripe_payment_method_id" db:"stripe_payment_method_id"`
	StripeReceiptURL                    *string    `json:"stripe_receipt_url" db:"stripe_receipt_url"`
	StripeChargeID                      *string    `json:"stripe_charge_id" db:"stripe_charge_id"`
	EnergyPricePerKWH                   *float64   `json:"energy_price_per_kwh" db:"energy_price_per_kwh"`
	ChargingTimePricePerHour            *float64   `json:"charging_time_price_per_hour" db:"charging_time_price_per_hour"`
	ParkingTimePricePerHour             *float64   `json:"parking_time_price_per_hour" db:"parking_time_price_per_hour"`
	TransactionFeeThresholdCents        *int64     `json:"transaction_fee_threshold_cents" db:"transaction_fee_threshold_cents"`
	TransactionFeeBelowThresholdCents   *int64     `json:"transaction_fee_below_threshold_cents" db:"transaction_fee_below_threshold_cents"`
	TransactionFeeAboveThresholdPercent *int64     `json:"transaction_fee_above_threshold_percent" db:"transaction_fee_above_threshold_percent"`
	StripeReceiptNumber                 *string    `json:"stripe_receipt_number" db:"stripe_receipt_number"`
	TotalCPOSplitCents                  *int64     `json:"total_cpo_split_cents" db:"total_cpo_split_cents"`
	TotalPlatformFeeCents               *int64     `json:"total_platform_fee_cents" db:"total_platform_fee_cents"`
	TotalStripeFeeCents                 *int64     `json:"total_stripe_fee_cents" db:"total_stripe_fee_cents"`
	TotalTaxCents                       *int64     `json:"total_tax_cents" db:"total_tax_cents"`
	TotalCPONetCents                    *int64     `json:"total_cpo_net_cents" db:"total_cpo_net_cents"`
	CardBrand                           *string    `json:"card_brand" db:"card_brand"`
	CardLast4                           *string    `json:"card_last4" db:"card_last4"`
	ChargerDisconnectedAt               *time.Time `json:"charger_disconnected_at" db:"charger_disconnected_at"`
	PayableParkingChargeCents           *int64     `json:"payable_parking_charge_cents" db:"payable_parking_charge_cents"`
	StripeTaxCents                      *int64     `json:"stripe_tax_cents" db:"stripe_tax_cents"`
	Status                              *string    `json:"status" db:"status"`
	TotalAmountCents                    *int64     `json:"total_amount_cents" db:"total_amount_cents"`
	FreePeriodMinutes                   *int64     `json:"free_period_minutes" db:"free_period_minutes"`
	IsLocal                             *bool      `json:"is_local" db:"is_local"`
	TimeSuspendedEv                     *int64     `json:"time_suspended_ev" db:"time_suspended_ev"`
	TimeSuspendedEvse                   *int64     `json:"time_suspended_evse" db:"time_suspended_evse"`
	Reason                              *string    `json:"reason" db:"reason"`
	History                             *any       `json:"history" db:"history"`
	SessionType                         *string    `json:"session_type" db:"session_type"`
	ChargingStationNeedsUnlockFirst     *bool      `json:"charging_station_needs_unlock_first" db:"charging_station_needs_unlock_first"`
}

func SelectVUserSessionsPublic(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*VUserSessionsPublic, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM v_user_sessions_public%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*VUserSessionsPublic, 0)
	for rows.Next() {
		rowCount++

		var item VUserSessionsPublic
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectVUserSessionsPublic(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectVUserSessionsPublic(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var VUsersTable = "v_users"
var VUsersIDColumn = "id"
var VUsersEmailColumn = "email"
var VUsersEnabledColumn = "enabled"
var VUsersDefaultTokenIDColumn = "default_token_id"
var VUsersCreatedAtColumn = "created_at"
var VUsersUpdatedAtColumn = "updated_at"
var VUsersDeletedAtColumn = "deleted_at"
var VUsersActiveUserSubscriptionIDColumn = "active_user_subscription_id"
var VUsersNameColumn = "name"
var VUsersDOBColumn = "dob"
var VUsersNationalityColumn = "nationality"
var VUsersCountryColumn = "country"
var VUsersAddressColumn = "address"
var VUsersCompanyNameColumn = "company_name"
var VUsersCompanyWebsiteColumn = "company_website"
var VUsersCompanyCountryColumn = "company_country"
var VUsersCompanyAddressColumn = "company_address"
var VUsersCompanyLicenseNumberColumn = "company_license_number"
var VUsersCPOTypeColumn = "cpo_type"
var VUsersCPOReasonColumn = "cpo_reason"
var VUsersStripeCustomerIDColumn = "stripe_customer_id"
var VUsersFirstNameColumn = "first_name"
var VUsersLastNameColumn = "last_name"
var VUsersPhoneColumn = "phone"
var VUsersCountryDialCodeColumn = "country_dial_code"
var VUsersImageColumn = "image"
var VUsersStripeAccountIDColumn = "stripe_account_id"
var VUsersStripeAccountTypeColumn = "stripe_account_type"
var VUsersRoleIDColumn = "role_id"
var VUsersNotificationTokenColumn = "notification_token"
var VUsersEmailValidateCodeColumn = "email_validate_code"
var VUsersCustomerIDColumn = "customer_id"
var VUsersDiscountIDColumn = "discount_id"
var VUsersIsEmailVerifiedColumn = "is_email_verified"
var VUsersManagerUserIDColumn = "manager_user_id"
var VUsersColumns = []string{"id", "email", "enabled", "default_token_id", "created_at", "updated_at", "deleted_at", "active_user_subscription_id", "name", "dob", "nationality", "country", "address", "company_name", "company_website", "company_country", "company_address", "company_license_number", "cpo_type", "cpo_reason", "stripe_customer_id", "first_name", "last_name", "phone", "country_dial_code", "image", "stripe_account_id", "stripe_account_type", "role_id", "notification_token", "email_validate_code", "customer_id", "discount_id", "is_email_verified", "manager_user_id"}

type VUsers struct {
	ID                       *uuid.UUID `json:"id" db:"id"`
	Email                    *string    `json:"email" db:"email"`
	Enabled                  *bool      `json:"enabled" db:"enabled"`
	DefaultTokenID           *string    `json:"default_token_id" db:"default_token_id"`
	CreatedAt                *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt                *time.Time `json:"deleted_at" db:"deleted_at"`
	ActiveUserSubscriptionID *uuid.UUID `json:"active_user_subscription_id" db:"active_user_subscription_id"`
	Name                     *string    `json:"name" db:"name"`
	DOB                      *time.Time `json:"dob" db:"dob"`
	Nationality              *string    `json:"nationality" db:"nationality"`
	Country                  *string    `json:"country" db:"country"`
	Address                  *string    `json:"address" db:"address"`
	CompanyName              *string    `json:"company_name" db:"company_name"`
	CompanyWebsite           *string    `json:"company_website" db:"company_website"`
	CompanyCountry           *string    `json:"company_country" db:"company_country"`
	CompanyAddress           *string    `json:"company_address" db:"company_address"`
	CompanyLicenseNumber     *string    `json:"company_license_number" db:"company_license_number"`
	CPOType                  *string    `json:"cpo_type" db:"cpo_type"`
	CPOReason                *string    `json:"cpo_reason" db:"cpo_reason"`
	StripeCustomerID         *string    `json:"stripe_customer_id" db:"stripe_customer_id"`
	FirstName                *string    `json:"first_name" db:"first_name"`
	LastName                 *string    `json:"last_name" db:"last_name"`
	Phone                    *string    `json:"phone" db:"phone"`
	CountryDialCode          *string    `json:"country_dial_code" db:"country_dial_code"`
	Image                    *string    `json:"image" db:"image"`
	StripeAccountID          *string    `json:"stripe_account_id" db:"stripe_account_id"`
	StripeAccountType        *string    `json:"stripe_account_type" db:"stripe_account_type"`
	RoleID                   *int64     `json:"role_id" db:"role_id"`
	NotificationToken        *string    `json:"notification_token" db:"notification_token"`
	EmailValidateCode        *string    `json:"email_validate_code" db:"email_validate_code"`
	CustomerID               *int64     `json:"customer_id" db:"customer_id"`
	DiscountID               *uuid.UUID `json:"discount_id" db:"discount_id"`
	IsEmailVerified          *bool      `json:"is_email_verified" db:"is_email_verified"`
	ManagerUserID            *uuid.UUID `json:"manager_user_id" db:"manager_user_id"`
}

func SelectVUsers(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*VUsers, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM v_users%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*VUsers, 0)
	for rows.Next() {
		rowCount++

		var item VUsers
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectVUsers(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectVUsers(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var VUsersRegisteredOverTimeTable = "v_users_registered_over_time"
var VUsersRegisteredOverTimeTimeColumn = "time"
var VUsersRegisteredOverTimeCountColumn = "count"
var VUsersRegisteredOverTimeColumns = []string{"time", "count"}

type VUsersRegisteredOverTime struct {
	Time  *time.Time `json:"time" db:"time"`
	Count *float64   `json:"count" db:"count"`
}

func SelectVUsersRegisteredOverTime(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*VUsersRegisteredOverTime, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM v_users_registered_over_time%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*VUsersRegisteredOverTime, 0)
	for rows.Next() {
		rowCount++

		var item VUsersRegisteredOverTime
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectVUsersRegisteredOverTime(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectVUsersRegisteredOverTime(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}

var VoltguardEventsTable = "voltguard_events"
var VoltguardEventsIDColumn = "id"
var VoltguardEventsErrorCodeColumn = "error_code"
var VoltguardEventsEventTypeColumn = "event_type"
var VoltguardEventsEventMessageColumn = "event_message"
var VoltguardEventsEventTimestampColumn = "event_timestamp"
var VoltguardEventsOcppServerLogIDColumn = "ocpp_server_log_id"
var VoltguardEventsIhomerLogIDColumn = "ihomer_log_id"
var VoltguardEventsIsV2LogColumn = "is_v2_log"
var VoltguardEventsLastSeenTimestampColumn = "last_seen_timestamp"
var VoltguardEventsChargerIDColumn = "charger_id"
var VoltguardEventsConnectorIDColumn = "connector_id"
var VoltguardEventsLocationIDColumn = "location_id"
var VoltguardEventsUserIDColumn = "user_id"
var VoltguardEventsCreatedAtColumn = "created_at"
var VoltguardEventsUpdatedAtColumn = "updated_at"
var VoltguardEventsDeletedAtColumn = "deleted_at"
var VoltguardEventsColumns = []string{"id", "error_code", "event_type", "event_message", "event_timestamp", "ocpp_server_log_id", "ihomer_log_id", "is_v2_log", "last_seen_timestamp", "charger_id", "connector_id", "location_id", "user_id", "created_at", "updated_at", "deleted_at"}

type VoltguardEvents struct {
	ID                uuid.UUID  `json:"id" db:"id"`
	ErrorCode         string     `json:"error_code" db:"error_code"`
	EventType         string     `json:"event_type" db:"event_type"`
	EventMessage      *string    `json:"event_message" db:"event_message"`
	EventTimestamp    time.Time  `json:"event_timestamp" db:"event_timestamp"`
	OcppServerLogID   *uuid.UUID `json:"ocpp_server_log_id" db:"ocpp_server_log_id"`
	IhomerLogID       *string    `json:"ihomer_log_id" db:"ihomer_log_id"`
	IsV2Log           bool       `json:"is_v2_log" db:"is_v2_log"`
	LastSeenTimestamp *time.Time `json:"last_seen_timestamp" db:"last_seen_timestamp"`
	ChargerID         string     `json:"charger_id" db:"charger_id"`
	ConnectorID       *int64     `json:"connector_id" db:"connector_id"`
	LocationID        uuid.UUID  `json:"location_id" db:"location_id"`
	UserID            uuid.UUID  `json:"user_id" db:"user_id"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at" db:"deleted_at"`
}

func SelectVoltguardEvents(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]*VoltguardEvents, error) {
	mu.RLock()
	debug := actualDebug
	mu.RUnlock()

	var buildStart int64
	var buildStop int64
	var execStart int64
	var execStop int64
	var scanStart int64
	var scanStop int64

	var sql string
	var columnCount int = len(columns)
	var rowCount int64
	defer func() {
		if !debug {
			return
		}

		buildDuration := 0.0
		execDuration := 0.0
		scanDuration := 0.0

		if buildStop > 0 {
			buildDuration = float64(buildStop-buildStart) * 1e-9
		}

		if execStop > 0 {
			execDuration = float64(execStop-execStart) * 1e-9
		}

		if scanStop > 0 {
			scanDuration = float64(scanStop-scanStart) * 1e-9
		}

		if debug {
			logger.Printf(
				"selected %v columns, %v rows; %.3f seconds to build, %.3f seconds to execute, %.3f seconds to scan; sql:\n%v",
				columnCount, rowCount, buildDuration, execDuration, scanDuration, sql,
			)
		}
	}()

	buildStart = time.Now().UnixNano()

	where := strings.TrimSpace(strings.Join(wheres, " AND "))
	if len(where) > 0 {
		where = " WHERE " + where
	}

	sql = fmt.Sprintf(
		"SELECT %v FROM voltguard_events%v",
		strings.Join(columns, ", "),
		where,
	)

	if orderBy != nil {
		actualOrderBy := strings.TrimSpace(*orderBy)
		if len(actualOrderBy) > 0 {
			sql += fmt.Sprintf(" ORDER BY %v", actualOrderBy)
		}
	}

	if limit != nil && *limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %v", *limit)
	}

	buildStop = time.Now().UnixNano()

	execStart = time.Now().UnixNano()

	selectCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	rows, err := db.QueryxContext(selectCtx, sql)
	if err != nil {
		return nil, err
	}

	execStop = time.Now().UnixNano()

	scanStart = time.Now().UnixNano()

	items := make([]*VoltguardEvents, 0)
	for rows.Next() {
		rowCount++

		var item VoltguardEvents
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		// this is where any post-scan processing would appear (if required)

		items = append(items, &item)
	}

	scanStop = time.Now().UnixNano()

	return items, nil
}

func genericSelectVoltguardEvents(ctx context.Context, db *sqlx.DB, columns []string, orderBy *string, limit *int, offset *int, wheres ...string) ([]any, error) {
	items, err := SelectVoltguardEvents(ctx, db, columns, orderBy, limit, offset, wheres...)
	if err != nil {
		return nil, err
	}

	genericItems := make([]any, 0)
	for _, item := range items {
		genericItems = append(genericItems, item)
	}

	return genericItems, nil
}
