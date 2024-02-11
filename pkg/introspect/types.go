package introspect

type Column struct {
	Name              string  `json:"column"`
	DataType          string  `json:"datatype"`
	TableName         string  `json:"table"`
	Pos               int64   `json:"pos"`
	TypeID            string  `json:"typeid"`
	TypeLen           int64   `json:"typelen"`
	TypeMod           int64   `json:"typemod"`
	NotNull           bool    `json:"notnull"`
	HasDefault        bool    `json:"hasdefault"`
	HasMissing        bool    `json:"hasmissing"`
	ForeignTableName  *string `json:"ftable"`
	ForeignColumnName *string `json:"fcolumn"`
	ParentID          string  `json:"parent_id"`
	ForeignTable      *Table  `json:"-"`
	ForeignColumn     *Column `json:"-"`
	ParentTable       *Table  `json:"-"`
	ZeroType          any     `json:"-"`
	TypeTemplate      string  `json:"-"`
}

type Table struct {
	Name               string             `json:"tablename"`
	OID                string             `json:"oid"`
	Schema             string             `json:"schema"`
	RelTuples          float64            `json:"reltuples"`
	RelKind            string             `json:"relkind"`
	RelAM              string             `json:"relam"`
	RelACL             any                `json:"relacl"` // TODO: figure out this (probable) struct as required
	RelType            string             `json:"reltype"`
	RelOwner           string             `json:"relowner"`
	RelHasIndex        bool               `json:"relhasindex"`
	Columns            []*Column          `json:"columns"`
	ColumnByName       map[string]*Column `json:"-"`
	ForeignTableByName map[string]*Table  `json:"-"` // WARNING: only takes holds the last in the event a table with more than one relationship to the same foreign table
}
