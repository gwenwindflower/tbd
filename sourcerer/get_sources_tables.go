package sourcerer

import (
	"context"
	"fmt"
	"log"

	"github.com/gwenwindflower/tbd/shared"

	dbsql "github.com/databricks/databricks-sql-go"
	"google.golang.org/api/iterator"
)

func (sfc *SfConn) GetSourceTables(ctx context.Context) (shared.SourceTables, error) {
	ts := shared.SourceTables{}
	defer sfc.Cancel()
	// TODO: why doesn't this work?
	// q := `SELECT table_name FROM information_schema.tables WHERE table_schema = '?'`
	q := fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE table_schema = '%s'", sfc.Schema)
	rows, err := sfc.Db.QueryContext(ctx, q)
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
	q := "SELECT table_name FROM information_schema.tables WHERE table_schema = '?'"
	rows, err := dc.Db.QueryContext(ctx, q, dc.Schema)
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
	q := "SELECT table_name FROM information_schema.tables WHERE table_schema = '$1'"
	rows, err := pgc.Db.QueryContext(ctx, q, pgc.Schema)
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

func (dbxc *DbxConn) GetSourceTables(ctx context.Context) (shared.SourceTables, error) {
	ts := shared.SourceTables{}
	defer dbxc.Cancel()
	q := "SELECT table_name FROM information_schema.tables WHERE table_schema = :p_schema"
	rows, err := dbxc.Db.QueryContext(ctx, q, dbsql.Parameter{Name: "p_schema", Value: dbxc.Schema})
	if err != nil {
		log.Fatalf("Error fetching tables: %v\n", err)
	}
	defer rows.Close()
	for rows.Next() {
		var table shared.SourceTable
		if err := rows.Scan(&table.Name); err != nil {
			log.Fatalf("Error scanning tables: %v\n", err)
		}
		table.Schema = dbxc.Schema
		ts.SourceTables = append(ts.SourceTables, table)
	}
	return ts, nil
}
