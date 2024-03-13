package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"tbd/sourcerer"
	"time"

	"github.com/charmbracelet/huh/spinner"
	_ "github.com/snowflakedb/gosnowflake"
)

type DbtProfile struct {
	Target  string `yaml:"target"`
	Outputs map[string]struct {
		ConnType      string `yaml:"type"`
		Account       string `yaml:"account"`
		User          string `yaml:"user"`
		Role          string `yaml:"role"`
		Authenticator string `yaml:"authenticator"`
		Database      string `yaml:"database"`
		Schema        string `yaml:"schema"`
		Threads       int    `yaml:"threads"`
	} `yaml:"outputs"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	formResponse := Forms()
	if !formResponse.Confirm {
		// TODO: Read up on error types in Go,
		// do a pass on the full codebase to make them correct + consistent
		log.Fatal("⛔ User cancelled.")
	}
	connectionDetails := SetConnectionDetails(formResponse)

	var (
		dbElapsed         float64
		processingElapsed float64
	)
	s := spinner.New()
	s.Action(func() {
		connectionStart := time.Now()
		buildDir := formResponse.BuildDir

		tables, err := sourcerer.GetSources(ctx, connectionDetails)
		if err != nil {
			log.Fatalf("Error getting sources: %v", err)
		}

		dbElapsed = time.Since(connectionStart).Seconds()
		// End of database interaction, start of processing
		processingStart := time.Now()

		if formResponse.GenerateDescriptions {
			GenerateColumnDescriptions(tables)
		}
		PrepBuildDir(buildDir)
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			WriteYAML(tables, buildDir)
		}()
		go func() {
			defer wg.Done()
			WriteStagingModels(tables, buildDir)
		}()
		wg.Wait()
		processingElapsed = time.Since(processingStart).Seconds()
	}).Title("🏎️✨ Generating YAML and SQL files...").Run()
	fmt.Printf("🏁 Done in %.1fs getting data from the db and %.1fs processing! ", dbElapsed, processingElapsed)
	fmt.Println("Your YAML and SQL files are in the build directory.")
}
