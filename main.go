package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/briandowns/spinner"
	_ "github.com/snowflakedb/gosnowflake"
)

// Connection details for Snowflake
var (
	dbUsername = os.Getenv("SNOWFLAKE_SANDBOX_USER")
	dbAccount  = os.Getenv("SNOWFLAKE_SANDBOX_ACCOUNT")
	dbSchema   = "DBT_WINNIE"
	dbDatabase = "ANALYTICS"
)

// Type definitions for the YAML file
type Column struct {
	Name     string `yaml:"name"`
	DataType string `yaml:"data_type"`
}

type SourceTable struct {
	Name    string   `yaml:"name"`
	Columns []Column `yaml:"columns"`
}

type SourceTables struct {
	SourceTables []SourceTable `yaml:"sources"`
}

func main() {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	log.Println("Writing tables and columns to sources YAML...")

	databaseType := "snowflake"
	connStr := fmt.Sprintf("%s@%s/%s/%s?authenticator=externalbrowser", dbUsername, dbAccount, dbDatabase, dbSchema)
	ctx, db, err := ConnectToDB(connStr, databaseType)
	if err != nil {
		log.Fatal(err)
	}
	tables, err := GetTables(db, ctx)
	if err != nil {
		log.Fatal(err)
	}
	PutColumnsOnTables(db, ctx, tables)
	WriteYAML(tables)

	log.Println("Writing staging models...")
	WriteStagingModels(tables)
	s.Stop()
}
