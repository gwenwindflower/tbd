package sourcerer

import (
	"context"
	"fmt"
	"log"
	"strings"
	"tbd/shared"
)

func GetSources(ctx context.Context, connectionDetails shared.ConnectionDetails) (shared.SourceTables, error) {
	tables := shared.SourceTables{}

	switch connectionDetails.Warehouse {
	case "snowflake":
		{
			dbConn := SnowflakeConnection{}
			db, cancel, err := dbConn.ConnectToDB(ctx, connectionDetails)
			defer cancel()
			if err != nil {
				log.Fatalf("couldn't connect to database: %v", err)
			}
			rows, err := db.QueryContext(ctx, fmt.Sprintf("SELECT table_name FROM information_schema.tables where table_schema = '%s'", strings.ToUpper(connectionDetails.Schema)))
			if err != nil {
				return tables, err
			}
			defer rows.Close()
			for rows.Next() {
				var table shared.SourceTable
				if err := rows.Scan(&table.Name); err != nil {
					return tables, err
				}
				tables.SourceTables = append(tables.SourceTables, table)
			}

			PutColumnsOnTables(ctx, db, tables, connectionDetails)
		}
	}
	return tables, nil
}
