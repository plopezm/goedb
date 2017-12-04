package dialect

import "github.com/plopezm/goedb/database/models"

// Dialect represents a database dialect
type Dialect interface {
	GetModel(name string) (models.Table, bool)
	SetModel(name string, table models.Table)
	DeleteModel(name string)
	Create(table models.Table) string
	Insert(table models.Table, instance interface{}) (string, error)
	First(table models.Table, where string, instance interface{}) (string, error)
	Find(table models.Table, where string, instance interface{}) (string, error)
	Update(table models.Table, instance interface{}) (string, error)
	Delete(table models.Table, where string, instance interface{}) (string, error)
	Drop(tableName string) string
}

// GetDialect returns the dialect depending on the driver used
func GetDialect(driver string) (dialect Dialect) {
	switch driver {
	case "sqlite3", "postgres", "pgx":
		dialect = GetSQLDialectInstance(driver)
	default:
		dialect = GetSQLDialectInstance(driver)
	}
	return dialect
}
