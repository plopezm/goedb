package metadata

import (
	"errors"
	"reflect"
	"strings"
)

// Models contains every model migrated
var Models map[string]GoedbTable

func init() {
	Models = make(map[string]GoedbTable)
}

// GoedbTable represents the metadata of a table
type GoedbTable struct {
	Name           string
	Columns        []GoedbColumn
	PrimaryKeyName string       `json:"-"`
	PrimaryKeyType reflect.Kind `json:"-"`
}

// GoedbColumn represents the metadata of a column
type GoedbColumn struct {
	Title               string
	ColumnType          reflect.Kind
	ColumnTypeName      string
	PrimaryKey          bool
	Unique              bool
	ForeignKey          bool
	ForeignKeyReference string
	AutoIncrement       bool
	IsComplex           bool
}

// GetType returns the type of a struct
func GetType(i interface{}) reflect.Type {
	typ := reflect.TypeOf(i)

	// if a pointer to a struct is passed, get the type of the dereferenced object
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ.Kind() == reflect.Slice {
		typ = typ.Elem()
	}

	return typ
}

// GetValue returns the value of a struct
func GetValue(i interface{}) reflect.Value {
	val := reflect.ValueOf(i)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	return val
}

func tagAttributeExists(tag reflect.StructTag, attribute string) bool {

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

// GetGoedbTagTypeAndValue returns the tag and the value of a struct
func GetGoedbTagTypeAndValue(instanceType reflect.Type, instanceValue reflect.Value, goedbTag string) (reflect.Type, reflect.Value, error) {
	for i := 0; i < instanceType.NumField(); i++ {
		field := instanceType.Field(i)
		value := instanceValue.Field(i)
		if tagAttributeExists(field.Tag, goedbTag) {
			return field.Type, value, nil
		}
	}
	return nil, reflect.Value{}, errors.New(" Goedb:" + goedbTag + " not found")
}

// GetGoedbTagTypeAndValueOfIndexField returns the type and the value of a index field
func GetGoedbTagTypeAndValueOfIndexField(instanceType reflect.Type, instanceValue reflect.Value, goedbTag string, fieldID int) (reflect.Type, reflect.Value, error) {
	fieldType := instanceType.Field(fieldID).Type
	fieldValue := instanceValue.Field(fieldID)

	return GetGoedbTagTypeAndValue(fieldType, fieldValue, goedbTag)
}

func processColumnType(column *GoedbColumn, columnType reflect.Type, columnValue reflect.Value) error {

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
	return nil
}

// ParseModel generates a GoedbTable, the model of a struct
func ParseModel(entity interface{}) GoedbTable {
	entityType := GetType(entity)
	entityValue := GetValue(entity)

	table := GoedbTable{}
	table.Name = entityType.Name()
	table.Columns = make([]GoedbColumn, 0)

	for i := 0; i < entityType.NumField(); i++ {
		tablecol := GoedbColumn{}
		tablecol.Title = entityType.Field(i).Name
		processColumnType(&tablecol, entityType.Field(i).Type, entityValue)

		if tag, ok := entityType.Field(i).Tag.Lookup("goedb"); ok {
			params := strings.Split(tag, ",")
			for _, val := range params {
				switch val {
				case "pk":
					tablecol.PrimaryKey = true
					table.PrimaryKeyName = tablecol.Title
					table.PrimaryKeyType = tablecol.ColumnType
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

func getSubStructAddresses(slice *[]interface{}, value reflect.Value) {
	for j := 0; j < value.NumField(); j++ {
		subField := value.Field(j)
		if subField.Kind() == reflect.Struct {
			getSubStructAddresses(slice, subField)
			continue
		}
		*slice = append(*slice, subField.Addr().Interface())
	}
}

// StructToSliceOfAddresses returns a slice with the addresses of each struct field,
// so any modification on the slide will modify the source struct fields
func StructToSliceOfAddresses(structPtr interface{}) []interface{} {

	var fieldArr reflect.Value
	if _, ok := structPtr.(reflect.Value); ok {
		fieldArr = structPtr.(reflect.Value)
	} else {
		fieldArr = reflect.ValueOf(structPtr).Elem()
	}

	if fieldArr.Kind() == reflect.Ptr {
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
