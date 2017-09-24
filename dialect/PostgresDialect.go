package dialect

import (
	"errors"
	"github.com/plopezm/goedb/metadata"
	"reflect"
)

type PostgresDialect struct {
}

func (dialect *PostgresDialect) GetSQLColumnModel(value metadata.GoedbColumn) (string, string, string, error) {
	var pksFound string
	var constraints string
	column := value.Title

	switch value.ColumnType {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint:
		if value.AutoIncrement {
			column += " SERIAL"
		} else {
			column += " INTEGER"
		}
	case reflect.Int64, reflect.Uint64:
		if value.AutoIncrement {
			column += " SERIAL"
		} else {
			column += " BIGINT"
		}
	case reflect.Float32, reflect.Float64:
		column += " FLOAT"
	case reflect.Bool:
		column += " BOOLEAN"
	case reflect.String:
		column += " VARCHAR"
	default:
		return "", "", "", errors.New("Type unknown")
	}

	if value.Unique {
		column += " UNIQUE"
	}

	if value.PrimaryKey {
		pksFound += value.Title + ","
	}

	if value.ForeignKey {
		constraints += ", FOREIGN KEY (" + value.Title + ") REFERENCES " + value.ForeignKeyReference + " ON DELETE CASCADE"
	}
	column += ","
	return column, pksFound, constraints, nil
}
