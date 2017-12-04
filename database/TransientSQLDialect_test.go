package database

import (
	"reflect"
	"testing"
)

func Test_getTransientSQLCreateColumn(t *testing.T) {
	type args struct {
		value Column
	}
	tests := [...]struct {
		name              string
		args              args
		wantSQLColumnLine string
		wantPrimaryKey    string
		wantConstraints   string
		wantErr           bool
	}{
		// TODO: Add test cases.
		{
			name: "TestColumnPrimaryKeyAutoincrement",
			args: args{
				value: Column{
					Title:         "PKColumn",
					PrimaryKey:    true,
					ColumnType:    reflect.Uint64,
					AutoIncrement: true,
				},
			},
			wantSQLColumnLine: "PKColumn BIGINT PRIMARY KEY AUTOINCREMENT,",
			wantPrimaryKey:    "",
			wantConstraints:   "",
		},
		{
			name: "TestColumnPrimaryKey",
			args: args{
				value: Column{
					Title:      "PKColumn",
					PrimaryKey: true,
					ColumnType: reflect.Int,
				},
			},
			wantSQLColumnLine: "PKColumn INTEGER,",
			wantPrimaryKey:    "PKColumn,",
			wantConstraints:   "",
		},
		{
			name: "TestColumnPrimaryKeyWithConstraints",
			args: args{
				value: Column{
					Title:      "PKColumn",
					PrimaryKey: true,
					ColumnType: reflect.Int,
					ForeignKey: ForeignKey{IsForeignKey: true, ForeignKeyTableReference: "OtherTable", ForeignKeyColumnReference: "OtherTablePK"},
				},
			},
			wantSQLColumnLine: "PKColumn INTEGER,",
			wantPrimaryKey:    "PKColumn,",
			wantConstraints:   ", FOREIGN KEY (PKColumn) REFERENCES OtherTable(OtherTablePK) ON DELETE CASCADE",
		},
		{
			name: "TestErrorTypeNotFound",
			args: args{
				value: Column{
					Title:      "PKColumn",
					PrimaryKey: true,
					ColumnType: reflect.Struct,
					ForeignKey: ForeignKey{IsForeignKey: true, ForeignKeyTableReference: "OtherTable", ForeignKeyColumnReference: "OtherTablePK"},
				},
			},
			wantErr: true,
		},
		{
			name: "TestStringColumnUnique",
			args: args{
				value: Column{
					Title:      "UniqueColumnString",
					ColumnType: reflect.String,
					Unique:     true,
				},
			},
			wantSQLColumnLine: "UniqueColumnString VARCHAR UNIQUE,",
		},
		{
			name: "TestNormalStringColumn",
			args: args{
				value: Column{
					Title:      "NormalColumnString",
					ColumnType: reflect.String,
				},
			},
			wantSQLColumnLine: "NormalColumnString VARCHAR,",
		},
		{
			name: "TestNormalStringColumn",
			args: args{
				value: Column{
					Title:      "NormalColumnString",
					ColumnType: reflect.Float64,
				},
			},
			wantSQLColumnLine: "NormalColumnString FLOAT,",
		},
		{
			name: "TestNormalStringColumn",
			args: args{
				value: Column{
					Title:      "NormalColumnString",
					ColumnType: reflect.Bool,
				},
			},
			wantSQLColumnLine: "NormalColumnString BOOLEAN,",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantSQLColumnLine, gotPrimaryKey, gotConstraints, err := getTransientSQLCreateColumn(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTransientSQLCreateColumn() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if wantSQLColumnLine != tt.wantSQLColumnLine {
				t.Errorf("getTransientSQLCreateColumn() wantSQLColumnLine = %v, want = %v", wantSQLColumnLine, tt.wantSQLColumnLine)
			}
			if gotPrimaryKey != tt.wantPrimaryKey {
				t.Errorf("getTransientSQLCreateColumn() gotPrimaryKey = %v, want = %v", gotPrimaryKey, tt.wantPrimaryKey)
			}
			if gotConstraints != tt.wantConstraints {
				t.Errorf("getTransientSQLCreateColumn() gotConstraints = %v, want = %v", gotConstraints, tt.wantConstraints)
			}
		})
	}
}

