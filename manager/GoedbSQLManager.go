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
	metadata.Models[table.Name] = table
	sqltab := getSQLTableModel(table)
	print("CREATE: "+sqltab+"\n")
	_, err := sqld.db.Exec(sqltab)
	return err
}

func (sqld *GoedbSQLDriver) Model(i interface{}) (metadata.GoedbTable, error){
	var q metadata.GoedbTable
	if q, ok := metadata.Models[metadata.GetType(i).Name()]; ok{
		return q, nil
	}
	return q, errors.New("Model not found")
}

func (sqld *GoedbSQLDriver) Insert(instance interface{})(GoedbResult, error){
	var result sql.Result
	var goedbres GoedbResult

	model,err := sqld.Model(instance)
	if err != nil {
		return goedbres, err
	}
	columns, values, err := getColumnsAndValues(model, instance)
	if err != nil {
		return GoedbResult{}, err
	}
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


func referenceSQLEntity(from *string, query *string, constraints *string, referencedTable metadata.GoedbTable){
	*from += referencedTable.Name+","
	for _, referencedColumn := range referencedTable.Columns{

		if !referencedColumn.IsComplex {
			*query += referencedTable.Name+"."+referencedColumn.Title + ","
			continue
		}
		*constraints += " AND " + referencedTable.Name+"."+referencedColumn.Title + " = " + metadata.Models[referencedColumn.ColumnTypeName].Name+"."+metadata.Models[referencedColumn.ColumnTypeName].PrimaryKeyName
		referenceSQLEntity(from, query, constraints, metadata.Models[referencedColumn.ColumnTypeName])
	}
}

func generateSQLQuery(model metadata.GoedbTable) (query string, constraints string){
	query = "SELECT "
	from := " FROM "+model.Name+","
	constraints = ""

	for _, column := range model.Columns {
		if !column.IsComplex {
			query += model.Name+"."+column.Title + ","
			continue
		}
		referencedTable := metadata.Models[column.ColumnTypeName]
		constraints += model.Name+"."+column.Title + " = " + referencedTable.Name+"."+referencedTable.PrimaryKeyName
		referenceSQLEntity(&from, &query, &constraints, referencedTable)
	}
	//Removing the last ','
	query = query[:len(query)-1] + from[:len(from)-1] + " WHERE "
	return query, constraints
}

func (sqld *GoedbSQLDriver) First(instance interface{}, where string) (error){
	model,err := sqld.Model(instance)
	if err != nil {
		return err
	}
	sql, constraints := generateSQLQuery(model)
	if where == "" {
		pkc, pkv, err := getPKs(model, instance)
		if err != nil {
			return errors.New("Error getting primary key")
		}
		sql += pkc + "=" + pkv
	}else{
		sql += where
	}

	sql += " AND " + constraints

	println("[First method]: QUERY: "+sql)
	rows, err := sqld.db.Query(sql)
	if err != nil{
		return err
	}
	defer rows.Close()

	instanceValuesAddresses := metadata.StructToSliceOfAddresses(instance)

	if !rows.Next() {
		return errors.New("Record not found")
	}

	rows.Scan(instanceValuesAddresses...)
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

		entityFieldsAsSlice := metadata.StructToSliceOfAddresses(entityPtr)
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
	delete(metadata.Models, name)
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
	case  reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int,  reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint:
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
			case  reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int, reflect.Int64, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint, reflect.Uint64:
				return gt.Columns[i].Title, strconv.FormatInt(v.Int(), 10), nil
			case reflect.Float32, reflect.Float64:
				return gt.Columns[i].Title, strconv.FormatFloat(v.Float(), 'f', 6, 64), nil
			case reflect.Bool:
				if v.Bool() {
					return gt.Columns[i].Title, "1", nil
				}else{
					return gt.Columns[i].Title, "0", nil
				}
			case reflect.String:
				return gt.Columns[i].Title, "'"+v.String()+"'", nil
			}
		}
	}

	return "", "", errors.New("No PK found")
}

/*
	Returns columns names and values for inserting values
 */
func getColumnsAndValues(metatable metadata.GoedbTable, instance interface{}) (string, string, error){
	strCols := ""
	strValues := ""

	instanceType := metadata.GetType(instance)
	intanceValue := metadata.GetValue(instance)

	for i:=0;i<len(metatable.Columns);i++ {
		var value reflect.Value
		if metatable.Columns[i].IsComplex {
			var err error
			_, value, err =  metadata.GetGoedbTagTypeAndValueOfIndexField(instanceType, intanceValue, "pk", i)
			if err != nil {
				return "", "", err
			}
		}else{
			value = intanceValue.Field(i)
		}

		switch metatable.Columns[i].ColumnType {
		case  reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int, reflect.Int64, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint, reflect.Uint64:
			strValues += strconv.FormatInt(value.Int(), 10)+","
		case reflect.Float32, reflect.Float64:
			strValues += strconv.FormatFloat(value.Float(), 'f', 6, 64)+","
		case reflect.Bool:
			if value.Bool() {
				strValues += "1,"
			}else{
				strValues += "0,"
			}
		case reflect.String:
			strValues += "'"+ value.String()+"',"
		}
		strCols += metatable.Columns[i].Title + ","
	}

	return strCols[:len(strCols)-1], strValues[:len(strValues)-1], nil
}



