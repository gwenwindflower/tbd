package main

import (
	"fmt"
	"log"
	"os"
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
)

// Type definitions for the YAML file
type Column struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	DataType    string `yaml:"data_type"`
}

type SourceTable struct {
	DataTypeGroups map[string][]Column `yaml:"-"`
	Name           string              `yaml:"name"`
	Columns        []Column            `yaml:"columns"`
}

type SourceTables struct {
	SourceTables []SourceTable `yaml:"sources"`
}

func main() {
	useForm := true
	if useForm {
		Forms()
	} else {
		warehouse = "snowflake"
		dbUsername = os.Getenv("SNOWFLAKE_SANDBOX_USER")
		dbAccount = os.Getenv("SNOWFLAKE_SANDBOX_ACCOUNT")
		dbDatabase = "ANALYTICS"
		dbSchema = "JAFFLE_SHOP_RAW"
		buildDir = "build"
		generateDescriptions = true
		groqKeyEnvVar = "GROQ_API_KEY"
		confirm = true
	}
	if !confirm {
		log.Fatal("⛔ User cancelled.")
	}
	databaseType := warehouse
	if warehouse == "snowflake" {
		dbAccount = strings.ToUpper(dbAccount)
		dbUsername = strings.ToUpper(dbUsername)
		dbSchema = strings.ToUpper(dbSchema)
		dbDatabase = strings.ToUpper(dbDatabase)
	}

	s := spinner.New()

	connStr := fmt.Sprintf("%s@%s/%s/%s?authenticator=externalbrowser", dbUsername, dbAccount, dbDatabase, dbSchema)
	var (
		connectionElapsed float64
		processingElapsed float64
	)
	s.Action(func() {
		connectionStart := time.Now()
		ctx, db, err := ConnectToDB(connStr, databaseType)
		if err != nil {
			log.Fatal(err)
		}
		tables, err := GetTables(db, ctx)
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
