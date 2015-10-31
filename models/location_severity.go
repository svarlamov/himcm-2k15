package models
import "github.com/svarlamov/himcm2k15/utils"

type LocationSeverities map[string]int64

func MakeLocationSeverities(raw [][]string) LocationSeverities {
	severities := make(LocationSeverities)
	for ind, row := range raw {
		if ind == 0 {
			continue
		}
		severities[row[0]] = utils.SToIP(row[1])
	}
	return severities
}