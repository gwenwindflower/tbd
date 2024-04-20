package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gwenwindflower/tbd/shared"
)

func SetConnectionDetails(fr FormResponse, ps DbtProfiles) shared.ConnectionDetails {
	var cd shared.ConnectionDetails
	if fr.UseDbtProfile {
		profile, err := GetDbtProfile(fr.DbtProfileName, ps)
		if err != nil {
			log.Fatalf("Could not get dbt profile %v\n", err)
		}
		switch profile.Outputs[fr.DbtProfileOutput].ConnType {
		case "snowflake":
			{
				cd = shared.ConnectionDetails{
					ConnType: profile.Outputs[fr.DbtProfileOutput].ConnType,
					Username: profile.Outputs[fr.DbtProfileOutput].User,
					Account:  profile.Outputs[fr.DbtProfileOutput].Account,
					Database: profile.Outputs[fr.DbtProfileOutput].Database,
					Schema:   fr.Schema,
				}
			}
		case "bigquery":
			{
				cd = shared.ConnectionDetails{
					ConnType: profile.Outputs[fr.DbtProfileOutput].ConnType,
					Project:  profile.Outputs[fr.DbtProfileOutput].Project,
					Dataset:  fr.Schema,
				}
			}
		case "duckdb":
			{
				cd = shared.ConnectionDetails{
					ConnType: profile.Outputs[fr.DbtProfileOutput].ConnType,
					Path:     profile.Outputs[fr.DbtProfileOutput].Path,
					Database: profile.Outputs[fr.DbtProfileOutput].Database,
					Schema:   fr.Schema,
				}
			}
		case "postgres":
			var sslMode string
			if profile.Outputs[fr.DbtProfileOutput].SslMode == "" {
				sslMode = "disable"
			} else {
				sslMode = profile.Outputs[fr.DbtProfileOutput].SslMode
			}
			{
				cd = shared.ConnectionDetails{
					ConnType: profile.Outputs[fr.DbtProfileOutput].ConnType,
					Host:     profile.Outputs[fr.DbtProfileOutput].Host,
					Port:     profile.Outputs[fr.DbtProfileOutput].Port,
					Username: profile.Outputs[fr.DbtProfileOutput].User,
					Password: profile.Outputs[fr.DbtProfileOutput].Password,
					Database: profile.Outputs[fr.DbtProfileOutput].Database,
					SslMode:  sslMode,
					Schema:   fr.Schema,
				}
			}
		default:
			{
				log.Fatalf("Unsupported connection type %v\n", profile.Outputs[fr.DbtProfileOutput].ConnType)
			}
		}
	} else {
		switch fr.Warehouse {
		case "snowflake":
			{
				cd = shared.ConnectionDetails{
					ConnType: fr.Warehouse,
					Username: fr.Username,
					Account:  fr.Account,
					Schema:   fr.Schema,
					Database: fr.Database,
				}
			}
		case "bigquery":
			{
				cd = shared.ConnectionDetails{
					ConnType: fr.Warehouse,
					Project:  fr.Project,
					Dataset:  fr.Dataset,
				}
			}
		case "duckdb":
			wd, err := os.Getwd()
			if err != nil {
				log.Fatalf("Could not get working directory: %v\n", err)
			}
			{
				cd = shared.ConnectionDetails{
					ConnType: fr.Warehouse,
					Path: func() string {
						if fr.Path == "md:" {
							return "md:"
						}
						return filepath.Join(wd, fr.Path)
					}(),
					Database: fr.Database,
					Schema:   fr.Schema,
				}
			}
		case "postgres":
			port, err := strconv.Atoi(fr.Port)
			if err != nil || port < 1000 || port > 9999 {
				log.Fatalf("Port must be a 4-digit number\n")
			}
			{
				cd = shared.ConnectionDetails{
					ConnType: fr.Warehouse,
					Host:     fr.Host,
					Port:     port,
					Username: fr.Username,
					Password: fr.Password,
					Database: fr.Database,
					Schema:   fr.Schema,
					SslMode:  fr.SslMode,
				}
			}
		default:
			{
				log.Fatalf("Unsupported connection type %v\n", fr.Warehouse)
			}
		}
	}
	return cd
}
