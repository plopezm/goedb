package dialect

import (
	"errors"
	"github.com/plopezm/goedb/metadata"
	"reflect"
)

type SQLiteDialect struct {
}

func (dialect *SQLiteDialect) GetSQLCreate(table metadata.GoedbTable) string {
	columns := ""
	pksFound := ""
	constraints := ""

	for _, value := range table.Columns {
		columnModel, pksColModel, constModel, err := dialect.getSQLColumnModel(value)
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

func (dialect *SQLiteDialect) GetSQLDelete(table metadata.GoedbTable, where string, instance interface{}) (string, error) {
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

func (dialect *SQLiteDialect) GetSQLInsert(table metadata.GoedbTable, instance interface{}) (string, error) {
	columns, values, err := getColumnsAndValues(table, instance)
	if err != nil {
		return "", err
	}
	sql := "INSERT INTO " + table.Name + " (" + columns + ") values(" + values + ")"
	return sql, nil
}

func (dialect *SQLiteDialect) GetFirstQuery(table metadata.GoedbTable, where string, instance interface{}) (string, error) {
	sql, relationContraints := dialect.generateSQLQuery(table)
	if where == "" {
		pkc, pkv, err := getPKs(table, instance)
		if err != nil {
			return "", errors.New("Error getting primary key")
		}
		sql += " WHERE " + table.Name + "." + pkc + "=" + pkv
	} else {
		sql += " WHERE " + where
	}
	//contraints are generated by relations between objects
	sql += relationContraints

	return sql, nil
}

func (dialect *SQLiteDialect) GetFindQuery(table metadata.GoedbTable, where string) (string, error) {
	//SQL generated by entity
	sql, relationContraints := dialect.generateSQLQuery(table)

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

func (dialect *SQLiteDialect) GetDropTableQuery(table metadata.GoedbTable) string {
	return "DROP TABLE " + table.Name
}

func (dialect *SQLiteDialect) generateSQLQuery(model metadata.GoedbTable) (query string, constraints string) {
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

func (dialect *SQLiteDialect) getSQLColumnModel(value metadata.GoedbColumn) (string, string, string, error) {
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