func TestTransientSQLDialect_Create(t *testing.T) {
	type args struct {
		table Table
	}
	tests := [...]struct {
		name    string
		dialect *TransientSQLDialect
		args    args
		want    string
	}{
		// TODO: Add test cases.
		{
			name:    "TestCreateTable",
			dialect: new(TransientSQLDialect),
			args: args{
				table: Table{
					Name: "Table1",
					Columns: []Column{
						{
							Title:         "PKColumn",
							PrimaryKey:    true,
							ColumnType:    reflect.Uint64,
							AutoIncrement: true,
						},
						{
							Title:      "NormalColumnString",
							ColumnType: reflect.String,
						},
					},
				},
			},
			want: "CREATE TABLE Table1 (PKColumn BIGINT PRIMARY KEY AUTOINCREMENT,NormalColumnString VARCHAR)",
		},
		{
			name:    "TestCreateTableWithPrimaryKeys",
			dialect: new(TransientSQLDialect),
			args: args{
				table: Table{
					Name: "Table1",
					Columns: []Column{
						{
							Title:      "PKColumn1",
							PrimaryKey: true,
							ColumnType: reflect.Uint64,
						},
						{
							Title:      "PKColumn2",
							PrimaryKey: true,
							ColumnType: reflect.Uint64,
						},
						{
							Title:      "NormalColumnString",
							ColumnType: reflect.String,
						},
						{
							Title:      "NormalColumnString",
							ColumnType: reflect.Struct,
						},
					},
				},
			},
			want: "CREATE TABLE Table1 (PKColumn1 BIGINT,PKColumn2 BIGINT,NormalColumnString VARCHAR, PRIMARY KEY (PKColumn1,PKColumn2))",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dialect := &TransientSQLDialect{}
			if got := dialect.Create(tt.args.table); got != tt.want {
				t.Errorf("TransientSQLDialect.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getGoedbTableTest1() Table {

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

	return ParseModel(&TestTableWithFK{})
}

func getGoedbTableTest2() Table {

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

	return ParseModel(&TestTableWithFK2{})
}

func getGoedbTableMapTest() (modelMap map[string]Table) {
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
	modelMap = make(map[string]Table)
	modelMap["TestTable"] = ParseModel(&TestTable{})
	modelMap["TestTableWithFK"] = ParseModel(&TestTableWithFK{})
	modelMap["TestTableWithFK2"] = ParseModel(&TestTableWithFK2{})
	return modelMap
}

func Test_generateSQLQuery(t *testing.T) {
	type args struct {
		table    Table
		modelMap map[string]Table
	}
	tests := []struct {
		name            string
		args            args
		wantQuery       string
		wantConstraints string
		wantErr         bool
	}{
		// TODO: Add test cases.
		{
			name: "TestGenerateSQLQuery",
			args: args{
				table:    getGoedbTableTest1(),
				modelMap: getGoedbTableMapTest(),
			},
			wantQuery:       "SELECT TestTableWithFK.Name,TestTable.ID,TestTable.Name,TestTableWithFK.Desc FROM TestTableWithFK,TestTable",
			wantConstraints: " AND TestTableWithFK.TestTableName = TestTable.Name",
		},
		{
			name: "TestGenerateSQLQueryMoreThanOneStructAsDependency",
			args: args{
				table:    getGoedbTableTest2(),
				modelMap: getGoedbTableMapTest(),
			},
			wantQuery:       "SELECT TestTableWithFK2.Name,TestTableWithFK.Name,TestTable.ID,TestTable.Name,TestTableWithFK.Desc FROM TestTableWithFK2,TestTableWithFK,TestTable",
			wantConstraints: " AND TestTableWithFK2.TestTableWithFKName = TestTableWithFK.Name AND TestTableWithFK.TestTableName = TestTable.Name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery, gotConstraints, err := generateSQLQuery(tt.args.table, tt.args.modelMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateSQLQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotQuery != tt.wantQuery {
				t.Errorf("generateSQLQuery() gotQuery = %v, want %v", gotQuery, tt.wantQuery)
			}
			if gotConstraints != tt.wantConstraints {
				t.Errorf("generateSQLQuery() gotConstraints = %v, want %v", gotConstraints, tt.wantConstraints)
			}
		})
	}
}

func getGoedbTableTest1Value() interface{} {

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

	return &TestTableWithFK{
		Name:          "TestTableWithFK-Name",
		TestTableName: TestTable{ID: 1, Name: "TestTableName-Name-ID"},
		Ignorable:     true,
		Desc:          "testing description",
	}
}

func Test_getPrimaryKeysAndValues(t *testing.T) {
	type args struct {
		gt  Table
		obj interface{}
	}
	tests := []struct {
		name            string
		args            args
		wantColumnName  []string
		wantColumnValue []string
		wantErr         bool
	}{
		// TODO: Add test cases.
		{
			name: "GetPrimaryKeysValues",
			args: args{
				gt:  getGoedbTableTest1(),
				obj: getGoedbTableTest1Value(),
			},
			wantColumnName:  []string{"Name", "TestTableName"},
			wantColumnValue: []string{"'TestTableWithFK-Name'", "'TestTableName-Name-ID'"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotColumnName, gotColumnValue, err := getPrimaryKeysAndValues(tt.args.gt, tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("getPrimaryKeysAndValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotColumnName, tt.wantColumnName) {
				t.Errorf("getPrimaryKeysAndValues() gotColumnName = %v, want %v", gotColumnName, tt.wantColumnName)
			}
			if !reflect.DeepEqual(gotColumnValue, tt.wantColumnValue) {
				t.Errorf("getPrimaryKeysAndValues() gotColumnValue = %v, want %v", gotColumnValue, tt.wantColumnValue)
			}
		})
	}
}

func Test_getColumnsAndValues(t *testing.T) {
	type args struct {
		table    Table
		instance interface{}
	}
	tests := []struct {
		name        string
		args        args
		wantColumns []string
		wantValues  []string
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			name: "GetColumnsAndValuesSQL",
			args: args{
				table:    getGoedbTableTest1(),
				instance: getGoedbTableTest1Value(),
			},
			wantColumns: []string{"Name", "TestTableName", "Desc"},
			wantValues:  []string{"'TestTableWithFK-Name'", "'TestTableName-Name-ID'", "'testing description'"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotColumns, gotValues, err := getColumnsAndValues(tt.args.table, tt.args.instance)
			if (err != nil) != tt.wantErr {
				t.Errorf("getColumnsAndValuesSQL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotColumns, tt.wantColumns) {
				t.Errorf("getColumnsAndValuesSQL() gotColumns = %v, want %v", gotColumns, tt.wantColumns)
			}
			if !reflect.DeepEqual(gotValues, tt.wantValues) {
				t.Errorf("getColumnsAndValuesSQL() gotValues = %v, want %v", gotValues, tt.wantValues)
			}
		})
	}
}

func TestTransientSQLDialect_First(t *testing.T) {
	type args struct {
		table    Table
		where    string
		instance interface{}
	}
	tests := []struct {
		name    string
		dialect *TransientSQLDialect
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "TransientSQLDialect_First_WithoutWhere",
			args: args{
				table:    getGoedbTableTest1(),
				instance: getGoedbTableTest1Value(),
			},
			want:    "SELECT TestTableWithFK.Name,TestTable.ID,TestTable.Name,TestTableWithFK.Desc FROM TestTableWithFK,TestTable WHERE TestTableWithFK.Name='TestTableWithFK-Name' AND TestTableWithFK.TestTableName='TestTableName-Name-ID' AND TestTableWithFK.TestTableName = TestTable.Name",
			wantErr: false,
			dialect: &TransientSQLDialect{Models: getGoedbTableMapTest()},
		},
		{
			name: "TransientSQLDialect_First_WithWhere",
			args: args{
				table:    getGoedbTableTest1(),
				instance: getGoedbTableTest1Value(),
				where:    "TestTableWithFK.Desc = 'description1'",
			},
			want:    "SELECT TestTableWithFK.Name,TestTable.ID,TestTable.Name,TestTableWithFK.Desc FROM TestTableWithFK,TestTable WHERE TestTableWithFK.Desc = 'description1' AND TestTableWithFK.TestTableName = TestTable.Name",
			wantErr: false,
			dialect: &TransientSQLDialect{Models: getGoedbTableMapTest()},
		},
		{
			name: "TransientSQLDialect_First_NoModelFound",
			args: args{
				table:    getGoedbTableTest1(),
				instance: getGoedbTableTest1Value(),
			},
			want:    "",
			wantErr: true,
			dialect: &TransientSQLDialect{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dialect := tt.dialect
			got, err := dialect.First(tt.args.table, tt.args.where, tt.args.instance)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransientSQLDialect.First() error = [%v], wantErr [%v]", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TransientSQLDialect.First() = [%v], want [%v]", got, tt.want)
			}
		})
	}
}

