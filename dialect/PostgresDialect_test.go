package dialect

import (
	"testing"
	"github.com/plopezm/goedb/metadata"
	"reflect"
	"github.com/stretchr/testify/assert"
)

func TestGetSQLColumnModelPrimaryKey (t *testing.T){
	dialect := new(PostgresDialect)
	var column metadata.GoedbColumn

	column.Title = "primaryKeyColumn"
	column.AutoIncrement = true
	column.PrimaryKey = true
	column.ColumnType = reflect.Int32

	columns, pksFound, constraints, err := dialect.GetSQLColumnModel(column)
	t.Log(columns)
	t.Log(pksFound)
	t.Log(constraints)
	t.Log(err)
	assert.Equal(t, "primaryKeyColumn SERIAL,", columns)
	assert.Equal(t, "primaryKeyColumn,", pksFound)
	assert.Equal(t, "", constraints)
	assert.Nil(t, err)
}

func TestGetSQLColumnModelFloat (t *testing.T){
	dialect := new(PostgresDialect)
	var column metadata.GoedbColumn

	column.Title = "floatColumn"
	column.ColumnType = reflect.Float32

	columns, pksFound, constraints, err := dialect.GetSQLColumnModel(column)
	t.Log(columns)
	t.Log(pksFound)
	t.Log(constraints)
	t.Log(err)
	assert.Equal(t, "floatColumn FLOAT,", columns)
	assert.Equal(t, "", pksFound)
	assert.Equal(t, "", constraints)
	assert.Nil(t, err)
}

func TestGetSQLColumnModelBool (t *testing.T){
	dialect := new(PostgresDialect)
	var column metadata.GoedbColumn

	column.Title = "boolColumn"
	column.ColumnType = reflect.Bool

	columns, pksFound, constraints, err := dialect.GetSQLColumnModel(column)
	t.Log(columns)
	t.Log(pksFound)
	t.Log(constraints)
	t.Log(err)
	assert.Equal(t, "boolColumn BOOLEAN,", columns)
	assert.Equal(t, "", pksFound)
	assert.Equal(t, "", constraints)
	assert.Nil(t, err)
}

func TestGetSQLColumnModelString (t *testing.T){
	dialect := new(PostgresDialect)
	var column metadata.GoedbColumn

	column.Title = "stringColumn"
	column.ColumnType = reflect.String

	columns, pksFound, constraints, err := dialect.GetSQLColumnModel(column)
	t.Log(columns)
	t.Log(pksFound)
	t.Log(constraints)
	t.Log(err)
	assert.Equal(t, "stringColumn VARCHAR,", columns)
	assert.Equal(t, "", pksFound)
	assert.Equal(t, "", constraints)
	assert.Nil(t, err)
}

func TestGetSQLColumnModelBigInt (t *testing.T){
	dialect := new(PostgresDialect)
	var column metadata.GoedbColumn

	column.Title = "bigIntColumn"
	column.ColumnType = reflect.Int64

	columns, pksFound, constraints, err := dialect.GetSQLColumnModel(column)
	t.Log(columns)
	t.Log(pksFound)
	t.Log(constraints)
	t.Log(err)
	assert.Equal(t, "bigIntColumn BIGINT,", columns)
	assert.Equal(t, "", pksFound)
	assert.Equal(t, "", constraints)
	assert.Nil(t, err)
}

func TestGetSQLColumnModelBigIntPrimaryKey (t *testing.T){
	dialect := new(PostgresDialect)
	var column metadata.GoedbColumn

	column.Title = "bigIntColumn"
	column.ColumnType = reflect.Int64
	column.PrimaryKey = true

	columns, pksFound, constraints, err := dialect.GetSQLColumnModel(column)
	t.Log(columns)
	t.Log(pksFound)
	t.Log(constraints)
	t.Log(err)
	assert.Equal(t, "bigIntColumn BIGINT,", columns)
	assert.Equal(t, "bigIntColumn,", pksFound)
	assert.Equal(t, "", constraints)
	assert.Nil(t, err)
}

func TestGetSQLColumnModelBigIntPrimaryKeyAutoincrement (t *testing.T){
	dialect := new(PostgresDialect)
	var column metadata.GoedbColumn

	column.Title = "bigIntColumn"
	column.ColumnType = reflect.Int64
	column.PrimaryKey = true
	column.AutoIncrement = true

	columns, pksFound, constraints, err := dialect.GetSQLColumnModel(column)
	t.Log(columns)
	t.Log(pksFound)
	t.Log(constraints)
	t.Log(err)
	assert.Equal(t, "bigIntColumn SERIAL,", columns)
	assert.Equal(t, "bigIntColumn,", pksFound)
	assert.Equal(t, "", constraints)
	assert.Nil(t, err)
}
