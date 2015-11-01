package main

import (
	"github.com/svarlamov/himcm2k15/models"
	"github.com/svarlamov/himcm2k15/utils"
)

func main() {
	masterSet, err := utils.ParseCSVFile("./resources/master_set_with_districts.csv")
	if err != nil {
		panic(err)
	}
	coefficientSet, err := utils.ParseCSVFile("./resources/coefficients.csv")
	if err != nil {
		panic(err)
	}
	locationsSevsSet, err := utils.ParseCSVFile("./resources/location_severities.csv")
	if err != nil {
		panic(err)
	}
	rawXYVals, err := utils.ParseCSVFile("./resources/x_y_lookup.csv")
	if err != nil {
		panic(err)
	}
	rawScoreRanges, err := utils.ParseCSVFile("./resources/score_ranges.csv")
	if err != nil {
		panic(err)
	}
	xyVals := make(models.XYLookupTable)
	xyVals.Populate(rawXYVals)
	rawCoeffs := models.MakeRawCoefficients(coefficientSet)
	dCoeffs := rawCoeffs.ConvertToMap()
	locationSevs := models.MakeLocationSeverities(locationsSevsSet)
	scoreRanges := make(models.ScoreRanges, len(rawScoreRanges)-1)
	scoreRanges.Populate(rawScoreRanges)
	crimes := models.MarshalCrimes(masterSet)
	params := models.ScorerParameters{
		DConst:           70,
		DCoeffs:          dCoeffs,
		LocationSevConst: 20,
		LocationSevs:     locationSevs,
		DomesticConst:    5,
		ArrestedConst:    5,
		XYValues:         xyVals,
		ScoreRanges:      scoreRanges,
	}
	crimes.ScoreCrimes(&params)
	utils.WriteCSVToFile("./resources/output.csv", crimes.MakeCSVStr())
	perDistrictPerRange := crimes.GetSumsPerDistrictPerRange()
	utils.WriteCSVToFile("./resources/crimes_per_dist_per_range.csv", models.MakeSumsPerDistrictPerRangeCSV(perDistrictPerRange))
	sumsPerRange := crimes.GetSumsPerRange()
	utils.WriteCSVToFile("./resources/sums_per_range.csv", models.MakeSumsPerRangeCSV(sumsPerRange))
}