func TestTransientSQLDialect_Find(t *testing.T) {
	type fields struct {
		Models map[string]Table
	}
	type args struct {
		table    Table
		where    string
		instance interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "TransientSQLDialect_Find_WithoutWhere",
			args: args{
				table:    getGoedbTableTest1(),
				instance: getGoedbTableTest1Value(),
			},
			want:    "SELECT TestTableWithFK.Name,TestTable.ID,TestTable.Name,TestTableWithFK.Desc FROM TestTableWithFK,TestTable WHERE TestTableWithFK.TestTableName = TestTable.Name",
			wantErr: false,
			fields: fields{
				Models: getGoedbTableMapTest(),
			},
		},
		{
			name: "TransientSQLDialect_Find_WithWhere",
			args: args{
				table:    getGoedbTableTest1(),
				instance: getGoedbTableTest1Value(),
				where:    "TestTableWithFK.Desc = 'description1'",
			},
			want:    "SELECT TestTableWithFK.Name,TestTable.ID,TestTable.Name,TestTableWithFK.Desc FROM TestTableWithFK,TestTable WHERE TestTableWithFK.Desc = 'description1' AND TestTableWithFK.TestTableName = TestTable.Name",
			wantErr: false,
			fields: fields{
				Models: getGoedbTableMapTest(),
			},
		},
		{
			name: "TransientSQLDialect_Find_WithModelNotFound",
			args: args{
				table:    getGoedbTableTest1(),
				instance: getGoedbTableTest1Value(),
				where:    "TestTableWithFK.Desc = 'description1'",
			},
			want:    "SELECT TestTableWithFK.Name,TestTable.ID,TestTable.Name,TestTableWithFK.Desc FROM TestTableWithFK,TestTable WHERE TestTableWithFK.Desc = 'description1' AND TestTableWithFK.TestTableName = TestTable.Name",
			wantErr: false,
			fields: fields{
				Models: getGoedbTableMapTest(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dialect := &TransientSQLDialect{
				Models: tt.fields.Models,
			}
			got, err := dialect.Find(tt.args.table, tt.args.where, tt.args.instance)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransientSQLDialect.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TransientSQLDialect.Find() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransientSQLDialect_Update(t *testing.T) {
	type fields struct {
		Models map[string]Table
	}
	type args struct {
		table    Table
		instance interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "TransientSQLDialect_Update",
			args: args{
				table:    getGoedbTableTest1(),
				instance: getGoedbTableTest1Value(),
			},
			want:    "UPDATE TestTableWithFK SET Name = 'TestTableWithFK-Name',TestTableName = 'TestTableName-Name-ID',Desc = 'testing description' WHERE TestTableWithFK.Name='TestTableWithFK-Name' AND TestTableWithFK.TestTableName='TestTableName-Name-ID'",
			wantErr: false,
			fields: fields{
				Models: getGoedbTableMapTest(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dialect := &TransientSQLDialect{
				Models: tt.fields.Models,
			}
			got, err := dialect.Update(tt.args.table, tt.args.instance)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransientSQLDialect.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TransientSQLDialect.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransientSQLDialect_Delete(t *testing.T) {
	type fields struct {
		Models map[string]Table
	}
	type args struct {
		table    Table
		where    string
		instance interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "TransientSQLDialect_Delete_WithoutWhere",
			args: args{
				table:    getGoedbTableTest1(),
				instance: getGoedbTableTest1Value(),
			},
			want:    "DELETE FROM TestTableWithFK WHERE Name='TestTableWithFK-Name' AND TestTableName='TestTableName-Name-ID'",
			wantErr: false,
			fields: fields{
				Models: getGoedbTableMapTest(),
			},
		},
		{
			name: "TransientSQLDialect_Delete_WithWhere",
			args: args{
				table:    getGoedbTableTest1(),
				instance: getGoedbTableTest1Value(),
				where:    "Name='TestTableWithFK-Name'",
			},
			want:    "DELETE FROM TestTableWithFK WHERE Name='TestTableWithFK-Name'",
			wantErr: false,
			fields: fields{
				Models: getGoedbTableMapTest(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dialect := &TransientSQLDialect{
				Models: tt.fields.Models,
			}
			got, err := dialect.Delete(tt.args.table, tt.args.where, tt.args.instance)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransientSQLDialect.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TransientSQLDialect.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransientSQLDialect_Drop(t *testing.T) {
	type fields struct {
		Models map[string]Table
	}
	type args struct {
		tableName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
		{
			name: "TransientSQLDialect_Drop",
			args: args{
				tableName: "TableToRemove",
			},
			want: "DROP TABLE TableToRemove",
			fields: fields{
				Models: getGoedbTableMapTest(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dialect := &TransientSQLDialect{
				Models: tt.fields.Models,
			}
			if got := dialect.Drop(tt.args.tableName); got != tt.want {
				t.Errorf("TransientSQLDialect.Drop() = %v, want %v", got, tt.want)
			}
		})
	}
}
