package manager

import (
	"database/sql"
	"errors"
	"github.com/plopezm/goedb/metadata"
	"reflect"
	"github.com/jmoiron/sqlx"
	"github.com/plopezm/goedb/dialect"
)

// GoedbSQLDriver constains the database connection
type GoedbSQLDriver struct {
	db *sqlx.DB
	Dialect dialect.Dialect
}

// Open creates the connection with the database
// **DON'T open a connection**
// This will be managed by goedb
func (sqld *GoedbSQLDriver) Open(driver string, params string) error {
	db, err := sqlx.Connect(driver, params)
	if err != nil {
		return err
	}

	err = db.Ping() // Send a ping to make sure the database connection is alive.
	if err != nil {
		db.Close()
		return err
	}

	sqld.db = db

	if driver == "sqlite3" {
		sqld.db.Exec("PRAGMA foreign_keys = ON")
	}
	return nil

}

// Close finishes the connection
func (sqld *GoedbSQLDriver) Close() error {
	if sqld.db == nil {
		return errors.New("DB is closed")
	}
	return sqld.db.Close()
}

// Migrate creates the table in the database
func (sqld *GoedbSQLDriver) Migrate(i interface{}) error {
	sqld.DropTable(i)
	table := metadata.ParseModel(i)
	metadata.Models[table.Name] = table
	sqltab := sqld.Dialect.GetSQLCreate(table)
	_, err := sqld.db.Exec(sqltab)
	return err
}

// Model returns the metadata of each structure migrated
func (sqld *GoedbSQLDriver) Model(i interface{}) (metadata.GoedbTable, error) {
	var q metadata.GoedbTable
	if q, ok := metadata.Models[metadata.GetType(i).Name()]; ok {
		return q, nil
	}
	return q, errors.New("Model not found")
}

// Insert creates a new row with the object in the database (it must be migrated)
func (sqld *GoedbSQLDriver) Insert(instance interface{}) (GoedbResult, error) {
	var goedbres GoedbResult
	var result sql.Result
	model, err := sqld.Model(instance)
	if err != nil {
		return goedbres, err
	}
	//columns, values, err := getColumnsAndValues(model, instance)
	//if err != nil {
	//	return GoedbResult{}, err
	//}
	//sql := "INSERT INTO " + model.Name + " (" + columns + ") values(" + values + ")"

	sql, err := sqld.Dialect.GetSQLInsert(model, instance)
	if err != nil{
		return goedbres, err
	}
	result, err = sqld.db.Exec(sql)
	if err != nil {
		return goedbres, err
	}

	goedbres.NumRecordsAffected, _ = result.RowsAffected()
	return goedbres, nil
}

// Remove removes a row with the object in the database (it must be migrated)
func (sqld *GoedbSQLDriver) Remove(i interface{}, where string, params map[string]interface{}) (GoedbResult, error) {
	var result sql.Result
	var goedbres GoedbResult

	model, err := sqld.Model(i)
	if err != nil {
		return goedbres, err
	}

	//sql := "DELETE FROM " + model.Name + " WHERE "
	//if where == "" {
	//	pkc, pkv, err := getPKs(model, i)
	//	if err != nil {
	//		return goedbres, err
	//	}
	//	sql += pkc + "=" + pkv
	//} else {
	//	sql += where
	//}

	sql, err := sqld.Dialect.GetSQLDelete(model, where, i)
	if err != nil {
		return goedbres,err
	}

	//stmt, err := sqld.db.Prepare(sql)
	//if err != nil {
	//	return goedbres, err
	//}
	//defer stmt.Close()
	//
	//result, err = stmt.Exec(params...)
	//if err != nil {
	//	return goedbres, err
	//}
	result, err = sqld.db.NamedExec(sql, params)
	goedbres.NumRecordsAffected, _ = result.RowsAffected()
	return goedbres, err
}

//func referenceSQLEntity(from *string, query *string, constraints *string, referencedTable metadata.GoedbTable) {
//	*from += referencedTable.Name + ","
//	for _, referencedColumn := range referencedTable.Columns {
//
//		if !referencedColumn.IsComplex {
//			*query += referencedTable.Name + "." + referencedColumn.Title + ","
//			continue
//		}
//		*constraints += " AND " + referencedTable.Name + "." + referencedColumn.Title + " = " + metadata.Models[referencedColumn.ColumnTypeName].Name + "." + metadata.Models[referencedColumn.ColumnTypeName].PrimaryKeyName
//		referenceSQLEntity(from, query, constraints, metadata.Models[referencedColumn.ColumnTypeName])
//	}
//}

