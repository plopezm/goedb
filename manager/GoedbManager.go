package manager

import (
	"github.com/plopezm/goedb/metadata"
	"database/sql"
)

type GoedbResult struct {
	NumRecordsAffected int64
}

type EntityManager interface {
	Open(driver string, params string) error
	Close() error
	Migrate(i interface{}) error
	DropTable(i interface{}) error
	Model(i interface{})(metadata.GoedbTable, error)
	Insert(i interface{}) (GoedbResult, error)
	Remove(i interface{}) (GoedbResult, error)
	First(i interface{}, where string, params ...interface{}) error
	Find(i interface{}, where string, params ...interface{}) error
	TxBegin() (*sql.Tx, error)
}
