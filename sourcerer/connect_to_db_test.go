package sourcerer

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gwenwindflower/tbd/shared"
)

func TestConnectToDbSnowflake(t *testing.T) {
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
	sfc, ok := conn.(*SfConn)
	if !ok {
		t.Errorf("conn not of type SfConn: %v", err)
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	sfc.Db = db
	defer sfc.Db.Close()
	mock.ExpectBegin()
	if _, err := sfc.Db.Begin(); err != nil {
		t.Errorf("error '%s' was not expected, while pinging db", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestConnectToDbPostgres(t *testing.T) {
	cd := shared.ConnectionDetails{
		ConnType: "postgres",
		Host:     "localhost",
		Port:     5432,
		Username: "frodo",
		Password: "0nering",
		Database: "shire",
		Schema:   "hobbiton",
	}
	conn, err := GetConn(cd)
	if err != nil {
		t.Errorf("GetConn failed: %v", err)
	}
	pgc, ok := conn.(*PgConn)
	if !ok {
		t.Errorf("conn not of type PgConn: %v", err)
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	pgc.Db = db
	defer pgc.Db.Close()
	mock.ExpectBegin()
	if _, err := pgc.Db.Begin(); err != nil {
		t.Errorf("error '%s' was not expected, while pinging db", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
