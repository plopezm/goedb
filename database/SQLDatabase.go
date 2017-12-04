package database

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"

	"github.com/jmoiron/sqlx"
	"github.com/plopezm/goedb/config"
	"github.com/plopezm/goedb/database/dialect"
	"github.com/plopezm/goedb/database/models"
)

//SQLDatabase is the implementation of SQL for a Database interface
type SQLDatabase struct {
	db         *sqlx.DB
	Dialect    dialect.Dialect
	Datasource config.Datasource
}

// SetSchema sets the schema as default schema for a datasource
func (sqld *SQLDatabase) SetSchema(schema string) (sql.Result, error) {
	sql := "SET search_path TO " + schema
	return sqld.db.Exec(sql)
}

// Open creates the connection with the database
// **DON'T open a connection**
// This will be managed by goedb
func (sqld *SQLDatabase) Open(driver string, params string, schema string) error {
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
	if len(schema) > 0 {
		sqld.SetSchema(schema)
	}
	return nil
}

// Close finishes the connection
func (sqld *SQLDatabase) Close() error {
	if sqld.db == nil {
		return errors.New("DB is closed")
	}
	return sqld.db.Close()
}

// GetDBConnection returns the DB connection as *sqlx.DB.
// This method can be used if you wanna perform some query manually
func (sqld *SQLDatabase) GetDBConnection() *sqlx.DB {
	return sqld.db
}

// Model returns the metadata of each structure migrated
func (sqld *SQLDatabase) Model(i interface{}) (models.Table, error) {
	var table models.Table
	if table, ok := sqld.Dialect.GetModel(models.GetType(i).Name()); ok {
		return table, nil
	}
	return table, errors.New("Model not found")
}

// Migrate creates the table in the database
func (sqld *SQLDatabase) Migrate(i interface{}, autoCreate bool, dropIfExists bool) (err error) {
	if dropIfExists {
		sqld.DropTable(i)
	}
	table := models.ParseModel(i)
	sqld.Dialect.SetModel(table.Name, table)
	if autoCreate {
		sqltab := sqld.Dialect.Create(table)
		_, err = sqld.db.Exec(sqltab)
	}
	return err
}

// Insert creates a new row with the object in the database (it must be migrated)
func (sqld *SQLDatabase) Insert(instance interface{}) (goedbres models.Result, err error) {
	var result sql.Result
	model, err := sqld.Model(instance)
	if err != nil {
		return goedbres, err
	}

	sql, err := sqld.Dialect.Insert(model, instance)
	if err != nil {
		return goedbres, err
	}
	result, err = sqld.db.Exec(sql)
	if err != nil {
		return goedbres, err
	}

	goedbres.NumRecordsAffected, _ = result.RowsAffected()
	goedbres.LastInsertId, _ = result.LastInsertId()
	return goedbres, nil
}

// Update updates an object using its primery key
func (sqld *SQLDatabase) Update(instance interface{}) (goedbres models.Result, err error) {
	var result sql.Result
	model, err := sqld.Model(instance)
	if err != nil {
		return goedbres, err
	}

	sql, err := sqld.Dialect.Update(model, instance)
	fmt.Println("UPDATE SQL: ", sql)
	if err != nil {
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
func (sqld *SQLDatabase) Remove(i interface{}, where string, params map[string]interface{}) (goedbres models.Result, err error) {
	model, err := sqld.Model(i)
	if err != nil {
		return goedbres, err
	}

	sql, err := sqld.Dialect.Delete(model, where, i)
	if err != nil {
		return goedbres, err
	}

	result, err := sqld.db.NamedExec(sql, params)
	goedbres.NumRecordsAffected, _ = result.RowsAffected()
	return goedbres, err
}

// First returns the first record found
func (sqld *SQLDatabase) First(instance interface{}, where string, params map[string]interface{}) error {
	model, err := sqld.Model(instance)
	if err != nil {
		return err
	}

	sql, err := sqld.Dialect.First(model, where, instance)
	if err != nil {
		return err
	}
	rows, err := sqld.db.NamedQuery(sql, params)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		instanceValuesAddresses := models.StructToSliceOfAddressesWithRules(instance, sqld.Dialect.GetModel)
		err = rows.Scan(instanceValuesAddresses...)
	} else {
		err = errors.New("Not found")
	}
	return err
}

// NativeFirst returns the first record found
func (sqld *SQLDatabase) NativeFirst(instance interface{}, sql string, params map[string]interface{}) error {
	rows, err := sqld.db.NamedQuery(sql, params)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		instanceValuesAddresses := models.StructToSliceOfAddresses(instance)
		err = rows.Scan(instanceValuesAddresses...)
	} else {
		err = errors.New("Not found")
	}
	return err
}

// Find returns all records found
func (sqld *SQLDatabase) Find(instance interface{}, where string, params map[string]interface{}) error {

	if reflect.TypeOf(instance).Elem().Kind() != reflect.Slice {
		return errors.New("The intput value is not a pointer of a slice")
	}

	model, err := sqld.Model(instance)
	if err != nil {
		return err
	}

	sql, err := sqld.Dialect.Find(model, where, instance)
	if err != nil {
		return err
	}

	rows, err := sqld.db.NamedQuery(sql, params)
	if err != nil {
		return err
	}
	defer rows.Close()

	//Creates a new pointer with the same type that resultEntitySlice
	slicePtr := reflect.ValueOf(instance)
	//it gets the value of the slice pointer
	slice := reflect.Indirect(slicePtr)

	entityType := models.GetType(instance)

	if !rows.Next() {
		return errors.New("Records not found")
	}

	for {
		entityPtr := reflect.New(entityType)

		entityFieldsAsSlice := models.StructToSliceOfAddressesWithRules(entityPtr, sqld.Dialect.GetModel)
		rows.Scan(entityFieldsAsSlice...)

		slice.Set(reflect.Append(slice, entityPtr.Elem()))

		if !rows.Next() {
			break
		}
	}

	return nil
}

// NativeFind returns all records found
func (sqld *SQLDatabase) NativeFind(resultEntitySlice interface{}, sql string, params map[string]interface{}) error {

	if reflect.TypeOf(resultEntitySlice).Elem().Kind() != reflect.Slice {
		return errors.New("The intput value is not a pointer of a slice")
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

	entityType := models.GetType(resultEntitySlice)

	if !rows.Next() {
		return errors.New("Records not found")
	}

	for {
		entityPtr := reflect.New(entityType)

		entityFieldsAsSlice := models.StructToSliceOfAddresses(entityPtr)
		rows.Scan(entityFieldsAsSlice...)

		slice.Set(reflect.Append(slice, entityPtr.Elem()))

		if !rows.Next() {
			break
		}
	}
	return nil
}

// DropTable removes a table from the database
func (sqld *SQLDatabase) DropTable(i interface{}) error {
	typ := models.GetType(i)
	name := typ.Name()

	table, ok := sqld.Dialect.GetModel(name)
	if !ok {
		return errors.New("Model not found")
	}
	sql := sqld.Dialect.Drop(table.Name)

	_, err := sqld.db.Exec(sql)
	if err != nil {
		return err
	}
	sqld.Dialect.DeleteModel(name)
	return nil
}

// TxBegin is used to set a transaction
func (sqld *SQLDatabase) TxBegin() (*sql.Tx, error) {
	return sqld.db.Begin()
}
