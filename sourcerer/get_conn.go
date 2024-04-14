package sourcerer

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/gwenwindflower/tbd/shared"

	"cloud.google.com/go/bigquery"
)

type DbConn interface {
	ConnectToDB(ctx context.Context) error
	GetSources(ctx context.Context) (shared.SourceTables, error)
	GetColumns(ctx context.Context, t shared.SourceTable) ([]shared.Column, error)
	PutColumnsOnTables(ctx context.Context, tables shared.SourceTables)
}

type SfConn struct {
	Account  string
	Username string
	Database string
	Schema   string
	Db       *sql.DB
	Cancel   context.CancelFunc
}

type BqConn struct {
	Project string
	Dataset string
	Bq      *bigquery.Client
	Cancel  context.CancelFunc
}

type DuckConn struct {
	Path     string
	Database string
	Schema   string
	Db       *sql.DB
	Cancel   context.CancelFunc
}

func GetConn(cd shared.ConnectionDetails) (DbConn, error) {
	switch cd.ConnType {
	case "snowflake":
		// TODO: Why do I need to use a pointer here?
		return &SfConn{
			Account:  strings.ToUpper(cd.Account),
			Username: strings.ToUpper(cd.Username),
			Database: strings.ToUpper(cd.Database),
			Schema:   strings.ToUpper(cd.Schema),
		}, nil
	case "bigquery":
		return &BqConn{
			Project: cd.Project,
			Dataset: cd.Dataset,
		}, nil
	case "duckdb":
		{
			wd, err := os.Getwd()
			if err != nil {
				return nil, err
			}
			p := filepath.Join(wd, cd.Path)
			return &DuckConn{
				Path:     p,
				Database: cd.Database,
				Schema:   cd.Schema,
			}, nil
		}
	default:
		return nil, errors.New("unsupported connection type")
	}
}
