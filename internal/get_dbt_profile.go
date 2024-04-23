package internal

import (
	"fmt"
)

func GetDbtProfile(pn string, ps DbtProfiles) (DbtProfile, error) {
	if p, ok := ps[pn]; ok {
		return p, nil
	} else {
		return DbtProfile{}, fmt.Errorf("no profile named %s", pn)
	}
}
