package sourcerer

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gwenwindflower/tbd/shared"
)

func TestGetColumnsSnowflake(t *testing.T) {
	t.SkipNow()
	ctx := context.Background()
	st := shared.SourceTable{
		Name: "table1",
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
	SfConn.Db = db
	defer SfConn.Db.Close()
	q := fmt.Sprintf("SELECT column_name, data_type FROM information_schema.columns WHERE table_schema = '%s' AND table_name = '%s'", SfConn.Schema, st.Name)
	mock.ExpectQuery(q).WillReturnRows(sqlmock.NewRows([]string{"column_name", "data_type"}).AddRow("column1", "varchar").AddRow("column2", "varchar").AddRow("column3", "int"))
	cols, err := SfConn.GetColumns(ctx, st)
	if err != nil {
		t.Errorf("GetColumns failed: %v", err)
	}
	if len(cols) != 1 {
		t.Errorf("GetColumns failed: expected 1 column, got %d", len(cols))
	}
	if cols[0].Name != "column1" {
		t.Errorf("GetColumns failed: expected column name %s, got %s", "column1", cols[0].Name)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
