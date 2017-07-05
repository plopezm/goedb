package drivers

import "reflect"

type GoedbTable struct{
	Name    string
	Columns []GoedbColumn
	Model  	reflect.Type		`json:"-"`
}

type GoedbColumn struct{
	Title   	string
	Ctype   	string
	Pk      	bool
	Unique  	bool
	Fk      	bool
	Fkref   	string
	Autoinc 	bool
}

type GoedbResult struct {
	NumRecordsAffected int64
}

type GoedbDriver interface {
	Open(driver string, params string) error
	Close() error
	Migrate(i interface{}) error
	DropTable(i interface{}) error
	Model(i interface{})(GoedbTable, error)
	Insert(i interface{}) (GoedbResult, error)
	Remove(i interface{}) (GoedbResult, error)
	First(i interface{}, params string) error
	Find(i interface{}, params string) error
}
