package main

import (
	"context"
	"database/sql"
	"fmt"
)

func GetTables(db *sql.DB, ctx context.Context, schema string) (SourceTables, error) {
	tables := SourceTables{}

	rows, err := db.QueryContext(ctx, fmt.Sprintf("SELECT table_name FROM information_schema.tables where table_schema = '%s'", schema))
	if err != nil {
		return tables, err
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
