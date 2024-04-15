package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/huh/spinner"
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
	s := spinner.New()
	err := s.Action(func() {
		e.DbStart = time.Now()

		bd := formResponse.BuildDir

		dbc, err := sourcerer.GetConn(cd)
		if err != nil {
			log.Fatalf("Error getting connection: %v\n", err)
		}
		ts, err := dbc.GetSources(ctx)
		if err != nil {
			log.Fatalf("Error getting sources: %v\n", err)
		}

		e.DbElapsed = time.Since(e.DbStart).Seconds()
		// End of database interaction, start of processing
		e.ProcessingStart = time.Now()

		if formResponse.GenerateDescriptions {
			GenerateColumnDescriptions(ts)
		}
		PrepBuildDir(bd)
		err = WriteFiles(ts, bd)
		if err != nil {
			log.Fatalf("Error writing files: %v\n", err)
		}
	}).Title("🏎️✨ Generating YAML and SQL files...").Run()
	if err != nil {
		log.Fatalf("Error running spinner action: %v\n", err)
	}
	e.ProcessingElapsed = time.Since(e.ProcessingStart).Seconds()
	fmt.Printf("🏁 Done in %.1fs fetching data and %.1fs writing files! ", e.DbElapsed, e.ProcessingElapsed)
	fmt.Println("Your YAML and SQL files are in the build directory.")
}
