package main

import (
	"fmt"
	"sync"

	"github.com/gwenwindflower/tbd/shared"
	"github.com/schollz/progressbar/v3"
)

func InferColumnFields(llm Llm, ts shared.SourceTables) error {
	var wg sync.WaitGroup
	semaphore, limiter := llm.GetRateLimiter()
	defer limiter.Stop()

	bar := progressbar.NewOptions(len(ts.SourceTables),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(30),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetDescription("🤖📝"),
	)
	for i := range ts.SourceTables {
		for j := range ts.SourceTables[i].Columns {

			semaphore <- struct{}{}
			<-limiter.C

			wg.Add(1)
			go func(i, j int) error {
				defer wg.Done()
				defer func() { <-semaphore }()

				tableName := ts.SourceTables[i].Name
				columnName := ts.SourceTables[i].Columns[j].Name
				descPrompt := fmt.Sprintf(DESC_PROMPT, tableName, columnName)
				testsPrompt := fmt.Sprintf(TESTS_PROMPT, tableName, columnName)
				err := llm.SetDescription(descPrompt, ts, i, j)
				if err != nil {
					return fmt.Errorf("error setting description: %v", err)
				}
				err = llm.SetTests(testsPrompt, ts, i, j)
				if err != nil {
					return fmt.Errorf("error setting tests: %v", err)
				}
				return nil
			}(i, j)
		}
		bar.Add(1)
	}
	wg.Wait()
	return nil
}
