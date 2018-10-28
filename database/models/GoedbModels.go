package models

import "reflect"

// Table represents the metadata of a table
type Table struct {
	Name          string
	Columns       []Column
	MappedColumns []Column
	PrimaryKeys   []PrimaryKey
}

//PrimaryKey contains the name and the type of a primary key
type PrimaryKey struct {
	Name string
	Type reflect.Kind
}

//ForeignKey contains the table and column reference of a ForeignKey
type ForeignKey struct {
	IsForeignKey              bool
	ForeignKeyTableReference  string
	ForeignKeyColumnReference string
}

type MappedByField struct {
	TargetTableName string
	TargetTablePK   string
}

// Column represents the metadata of a column
type Column struct {
	Title          string
	ColumnType     reflect.Kind
	ColumnTypeName string
	PrimaryKey     bool
	Unique         bool
	ForeignKey     ForeignKey
	AutoIncrement  bool
	IsComplex      bool
	Ignore         bool
	IsMapped       bool
	MappedBy       MappedByField
}

// Result is the result for some operation in database
type Result struct {
	NumRecordsAffected int64
	LastInsertId       int64
}
