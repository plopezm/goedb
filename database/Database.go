package database

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/plopezm/goedb/database/models"
)

// EntityManager is the manager used to interact with the database
type EntityManager interface {
	SetSchema(schema string) (sql.Result, error)
	Open(driver string, params string, schema string) error
	Close() error
	GetDBConnection() *sqlx.DB
	Migrate(i interface{}, autoCreate bool, dropIfExists bool) error
	DropTable(i interface{}) error
	Model(i interface{}) (models.Table, error)
	Insert(i interface{}) (models.Result, error)
	Update(i interface{}) (models.Result, error)
	Remove(i interface{}, where string, params map[string]interface{}) (models.Result, error)
	First(i interface{}, where string, params map[string]interface{}) error
	Find(i interface{}, where string, params map[string]interface{}) error
	NativeFirst(i interface{}, query string, params map[string]interface{}) error
	NativeFind(i interface{}, query string, params map[string]interface{}) error
	TxBegin() (*sql.Tx, error)
}
