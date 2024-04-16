package main

import (
	"errors"
	"sync"

	"github.com/gwenwindflower/tbd/shared"
)

func WriteFiles(ts shared.SourceTables, bd string, prefix string) error {
	if len(ts.SourceTables) == 0 {
		return errors.New("no tables to write")
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		WriteYAML(ts, bd)
	}()
	go func() {
		defer wg.Done()
		WriteStagingModels(ts, bd, prefix)
	}()
	wg.Wait()
	return nil
}
