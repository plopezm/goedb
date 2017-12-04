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
