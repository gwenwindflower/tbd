package main

import (
	"context"
	"database/sql"
	"log"
	"sync"
)

func PutColumnsOnTables(db *sql.DB, ctx context.Context, tables SourceTables) {
	mutex := sync.Mutex{}

	var wg sync.WaitGroup
	wg.Add(len(tables.SourceTables))

	for i := range tables.SourceTables {
		go func(i int) {
			defer wg.Done()

			columns, err := GetColumns(db, ctx, tables.SourceTables[i])
			if err != nil {
				log.Printf("Error fetching columns for table %s: %v\n", tables.SourceTables[i].Name, err)
				return
			}

			mutex.Lock()
			tables.SourceTables[i].Columns = columns
			mutex.Unlock()
		}(i)
	}
	wg.Wait()
}
