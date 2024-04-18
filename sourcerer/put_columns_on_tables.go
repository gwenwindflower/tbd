package sourcerer

import (
	"context"
	"fmt"
	"regexp"
	"sync"

	"github.com/gwenwindflower/tbd/shared"
	"github.com/schollz/progressbar/v3"
)

func PutColumnsOnTables(ctx context.Context, ts shared.SourceTables, dbc DbConn) error {
	dataTypeGroupMap := map[string]string{
		"(text|char|varchar)":                                "text",
		"(float|int|num|number|bigint|float32|float64|int8)": "numbers",
		"(bool|boolean|bit)":                                 "booleans",
		"(json|struct)":                                      "json",
		"(date|datetime)":                                    "datetimes",
		"(timestamp|timestamptz|timestampntz|timestampltz)":  "timestamps",
	}
	mutex := sync.Mutex{}

	bar := progressbar.NewOptions(len(ts.SourceTables),
		progressbar.OptionSetWidth(5),
		progressbar.OptionShowCount(),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetDescription("[magenta]🏎️✨ Getting warehouse metadata...[reset]"),
	)
	var wg sync.WaitGroup
	wg.Add(len(ts.SourceTables))
	for i := range ts.SourceTables {
		go func(i int) error {
			defer wg.Done()

			columns, err := dbc.GetColumns(ctx, ts.SourceTables[i])
			if err != nil {
				return err
			}

			mutex.Lock()
			ts.SourceTables[i].Columns = columns
			ts.SourceTables[i].DataTypeGroups = make(map[string][]shared.Column)
			// Create a map of data types groups to hold column slices by data type
			// This lets us group columns by their data type e.g. in templates
			for j := range ts.SourceTables[i].Columns {
				for k, v := range dataTypeGroupMap {
					r, _ := regexp.Compile(fmt.Sprintf(`(?i).*%s.*`, k))
					if r.MatchString(ts.SourceTables[i].Columns[j].DataType) {
						ts.SourceTables[i].DataTypeGroups[v] = append(ts.SourceTables[i].DataTypeGroups[v], ts.SourceTables[i].Columns[j])
					}
				}
			}
			bar.Add(1)
			mutex.Unlock()
			return nil
		}(i)
	}
	wg.Wait()
	return nil
}
