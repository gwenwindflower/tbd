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
			log.Fatal(err)
		}
		connectionDetails = shared.ConnectionDetails{
			Warehouse: profile.Outputs[formResponse.DbtProfile].Warehouse,
			Username:  profile.Outputs[formResponse.DbtProfile].User,
			Account:   profile.Outputs[formResponse.DbtProfile].Account,
			Schema:    profile.Outputs[formResponse.DbtProfile].Schema,
			Database:  profile.Outputs[formResponse.DbtProfile].Database,
		}
	} else {
		connectionDetails = shared.ConnectionDetails{
			Warehouse: formResponse.Warehouse,
			Username:  formResponse.Username,
			Account:   formResponse.Account,
			Schema:    formResponse.Schema,
			Database:  formResponse.Database,
		}
	}
	return connectionDetails
}
