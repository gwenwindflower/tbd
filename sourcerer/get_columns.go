package sourcerer

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"tbd/shared"
)

func GetColumns(db *sql.DB, ctx context.Context, table shared.SourceTable, connectionDetails shared.ConnectionDetails) ([]shared.Column, error) {
	var columns []shared.Column

	query := fmt.Sprintf("SELECT column_name, data_type FROM information_schema.columns WHERE table_schema = '%s' AND table_name = '%s'", connectionDetails.Schema, table.Name)
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Fatalf("Error fetching columns for table %s: %v\n", table.Name, err)
	}
	defer rows.Close()

	for rows.Next() {
		column := shared.Column{}
		if err := rows.Scan(&column.Name, &column.DataType); err != nil {
			return nil, err
		}
		columns = append(columns, column)
	}

	return columns, nil
}
