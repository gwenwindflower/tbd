package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gwenwindflower/tbd/sourcerer"
)

type Elapsed struct {
	DbStart           time.Time
	DbElapsed         float64
	ProcessingStart   time.Time
	ProcessingElapsed float64
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	formResponse := Forms()
	if !formResponse.Confirm {
		log.Fatal("⛔ User cancelled.")
	}
	cd := SetConnectionDetails(formResponse)

	e := Elapsed{}
	e.DbStart = time.Now()

	bd := formResponse.BuildDir
	err := PrepBuildDir(bd)
	if err != nil {
		log.Fatalf("Error preparing build directory: %v\n", err)
	}
	dbc, err := sourcerer.GetConn(cd)
	if err != nil {
		log.Fatalf("Error getting database connection: %v\n", err)
	}
	err = dbc.ConnectToDb(ctx)
	if err != nil {
		log.Fatalf("Error connecting to database: %v\n", err)
	}
	ts, err := dbc.GetSourceTables(ctx)
	if err != nil {
		log.Fatalf("Error getting sources: %v\n", err)
	}
	err = sourcerer.PutColumnsOnTables(ctx, ts, dbc)
	if err != nil {
		log.Fatalf("Error putting columns on tables: %v\n", err)
	}

	e.DbElapsed = time.Since(e.DbStart).Seconds()
	// End of database interaction, start of processing
	e.ProcessingStart = time.Now()

	if formResponse.GenerateDescriptions {
		GenerateColumnDescriptions(ts)
	}
	if formResponse.CreateProfile {
		WriteProfile(cd, bd)
	}
	if formResponse.ScaffoldProject {
		s, err := WriteScaffoldProject(cd, bd, formResponse.ProjectName)
		if err != nil {
			log.Fatalf("Error scaffolding project: %v\n", err)
		}
		bd = s
	}
	err = WriteFiles(ts, bd, formResponse.Prefix)
	if err != nil {
		log.Fatalf("Error writing files: %v\n", err)
	}
	e.ProcessingElapsed = time.Since(e.ProcessingStart).Seconds()
	fmt.Printf("\n🏁 Done in %.1fs fetching data and %.1fs writing files!\nYour YAML and SQL files are in the %s directory.", e.DbElapsed, e.ProcessingElapsed, formResponse.BuildDir)
}
