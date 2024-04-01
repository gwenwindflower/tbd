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

	switch connectionDetails.ConnType {
	case "snowflake":
		{
			dbConn := SnowflakeConnection{}
			// Snowflake requires uppercase for all identifiers
			connectionDetails.Username = strings.ToUpper(connectionDetails.Username)
			connectionDetails.Account = strings.ToUpper(connectionDetails.Account)
			connectionDetails.Database = strings.ToUpper(connectionDetails.Database)
			connectionDetails.Schema = strings.ToUpper(connectionDetails.Schema)

			db, cancel, err := dbConn.ConnectToDB(ctx, connectionDetails)
			defer cancel()
			if err != nil {
				log.Fatalf("Couldn't connect to database: %v\n", err)
			}
			rows, err := db.QueryContext(ctx, fmt.Sprintf("SELECT table_name FROM information_schema.tables where table_schema = '%s'", connectionDetails.Schema))
			if err != nil {
				log.Fatalf("Error fetching tables: %v\n", err)
			}
			defer rows.Close()
			for rows.Next() {
				var table shared.SourceTable
				if err := rows.Scan(&table.Name); err != nil {
					log.Fatalf("Error scanning tables: %v\n", err)
				}
				tables.SourceTables = append(tables.SourceTables, table)
			}
			PutColumnsOnTables(ctx, db, tables, connectionDetails)
		}
	default:
		{
			return tables, fmt.Errorf("unsupported warehouse: %s", connectionDetails.ConnType)
		}
	}
	return tables, nil
}
