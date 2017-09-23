package dialect

import (
	"github.com/plopezm/goedb/metadata"
)

type Dialect interface {

	GetSQLCreate(table metadata.GoedbTable) string
	GetSQLDelete(table metadata.GoedbTable, where string, instance interface{}) (string, error)
	GetSQLInsert(table metadata.GoedbTable, instance interface{}) (string, error)
	GetFirstQuery(table metadata.GoedbTable, where string, instance interface{}) (string, error)
	GetFindQuery(table metadata.GoedbTable, where string) (string, error)
}

func GetDialect(driver string) Dialect{
	switch driver {
	case "sqlite3":
		return new(SQLiteDialect)
	default:
		return new(SQLiteDialect)
	}
}