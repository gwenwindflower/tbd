package main

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestGetTables(t *testing.T) {
	dbUsername = os.Getenv("SNOWFLAKE_SANDBOX_USERNAME")
	dbAccount = os.Getenv("SNOWFLAKE_SANDBOX_ACCOUNT")
	dbDatabase = "ANALYTICS"
	dbSchema = "DBT_WINNIE"
	databaseType := "snowflake"
	connStr := fmt.Sprintf("%s@%s/%s/%s?authenticator=externalbrowser", dbUsername, dbAccount, dbDatabase, dbSchema)
	ctx, db, err := ConnectToDB(connStr, databaseType)
	if err != nil {
		log.Fatal(err)
	}
	tables, err := GetTables(db, ctx)
	if err != nil {
		t.Errorf("Error getting tables: %v", err)
	}
	if len(tables.SourceTables) == 0 {
		t.Errorf("No tables found")
	} else {
		t.Logf("%v tables found.", len(tables.SourceTables))
	}
}
