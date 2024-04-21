package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/fatih/color"
	"github.com/gwenwindflower/tbd/shared"
	"github.com/schollz/progressbar/v3"
)

func InferColumnFields(llm Llm, ts shared.SourceTables) {
	var wg sync.WaitGroup
	semaphore, limiter := llm.GetRateLimiter()
	defer limiter.Stop()

	bar := progressbar.NewOptions(countColumns(ts),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(30),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionOnCompletion(func() {
			color.HiGreen("\nColumn config generated.")
		}),
		progressbar.OptionSetDescription("🤖📝"),
	)
	for i := range ts.SourceTables {
		for j := range ts.SourceTables[i].Columns {

			semaphore <- struct{}{}
			<-limiter.C

			wg.Add(1)
			go func(i, j int) {
				defer wg.Done()
				defer func() { <-semaphore }()

				tableName := ts.SourceTables[i].Name
				columnName := ts.SourceTables[i].Columns[j].Name
				descPrompt := fmt.Sprintf(DESC_PROMPT, tableName, columnName)
				testsPrompt := fmt.Sprintf(TESTS_PROMPT, tableName, columnName)
				err := llm.SetDescription(descPrompt, ts, i, j)
				if err != nil {
					log.Fatalf("Error generating descriptions: %v\n", err)
				}
				err = llm.SetTests(testsPrompt, ts, i, j)
				if err != nil {
					log.Fatalf("Error generating tests: %v\n", err)
				}
				bar.Add(1)
			}(i, j)
		}
	}
	wg.Wait()
}

func countColumns(ts shared.SourceTables) int {
	c := 0
	for _, t := range ts.SourceTables {
		c += len(t.Columns)
	}
	return c
}
