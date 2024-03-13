package sourcerer

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"sync"
	"tbd/shared"
)

func PutColumnsOnTables(ctx context.Context, db *sql.DB, tables shared.SourceTables, connectionDetails shared.ConnectionDetails) {
	mutex := sync.Mutex{}

	var wg sync.WaitGroup
	wg.Add(len(tables.SourceTables))

	dataTypeGroupMap := map[string]string{
		"(text|char)":     "text",
		"(float|int|num)": "numbers",
		"(bool|bit)":      "booleans",
		"json":            "json",
		"date":            "datetimes",
		"timestamp":       "timestamps",
	}

	for i := range tables.SourceTables {
		go func(i int) {
			defer wg.Done()

			columns, err := GetColumns(db, ctx, tables.SourceTables[i], connectionDetails)
			if err != nil {
				log.Printf("Error fetching columns for table %s: %v\n", tables.SourceTables[i].Name, err)
				return
			}

			mutex.Lock()
			tables.SourceTables[i].Columns = columns
			tables.SourceTables[i].DataTypeGroups = make(map[string][]shared.Column)
			// Create a map of data types groups to hold column slices by data type
			// This lets us group columns by their data type e.g. in templates
			for j := range tables.SourceTables[i].Columns {
				for k, v := range dataTypeGroupMap {
					r, _ := regexp.Compile(fmt.Sprintf(`(?i).*%s.*`, k))
					if r.MatchString(tables.SourceTables[i].Columns[j].DataType) {
						tables.SourceTables[i].DataTypeGroups[v] = append(tables.SourceTables[i].DataTypeGroups[v], tables.SourceTables[i].Columns[j])
					}
				}
			}
			mutex.Unlock()
		}(i)
	}
	wg.Wait()
}
