package dialect

import (
	"reflect"
	"testing"

	"github.com/plopezm/goedb/database/models"
)

func TestPostgresSpecifics_GetSQLCreateTableColumn(t *testing.T) {
	type args struct {
		value models.Column
	}
	tests := []struct {
		name              string
		dialect           *PostgresDialect
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
			wantSQLColumnLine: "PKColumn SERIAL,",
			wantPrimaryKey:    "PKColumn,",
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
			dialect := &PostgresDialect{}
			got, got1, got2, err := dialect.GetSQLCreateTableColumn(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostgresDialect.GetSQLCreateTableColumn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantSQLColumnLine {
				t.Errorf("PostgresDialect.GetSQLCreateTableColumn() got = %v, want %v", got, tt.wantSQLColumnLine)
			}
			if got1 != tt.wantPrimaryKey {
				t.Errorf("PostgresDialect.GetSQLCreateTableColumn() got1 = %v, want %v", got1, tt.wantPrimaryKey)
			}
			if got2 != tt.wantConstraints {
				t.Errorf("PostgresDialect.GetSQLCreateTableColumn() got2 = %v, want %v", got2, tt.wantConstraints)
			}
		})
	}
}
