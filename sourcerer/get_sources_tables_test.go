package sourcerer

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gwenwindflower/tbd/shared"
)

func TestGetSourceTablesSnowflake(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
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
	SfConn.Db = db
	SfConn.Cancel = cancel
	defer SfConn.Db.Close()
	q := fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE table_schema = '%s'", SfConn.Schema)
	mock.ExpectQuery(q).WillReturnRows(sqlmock.NewRows([]string{"table_name"}).AddRow("table1").AddRow("table2"))
	ts, err := SfConn.GetSourceTables(ctx)
	if err != nil {
		t.Errorf("GetSources failed: %v", err)
	}
	if len(ts.SourceTables) != 2 {
		t.Errorf("GetSources failed: expected 2 sources, got %d", len(ts.SourceTables))
	}
	if ts.SourceTables[0].Name != "table1" {
		t.Errorf("GetSources failed: expected source name %s, got %s", "table1", ts.SourceTables[0].Name)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
