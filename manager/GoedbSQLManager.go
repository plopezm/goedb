package manager

import (
	"database/sql"
	"reflect"
	"errors"
	"strconv"
	"goedb/metadata"
)

type GoedbSQLDriver struct{
	db     *sql.DB
	tables map[string]metadata.GoedbTable
}

func (sqld *GoedbSQLDriver) Open(driver string, params string) error{
	db, err := sql.Open(driver, params)
	if err != nil {
		return err
	}

	err = db.Ping() // Send a ping to make sure the database connection is alive.
	if err != nil {
		db.Close()
		return err
	}

	sqld.db = db
	if sqld.tables == nil {
		sqld.tables = make(map[string]metadata.GoedbTable)
	}

	if driver == "sqlite3" {
		sqld.db.Exec("PRAGMA foreign_keys = ON")
	}
	return nil

}

func (sqld *GoedbSQLDriver) Close() error{
	if sqld.db == nil {
		return errors.New("DB is closed")
	}
	return sqld.db.Close()
}

func (sqld *GoedbSQLDriver) Migrate(i interface{}) (error){
	sqld.DropTable(i)
	table := metadata.ParseModel(i)
	sqld.tables[table.Name] = table
	sqltab := getSQLTableModel(table)
	_, err := sqld.db.Exec(sqltab)
	return err
}

func (sqld *GoedbSQLDriver) Model(i interface{}) (metadata.GoedbTable, error){
	var q metadata.GoedbTable
	if q, ok := sqld.tables[metadata.GetType(i).Name()]; ok{
		return q, nil
	}
	return q, errors.New("Model not found")
}

func (sqld *GoedbSQLDriver) Insert(i interface{})(GoedbResult, error){
	var result sql.Result
	var goedbres GoedbResult

	model,err := sqld.Model(i)
	if err != nil {
		return goedbres, err
	}
	columns, values := getColumnsAndValues(model, i)
	sql := "INSERT INTO "+model.Name +" ("+columns+") values("+values+")"
	result, err = sqld.db.Exec(sql)
	if err != nil {
		return goedbres, err
	}

	goedbres.NumRecordsAffected, _ = result.RowsAffected()
	return goedbres, nil
}

func (sqld *GoedbSQLDriver) Remove(i interface{})(GoedbResult, error){
	var result sql.Result
	var goedbres GoedbResult

	model,err := sqld.Model(i)
	if err != nil {
		return goedbres, err
	}

	pkc, pkv, err := getPKs(model, i)
	if err != nil {
		return goedbres, err
	}

	sql := "DELETE FROM "+model.Name +" WHERE "+pkc+ "=" + pkv
	result, err = sqld.db.Exec(sql)
	if err != nil {
		return goedbres, err
	}

	goedbres.NumRecordsAffected, _ = result.RowsAffected()
	return goedbres, err
}

func (sqld *GoedbSQLDriver) First(i interface{}, where string) (error){
	model,err := sqld.Model(i)
	if err != nil {
		return err
	}

	sql := "SELECT * FROM " + model.Name + " WHERE "
	if where == "" {
		pkc, pkv, err := getPKs(model, i)
		if err != nil {
			return errors.New("Error getting primary key")
		}
		sql += pkc + "=" + pkv
	}else{
		sql += where
	}

	rows, err := sqld.db.Query(sql)
	if err != nil{
		return err
	}
	defer rows.Close()

	valuePtrs := structToSliceOfFieldAddress(i)

	if !rows.Next() {
		return errors.New("Record not found")
	}

	rows.Scan(valuePtrs...)
	return nil
}

func (sqld *GoedbSQLDriver) Find(resultEntitySlice interface{}, where string) error{
	model,err := sqld.Model(resultEntitySlice)
	if err != nil {
		return err
	}

	var sqlQuery string

	if where != "" {
		where = " WHERE "+ where
	}

	sqlQuery = "SELECT * FROM " + model.Name + where
	rows, err := sqld.db.Query(sqlQuery)
	if err != nil {
		return err
	}

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

		entityFieldsAsSlice := structToSliceOfFieldAddress(entityPtr)
		rows.Scan(entityFieldsAsSlice...)

		slice.Set(reflect.Append(slice, entityPtr.Elem()))

		if !rows.Next() { break }
	}

	return nil
}

func (sqld *GoedbSQLDriver) DropTable(i interface{}) error{
	typ := metadata.GetType(i)
	name := typ.Name()

	_, err := sqld.db.Exec("DROP TABLE "+name)
	if err != nil {
		return err
	}
	delete(sqld.tables, name)
	return nil
}


