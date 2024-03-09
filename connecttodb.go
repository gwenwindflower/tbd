package main

import (
	"context"
	"database/sql"
	"log"
)

func ConnectToDB(connStr string, databaseType string) (ctx context.Context, db *sql.DB, err error) {
	ctx = context.Background()
	db, err = sql.Open(databaseType, connStr)
	if err != nil {
		log.Fatal(err)
	}
	return ctx, db, err
}
