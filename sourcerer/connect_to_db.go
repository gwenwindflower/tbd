package sourcerer

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/bigquery"
	_ "github.com/databricks/databricks-sql-go"
	_ "github.com/lib/pq"
	_ "github.com/marcboeker/go-duckdb"
	_ "github.com/snowflakedb/gosnowflake"
)

func (sfc *SfConn) ConnectToDb(ctx context.Context) (err error) {
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

func (bqc *BqConn) ConnectToDb(ctx context.Context) (err error) {
	_, bqc.Cancel = context.WithTimeout(ctx, 1*time.Minute)
	defer bqc.Cancel()
	bqc.Bq, err = bigquery.NewClient(ctx, bqc.Project)
	if err != nil {
		log.Fatalf("Could not connect to BigQuery %v\n", err)
	}
	return err
}

func (dc *DuckConn) ConnectToDb(ctx context.Context) (err error) {
	_, dc.Cancel = context.WithTimeout(ctx, 1*time.Minute)
	defer dc.Cancel()
	if dc.Path == "md:" {
		dc.Db, err = sql.Open("duckdb", "md:")
		if err != nil {
			log.Fatalf("Could not connect to DuckDB %v\n", err)
		}
		return err
	}
	if _, err := os.Stat(dc.Path); os.IsNotExist(err) {
		log.Fatalf("Path does not exist: %v\n", err)
	}
	dc.Db, err = sql.Open("duckdb", dc.Path)
	if err != nil {
		log.Fatalf("Could not connect to DuckDB %v\n", err)
	}
	return err
}

func (pgc *PgConn) ConnectToDb(ctx context.Context) (err error) {
	_, pgc.Cancel = context.WithTimeout(ctx, 1*time.Minute)
	defer pgc.Cancel()
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", pgc.Host, pgc.Port, pgc.Username, pgc.Password, pgc.Database, pgc.SslMode)
	pgc.Db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Could not connect to Postgres %v\n", err)
	}
	return err
}

func (dbxc *DbxConn) ConnectToDb(ctx context.Context) (err error) {
	_, dbxc.Cancel = context.WithTimeout(ctx, 1*time.Minute)
	defer dbxc.Cancel()
	connStr := fmt.Sprintf("token:%s@%s:443%s?catalog=%s&schema=%s", dbxc.Token, dbxc.Host, dbxc.HttpPath, dbxc.Catalog, dbxc.Schema)
	dbxc.Db, err = sql.Open("databricks", connStr)
	if err != nil {
		log.Fatalf("Could not connect to Databricks %v\n", err)
	}
	return err
}
