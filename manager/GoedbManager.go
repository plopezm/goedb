package manager

import (
	"database/sql"
	"github.com/plopezm/goedb/metadata"
)

// GoedbResult is the result for some operation in database
type GoedbResult struct {
	NumRecordsAffected int64
}

// EntityManager is the manager used to interact with the database
type EntityManager interface {
	Open(driver string, params string) error
	Close() error
	Migrate(i interface{}) error
	DropTable(i interface{}) error
	Model(i interface{}) (metadata.GoedbTable, error)
	Insert(i interface{}) (GoedbResult, error)
	Remove(i interface{}, where string, params ...interface{}) (GoedbResult, error)
	First(i interface{}, where string, params ...interface{}) error
	Find(i interface{}, where string, params ...interface{}) error
	TxBegin() (*sql.Tx, error)
}
