package models

import (
	"fmt"
	"github.com/svarlamov/himcm2k15/utils"
)

type ScorerParameters struct {
	DConst  int64
	DCoeffs map[string]int64
	LocationSevConst int64
	LocationSevs LocationSeverities
	DomesticConst int64
	ArrestedConst int64
}

func (params *ScorerParameters) ScoreCrime(crime Crime) int64 {
	dCoeff := params.DCoeffs[fmt.Sprintf("%s - %s", crime.PrimaryDesc, crime.SecondaryDesc)]
	if dCoeff == 0 {
		fmt.Println(fmt.Sprintf("%s - %s", crime.PrimaryDesc, crime.SecondaryDesc))
		panic(fmt.Sprintln("Could not find coefficient for,", crime.PrimaryDesc, "and", crime.SecondaryDesc))
	}
	locSev := params.LocationSevs[crime.LocationDesc]
	var arrestedCoeff int64 = 0
	if crime.Arrest == true {
		arrestedCoeff = 1
	}
	var domesticCoeff int64 = 0
	if crime.Domestic != true {
		domesticCoeff = 1
	}
	score := (params.DConst * dCoeff) + (params.LocationSevConst + locSev) + (params.ArrestedConst * arrestedCoeff) + (params.DomesticConst * domesticCoeff)
	return score
}

type Crime struct {
	CaseId string
	DateOfOcc string
	PrimaryDesc string
	SecondaryDesc string
	LocationDesc string
	Arrest bool
	Domestic bool
	Beat int64
	Score int64
}

type Crimes []Crime

func (crimes Crimes) ScoreCrimes(params *ScorerParameters) {
	for ind, crime := range crimes {
		crimes[ind].Score = params.ScoreCrime(crime)
	}
}

func (crimes Crimes) MakeCSVStr() string {
	rawExp := `"CASE #","DATE  OF OCCURRENCE","PRIMARY DESCRIPTION","SECONDARY DESCRIPTION","LOCATION DESCRIPTION","ARREST","DOMESTIC","BEAT","SCORE"`
	rawExp += "\n"
	for ind, crime := range crimes {
		rawExp += fmt.Sprint(`"`, crime.CaseId, `"`, ",", `"`, crime.DateOfOcc, `"`, ",", `"`, crime.PrimaryDesc, `"`, ",", `"`, crime.SecondaryDesc, `"`, ",", `"`, crime.LocationDesc, `"`, ",", `"`, crime.Arrest, `"`, ",", `"`, crime.Domestic, `"`, ",", `"`, crime.Beat, `"`, ",", `"`, crime.Score, `"`)
		if ind < len(crimes) - 1 {
			rawExp += "\n"
		}
	}
	return rawExp
}

func MarshalCrimes(raw [][]string) Crimes {
	crimes := make(Crimes, len(raw) - 1)
	for ind, row := range raw {
		if ind == 0 {
			continue
		}
		ind--
		crimes[ind].CaseId = row[0]
		crimes[ind].DateOfOcc = row[1]
		crimes[ind].PrimaryDesc = row[2]
		crimes[ind].SecondaryDesc = row[3]
		crimes[ind].LocationDesc = row[4]
		tempDidArrest := row[5]
		if tempDidArrest == "Y" {
			crimes[ind].Arrest  = true
		} else if tempDidArrest == "N" {
			crimes[ind].Arrest = false
		} else {
			panic(fmt.Sprintln("Could not find arrest value for,", tempDidArrest))
		}
		tempWasDom := row[6]
		if tempWasDom == "Y" {
			crimes[ind].Domestic  = true
		} else if tempWasDom == "N" {
			crimes[ind].Domestic = false
		} else {
			panic(fmt.Sprintln("Could not find arrest value for,", tempWasDom))
		}
		crimes[ind].Beat = utils.SToIP(row[7])
		temp := fmt.Sprintln(crimes[ind].PrimaryDesc, crimes[ind].SecondaryDesc)
		if len(temp) < 1 {
			fmt.Println(crimes[ind].PrimaryDesc, crimes[ind].SecondaryDesc)
		}
		if len(row[2]) < 1 || len(row[3]) < 1 {
			fmt.Println("Failed with,", row[2], row[3])
			panic("failed parsing crimes")
		}
	}
	return crimes
}