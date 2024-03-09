package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/snowflakedb/gosnowflake"
	"gopkg.in/yaml.v2"
)

var (
	dbUsername = os.Getenv("SNOWFLAKE_SANDBOX_USER")
	dbPassword = os.Getenv("SNOWFLAKE_SANDBOX_PASSWORD")
	dbAccount  = os.Getenv("SNOWFLAKE_SANDBOX_ACCOUNT")
	dbSchema   = "DBT_WINNIE"
	dbDatabase = "ANALYTICS"
)

type Source struct {
	Name    string   `yaml:"name"`
	Columns []string `yaml:"columns"`
}

type Sources struct {
	Sources []Source `yaml:"sources"`
}

func main() {
	connStr := fmt.Sprintf("%s:%s@%s/%s/%s", dbUsername, dbPassword, dbAccount, dbDatabase, dbSchema)

	ctx := context.Background()
	db, err := sql.Open("snowflake", connStr)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.QueryContext(ctx, "SELECT table_name FROM information_schema.tables where table_schema = 'DBT_WINNIE'")
	if err != nil {
		log.Fatal(err)
	}

	tables := []Source{}
	for rows.Next() {
		var table Source
		if err := rows.Scan(&table); err != nil {
			log.Fatal(err)
		}
		tables = append(tables, table)
	}
	rows.Close()
	var wg sync.WaitGroup
	wg.Add(len(tables))

	tableData := make(map[string][]string)
	mutex := sync.Mutex{}

	for _, table := range tables {
		go func(table string) {
			defer wg.Done()

			columns, err := getColumnsForTable(db, ctx, table)
			if err != nil {
				log.Printf("Error fetching columns for table %s: %v\n", table, err)
				return
			}

			mutex.Lock()
			tableData[table] = columns
			mutex.Unlock()
		}(table)
	}

	wg.Wait()

	yamlData, err := yaml.Marshal(tableData)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(yamlData))
}

func getColumnsForTable(db *sql.DB, ctx context.Context, table string) ([]string, error) {
	var columns []string

	query := fmt.Sprintf("SELECT column_name FROM information_schema.columns WHERE table_name = '%s' AND table_schema = '%s'", table, dbSchema)
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var columnName string
		if err := rows.Scan(&columnName); err != nil {
			return nil, err
		}
		columns = append(columns, columnName)
	}

	return columns, nil
}
