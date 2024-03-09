package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

func GetTables(db *sql.DB, ctx context.Context) (SourceTables, error) {
	tables := SourceTables{}

	rows, err := db.QueryContext(ctx, fmt.Sprintf("SELECT table_name FROM information_schema.tables where table_schema = '%s'", dbSchema))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var table SourceTable
		if err := rows.Scan(&table.Name); err != nil {
			return tables, err
		}
		tables.SourceTables = append(tables.SourceTables, table)
	}

	return tables, nil
}
