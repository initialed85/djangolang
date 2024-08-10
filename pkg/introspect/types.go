package introspect

type Column struct {
	Name               string  `json:"column"`
	DataType           string  `json:"datatype"`
	TableName          string  `json:"table"`
	Pos                int64   `json:"pos"`
	TypeID             string  `json:"typeid"`
	TypeLen            int64   `json:"typelen"`
	TypeMod            int64   `json:"typemod"`
	NotNull            bool    `json:"notnull"`
	HasDefault         bool    `json:"hasdefault"`
	HasMissing         bool    `json:"hasmissing"`
	IsPrimaryKey       bool    `json:"ispkey"`
	ForeignTableName   *string `json:"ftable"`
	ForeignColumnName  *string `json:"fcolumn"`
	ParentID           string  `json:"parent_id"`
	ForeignTable       *Table  `json:"-"`
	ForeignColumn      *Column `json:"-"`
	ParentTable        *Table  `json:"-"`
	ZeroType           any     `json:"zero_type"`
	QueryTypeTemplate  string  `json:"query_type_template"`
	StreamTypeTemplate string  `json:"stream_type_template"`
	TypeTemplate       string  `json:"type_template"`
}

type Table struct {
	Name                string             `json:"tablename"`
	OID                 string             `json:"oid"`
	Schema              string             `json:"schema"`
	RelTuples           float64            `json:"reltuples"`
	RelKind             string             `json:"relkind"`
	RelAM               string             `json:"relam"`
	RelACL              any                `json:"relacl"` // TODO: figure out this (probable) struct as required
	RelType             string             `json:"reltype"`
	RelOwner            string             `json:"relowner"`
	RelHasIndex         bool               `json:"relhasindex"`
	Columns             []*Column          `json:"columns"`
	ColumnByName        map[string]*Column `json:"-"`
	PrimaryKeyColumn    *Column            `json:"-"`
	ForeignTables       []*Table           `json:"-"`
	ReferencedByColumns []*Column          `json:"-"`
}
