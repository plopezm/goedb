package dialect

import "github.com/plopezm/goedb/database/models"

//DBAccess is a small change in a dbaccess, it will be used for similar databases
type Dialect interface {
	GetSQLCreateTableColumn(value models.Column) (sqlColumnLine string, primaryKey string, constraints string, err error)
}
