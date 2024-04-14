package sourcerer

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"sync"
	"tbd/shared"
)

func (sfc *SfConn) PutColumnsOnTables(ctx context.Context, tables shared.SourceTables) {
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

			columns, err := sfc.GetColumns(ctx, tables.SourceTables[i])
			if err != nil {
				log.Fatalf("Error fetching columns for table %s: %v\n", tables.SourceTables[i].Name, err)
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

func (bqc *BqConn) PutColumnsOnTables(ctx context.Context, tables shared.SourceTables) {
	mutex := sync.Mutex{}
	var wg sync.WaitGroup
	wg.Add(len(tables.SourceTables))
	dataTypeGroupMap := map[string]string{
		"(string)":    "text",
		"(float|int)": "numbers",
		"(bool)":      "booleans",
		"(json)":      "json",
		"(date)":      "datetimes",
		"(timestamp)": "timestamps",
	}
	for i := range tables.SourceTables {
		go func(i int) {
			defer wg.Done()
			columns, err := bqc.GetColumns(ctx, tables.SourceTables[i])
			if err != nil {
				log.Fatalf("Error fetching columns for table %s: %v\n", tables.SourceTables[i].Name, err)
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