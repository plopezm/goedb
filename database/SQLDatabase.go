package database

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/plopezm/goedb/config"
	"github.com/plopezm/goedb/metadata"
)

//SQLDatabase is the implementation of SQL for a Database interface
type SQLDatabase struct {
	db         *sqlx.DB
	Dialect    Dialect
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

// Migrate creates the table in the database
func (sqld *SQLDatabase) Migrate(i interface{}, autoCreate bool, dropIfExists bool) (err error) {
	if dropIfExists {
		sqld.DropTable(i)
	}
	table := ParseModel(i)
	Models[table.Name] = table
	if autoCreate {
		sqltab := sqld.Dialect.Create(table)
		_, err = sqld.db.Exec(sqltab)
	}
	return err
}

// DropTable removes a table from the database
func (sqld *SQLDatabase) DropTable(i interface{}) error {
	typ := metadata.GetType(i)
	name := typ.Name()

	sql := sqld.Dialect.Drop(Models[name].Name)

	_, err := sqld.db.Exec(sql)
	if err != nil {
		return err
	}
	delete(metadata.Models, name)
	return nil
}
