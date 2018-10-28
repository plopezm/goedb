package models

import (
	"reflect"
	"testing"
)

type testtable struct {
	ID     uint64            `goedb:"pk,autoincrement"`
	Name   string            `goedb:"unique"`
	Childs []testtablewithfk `goedb:"mappedBy=testtablewithfk(Name)"`
}

type testtablewithfk struct {
	Name          string    `goedb:"pk"`
	TestTableName testtable `goedb:"pk,fk=testtable(Name)"`
	Ignorable     bool      `goedb:"ignore"`
	Desc          string
}

type testtablewithfk2 struct {
	Name                string          `goedb:"pk"`
	TestTableWithFKName testtablewithfk `goedb:"pk,fk=testtablewithfk(Name)"`
}

func TestParseModel(t *testing.T) {
	type args struct {
		entity interface{}
	}
	tests := []struct {
		name string
		args args
		want Table
	}{
		// TODO: Add test cases.
		{
			name: "TestParseModelTestTable",
			args: args{
				entity: &testtable{},
			},
			want: Table{
				Name: "testtable",
				Columns: []Column{
					{
						Title:          "ID",
						AutoIncrement:  true,
						PrimaryKey:     true,
						ColumnType:     reflect.Uint64,
						ColumnTypeName: "uint64",
					},
					{
						Title:          "Name",
						Unique:         true,
						ColumnType:     reflect.String,
						ColumnTypeName: "string",
					},
				},
				MappedColumns: []Column{
					{
						Title:      "Childs",
						ColumnType: reflect.Slice,
						MappedBy: MappedByField{
							TargetTablePK:   "Name",
							TargetTableName: "testtablewithfk",
						},
						IsMapped: true,
					},
				},
				PrimaryKeys: []PrimaryKey{
					{
						Name: "ID", Type: reflect.Uint64,
					},
					{
						Name: "Name", Type: reflect.String,
					},
				},
			},
		},
		{
			name: "TestParseModelForeignKey",
			args: args{
				entity: &testtablewithfk{},
			},
			want: Table{
				Name: "testtablewithfk",
				Columns: []Column{
					{
						Title:          "Name",
						PrimaryKey:     true,
						ColumnType:     reflect.String,
						ColumnTypeName: "string",
					},
					{
						Title:          "TestTableName",
						PrimaryKey:     true,
						ColumnType:     reflect.String,
						ColumnTypeName: "testtable",
						ForeignKey:     ForeignKey{IsForeignKey: true, ForeignKeyTableReference: "testtable", ForeignKeyColumnReference: "Name"},
						IsComplex:      true,
					},
					{
						Title:          "Ignorable",
						Ignore:         true,
						ColumnType:     reflect.Bool,
						ColumnTypeName: "bool",
					},
					{
						Title:          "Desc",
						ColumnType:     reflect.String,
						ColumnTypeName: "string",
					},
				},
				PrimaryKeys: []PrimaryKey{
					{
						Name: "Name", Type: reflect.String,
					},
					{
						Name: "TestTableName", Type: reflect.String,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseModel(tt.args.entity); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseModel() =\n %v \n want\n %v", got, tt.want)
			}
		})
	}
}

func getGoedbTableTestParser() interface{} {
	return &testtablewithfk2{
		Name: "ExampleMultiStruct",
		TestTableWithFKName: testtablewithfk{
			Name:          "testtablewithfk-Name",
			TestTableName: testtable{ID: 1, Name: "TestTableName-Name-ID"},
			Ignorable:     true,
			Desc:          "testing description",
		},
	}
}

func MockGetModel(name string) (Table, bool) {
	if name == "testtablewithfk2" {
		return ParseModel(&testtablewithfk2{
			Name: "ExampleMultiStruct",
			TestTableWithFKName: testtablewithfk{
				Name:          "testtablewithfk-Name",
				TestTableName: testtable{ID: 1, Name: "TestTableName-Name-ID"},
				Ignorable:     true,
				Desc:          "testing description",
			},
		}), true
	} else if name == "testtablewithfk" {
		return ParseModel(&testtablewithfk{
			Name:          "testtablewithfk-Name",
			TestTableName: testtable{ID: 1, Name: "TestTableName-Name-ID"},
			Ignorable:     true,
			Desc:          "testing description",
		}), true
	} else {
		return ParseModel(&testtable{ID: 1, Name: "TestTableName-Name-ID"}), true
	}
}

func TestStructToSliceOfAddressesWithRules(t *testing.T) {
	type args struct {
		structPtr interface{}
		GetModel  func(name string) (Table, bool)
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		// TODO: Add test cases.
		{
			name: "TestStructToSliceOfAddressesWithRules_MultipleStructs",
			args: args{structPtr: getGoedbTableTestParser(), GetModel: MockGetModel},
			want: []interface{}{"addr1", "addr2", "addr3", "addr4", "addr5", "addr6"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StructToSliceOfAddressesWithRules(tt.args.structPtr, tt.args.GetModel); len(got) != len(tt.want) {
				t.Errorf("StructToSliceOfAddressesWithRules() = %v, want %v", got, tt.want)
			}
		})
	}
}