/* ======================================
	    Support functions
   ====================================== */

func getSQLColumnModel(value metadata.GoedbColumn) (string, string, string, error){
	var pksFound string
	var constraints string
	column := value.Title

	switch value.ColumnType {
	case "char":
		column += " CHARACTER"
	case  "int8", "int16", "int32", "int",  "uint8", "uint16", "uint32", "uint":
		column += " INTEGER"
	case "int64", "uint64":
		column += " BIGINT"
	case "float32", "float64":
		column += " FLOAT"
	case "bool":
		column += " BOOLEAN"
	case "string":
		column += " VARCHAR"
	default:
		return "","","",errors.New("Type unknown")
	}

	if value.Unique {
		column += " UNIQUE"
	}

	if value.AutoIncrement {
		column += " AUTOINCREMENT"
	}

	if value.PrimaryKey {
		pksFound += value.Title+","
	}

	if value.ForeignKey {
		constraints += ", FOREIGN KEY ("+value.Title +") REFERENCES "+value.ForeignKeyReference +" ON DELETE CASCADE"
	}
	column += ","

	return column, pksFound, constraints, nil
}

func getSQLTableModel(table metadata.GoedbTable) (string){
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
		constraints += ", PRIMARY KEY ("+ pksFound +")"
	}

	lastColumnIndex := len(columns)
	return "CREATE TABLE "+table.Name +" (" +columns[:lastColumnIndex-1] + constraints+")"
}

func getPKs(gt metadata.GoedbTable, obj interface{}) (string, string, error){
	val := reflect.ValueOf(obj)

	if val.Kind() == reflect.Ptr{
		val = val.Elem()
	}

	for i:=0;i<len(gt.Columns);i++ {
		v := val.Field(i)
		if gt.Columns[i].PrimaryKey {
			switch gt.Columns[i].ColumnType {
			case  "int8", "int16", "int32", "int",  "uint8", "uint16", "uint32", "uint", "int64", "uint64":
				return gt.Columns[i].Title, strconv.FormatInt(v.Int(), 10), nil
			case "float32", "float64":
				return gt.Columns[i].Title, strconv.FormatFloat(v.Float(), 'f', 6, 64), nil
			case "bool":
				if v.Bool() {
					return gt.Columns[i].Title, "1", nil
				}else{
					return gt.Columns[i].Title, "0", nil
				}
			case "string","char":
				return gt.Columns[i].Title, "'"+v.String()+"'", nil
			}
		}
	}

	return "", "", errors.New("No PK found")
}

/*
	Returns a slice with the addresses of each struct field,
	so any modification on the slide will modify the source struct fields
 */
func structToSliceOfFieldAddress(structPtr interface{}) []interface{} {

	var fieldArr reflect.Value
	if _, ok  := structPtr.(reflect.Value); ok{
		fieldArr = structPtr.(reflect.Value)
	}else{
		fieldArr = reflect.ValueOf(structPtr).Elem()
	}

	if fieldArr.Kind() == reflect.Ptr{
		fieldArr = fieldArr.Elem()
	}

	fieldAddrArr := make([]interface{}, fieldArr.NumField())

	for i := 0; i < fieldArr.NumField(); i++ {
		f := fieldArr.Field(i)
		fieldAddrArr[i] = f.Addr().Interface()
	}

	return fieldAddrArr
}

/*
	Returns columns names and values for inserting values
 */
func getColumnsAndValues(gt metadata.GoedbTable, obj interface{}) (string, string){
	strCols := ""
	strValues := ""

	val := reflect.ValueOf(obj)

	if val.Kind() == reflect.Ptr{
		val = val.Elem()
	}

	for i:=0;i<len(gt.Columns);i++ {
		v := val.Field(i)
		switch gt.Columns[i].ColumnType {
		case  "int8", "int16", "int32", "int",  "uint8", "uint16", "uint32", "uint", "int64", "uint64":
			strValues += strconv.FormatInt(v.Int(), 10)+","
		case "float32", "float64":
			strValues += strconv.FormatFloat(v.Float(), 'f', 6, 64)+","
		case "bool":
			if v.Bool() {
				strValues += "1,"
			}else{
				strValues += "0,"
			}
		case "string","char":
			strValues += "'"+v.String()+"',"
		}
		strCols += gt.Columns[i].Title + ","
	}

	return strCols[:len(strCols)-1], strValues[:len(strValues)-1]
}



