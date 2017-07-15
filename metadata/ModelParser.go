package metadata

import (
	"reflect"
	"strings"
	"errors"
)

type GoedbTable struct{
	Name    string
	Columns []GoedbColumn
	Model  	reflect.Type		`json:"-"`
}

type GoedbColumn struct{
	Title               string
	ColumnType          string
	PrimaryKey          bool
	Unique              bool
	ForeignKey          bool
	ForeignKeyReference string
	AutoIncrement       bool
}

func GetType(i interface{}) (reflect.Type){
	typ := reflect.TypeOf(i)

	// if a pointer to a struct is passed, get the type of the dereferenced object
	if typ.Kind() == reflect.Ptr{
		typ = typ.Elem()
	}

	if typ.Kind() == reflect.Slice{
		typ = typ.Elem()
	}

	return typ
}

func tagAttributeExists(tag reflect.StructTag, attribute string) bool{

	if tag, ok := tag.Lookup("goedb"); ok {
		params := strings.Split(tag, ",")
		for _, val := range params {
			if val == "pk" {
				return true
			}
		}
	}
	return false

}

func getPrimaryKeyType(structType reflect.Type) (reflect.Type, error) {

	for i:=0;i< structType.NumField();i++ {
		fieldValue := structType.Field(i)
		if tagAttributeExists(fieldValue.Tag, "pk"){
			return fieldValue.Type, nil
		}
	}

	return nil, errors.New("Primary key not found")
}

func processColumnType(column *GoedbColumn, columnType reflect.Type) error{
	if columnType.Kind() != reflect.Struct {
		column.ColumnType = columnType.Name()
		return nil
	}
	//TODO: Parse struct and get primary-key
	primaryKeyType, err := getPrimaryKeyType(columnType)
	if err != nil {
		return err
	}
	column.ColumnType = primaryKeyType.Name()
	return nil
}

func ParseModel(entity interface{}) (GoedbTable){
	entityType := reflect.TypeOf(entity)

	// if a pointer to a struct is passed, get the type of the dereferenced object
	if entityType.Kind() == reflect.Ptr{
		entityType = entityType.Elem()
	}

	table := GoedbTable{}
	table.Name = entityType.Name()
	table.Columns = make([]GoedbColumn, 0)

	for i:=0;i< entityType.NumField();i++ {
		tablecol := GoedbColumn{}
		tablecol.Title = entityType.Field(i).Name
		processColumnType(&tablecol, entityType.Field(i).Type)
		//tablecol.ColumnType = entityType.Field(i).Type.Name()

		if tag, ok := entityType.Field(i).Tag.Lookup("goedb"); ok {
			params := strings.Split(tag, ",")
			for _, val := range params {
				switch val {
				case "pk":
					tablecol.PrimaryKey = true
				case "autoincrement":
					tablecol.AutoIncrement = true
				case "unique":
					tablecol.Unique = true
				default:
					if strings.Contains(val, "fk=") {
						tablecol.ForeignKey = true
						tablecol.ForeignKeyReference = strings.Split(val, "=")[1]
					}
				}
			}
		}
		table.Columns = append(table.Columns, tablecol)
	}

	return table
}