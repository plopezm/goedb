package database

import (
	"errors"
	"reflect"
	"strings"
)

// Models contains every model migrated
var Models map[string]Table

// Table represents the metadata of a table
type Table struct {
	Name        string
	Columns     []Column
	PrimaryKeys []PrimaryKey
}

type PrimaryKey struct {
	Name string
	Type reflect.Kind
}

type ForeignKey struct {
	IsForeignKey              bool
	ForeignKeyTableReference  string
	ForeignKeyColumnReference string
}

// Column represents the metadata of a column
type Column struct {
	Title          string
	ColumnType     reflect.Kind
	ColumnTypeName string
	PrimaryKey     bool
	Unique         bool
	ForeignKey     ForeignKey
	AutoIncrement  bool
	IsComplex      bool
	Ignore         bool
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
			if strings.Contains(attribute, val) {
				return true
			}
		}
	}
	return false
}

// GetGoedbTagTypeAndValueOfForeignKeyReference returns the tag and the value of a struct
func GetGoedbTagTypeAndValueOfForeignKeyReference(instanceType reflect.Type, instanceValue reflect.Value, goedbTag string, foreignKeyReference ForeignKey) (reflect.Type, reflect.Value, error) {
	for i := 0; i < instanceType.NumField(); i++ {
		field := instanceType.Field(i)
		value := instanceValue.Field(i)
		if tagAttributeExists(field.Tag, goedbTag) && foreignKeyReference.ForeignKeyColumnReference == field.Name {
			return field.Type, value, nil
		}
	}
	return nil, reflect.Value{}, errors.New(" Goedb:" + goedbTag + " not found")
}

// GetGoedbTagTypeAndValueOfIndexField returns the type and the value of a index field
/*func GetGoedbTagTypeAndValueOfIndexField(instanceType reflect.Type, instanceValue reflect.Value, goedbTag string, fieldID int) (reflect.Type, reflect.Value, error) {
	fieldType := instanceType.Field(fieldID).Type
	fieldValue := instanceValue.Field(fieldID)

	return GetGoedbTagTypeAndValue(fieldType, fieldValue, goedbTag)
}*/

func processColumnType(column *Column, columnType reflect.Type, columnValue reflect.Value) error {

	column.ColumnTypeName = columnType.Name()
	if columnType.Kind() != reflect.Struct {
		column.ColumnType = columnType.Kind()
		return nil
	}
	primaryKeyType, _, err := GetGoedbTagTypeAndValueOfForeignKeyReference(columnType, columnValue, "pk,unique", column.ForeignKey)
	if err != nil {
		return err
	}

	column.ColumnType = primaryKeyType.Kind()
	column.IsComplex = true
	return nil
}

// ParseModel generates a GoedbTable, the model of a struct
func ParseModel(entity interface{}) Table {
	entityType := GetType(entity)
	entityValue := GetValue(entity)

	table := Table{}
	table.Name = entityType.Name()
	table.Columns = make([]Column, 0)

	for i := 0; i < entityType.NumField(); i++ {
		tablecol := Column{}
		tablecol.Title = entityType.Field(i).Name

		if tag, ok := entityType.Field(i).Tag.Lookup("goedb"); ok {
			params := strings.Split(tag, ",")
			for _, val := range params {
				switch val {
				case "pk":
					tablecol.PrimaryKey = true
					//table.PrimaryKeyName = tablecol.Title
					//table.PrimaryKeyType = tablecol.ColumnType
				case "autoincrement":
					tablecol.AutoIncrement = true
				case "unique":
					tablecol.Unique = true
				case "ignore":
					tablecol.Ignore = true
				default:
					if strings.Contains(val, "fk=") {
						tablecol.ForeignKey.IsForeignKey = true
						//References are received in the following format: ReferencedTable(ReferencedColumn)
						fktag := val[3:]
						fksubtags := strings.Split(fktag, "(")
						tablecol.ForeignKey.ForeignKeyTableReference = fksubtags[0]
						tablecol.ForeignKey.ForeignKeyColumnReference = fksubtags[1][:len(fksubtags[1])-1]

					}
				}
			}
		}
		processColumnType(&tablecol, entityType.Field(i).Type, entityValue)
		if tablecol.PrimaryKey || tablecol.Unique {
			table.PrimaryKeys = append(table.PrimaryKeys, PrimaryKey{Name: tablecol.Title, Type: tablecol.ColumnType})
		}
		table.Columns = append(table.Columns, tablecol)
	}
	return table
}
