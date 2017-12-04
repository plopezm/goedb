package database

import (
	"reflect"
	"testing"

	"github.com/plopezm/goedb/database/specifics"

	"github.com/plopezm/goedb/database/models"
)

func TestSQLDialect_Create(t *testing.T) {
	type args struct {
		table models.Table
	}
	tests := [...]struct {
		name      string
		dialect   *SQLDialect
		specifics specifics.Specifics
		args      args
		want      string
	}{
		// TODO: Add test cases.
		{
			name: "TestCreateTable",
			dialect: &SQLDialect{
				Specifics: new(specifics.SQLiteSpecifics),
			},
			specifics: new(specifics.SQLiteSpecifics),
			args: args{
				table: models.Table{
					Name: "Table1",
					Columns: []models.Column{
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
			name: "TestCreateTableWithPrimaryKeys",
			dialect: &SQLDialect{
				Specifics: new(specifics.SQLiteSpecifics),
			},
			args: args{
				table: models.Table{
					Name: "Table1",
					Columns: []models.Column{
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
			dialect := tt.dialect
			if got := dialect.Create(tt.args.table); got != tt.want {
				t.Errorf("SQLDialect.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getGoedbTableTest1() models.Table {

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

func getGoedbTableTest2() models.Table {

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

func getGoedbTableMapTest() (modelMap map[string]models.Table) {
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
	modelMap = make(map[string]models.Table)
	modelMap["TestTable"] = ParseModel(&TestTable{})
	modelMap["TestTableWithFK"] = ParseModel(&TestTableWithFK{})
	modelMap["TestTableWithFK2"] = ParseModel(&TestTableWithFK2{})
	return modelMap
}

func Test_generateSQLQuery(t *testing.T) {
	type args struct {
		table    models.Table
		modelMap map[string]models.Table
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
		gt  models.Table
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
		table    models.Table
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

func TestSQLDialect_First(t *testing.T) {
	type args struct {
		table    models.Table
		where    string
		instance interface{}
	}
	tests := []struct {
		name    string
		dialect *SQLDialect
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "SQLDialect_First_WithoutWhere",
			args: args{
				table:    getGoedbTableTest1(),
				instance: getGoedbTableTest1Value(),
			},
			want:    "SELECT TestTableWithFK.Name,TestTable.ID,TestTable.Name,TestTableWithFK.Desc FROM TestTableWithFK,TestTable WHERE TestTableWithFK.Name='TestTableWithFK-Name' AND TestTableWithFK.TestTableName='TestTableName-Name-ID' AND TestTableWithFK.TestTableName = TestTable.Name",
			wantErr: false,
			dialect: &SQLDialect{Models: getGoedbTableMapTest()},
		},
		{
			name: "SQLDialect_First_WithWhere",
			args: args{
				table:    getGoedbTableTest1(),
				instance: getGoedbTableTest1Value(),
				where:    "TestTableWithFK.Desc = 'description1'",
			},
			want:    "SELECT TestTableWithFK.Name,TestTable.ID,TestTable.Name,TestTableWithFK.Desc FROM TestTableWithFK,TestTable WHERE TestTableWithFK.Desc = 'description1' AND TestTableWithFK.TestTableName = TestTable.Name",
			wantErr: false,
			dialect: &SQLDialect{Models: getGoedbTableMapTest()},
		},
		{
			name: "SQLDialect_First_NoModelFound",
			args: args{
				table:    getGoedbTableTest1(),
				instance: getGoedbTableTest1Value(),
			},
			want:    "",
			wantErr: true,
			dialect: &SQLDialect{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dialect := tt.dialect
			got, err := dialect.First(tt.args.table, tt.args.where, tt.args.instance)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLDialect.First() error = [%v], wantErr [%v]", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SQLDialect.First() = [%v], want [%v]", got, tt.want)
			}
		})
	}
}

func TestSQLDialect_Find(t *testing.T) {
	type fields struct {
		Models map[string]models.Table
	}
	type args struct {
		table    models.Table
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
			name: "SQLDialect_Find_WithoutWhere",
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
			name: "SQLDialect_Find_WithWhere",
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
			name: "SQLDialect_Find_WithModelNotFound",
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
			dialect := &SQLDialect{
				Models: tt.fields.Models,
			}
			got, err := dialect.Find(tt.args.table, tt.args.where, tt.args.instance)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLDialect.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SQLDialect.Find() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLDialect_Update(t *testing.T) {
	type fields struct {
		Models map[string]models.Table
	}
	type args struct {
		table    models.Table
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
			name: "SQLDialect_Update",
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
			dialect := &SQLDialect{
				Models: tt.fields.Models,
			}
			got, err := dialect.Update(tt.args.table, tt.args.instance)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLDialect.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SQLDialect.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLDialect_Delete(t *testing.T) {
	type fields struct {
		Models map[string]models.Table
	}
	type args struct {
		table    models.Table
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
			name: "SQLDialect_Delete_WithoutWhere",
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
			name: "SQLDialect_Delete_WithWhere",
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
			dialect := &SQLDialect{
				Models: tt.fields.Models,
			}
			got, err := dialect.Delete(tt.args.table, tt.args.where, tt.args.instance)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLDialect.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SQLDialect.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLDialect_Drop(t *testing.T) {
	type fields struct {
		Models map[string]models.Table
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
			name: "SQLDialect_Drop",
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
			dialect := &SQLDialect{
				Models: tt.fields.Models,
			}
			if got := dialect.Drop(tt.args.tableName); got != tt.want {
				t.Errorf("SQLDialect.Drop() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLDialect_Insert(t *testing.T) {
	type fields struct {
		Models map[string]models.Table
	}
	type args struct {
		table    models.Table
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
			name: "SQLDialect_Delete_WithWhere",
			args: args{
				table:    getGoedbTableTest1(),
				instance: getGoedbTableTest1Value(),
			},
			want:    "INSERT INTO TestTableWithFK (Name,TestTableName,Desc) values('TestTableWithFK-Name','TestTableName-Name-ID','testing description')",
			wantErr: false,
			fields: fields{
				Models: getGoedbTableMapTest(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dialect := &SQLDialect{
				Models: tt.fields.Models,
			}
			got, err := dialect.Insert(tt.args.table, tt.args.instance)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLDialect.Insert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SQLDialect.Insert() = %v, want %v", got, tt.want)
			}
		})
	}
}
