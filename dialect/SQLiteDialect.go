package dialect

import (
	"reflect"
	"strconv"
	"github.com/plopezm/goedb/metadata"
	"errors"
)

type SQLiteDialect struct{
}

func (dialect *SQLiteDialect) GetSQLCreate(table metadata.GoedbTable) string{
	columns := ""
	pksFound := ""
	constraints := ""

	for _, value := range table.Columns {
		columnModel, pksColModel, constModel, err := getSQLColumnModel(value)
		if err != nil {
			continue
		}
		columns += columnModel
		pksFound += pksColModel
		constraints += constModel
	}

	if len(pksFound) > 0 {
		pksFound = pksFound[:len(pksFound)-1]
		constraints += ", PRIMARY KEY (" + pksFound + ")"
	}

	lastColumnIndex := len(columns)
	return "CREATE TABLE " + table.Name + " (" + columns[:lastColumnIndex-1] + constraints + ")"
}

func (dialect *SQLiteDialect) GetSQLDelete(table metadata.GoedbTable, where string, instance interface{}) (string, error){
	sql := "DELETE FROM " + table.Name + " WHERE "
	if where == "" {
		pkc, pkv, err := getPKs(table, instance)
		if err != nil {
			return "", err
		}
		sql += pkc + "=" + pkv
	} else {
		sql += where
	}
	return sql, nil
}

func (dialect *SQLiteDialect) GetSQLInsert(table metadata.GoedbTable, instance interface{}) (string, error){
	columns, values, err := getColumnsAndValues(table, instance)
	if err != nil {
		return "", err
	}
	sql := "INSERT INTO " + table.Name + " (" + columns + ") values(" + values + ")"
	return sql, nil
}

func (dialect *SQLiteDialect) GetFirstQuery(table metadata.GoedbTable, where string, instance interface{}) (string, error){
	sql, relationContraints := generateSQLQuery(table)
	if where == "" {
		pkc, pkv, err := getPKs(table, instance)
		if err != nil {
			return "",errors.New("Error getting primary key")
		}
		sql += " WHERE " + table.Name + "." + pkc + "=" + pkv
	} else {
		sql += " WHERE " + where
	}
	//contraints are generated by relations between objects
	sql += relationContraints

	return sql, nil
}

func (dialect *SQLiteDialect) GetFindQuery(table metadata.GoedbTable, where string) (string, error){
	//SQL generated by entity
	sql, relationContraints := generateSQLQuery(table)

	if where != "" {
		//where clause
		sql += " WHERE " + where
		//contraints are generated by relations between objects
		sql += relationContraints
	} else if relationContraints != "" {
		sql += relationContraints
	}
	return sql, nil
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

func getSQLTableModel(table metadata.GoedbTable) string {
	columns := ""
	pksFound := ""
	constraints := ""

	for _, value := range table.Columns {
		columnModel, pksColModel, constModel, err := getSQLColumnModel(value)
		if err != nil {
			continue
		}
		columns += columnModel
		pksFound += pksColModel
		constraints += constModel
	}

	if len(pksFound) > 0 {
		pksFound = pksFound[:len(pksFound)-1]
		constraints += ", PRIMARY KEY (" + pksFound + ")"
	}

	lastColumnIndex := len(columns)
	return "CREATE TABLE " + table.Name + " (" + columns[:lastColumnIndex-1] + constraints + ")"
}

func getSQLColumnModel(value metadata.GoedbColumn) (string, string, string, error) {
	var pksFound string
	var constraints string
	column := value.Title

	switch value.ColumnType {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint:
		column += " INTEGER"
	case reflect.Int64, reflect.Uint64:
		column += " BIGINT"
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

	if value.PrimaryKey && value.AutoIncrement {
		column += " PRIMARY KEY AUTOINCREMENT"
	} else if value.PrimaryKey {
		pksFound += value.Title + ","
	}

	if value.ForeignKey {
		constraints += ", FOREIGN KEY (" + value.Title + ") REFERENCES " + value.ForeignKeyReference + " ON DELETE CASCADE"
	}
	column += ","
	return column, pksFound, constraints, nil
}

func generateSQLQuery(model metadata.GoedbTable) (query string, constraints string) {
	query = "SELECT "
	from := " FROM " + model.Name + ","
	constraints = ""

	for _, column := range model.Columns {
		if !column.IsComplex {
			query += model.Name + "." + column.Title + ","
			continue
		}
		referencedTable := metadata.Models[column.ColumnTypeName]
		constraints += " AND " + model.Name + "." + column.Title + " = " + referencedTable.Name + "." + referencedTable.PrimaryKeyName
		referenceSQLEntity(&from, &query, &constraints, referencedTable)
	}
	//Removing the last ','
	query = query[:len(query)-1] + from[:len(from)-1]
	return query, constraints
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
