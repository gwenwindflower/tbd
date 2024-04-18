package sourcerer

import (
	"testing"

	"github.com/gwenwindflower/tbd/shared"
)

func TestGetConnSnowflake(t *testing.T) {
	t.SkipNow()
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
	if SfConn.Account != "DUNEDAIN.SNOWFLAKECOMPUTING.COM" {
		t.Errorf("GetConn failed: Account is not correct")
	}
}

func TestGetConnBigQuery(t *testing.T) {
	cd := shared.ConnectionDetails{
		ConnType: "bigquery",
		Project:  "mirkwood",
		Dataset:  "hall_of_thranduil",
	}
	conn, err := GetConn(cd)
	if err != nil {
		t.Errorf("GetConn failed: %v", err)
	}
	if conn == nil {
		t.Errorf("GetConn failed: conn is nil")
	}
	BqConn, ok := conn.(*BqConn)
	if !ok {
		t.Errorf("GetConn failed: conn is not of type BqConn")
	}
	if BqConn.Dataset != "hall_of_thranduil" {
		t.Errorf("GetConn failed: Account is not correct")
	}
}

func TestGetConnDuckDB(t *testing.T) {
	cd := shared.ConnectionDetails{
		ConnType: "duckdb",
		Path:     "/path/to/duckdb.db",
		Database: "lothlorien",
		Schema:   "mallorn_trees",
	}
	conn, err := GetConn(cd)
	if err != nil {
		t.Errorf("GetConn failed: %v", err)
	}
	if conn == nil {
		t.Errorf("GetConn failed: conn is nil")
	}
	DuckConn, ok := conn.(*DuckConn)
	if !ok {
		t.Errorf("GetConn failed: conn is not of type DuckConn")
	}
	if DuckConn.Path != "/path/to/duckdb.db" {
		t.Errorf("GetConn failed: Account is not correct")
	}
}
