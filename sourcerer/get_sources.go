package sourcerer

import (
	"context"
	"fmt"
	"log"
	"tbd/shared"

	"google.golang.org/api/iterator"
)

func (sfc *SfConn) GetSources(ctx context.Context) (shared.SourceTables, error) {
	ts := shared.SourceTables{}

	err := sfc.ConnectToDB(ctx)
	defer sfc.Cancel()
	if err != nil {
		log.Fatalf("Couldn't connect to database: %v\n", err)
	}
	rows, err := sfc.Db.QueryContext(ctx, fmt.Sprintf("SELECT table_name FROM information_schema.tables where table_schema = '%s'", sfc.Schema))
	if err != nil {
		log.Fatalf("Error fetching tables: %v\n", err)
	}
	defer rows.Close()
	for rows.Next() {
		var table shared.SourceTable
		if err := rows.Scan(&table.Name); err != nil {
			log.Fatalf("Error scanning tables: %v\n", err)
		}
		ts.SourceTables = append(ts.SourceTables, table)
	}
	sfc.PutColumnsOnTables(ctx, ts)

	return ts, nil
}

func (bqc *BqConn) GetSources(ctx context.Context) (shared.SourceTables, error) {
	ts := shared.SourceTables{}
	err := bqc.ConnectToDB(ctx)
	defer bqc.Cancel()
	if err != nil {
		log.Fatalf("Couldn't connect to database: %v\n", err)
	}
	bqDataset := bqc.Bq.Dataset(bqc.Dataset)
	tableIter := bqDataset.Tables(ctx)
	for {
		table, err := tableIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Error fetching tables: %v\n", err)
		}
		ts.SourceTables = append(ts.SourceTables, shared.SourceTable{Name: table.TableID})
	}
	bqc.PutColumnsOnTables(ctx, ts)
	return ts, nil
}
