package sourcerer

import (
	"context"
	"fmt"
	"log"

	"github.com/gwenwindflower/tbd/shared"

	"google.golang.org/api/iterator"
)

func (sfc *SfConn) GetSourceTables(ctx context.Context) (shared.SourceTables, error) {
	ts := shared.SourceTables{}
	defer sfc.Cancel()
	rows, err := sfc.Db.QueryContext(ctx, fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE table_schema = '%s'", sfc.Schema))
	if err != nil {
		log.Fatalf("Error fetching tables: %v\n", err)
	}
	defer rows.Close()
	for rows.Next() {
		var table shared.SourceTable
		if err := rows.Scan(&table.Name); err != nil {
			log.Fatalf("Error scanning tables: %v\n", err)
		}
		table.Schema = sfc.Schema
		ts.SourceTables = append(ts.SourceTables, table)
	}
	return ts, nil
}

func (bqc *BqConn) GetSourceTables(ctx context.Context) (shared.SourceTables, error) {
	ts := shared.SourceTables{}
	defer bqc.Cancel()
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
		ts.SourceTables = append(ts.SourceTables, shared.SourceTable{Name: table.TableID, Schema: bqc.Dataset})
	}
	return ts, nil
}

func (dc *DuckConn) GetSourceTables(ctx context.Context) (shared.SourceTables, error) {
	ts := shared.SourceTables{}
	defer dc.Cancel()
	q := fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE table_schema = '%s'", dc.Schema)
	rows, err := dc.Db.QueryContext(ctx, q)
	if err != nil {
		log.Fatalf("Error fetching tables: %v\n", err)
	}
	defer rows.Close()
	for rows.Next() {
		var table shared.SourceTable
		if err := rows.Scan(&table.Name); err != nil {
			log.Fatalf("Error scanning tables: %v\n", err)
		}
		table.Schema = dc.Schema
		ts.SourceTables = append(ts.SourceTables, table)
	}
	return ts, nil
}

func (pgc *PgConn) GetSourceTables(ctx context.Context) (shared.SourceTables, error) {
	ts := shared.SourceTables{}
	defer pgc.Cancel()
	q := fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE table_schema = '%s'", pgc.Schema)
	rows, err := pgc.Db.QueryContext(ctx, q)
	if err != nil {
		log.Fatalf("Error fetching tables: %v\n", err)
	}
	defer rows.Close()
	for rows.Next() {
		var table shared.SourceTable
		if err := rows.Scan(&table.Name); err != nil {
			log.Fatalf("Error scanning tables: %v\n", err)
		}
		table.Schema = pgc.Schema
		ts.SourceTables = append(ts.SourceTables, table)
	}
	return ts, nil
}
