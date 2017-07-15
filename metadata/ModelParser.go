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
	IsComplexType		bool
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


func GetValue(i interface{}) (reflect.Value){
	val := reflect.ValueOf(i)

	if val.Kind() == reflect.Ptr{
		val = val.Elem()
	}

	return val
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
/*
	Returns the primary key type from a struct type
 */
func GetPrimaryKeyType(structType reflect.Type) (reflect.Type, error) {

	for i:=0;i< structType.NumField();i++ {
		field := structType.Field(i)
		if tagAttributeExists(field.Tag, "pk"){
			return field.Type, nil
		}
	}

	return nil, errors.New("Primary key not found")
}

func GetPrimaryKeyValue(instanceType reflect.Type, instanceValue reflect.Value) (reflect.Value, error) {

	for i:=0;i< instanceType.NumField();i++ {
		field := instanceType.Field(i)
		fieldValue := instanceValue.Field(i)
		if tagAttributeExists(field.Tag, "pk"){
			return fieldValue, nil
		}
	}

	return reflect.Value{}, errors.New("Primary key not found")
}

func GetForeignKeyValue(instanceType reflect.Type, instanceValue reflect.Value) (reflect.Value, error) {

	for i:=0;i< instanceType.NumField();i++ {
		field := instanceType.Field(i)
		fieldValue := instanceValue.Field(i)
		if tagAttributeExists(field.Tag, "fk"){
			return fieldValue, nil
		}
	}

	return reflect.Value{}, errors.New("Primary key not found")
}

/*
	Returns the primary key type and value from a struct type and value
 */
func GetForeignKeyItsPrimaryKeyValue(instanceType reflect.Type, instanceValue reflect.Value, fieldId int) (reflect.Value, error) {

	fieldType := instanceType.Field(fieldId).Type
	fieldValue := instanceValue.Field(fieldId)

	return GetForeignKeyValue(fieldType, fieldValue)
}

func processColumnType(column *GoedbColumn, columnType reflect.Type) error{
	if columnType.Kind() != reflect.Struct {
		column.ColumnType = columnType.Name()
		return nil
	}
	primaryKeyType, err := GetPrimaryKeyType(columnType)
	if err != nil {
		return err
	}
	column.ColumnType = primaryKeyType.Name()
	column.IsComplexType = true
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