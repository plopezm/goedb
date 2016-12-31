package goedb

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"reflect"
	"strings"
	"errors"
	"strconv"
)

type DB struct{
	db 	*sql.DB
	tables	map[string]GoedbTable
}


func NewGoeDB() (*DB){
	gdb := new(DB)
	gdb.tables = make(map[string]GoedbTable)
	return gdb
}

func (gdb *DB)Open(driver string, params string) (error){
	db, err := sql.Open(driver, params)
	if err != nil {
		return err
	}

	err = db.Ping() // Send a ping to make sure the database connection is alive.
	if err != nil {
		db.Close()
		return err
	}
	gdb.db = db
	return nil
}

func (gdb *DB) Close(){
	gdb.db.Close()
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

	for i:=0;i<len(gt.columns);i++ {
		v := val.Field(i)
		switch gt.columns[i].ctype {
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
		strCols += gt.columns[i].title + ","

	}

	return strCols[:len(strCols)-1], strValues[:len(strValues)-1]
}

func getPKs(gt GoedbTable, obj interface{}) (string, string, error){
	val := reflect.ValueOf(obj)

	if val.Kind() == reflect.Ptr{
		val = val.Elem()
	}

	for i:=0;i<len(gt.columns);i++ {
		v := val.Field(i)
		if gt.columns[i].pk {
			switch gt.columns[i].ctype {
			case  "int8", "int16", "int32", "int",  "uint8", "uint16", "uint32", "uint", "int64", "uint64":
				return gt.columns[i].title, strconv.FormatInt(v.Int(), 10), nil
			case "float32", "float64":
				return gt.columns[i].title, strconv.FormatFloat(v.Float(), 'f', 6, 64), nil
			case "bool":
				if v.Bool() {
					return gt.columns[i].title, "1", nil
				}else{
					return gt.columns[i].title, "0", nil
				}
			case "string","char":
				return gt.columns[i].title, "'"+v.String()+"'", nil
			}
		}
	}

	return "", "", errors.New("No PK found")
}

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

func parseModel(model interface{}) (GoedbTable){
	typ := reflect.TypeOf(model)

	// if a pointer to a struct is passed, get the type of the dereferenced object
	if typ.Kind() == reflect.Ptr{
		typ = typ.Elem()
	}

	table := GoedbTable{}
	table.name = typ.Name()
	table.columns = make([]GoedbColumn, 0)

	for i:=0;i<typ.NumField();i++ {
		tablecol := GoedbColumn{}
		tablecol.title = typ.Field(i).Name
		tablecol.ctype = typ.Field(i).Type.Name()

		if tag, ok := typ.Field(i).Tag.Lookup("goedb"); ok {
			params := strings.Split(tag, ",")
			for _, val := range params {
				switch val {
					case "pk":
						tablecol.pk = true
					case "autoincrement":
						tablecol.autoinc = true
					case "unique":
						tablecol.unique = true
					default:
						if strings.Contains(val, "fk=") {
							tablecol.fk = true
							tablecol.fkref = strings.Split(val, "=")[1]
						}
				}
			}
		}
		table.columns = append(table.columns, tablecol)
	}

	return table
}

func getSQLTableModel(table GoedbTable) (string){
	pksFound := ""
	columns := ""
	constraints := ""


	for _, value := range table.columns {
		columns += value.title

		switch value.ctype {
		case "char":
			columns += " CHARACTER"
		case  "int8", "int16", "int32", "int",  "uint8", "uint16", "uint32", "uint":
			columns += " INTEGER"
		case "int64", "uint64":
			columns += " BIGINT"
		case "float32", "float64":
			columns += " FLOAT"
		case "bool":
			columns += " BOOLEAN"
		case "string":
			columns += " VARCHAR"
		default:
			continue
		}

		if value.unique {
			columns += " UNIQUE"
		}

		if value.autoinc {
			columns += " AUTOINCREMENT"
		}

		if value.pk {
			pksFound += value.title+","
		}

		if value.fk {
			constraints += ", FOREIGN KEY ("+value.title +") REFERENCES "+value.fkref +" ON DELETE CASCADE"
		}
		columns += ","
	}
	if len(pksFound) > 0 {
		pksFound = pksFound[:len(pksFound)-1]
		constraints += ", PRIMARY KEY ("+ pksFound +")"
	}

	lastColumnIndex := len(columns)
	return "CREATE TABLE "+table.name +" (" +columns[:lastColumnIndex-1] + constraints+")"
}

func (gdb *DB) Migrate(i interface{}) (error){
	gdb.DropTable(i)
	table := parseModel(i)
	gdb.tables[table.name] = table
	sqltab := getSQLTableModel(table)
	_, err := gdb.db.Exec(sqltab)
	return err
}

func (gdb *DB) DropTable(i interface{}){
	typ := getType(i)

	gdb.db.Exec("DROP TABLE "+typ.Name())
}

func (gdb *DB) Model(i interface{}) (GoedbTable, error){
	var q GoedbTable
	if q, ok := gdb.tables[getType(i).Name()]; ok{
		return q, nil
	}
	return q, errors.New("Model not found")
}

func (gdb *DB) Insert(i interface{})(sql.Result, error){
	var result sql.Result

	model,err := gdb.Model(i)
	if err != nil {
		return result, err
	}

	columns, values := getColumnsAndValues(model, i)
	sql := "INSERT INTO "+model.name+" ("+columns+") values("+values+")"
	return gdb.db.Exec(sql)
}


func (gdb *DB) Remove(i interface{})(sql.Result, error){
	var result sql.Result

	model,err := gdb.Model(i)
	if err != nil {
		return result, err
	}

	pkc, pkv, err := getPKs(model, i)
	if err != nil {
		return nil,err
	}

	sql := "DELETE FROM "+model.name+" WHERE "+pkc+ "=" + pkv
	return gdb.db.Exec(sql)
}

func (gdb *DB) First(i interface{}, where string) (error){
	model,err := gdb.Model(i)
	if err != nil {
		return err
	}

	var sql string
	if where == "" {
		pkc, pkv, err := getPKs(model, i)
		if err != nil {
			return errors.New("Error getting primary key")
		}
		sql = "SELECT * FROM " + model.name + " WHERE " + pkc + "=" + pkv
	}else{
		sql = "SELECT * FROM " + model.name + " WHERE " + where
	}

	rows, err := gdb.db.Query(sql)
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

func (gdb *DB) Find(i interface{}, where string) error{
	model,err := gdb.Model(i)
	if err != nil {
		return err
	}

	var sql string

	if where == "" {
		sql = "SELECT * FROM " + model.name
	}else{
		sql = "SELECT * FROM " + model.name + " WHERE " + where
	}

	rows, err := gdb.db.Query(sql)
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