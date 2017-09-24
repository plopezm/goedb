package dialect

import (
	"github.com/plopezm/goedb/metadata"
	"reflect"
	"strconv"
	"errors"
)

type Dialect interface {
	GetSQLCreate(table metadata.GoedbTable) string
	GetSQLDelete(table metadata.GoedbTable, where string, instance interface{}) (string, error)
	GetSQLInsert(table metadata.GoedbTable, instance interface{}) (string, error)
	GetFirstQuery(table metadata.GoedbTable, where string, instance interface{}) (string, error)
	GetFindQuery(table metadata.GoedbTable, where string) (string, error)
	GetDropTableQuery(table metadata.GoedbTable) (string)
}

func GetDialect(driver string) Dialect{
	switch driver {
	case "sqlite3":
		return new(SQLiteDialect)
	case "postgres":
		return new(PostgresDialect)
	default:
		return new(SQLiteDialect)
	}
}

/*
	Returns columns names and values for inserting values
*/
func getColumnsAndValues(metatable metadata.GoedbTable, instance interface{}) (string, string, error) {
	strCols := ""
	strValues := ""

	instanceType := metadata.GetType(instance)
	intanceValue := metadata.GetValue(instance)

	for i := 0; i < len(metatable.Columns); i++ {
		var value reflect.Value

		if metatable.Columns[i].AutoIncrement {
			continue
		}

		if metatable.Columns[i].IsComplex {
			var err error
			_, value, err = metadata.GetGoedbTagTypeAndValueOfIndexField(instanceType, intanceValue, "pk", i)
			if err != nil {
				return "", "", err
			}
		} else {
			value = intanceValue.Field(i)
		}

		switch metatable.Columns[i].ColumnType {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int, reflect.Int64, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint, reflect.Uint64:
			strValues += strconv.FormatInt(value.Int(), 10) + ","
		case reflect.Float32, reflect.Float64:
			strValues += strconv.FormatFloat(value.Float(), 'f', 6, 64) + ","
		case reflect.Bool:
			if value.Bool() {
				strValues += "1,"
			} else {
				strValues += "0,"
			}
		case reflect.String:
			strValues += "'" + value.String() + "',"
		}
		strCols += metatable.Columns[i].Title + ","
	}

	return strCols[:len(strCols)-1], strValues[:len(strValues)-1], nil
}

func getPKs(gt metadata.GoedbTable, obj interface{}) (string, string, error) {
	val := reflect.ValueOf(obj)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < len(gt.Columns); i++ {
		v := val.Field(i)
		if gt.Columns[i].PrimaryKey {
			switch gt.Columns[i].ColumnType {
			case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int, reflect.Int64, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint, reflect.Uint64:
				return gt.Columns[i].Title, strconv.FormatInt(v.Int(), 10), nil
			case reflect.Float32, reflect.Float64:
				return gt.Columns[i].Title, strconv.FormatFloat(v.Float(), 'f', 6, 64), nil
			case reflect.Bool:
				if v.Bool() {
					return gt.Columns[i].Title, "1", nil
				}
				return gt.Columns[i].Title, "0", nil
			case reflect.String:
				return gt.Columns[i].Title, "'" + v.String() + "'", nil
			}
		}
	}
	return "", "", errors.New("No PK found")
}

func referenceSQLEntity(from *string, query *string, constraints *string, referencedTable metadata.GoedbTable) {
	*from += referencedTable.Name + ","
	for _, referencedColumn := range referencedTable.Columns {

		if !referencedColumn.IsComplex {
			*query += referencedTable.Name + "." + referencedColumn.Title + ","
			continue
		}
		*constraints += " AND " + referencedTable.Name + "." + referencedColumn.Title + " = " + metadata.Models[referencedColumn.ColumnTypeName].Name + "." + metadata.Models[referencedColumn.ColumnTypeName].PrimaryKeyName
		referenceSQLEntity(from, query, constraints, metadata.Models[referencedColumn.ColumnTypeName])
	}
}