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
	ProcessingStart   time.Time
	DbElapsed         float64
	ProcessingElapsed float64
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ps, err := FetchDbtProfiles()
	if err != nil {
		log.Fatalf("Error fetching dbt profiles: %v\n", err)
	}
	fr, err := Forms(ps)
	if err != nil {
		log.Fatalf("Error running form: %v\n", err)
	}
	if !fr.Confirm {
		log.Fatal("⛔ User cancelled.")
	}
	cd := SetConnectionDetails(fr, ps)

	e := Elapsed{}
	e.DbStart = time.Now()

	bd := fr.BuildDir
	err = PrepBuildDir(bd)
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

	if fr.GenerateDescriptions {
		GenerateColumnDescriptions(ts)
	}
	if fr.CreateProfile {
		WriteProfile(cd, bd)
	}
	if fr.ScaffoldProject {
		s, err := WriteScaffoldProject(cd, bd, fr.ProjectName)
		if err != nil {
			log.Fatalf("Error scaffolding project: %v\n", err)
		}
		bd = s
	}
	err = WriteFiles(ts, bd, fr.Prefix)
	if err != nil {
		log.Fatalf("Error writing files: %v\n", err)
	}
	e.ProcessingElapsed = time.Since(e.ProcessingStart).Seconds()
	fmt.Printf("\n🏁 Done in %.1fs fetching data and %.1fs writing files!\nYour YAML and SQL files are in the %s directory.", e.DbElapsed, e.ProcessingElapsed, fr.BuildDir)
}
