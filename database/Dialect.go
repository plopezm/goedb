package database

// GetDialect returns the dialect depending on the driver used
/*
func GetDialect(driver string) Dialect {
	switch driver {
	case "sqlite3":
		return new(SQLiteDialect)
	case "postgres", "pgx":
		return new(PostgresDialect)
	default:
		return new(SQLiteDialect)
	}
}
*/

// Dialect represents a database dialect
type Dialect interface {
	Create(table Table) string
	First(table Table, where string, instance interface{}) (string, error)
	Find(table Table, where string, instance interface{}) (string, error)
	Update(table Table, instance interface{}) (string, error)
	Delete(table Table, where string, instance interface{}) (string, error)
	Drop(tableName string) string
}

// GetDialect returns the dialect depending on the driver used
func GetDialect(driver string) Dialect {
	switch driver {
	case "sqlite3":
		return new(TransientSQLDialect)
	case "postgres", "pgx":
		return new(TransientSQLDialect)
	default:
		return new(TransientSQLDialect)
	}
}
