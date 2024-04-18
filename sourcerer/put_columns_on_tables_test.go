package sourcerer

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gwenwindflower/tbd/shared"
)

func TestPutColumnsOnTables(t *testing.T) {
	ctx := context.Background()
	ts := shared.SourceTables{
		SourceTables: []shared.SourceTable{
			{
				Name: "table1",
			},
		},
	}
	cd := shared.ConnectionDetails{
		ConnType: "snowflake",
		Account:  "dunedain.snowflakecomputing.com",
		Username: "aragorn",
		Database: "gondor",
		Schema:   "minas-tirith",
	}
	conn, err := GetConn(cd)
	if err != nil {
		t.Errorf("GetConn failed: %v", err)
	}
	if conn == nil {
		t.Errorf("GetConn failed: conn is nil")
	}
	SfConn, ok := conn.(*SfConn)
	if !ok {
		t.Errorf("GetConn failed: conn is not of type SfConn")
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	SfConn.Db = db
	rows := sqlmock.NewRows([]string{"column_name", "data_type"}).AddRow("column1", "text").AddRow("column2", "char").AddRow("COLUMN3", "int")
	mock.ExpectQuery("SELECT column_name, data_type FROM information_schema.columns WHERE table_schema = 'MINAS-TIRITH' AND table_name = 'table1'").WillReturnRows(rows)
	err = PutColumnsOnTables(ctx, ts, SfConn)
	if err != nil {
		t.Errorf("PutColumnsOnTables failed: %v", err)
	}
	if len(ts.SourceTables[0].Columns) != 3 {
		t.Errorf("PutColumnsOnTables failed: expected 3 columns, got %d", len(ts.SourceTables[0].Columns))
	}
	if ts.SourceTables[0].Columns[0].Name != "column1" {
		t.Errorf("PutColumnsOnTables failed: expected column name column1, got %s", ts.SourceTables[0].Columns[0].Name)
	}
	if ts.SourceTables[0].Columns[0].DataType != "text" {
		t.Errorf("PutColumnsOnTables failed: expected column data type text, got %s", ts.SourceTables[0].Columns[0].DataType)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %v", err)
	}
}
