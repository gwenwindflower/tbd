package main

import (
	"log"
	"tbd/shared"
)

func SetConnectionDetails(formResponse FormResponse) shared.ConnectionDetails {
	var connectionDetails shared.ConnectionDetails
	if formResponse.UseDbtProfile {
		profile, err := GetDbtProfile(formResponse.DbtProfile)
		if err != nil {
			log.Fatalf("Could not get dbt profile %v\n", err)
		}
		connectionDetails = shared.ConnectionDetails{
			ConnType: profile.Outputs[formResponse.DbtProfileOutput].ConnType,
			Username: profile.Outputs[formResponse.DbtProfileOutput].User,
			Account:  profile.Outputs[formResponse.DbtProfileOutput].Account,
			Schema:   formResponse.Schema,
			Database: profile.Outputs[formResponse.DbtProfileOutput].Database,
			Project:  profile.Outputs[formResponse.DbtProfileOutput].Project,
			Dataset:  formResponse.Schema,
		}
	} else {
		connectionDetails = shared.ConnectionDetails{
			ConnType: formResponse.Warehouse,
			Username: formResponse.Username,
			Account:  formResponse.Account,
			Schema:   formResponse.Schema,
			Database: formResponse.Database,
			Project:  formResponse.Project,
			Dataset:  formResponse.Dataset,
		}
	}
	return connectionDetails
}
