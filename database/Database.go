package database

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Result is the result for some operation in database
type Result struct {
	NumRecordsAffected int64
	LastInsertId       int64
}

// EntityManager is the manager used to interact with the database
type EntityManager interface {
	SetSchema(schema string) (sql.Result, error)
	Open(driver string, params string, schema string) error
	Close() error
	GetDBConnection() *sqlx.DB
	Migrate(i interface{}, autoCreate bool, dropIfExists bool) error
	DropTable(i interface{}) error
	Model(i interface{}) (Table, error)
	Insert(i interface{}) (Result, error)
	Update(i interface{}) (Result, error)
	Remove(i interface{}, where string, params map[string]interface{}) (Result, error)
	First(i interface{}, where string, params map[string]interface{}) error
	Find(i interface{}, where string, params map[string]interface{}) error
	NativeFirst(i interface{}, query string, params map[string]interface{}) error
	NativeFind(i interface{}, query string, params map[string]interface{}) error
	TxBegin() (*sql.Tx, error)
}
