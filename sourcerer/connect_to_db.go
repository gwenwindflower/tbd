package sourcerer

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"tbd/shared"
	"time"
)

type DBConnection interface {
	ConnectToDB() (*sql.DB, error)
}

type SnowflakeConnection struct {
	Username string
	Account  string
	Schema   string
	Database string
}

type BigQueryConnection struct {
	Project string
	Dataset string
}

func (sfc *SnowflakeConnection) ConnectToDB(ctx context.Context, connectionDetails shared.ConnectionDetails) (db *sql.DB, cancel context.CancelFunc, err error) {
	connStr := fmt.Sprintf(
		"%s@%s/%s/%s?authenticator=externalbrowser",
		strings.ToUpper(connectionDetails.Username),
		strings.ToUpper(connectionDetails.Account),
		strings.ToUpper(connectionDetails.Database),
		strings.ToUpper(connectionDetails.Schema),
	)

	_, cancel = context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()
	db, err = sql.Open("snowflake", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db, cancel, err
}
