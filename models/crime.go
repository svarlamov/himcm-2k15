package models

import (
	"fmt"
	"github.com/svarlamov/himcm2k15/utils"
)

type ScorerParameters struct {
	DConst           int64
	DCoeffs          map[string]int64
	LocationSevConst int64
	LocationSevs     LocationSeverities
	DomesticConst    int64
	ArrestedConst    int64
	XYValues         map[int64][2]int64
	ScoreRanges      ScoreRanges
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
	CaseId        string
	DateOfOcc     string
	PrimaryDesc   string
	SecondaryDesc string
	LocationDesc  string
	Arrest        bool
	Domestic      bool
	Beat          int64
	Score         int64
	XPlotValue    int64
	YPlotValue    int64
	District      int64
	RangeIndex    int64
}

type Crimes []Crime

func (crimes Crimes) ScoreCrimes(params *ScorerParameters) {
	fmt.Println(params.ScoreRanges)
	for ind, crime := range crimes {
		crimes[ind].Score = params.ScoreCrime(crime)
		crimes[ind].XPlotValue = params.XYValues[crimes[ind].Beat][0]
		crimes[ind].YPlotValue = params.XYValues[crimes[ind].Beat][1]
		// NOTE: -1 is a sentinel if a value does not fit within a valid range
		crimes[ind].RangeIndex = -1
		for innerInd, r := range params.ScoreRanges {
			if crimes[ind].Score > r.Lower && crimes[ind].Score < r.Upper {
				crimes[ind].RangeIndex = int64(innerInd)
				break
			}
		}
	}
}

func (crimes Crimes) MakeCSVStr() string {
	rawExp := fmt.Sprintln(`"CASE #","DATE OF OCCURRENCE","PRIMARY DESCRIPTION","SECONDARY DESCRIPTION","LOCATION DESCRIPTION","ARREST","DOMESTIC","BEAT","X PLOT VALUE","Y PLOT VALUE","DISTRICT","RANGE INDEX","SCORE"`)
	for ind, crime := range crimes {
		rawExp += fmt.Sprint(`"`, crime.CaseId, `"`, ",", `"`, crime.DateOfOcc, `"`, ",", `"`, crime.PrimaryDesc, `"`, ",", `"`, crime.SecondaryDesc, `"`, ",", `"`, crime.LocationDesc, `"`, ",", `"`, utils.BoolToYN(crime.Arrest), `"`, ",", `"`, utils.BoolToYN(crime.Domestic), `"`, ",", `"`, crime.Beat, `"`, ",", `"`, crime.XPlotValue, `"`, ",", `"`, crime.YPlotValue, `"`, ",", `"`, crime.District, `"`, ",", `"`, crime.RangeIndex, `"`, ",", `"`, crime.Score, `"`)
		if ind < len(crimes)-1 {
			rawExp += "\n"
		}
	}
	return rawExp
}

func MarshalCrimes(raw [][]string) Crimes {
	crimes := make(Crimes, len(raw)-1)
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
		if tempDidArrest == "Y" || tempDidArrest == "TRUE" || tempDidArrest == "true" {
			crimes[ind].Arrest = true
		} else if tempDidArrest == "N" || tempDidArrest == "FALSE" || tempDidArrest == "false" {
			crimes[ind].Arrest = false
		} else {
			panic(fmt.Sprintln("Could not find arrest value for,", tempDidArrest))
		}
		tempWasDom := row[6]
		if tempWasDom == "Y" || tempWasDom == "TRUE" || tempWasDom == "true" {
			crimes[ind].Domestic = true
		} else if tempWasDom == "N" || tempWasDom == "FALSE" || tempWasDom == "false" {
			crimes[ind].Domestic = false
		} else {
			panic(fmt.Sprintln("Could not find arrest value for,", tempWasDom))
		}
		crimes[ind].Beat = utils.SToIP(row[7])
		crimes[ind].District = utils.SToIP(row[8])
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

func (crimes Crimes) GetSumsPerDistrictPerRange() map[string]int64 {
	sumsPerDistPerRange := make(map[string]int64)
	for _, crime := range crimes {
		if crime.RangeIndex < 0 {
			continue
		}
		sumsPerDistPerRange[fmt.Sprint(crime.District, ",", crime.RangeIndex)]++
	}
	return sumsPerDistPerRange
}

func (crimes Crimes) GetSumsPerRange() []int64 {
	sums := make([]int64, 24)
	for _, crime := range crimes {
		if crime.RangeIndex < 0 {
			continue
		}
		sums[crime.RangeIndex]++
	}
	return sums
}

func MakeSumsPerRangeCSV(data []int64) string {
	rawExp := fmt.Sprintln(`"RANGE","SUM"`)
	for ind, sum := range data {
		rawExp += fmt.Sprint(ind, ",", sum, "\n")
	}
	return rawExp
}

func MakeSumsPerDistrictPerRangeCSV(data map[string]int64) string {
	rawExp := fmt.Sprintln(`"DISTRICT","RANGE","SUM"`)
	for key, sum := range data {
		rawExp += fmt.Sprint(key, ",", sum, "\n")
	}
	return rawExp
}