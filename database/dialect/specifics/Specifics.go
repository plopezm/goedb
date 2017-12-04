package specifics

import "github.com/plopezm/goedb/database/models"

//Specifics is a small change in a dialect, it will be used for similar databases
type Specifics interface {
	GetSQLCreateTableColumn(value models.Column) (sqlColumnLine string, primaryKey string, constraints string, err error)
}
