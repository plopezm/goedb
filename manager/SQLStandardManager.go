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
func (sqld *GoedbSQLDriver) Migrate(schema string, i interface{}) error {
	sqld.DropTable(i)
	table := metadata.ParseModel(schema, i)
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

	sql, err := sqld.Dialect.GetSQLDelete(model, where, i)
	if err != nil {
		return goedbres,err
	}

	result, err = sqld.db.NamedExec(sql, params)
	goedbres.NumRecordsAffected, _ = result.RowsAffected()
	return goedbres, err
}

// First returns the first record found
func (sqld *GoedbSQLDriver) First(instance interface{}, where string, params map[string]interface{}) error {
	model, err := sqld.Model(instance)
	if err != nil {
		return err
	}

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
