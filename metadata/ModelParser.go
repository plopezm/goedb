package metadata

import (
	"reflect"
	"strings"
	"errors"
)

var Models map[string]GoedbTable

func init(){
	Models = make(map[string]GoedbTable)
}

type GoedbTable struct{
	Name    string
	Columns []GoedbColumn
	Model  	reflect.Type		`json:"-"`
}

type GoedbColumn struct{
	Title               string
	ColumnType          reflect.Kind
	ColumnTypeName		string
	PrimaryKey          bool
	Unique              bool
	ForeignKey          bool
	ForeignKeyReference string
	AutoIncrement       bool
	IsComplex			bool
}

type GoedbComplexColumn struct {
	MappedFKValue				reflect.Value
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
			if val == attribute {
				return true
			}
		}
	}
	return false
}

func GetGoedbTagTypeAndValue(instanceType reflect.Type, instanceValue reflect.Value, goedbTag string) (reflect.Type, reflect.Value, error){
	for i:=0;i< instanceType.NumField();i++ {
		field := instanceType.Field(i)
		value := instanceValue.Field(i)
		if tagAttributeExists(field.Tag, goedbTag){
			return field.Type, value, nil
		}
	}
	return nil, reflect.Value{}, errors.New(" Goedb:"+goedbTag+" not found")
}

func GetGoedbTagTypeAndValueOfIndexField(instanceType reflect.Type, instanceValue reflect.Value, goedbTag string, fieldId int) (reflect.Type, reflect.Value, error){
	fieldType := instanceType.Field(fieldId).Type
	fieldValue := instanceValue.Field(fieldId)

	return GetGoedbTagTypeAndValue(fieldType, fieldValue, goedbTag)
}

func processColumnType(column *GoedbColumn, columnType reflect.Type, columnValue reflect.Value) error{

	column.ColumnTypeName = columnType.Name()
	if columnType.Kind() != reflect.Struct {
		column.ColumnType = columnType.Kind()
		return nil
	}
	primaryKeyType, _, err := GetGoedbTagTypeAndValue(columnType, columnValue, "pk")
	if err != nil {
		return err
	}

	column.ColumnType = primaryKeyType.Kind()
	column.IsComplex = true
	//column.ComplexColumn = new(GoedbComplexColumn)
	//column.ComplexColumn.MappedFKValue = primaryKeyValue
	//column.ComplexColumn.ReferencedStructName = columnType.Name()
	//column.ComplexColumn.ReferencedStructAttrNames = make([]string, columnType.NumField())
	//for i:=0;i < columnType.NumField();i++{
	//	column.ComplexColumn.ReferencedStructAttrNames[i] = columnType.Field(i).Name
	//}

	return nil
}

func ParseModel(entity interface{}) (GoedbTable){
	entityType := GetType(entity)
	entityValue := GetValue(entity)

	table := GoedbTable{}
	table.Name = entityType.Name()
	table.Columns = make([]GoedbColumn, 0)

	for i:=0;i< entityType.NumField();i++ {
		tablecol := GoedbColumn{}
		tablecol.Title = entityType.Field(i).Name
		processColumnType(&tablecol, entityType.Field(i).Type, entityValue)

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
	Models[table.Name] = table
	return table
}

func getSubStructAddresses(slice *[]interface{}, value reflect.Value){
	for j := 0; j < value.NumField(); j++ {
		subField := value.Field(j)
		if subField.Kind() == reflect.Struct {
			getSubStructAddresses(slice, subField)
			continue
		}
		*slice = append(*slice, subField.Addr().Interface())
	}
}

/*
	Returns a slice with the addresses of each struct field,
	so any modification on the slide will modify the source struct fields
 */
func StructToSliceOfAddresses(structPtr interface{}) []interface{} {

	var fieldArr reflect.Value
	if _, ok  := structPtr.(reflect.Value); ok{
		fieldArr = structPtr.(reflect.Value)
	}else{
		fieldArr = reflect.ValueOf(structPtr).Elem()
	}

	if fieldArr.Kind() == reflect.Ptr{
		fieldArr = fieldArr.Elem()
	}

	fieldAddrArr := make([]interface{}, 0)

	for i := 0; i < fieldArr.NumField(); i++ {
		f := fieldArr.Field(i)

		if f.Kind() == reflect.Struct {
			getSubStructAddresses(&fieldAddrArr, f)
			continue
		}
		fieldAddrArr = append(fieldAddrArr, f.Addr().Interface())
	}

	return fieldAddrArr
}