//func generateSQLQuery(model metadata.GoedbTable) (query string, constraints string) {
//	query = "SELECT "
//	from := " FROM " + model.Name + ","
//	constraints = ""
//
//	for _, column := range model.Columns {
//		if !column.IsComplex {
//			query += model.Name + "." + column.Title + ","
//			continue
//		}
//		referencedTable := metadata.Models[column.ColumnTypeName]
//		constraints += " AND " + model.Name + "." + column.Title + " = " + referencedTable.Name + "." + referencedTable.PrimaryKeyName
//		referenceSQLEntity(&from, &query, &constraints, referencedTable)
//	}
//	//Removing the last ','
//	query = query[:len(query)-1] + from[:len(from)-1]
//	return query, constraints
//}

// First returns the first record found
func (sqld *GoedbSQLDriver) First(instance interface{}, where string, params map[string]interface{}) error {
	model, err := sqld.Model(instance)
	if err != nil {
		return err
	}
	//sql, relationContraints := generateSQLQuery(model)
	//if where == "" {
	//	pkc, pkv, err := getPKs(model, instance)
	//	if err != nil {
	//		return errors.New("Error getting primary key")
	//	}
	//	sql += " WHERE " + model.Name + "." + pkc + "=" + pkv
	//} else {
	//	sql += " WHERE " + where
	//}
	//
	////contraints are generated by relations between objects
	//sql += relationContraints

	//stmt, err := sqld.db.Prepare(sql)
	//if err != nil {
	//	return err
	//}
	//defer stmt.Close()
	//
	//row := stmt.QueryRow(params...)
	//
	//instanceValuesAddresses := metadata.StructToSliceOfAddresses(instance)
	//
	//err = row.Scan(instanceValuesAddresses...)
	//return err

	sql, err := sqld.Dialect.GetFirstQuery(model, where, instance)
	if err != nil {
		return err
	}
	rows, err := sqld.db.NamedQuery(sql, params)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		instanceValuesAddresses := metadata.StructToSliceOfAddresses(instance)
		err = rows.Scan(instanceValuesAddresses...)
	}else{
		err = errors.New("Not found")
	}
	return err
}

// Find returns all records found
func (sqld *GoedbSQLDriver) Find(resultEntitySlice interface{}, where string, params map[string]interface{}) error {
	model, err := sqld.Model(resultEntitySlice)
	if err != nil {
		return err
	}
	////SQL generated by entity
	//sql, relationContraints := generateSQLQuery(model)
	//
	//if where != "" {
	//	//where clause
	//	sql += " WHERE " + where
	//	//contraints are generated by relations between objects
	//	sql += relationContraints
	//} else if relationContraints != "" {
	//	sql += relationContraints
	//}

	sql, err := sqld.Dialect.GetFindQuery(model, where)
	if err != nil {
		return err
	}

	rows, err := sqld.db.NamedQuery(sql, params)
	if err != nil {
		return err
	}
	defer rows.Close()

	//Creates a new pointer with the same type that resultEntitySlice
	slicePtr := reflect.ValueOf(resultEntitySlice)
	//it gets the value of the slice pointer
	slice := reflect.Indirect(slicePtr)

	entityType := metadata.GetType(resultEntitySlice)

	if !rows.Next() {
		return errors.New("Records not found")
	}

	for {
		entityPtr := reflect.New(entityType)

		entityFieldsAsSlice := metadata.StructToSliceOfAddresses(entityPtr)
		rows.Scan(entityFieldsAsSlice...)

		slice.Set(reflect.Append(slice, entityPtr.Elem()))

		if !rows.Next() {
			break
		}
	}

	return nil
}

// DropTable removes a table from the database
func (sqld *GoedbSQLDriver) DropTable(i interface{}) error {
	typ := metadata.GetType(i)
	name := typ.Name()

	table, err := sqld.Model(i)
	if err != nil {
		return err
	}

	sql := sqld.Dialect.GetDropTableQuery(table)

	_, err = sqld.db.Exec(sql)
	if err != nil {
		return err
	}
	delete(metadata.Models, name)
	return nil
}

// TxBegin is used to set a transaction
func (sqld *GoedbSQLDriver) TxBegin() (*sql.Tx, error) {
	return sqld.db.Begin()
}

/* ======================================
	    Support functions
   ====================================== */

