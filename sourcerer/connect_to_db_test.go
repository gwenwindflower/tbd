package sourcerer

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gwenwindflower/tbd/shared"
)

func TestConnectToDb(t *testing.T) {
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
	SfConn, ok := conn.(*SfConn)
	if !ok {
		t.Errorf("conn not of type SfConn: %v", err)
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	SfConn.Db = db
	defer SfConn.Db.Close()
	mock.ExpectBegin()
	if _, err := SfConn.Db.Begin(); err != nil {
		t.Errorf("error '%s' was not expected, while pinging db", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
