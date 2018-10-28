package dbaccess

import "github.com/plopezm/goedb/database/models"

// DatabaseAccess database access layer functions (could be a sql dbaccess or no-sql database)
type DatabaseAccess interface {
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

// GetDatabaseAccess returns the database depending on the driver used (could be a sql dbaccess or no-sql database)
func GetDatabaseAccess(driver string) (databaseAccess DatabaseAccess) {
	switch driver {
	case "sqlite3", "postgres", "pgx":
		databaseAccess = GetSQLDatabaseAccess(driver)
	default:
		databaseAccess = GetSQLDatabaseAccess(driver)
	}
	return databaseAccess
}
