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

func getGoedbTableTest() Table {

	type TestTable struct {
		ID   uint64 `goedb:"pk,autoincrement"`
		Name string `goedb:"unique"`
	}

	type TestTableWithFK struct {
		Name          string    `goedb:"pk"`
		TestTableName TestTable `goedb:"pk,fk=TestTable(Name)"`
		Ignorable     bool      `goedb:"ignore"`
	}

	return ParseModel(&TestTableWithFK{})
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
	}

	modelMap = make(map[string]Table)
	modelMap["TestTable"] = ParseModel(&TestTable{})
	modelMap["TestTableWithFK"] = ParseModel(&TestTableWithFK{})
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
				table:    getGoedbTableTest(),
				modelMap: getGoedbTableMapTest(),
			},
			wantQuery:       "SELECT TestTableWithFK.Name,TestTable.ID,TestTable.Name FROM TestTableWithFK,TestTable",
			wantConstraints: " AND TestTableWithFK.TestTableName = TestTable.Name",
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
