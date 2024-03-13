package main

import (
	"context"
	"fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestGetTables(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctx := context.TODO()

	rows := sqlmock.NewRows([]string{"table_name"}).
		AddRow("raw_orders").
		AddRow("raw_order_items").
		AddRow("raw_locations")

	schema := "test_schema"
	mock.ExpectQuery(fmt.Sprintf("SELECT table_name FROM information_schema.tables where table_schema = '%s'", schema)).WillReturnRows(rows)

	result, err := GetTables(db, ctx, schema)
	if err != nil {
		t.Fatalf("error was not expected while retrieving tables: %s", err)
	}

	if len(result.SourceTables) != 3 || result.SourceTables[0].Name != "raw_orders" || result.SourceTables[1].Name != "raw_order_items" {
		t.Errorf("result not match, got %+v", result)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetTablesNoTablesInSchema(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	ctx := context.TODO()

	schema := "test_schema"
	mock.NewColumn("table_name")
	mock.ExpectQuery(fmt.Sprintf("SELECT table_name FROM information_schema.tables where table_schema = '%s'", schema)).WillReturnError(err)
	_, err = GetTables(db, ctx, schema)
	if err == nil {
		t.Fatalf("error was expected while retrieving tables")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
