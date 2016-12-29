package goedb

type GoedbTable struct{
	name    	string
	columns 	[]GoedbColumn
}

type GoedbColumn struct{
	title   	string
	ctype   	string
	pk      	bool
	unique  	bool
	fk      	bool
	fkref   	string
	autoinc 	bool
}
