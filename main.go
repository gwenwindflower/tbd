package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/huh/spinner"
	_ "github.com/snowflakedb/gosnowflake"
)

// Connection details for Snowflake
var (
	confirm              bool
	warehouse            string
	dbUsername           string
	dbAccount            string
	dbDatabase           string
	dbSchema             string
	buildDir             string
	generateDescriptions bool
	groqKeyEnvVar        string
	useDbtProfile        bool
	dbtProfile           string
)

// Type definitions for the YAML file
type Column struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	DataType    string   `yaml:"data_type"`
	Tests       []string `yaml:"tests"`
}

type SourceTable struct {
	DataTypeGroups map[string][]Column `yaml:"-"`
	Name           string              `yaml:"name"`
	Columns        []Column            `yaml:"columns"`
}

type SourceTables struct {
	SourceTables []SourceTable `yaml:"sources"`
}

type DbtProfile struct {
	Target  string `yaml:"target"`
	Outputs map[string]struct {
		Warehouse     string `yaml:"type"`
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
	Forms()
	if !confirm {
		log.Fatal("⛔ User cancelled.")
	}
	if useDbtProfile {
		profile, err := GetDbtProfile(dbtProfile)
		if err != nil {
			log.Fatal(err)
		}
		warehouse = profile.Outputs["dev"].Warehouse
		dbUsername = strings.ToUpper(profile.Outputs["dev"].User)
		dbAccount = strings.ToUpper(profile.Outputs["dev"].Account)
		dbDatabase = strings.ToUpper(profile.Outputs["dev"].Database)
		dbSchema = strings.ToUpper(profile.Outputs["dev"].Schema)
	} else {
		dbAccount = strings.ToUpper(dbAccount)
		dbUsername = strings.ToUpper(dbUsername)
		dbSchema = strings.ToUpper(dbSchema)
		dbDatabase = strings.ToUpper(dbDatabase)
	}

	connStr := fmt.Sprintf("%s@%s/%s/%s?authenticator=externalbrowser", dbUsername, dbAccount, dbDatabase, dbSchema)

	var (
		connectionElapsed float64
		processingElapsed float64
	)
	s := spinner.New()
	s.Action(func() {
		connectionStart := time.Now()
		ctx, db, err := ConnectToDB(connStr, warehouse)
		if err != nil {
			log.Fatal(err)
		}
		tables, err := GetTables(db, ctx, dbSchema)
		if err != nil {
			log.Fatal(err)
		}
		connectionElapsed = time.Since(connectionStart).Seconds()
		processingStart := time.Now()
		PutColumnsOnTables(db, ctx, tables)
		if generateDescriptions {
			GenerateColumnDescriptions(tables)
		}
		CleanBuildDir(buildDir)
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
	fmt.Printf("🏁 Done in %.1fs getting data from the db and %.1fs processing! ", connectionElapsed, processingElapsed)
	fmt.Println("Your YAML and SQL files are in the build directory.")
}