//func getSQLColumnModel(value metadata.GoedbColumn) (string, string, string, error) {
//	var pksFound string
//	var constraints string
//	column := value.Title
//
//	switch value.ColumnType {
//	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint:
//		column += " INTEGER"
//	case reflect.Int64, reflect.Uint64:
//		column += " BIGINT"
//	case reflect.Float32, reflect.Float64:
//		column += " FLOAT"
//	case reflect.Bool:
//		column += " BOOLEAN"
//	case reflect.String:
//		column += " VARCHAR"
//	default:
//		return "", "", "", errors.New("Type unknown")
//	}
//
//	if value.Unique {
//		column += " UNIQUE"
//	}
//
//	if value.PrimaryKey && value.AutoIncrement {
//		column += " PRIMARY KEY AUTOINCREMENT"
//	} else if value.PrimaryKey {
//		pksFound += value.Title + ","
//	}
//
//	if value.ForeignKey {
//		constraints += ", FOREIGN KEY (" + value.Title + ") REFERENCES " + value.ForeignKeyReference + " ON DELETE CASCADE"
//	}
//	column += ","
//	return column, pksFound, constraints, nil
//}

//func getSQLTableModel(table metadata.GoedbTable) string {
//	columns := ""
//	pksFound := ""
//	constraints := ""
//
//	for _, value := range table.Columns {
//		columnModel, pksColModel, constModel, err := getSQLColumnModel(value)
//		if err != nil {
//			continue
//		}
//		columns += columnModel
//		pksFound += pksColModel
//		constraints += constModel
//	}
//
//	if len(pksFound) > 0 {
//		pksFound = pksFound[:len(pksFound)-1]
//		constraints += ", PRIMARY KEY (" + pksFound + ")"
//	}
//
//	lastColumnIndex := len(columns)
//	return "CREATE TABLE " + table.Name + " (" + columns[:lastColumnIndex-1] + constraints + ")"
//}

//func getPKs(gt metadata.GoedbTable, obj interface{}) (string, string, error) {
//	val := reflect.ValueOf(obj)
//
//	if val.Kind() == reflect.Ptr {
//		val = val.Elem()
//	}
//
//	for i := 0; i < len(gt.Columns); i++ {
//		v := val.Field(i)
//		if gt.Columns[i].PrimaryKey {
//			switch gt.Columns[i].ColumnType {
//			case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int, reflect.Int64, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint, reflect.Uint64:
//				return gt.Columns[i].Title, strconv.FormatInt(v.Int(), 10), nil
//			case reflect.Float32, reflect.Float64:
//				return gt.Columns[i].Title, strconv.FormatFloat(v.Float(), 'f', 6, 64), nil
//			case reflect.Bool:
//				if v.Bool() {
//					return gt.Columns[i].Title, "1", nil
//				}
//				return gt.Columns[i].Title, "0", nil
//			case reflect.String:
//				return gt.Columns[i].Title, "'" + v.String() + "'", nil
//			}
//		}
//	}
//
//	return "", "", errors.New("No PK found")
//}

///*
//	Returns columns names and values for inserting values
//*/
//func getColumnsAndValues(metatable metadata.GoedbTable, instance interface{}) (string, string, error) {
//	strCols := ""
//	strValues := ""
//
//	instanceType := metadata.GetType(instance)
//	intanceValue := metadata.GetValue(instance)
//
//	for i := 0; i < len(metatable.Columns); i++ {
//		var value reflect.Value
//
//		if metatable.Columns[i].AutoIncrement {
//			continue
//		}
//
//		if metatable.Columns[i].IsComplex {
//			var err error
//			_, value, err = metadata.GetGoedbTagTypeAndValueOfIndexField(instanceType, intanceValue, "pk", i)
//			if err != nil {
//				return "", "", err
//			}
//		} else {
//			value = intanceValue.Field(i)
//		}
//
//		switch metatable.Columns[i].ColumnType {
//		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int, reflect.Int64, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint, reflect.Uint64:
//			strValues += strconv.FormatInt(value.Int(), 10) + ","
//		case reflect.Float32, reflect.Float64:
//			strValues += strconv.FormatFloat(value.Float(), 'f', 6, 64) + ","
//		case reflect.Bool:
//			if value.Bool() {
//				strValues += "1,"
//			} else {
//				strValues += "0,"
//			}
//		case reflect.String:
//			strValues += "'" + value.String() + "',"
//		}
//		strCols += metatable.Columns[i].Title + ","
//	}
//
//	return strCols[:len(strCols)-1], strValues[:len(strValues)-1], nil
//}
