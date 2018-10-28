package dialect

import (
	"errors"
	"reflect"

	"github.com/plopezm/goedb/database/models"
)

//SQLite3Dialect contains a few functions that are different from standard sql dbaccess
type SQLite3Dialect struct {
}

// GetSQLCreateTableColumn returns the model of a column for SQLite3
func (specifics *SQLite3Dialect) GetSQLCreateTableColumn(value models.Column) (sqlColumnLine string, primaryKey string, constraints string, err error) {
	sqlColumnLine = value.Title

	switch value.ColumnType {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint:
		sqlColumnLine += " INTEGER"
	case reflect.Int64, reflect.Uint64:
		sqlColumnLine += " BIGINT"
	case reflect.Float32, reflect.Float64:
		sqlColumnLine += " FLOAT"
	case reflect.Bool:
		sqlColumnLine += " BOOLEAN"
	case reflect.String:
		sqlColumnLine += " VARCHAR"
	default:
		return "", "", "", errors.New("Type unknown")
	}

	if value.Unique {
		sqlColumnLine += " UNIQUE"
	}

	if value.PrimaryKey && value.AutoIncrement {
		sqlColumnLine += " PRIMARY KEY AUTOINCREMENT"
	} else if value.PrimaryKey {
		primaryKey += value.Title + ","
	}

	if value.ForeignKey.IsForeignKey {
		constraints += ", FOREIGN KEY (" + value.Title + ") REFERENCES " + value.ForeignKey.ForeignKeyTableReference + "(" + value.ForeignKey.ForeignKeyColumnReference + ")" + " ON DELETE CASCADE"
	}
	sqlColumnLine += ","
	return sqlColumnLine, primaryKey, constraints, nil
}
