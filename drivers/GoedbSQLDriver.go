package drivers

import (
	"database/sql"
	"reflect"
	"errors"
	"strconv"
	"strings"
)

type GoedbSQLDriver struct{
	db     *sql.DB
	tables map[string]GoedbTable
}

func init(){

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
		sqld.tables = make(map[string]GoedbTable)
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
	table := parseModel(i)
	sqld.tables[table.Name] = table
	sqltab := getSQLTableModel(table)
	_, err := sqld.db.Exec(sqltab)
	return err
}

func (sqld *GoedbSQLDriver) Model(i interface{}) (GoedbTable, error){
	var q GoedbTable
	if q, ok := sqld.tables[getType(i).Name()]; ok{
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

	var sql string
	if where == "" {
		pkc, pkv, err := getPKs(model, i)
		if err != nil {
			return errors.New("Error getting primary key")
		}
		sql = "SELECT * FROM " + model.Name + " WHERE " + pkc + "=" + pkv
	}else{
		sql = "SELECT * FROM " + model.Name + " WHERE " + where
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

func (sqld *GoedbSQLDriver) Find(i interface{}, where string) error{
	model,err := sqld.Model(i)
	if err != nil {
		return err
	}

	var sql string

	if where == "" {
		sql = "SELECT * FROM " + model.Name
	}else{
		sql = "SELECT * FROM " + model.Name + " WHERE " + where
	}

	rows, err := sqld.db.Query(sql)
	if err != nil {
		return err
	}

	slicePtr := reflect.ValueOf(i)
	slice := reflect.Indirect(slicePtr)

	slType := getType(i)

	if !rows.Next() {
		return errors.New("Records not found")
	}

	for {
		ptr := reflect.New(slType)

		valuePtrs := structToSliceOfFieldAddress(ptr)
		rows.Scan(valuePtrs...)

		slice.Set(reflect.Append(slice, ptr.Elem()))

		if !rows.Next() { break }
	}

	return nil
}

func (sqld *GoedbSQLDriver) DropTable(i interface{}) error{
	typ := getType(i)
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

func getType(i interface{}) (reflect.Type){
	typ := reflect.TypeOf(i)

	// if a pointer to a struct is passed, get the type of the dereferenced object
	if typ.Kind() == reflect.Ptr{
		typ = typ.Elem()
	}

	if typ.Kind() == reflect.Slice{
		typ = typ.Elem()
	}

	return typ
}

func getSQLColumnModel(value GoedbColumn) (string, string, string, error){
	var pksFound string
	var constraints string
	column := value.Title

	switch value.Ctype {
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
		return "","","",errors.New("Type unkown")
	}

	if value.Unique {
		column += " UNIQUE"
	}

	if value.Autoinc {
		column += " AUTOINCREMENT"
	}

	if value.Pk {
		pksFound += value.Title+","
	}

	if value.Fk {
		constraints += ", FOREIGN KEY ("+value.Title +") REFERENCES "+value.Fkref +" ON DELETE CASCADE"
	}
	column += ","

	return column, pksFound, constraints, nil
}

func getSQLTableModel(table GoedbTable) (string){
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

func getPKs(gt GoedbTable, obj interface{}) (string, string, error){
	val := reflect.ValueOf(obj)

	if val.Kind() == reflect.Ptr{
		val = val.Elem()
	}

	for i:=0;i<len(gt.Columns);i++ {
		v := val.Field(i)
		if gt.Columns[i].Pk {
			switch gt.Columns[i].Ctype {
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

func structToSliceOfFieldAddress(s interface{}) []interface{} {

	var fieldArr reflect.Value
	if _, ok  := s.(reflect.Value); ok{
		fieldArr = s.(reflect.Value)
	}else{
		fieldArr = reflect.ValueOf(s).Elem()
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

func getColumnsAndValues(gt GoedbTable, obj interface{}) (string, string){
	strCols := ""
	strValues := ""

	val := reflect.ValueOf(obj)

	if val.Kind() == reflect.Ptr{
		val = val.Elem()
	}

	for i:=0;i<len(gt.Columns);i++ {
		v := val.Field(i)
		switch gt.Columns[i].Ctype {
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

func parseModel(model interface{}) (GoedbTable){
	typ := reflect.TypeOf(model)

	// if a pointer to a struct is passed, get the type of the dereferenced object
	if typ.Kind() == reflect.Ptr{
		typ = typ.Elem()
	}

	table := GoedbTable{}
	table.Name = typ.Name()
	table.Columns = make([]GoedbColumn, 0)

	for i:=0;i<typ.NumField();i++ {
		tablecol := GoedbColumn{}
		tablecol.Title = typ.Field(i).Name
		tablecol.Ctype = typ.Field(i).Type.Name()

		if tag, ok := typ.Field(i).Tag.Lookup("goedb"); ok {
			params := strings.Split(tag, ",")
			for _, val := range params {
				switch val {
				case "pk":
					tablecol.Pk = true
				case "autoincrement":
					tablecol.Autoinc = true
				case "unique":
					tablecol.Unique = true
				default:
					if strings.Contains(val, "fk=") {
						tablecol.Fk = true
						tablecol.Fkref = strings.Split(val, "=")[1]
					}
				}
			}
		}
		table.Columns = append(table.Columns, tablecol)
	}

	return table
}