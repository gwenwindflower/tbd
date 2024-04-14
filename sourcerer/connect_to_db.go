package sourcerer

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/bigquery"
	_ "github.com/marcboeker/go-duckdb"
	_ "github.com/snowflakedb/gosnowflake"
)

func (sfc *SfConn) ConnectToDB(ctx context.Context) (err error) {
	connStr := fmt.Sprintf(
		"%s@%s/%s/%s?authenticator=externalbrowser",
		sfc.Username,
		sfc.Account,
		sfc.Database,
		sfc.Schema,
	)

	_, sfc.Cancel = context.WithTimeout(ctx, 1*time.Minute)
	defer sfc.Cancel()
	sfc.Db, err = sql.Open("snowflake", connStr)
	if err != nil {
		log.Fatalf("Could not connect to Snowflake %v\n", err)
	}
	return err
}

func (bqc *BqConn) ConnectToDB(ctx context.Context) (err error) {
	_, bqc.Cancel = context.WithTimeout(ctx, 1*time.Minute)
	defer bqc.Cancel()
	bqc.Bq, err = bigquery.NewClient(ctx, bqc.Project)
	if err != nil {
		log.Fatalf("Could not connect to BigQuery %v\n", err)
	}
	return err
}

func (dc *DuckConn) ConnectToDB(ctx context.Context) (err error) {
	_, dc.Cancel = context.WithTimeout(ctx, 1*time.Minute)
	defer dc.Cancel()
	if _, err := os.Stat(dc.Path); os.IsNotExist(err) {
		log.Fatalf("Path does not exist: %v\n", err)
	}
	dc.Db, err = sql.Open("duckdb", dc.Path)
	if err != nil {
		log.Fatalf("Could not connect to DuckDB %v\n", err)
	}
	return err
}
