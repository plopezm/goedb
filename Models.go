package goedb

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
