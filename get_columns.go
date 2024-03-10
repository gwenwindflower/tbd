package main

import (
	"context"
	"database/sql"
	"fmt"
)

func GetColumns(db *sql.DB, ctx context.Context, table SourceTable) ([]Column, error) {
	var columns []Column

	query := fmt.Sprintf("SELECT column_name, data_type FROM information_schema.columns WHERE table_schema = '%s' AND table_name = '%s'", dbSchema, table.Name)
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		column := Column{}
		if err := rows.Scan(&column.Name, &column.DataType); err != nil {
			return nil, err
		}
		columns = append(columns, column)
	}

	return columns, nil
}
