package specifics

import (
	"reflect"
	"testing"

	"github.com/plopezm/goedb/database/models"
)

func TestSQLiteSpecifics_GetSQLCreateTableColumn(t *testing.T) {
	type args struct {
		value models.Column
	}
	tests := []struct {
		name              string
		specifics         *SQLiteSpecifics
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
				value: models.Column{
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
				value: models.Column{
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
				value: models.Column{
					Title:      "PKColumn",
					PrimaryKey: true,
					ColumnType: reflect.Int,
					ForeignKey: models.ForeignKey{IsForeignKey: true, ForeignKeyTableReference: "OtherTable", ForeignKeyColumnReference: "OtherTablePK"},
				},
			},
			wantSQLColumnLine: "PKColumn INTEGER,",
			wantPrimaryKey:    "PKColumn,",
			wantConstraints:   ", FOREIGN KEY (PKColumn) REFERENCES OtherTable(OtherTablePK) ON DELETE CASCADE",
		},
		{
			name: "TestErrorTypeNotFound",
			args: args{
				value: models.Column{
					Title:      "PKColumn",
					PrimaryKey: true,
					ColumnType: reflect.Struct,
					ForeignKey: models.ForeignKey{IsForeignKey: true, ForeignKeyTableReference: "OtherTable", ForeignKeyColumnReference: "OtherTablePK"},
				},
			},
			wantErr: true,
		},
		{
			name: "TestStringColumnUnique",
			args: args{
				value: models.Column{
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
				value: models.Column{
					Title:      "NormalColumnString",
					ColumnType: reflect.String,
				},
			},
			wantSQLColumnLine: "NormalColumnString VARCHAR,",
		},
		{
			name: "TestNormalStringColumn",
			args: args{
				value: models.Column{
					Title:      "NormalColumnString",
					ColumnType: reflect.Float64,
				},
			},
			wantSQLColumnLine: "NormalColumnString FLOAT,",
		},
		{
			name: "TestNormalStringColumn",
			args: args{
				value: models.Column{
					Title:      "NormalColumnString",
					ColumnType: reflect.Bool,
				},
			},
			wantSQLColumnLine: "NormalColumnString BOOLEAN,",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			specifics := &SQLiteSpecifics{}
			gotSQLColumnLine, gotPrimaryKey, gotConstraints, err := specifics.GetSQLCreateTableColumn(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLiteSpecifics.GetSQLCreateTableColumn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSQLColumnLine != tt.wantSQLColumnLine {
				t.Errorf("SQLiteSpecifics.GetSQLCreateTableColumn() gotSqlColumnLine = %v, want %v", gotSQLColumnLine, tt.wantSQLColumnLine)
			}
			if gotPrimaryKey != tt.wantPrimaryKey {
				t.Errorf("SQLiteSpecifics.GetSQLCreateTableColumn() gotPrimaryKey = %v, want %v", gotPrimaryKey, tt.wantPrimaryKey)
			}
			if gotConstraints != tt.wantConstraints {
				t.Errorf("SQLiteSpecifics.GetSQLCreateTableColumn() gotConstraints = %v, want %v", gotConstraints, tt.wantConstraints)
			}
		})
	}
}
