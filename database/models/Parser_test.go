package models

import (
	"reflect"
	"testing"
)

func TestParseModel(t *testing.T) {
	type TestTable struct {
		ID   uint64 `goedb:"pk,autoincrement"`
		Name string `goedb:"unique"`
	}

	type TestTableWithFK struct {
		Name          string    `goedb:"pk"`
		TestTableName TestTable `goedb:"pk,fk=TestTable(Name)"`
		Ignorable     bool      `goedb:"ignore"`
	}

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
				entity: &TestTable{},
			},
			want: Table{
				Name: "TestTable",
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
				entity: &TestTableWithFK{},
			},
			want: Table{
				Name: "TestTableWithFK",
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
						ColumnTypeName: "TestTable",
						ForeignKey:     ForeignKey{IsForeignKey: true, ForeignKeyTableReference: "TestTable", ForeignKeyColumnReference: "Name"},
						IsComplex:      true,
					},
					{
						Title:          "Ignorable",
						Ignore:         true,
						ColumnType:     reflect.Bool,
						ColumnTypeName: "bool",
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

	type TestTable struct {
		ID   uint64 `goedb:"pk,autoincrement"`
		Name string `goedb:"unique"`
	}

	type TestTableWithFK struct {
		Name          string    `goedb:"pk"`
		TestTableName TestTable `goedb:"pk,fk=TestTable(Name)"`
		Ignorable     bool      `goedb:"ignore"`
		Desc          string
	}

	type TestTableWithFK2 struct {
		Name                string          `goedb:"pk"`
		TestTableWithFKName TestTableWithFK `goedb:"pk,fk=TestTableWithFK(Name)"`
	}

	return &TestTableWithFK2{
		Name: "ExampleMultiStruct",
		TestTableWithFKName: TestTableWithFK{
			Name:          "TestTableWithFK-Name",
			TestTableName: TestTable{ID: 1, Name: "TestTableName-Name-ID"},
			Ignorable:     true,
			Desc:          "testing description",
		},
	}
}

func MockGetModel(name string) (Table, bool) {
	type TestTable struct {
		ID   uint64 `goedb:"pk,autoincrement"`
		Name string `goedb:"unique"`
	}

	type TestTableWithFK struct {
		Name          string    `goedb:"pk"`
		TestTableName TestTable `goedb:"pk,fk=TestTable(Name)"`
		Ignorable     bool      `goedb:"ignore"`
		Desc          string
	}

	type TestTableWithFK2 struct {
		Name                string          `goedb:"pk"`
		TestTableWithFKName TestTableWithFK `goedb:"pk,fk=TestTableWithFK(Name)"`
	}

	if name == "TestTableWithFK2" {
		return ParseModel(&TestTableWithFK2{
			Name: "ExampleMultiStruct",
			TestTableWithFKName: TestTableWithFK{
				Name:          "TestTableWithFK-Name",
				TestTableName: TestTable{ID: 1, Name: "TestTableName-Name-ID"},
				Ignorable:     true,
				Desc:          "testing description",
			},
		}), true
	} else if name == "TestTableWithFK" {
		return ParseModel(&TestTableWithFK{
			Name:          "TestTableWithFK-Name",
			TestTableName: TestTable{ID: 1, Name: "TestTableName-Name-ID"},
			Ignorable:     true,
			Desc:          "testing description",
		}), true
	} else {
		return ParseModel(&TestTable{ID: 1, Name: "TestTableName-Name-ID"}), true
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
			want: []interface{}{"addr1", "addr2", "addr3", "addr4", "addr5"},
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
