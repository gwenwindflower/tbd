package sourcerer

import (
	"errors"
	"strings"
	"tbd/shared"
)

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
	default:
		return nil, errors.New("unsupported connection type")
	}
}
