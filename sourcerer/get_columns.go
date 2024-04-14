package sourcerer

import (
	"context"
	"fmt"
	"log"

	"github.com/gwenwindflower/tbd/shared"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

func (sfc *SfConn) GetColumns(ctx context.Context, t shared.SourceTable) ([]shared.Column, error) {
	var cs []shared.Column

	q := fmt.Sprintf("SELECT column_name, data_type FROM information_schema.columns WHERE table_schema = '%s' AND table_name = '%s'", sfc.Schema, t.Name)
	rows, err := sfc.Db.QueryContext(ctx, q)
	if err != nil {
		log.Fatalf("Error fetching columns for table %s: %v\n", t.Name, err)
	}
	defer rows.Close()

	for rows.Next() {
		c := shared.Column{}
		if err := rows.Scan(&c.Name, &c.DataType); err != nil {
			log.Fatalf("Error scanning columns for table %s: %v\n", t.Name, err)
		}
		cs = append(cs, c)
	}
	return cs, nil
}

func (bqc *BqConn) GetColumns(ctx context.Context, t shared.SourceTable) ([]shared.Column, error) {
	var cs []shared.Column
	qs := fmt.Sprintf("SELECT column_name, data_type FROM %s.%s.INFORMATION_SCHEMA.COLUMNS WHERE table_name = @table", bqc.Project, bqc.Dataset)
	q := bqc.Bq.Query(qs)
	q.Parameters = []bigquery.QueryParameter{
		{Name: "table", Value: t.Name},
	}
	it, err := q.Read(ctx)
	if err != nil {
		log.Fatalf("Error fetching columns for table %s: %v\n", t.Name, err)
	}
	for {
		var values []bigquery.Value
		err := it.Next(&values)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Error scanning columns for table %s: %v\n", t.Name, err)
		}
		c := shared.Column{
			Name:     values[0].(string),
			DataType: values[1].(string),
		}
		cs = append(cs, c)
	}
	return cs, nil
}

func (dc *DuckConn) GetColumns(ctx context.Context, t shared.SourceTable) ([]shared.Column, error) {
	var cs []shared.Column
	q := fmt.Sprintf("SELECT column_name, data_type FROM information_schema.columns WHERE table_schema = '%s' AND table_name = '%s'", dc.Schema, t.Name)
	rows, err := dc.Db.QueryContext(ctx, q)
	if err != nil {
		log.Fatalf("Error fetching columns for table %s: %v\n", t.Name, err)
	}
	defer rows.Close()
	for rows.Next() {
		c := shared.Column{}
		if err := rows.Scan(&c.Name, &c.DataType); err != nil {
			log.Fatalf("Error scanning columns for table %s: %v\n", t.Name, err)
		}
		cs = append(cs, c)
	}
	return cs, nil
}
