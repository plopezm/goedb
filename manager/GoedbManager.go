package manager

import (
	"goedb/metadata"
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
	First(i interface{}, params string) error
	Find(i interface{}, params string) error
}
