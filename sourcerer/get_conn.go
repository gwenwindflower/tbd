package sourcerer

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/gwenwindflower/tbd/shared"

	"cloud.google.com/go/bigquery"
)

type DbConn interface {
	ConnectToDb(ctx context.Context) error
	GetSourceTables(ctx context.Context) (shared.SourceTables, error)
	GetColumns(ctx context.Context, t shared.SourceTable) ([]shared.Column, error)
}

type SfConn struct {
	Db       *sql.DB
	Cancel   context.CancelFunc
	Account  string
	Username string
	Database string
	Schema   string
}

type BqConn struct {
	Bq      *bigquery.Client
	Cancel  context.CancelFunc
	Project string
	Dataset string
}

type DuckConn struct {
	Db       *sql.DB
	Cancel   context.CancelFunc
	Path     string
	Database string
	Schema   string
}

type PgConn struct {
	Db       *sql.DB
	Cancel   context.CancelFunc
	Host     string
	Username string
	Password string
	Database string
	Schema   string
	SslMode  string
	Port     int
}

func GetConn(cd shared.ConnectionDetails) (DbConn, error) {
	switch cd.ConnType {
	case "snowflake":
		{
			return &SfConn{
				Account:  strings.ToUpper(cd.Account),
				Username: strings.ToUpper(cd.Username),
				Database: strings.ToUpper(cd.Database),
				Schema:   strings.ToUpper(cd.Schema),
			}, nil
		}
	case "bigquery":
		{
			return &BqConn{
				Project: cd.Project,
				Dataset: cd.Dataset,
			}, nil
		}
	case "duckdb":
		{
			return &DuckConn{
				Path:     cd.Path,
				Database: cd.Database,
				Schema:   cd.Schema,
			}, nil
		}
	case "postgres":
		{
			return &PgConn{
				Host:     cd.Host,
				Port:     cd.Port,
				Username: cd.Username,
				Password: cd.Password,
				Database: cd.Database,
				Schema:   cd.Schema,
				SslMode:  cd.SslMode,
			}, nil
		}
	default:
		return nil, errors.New("unsupported connection type")
	}
}
