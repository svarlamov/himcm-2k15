package main

import (
	"github.com/svarlamov/himcm2k15/models"
	"github.com/svarlamov/himcm2k15/utils"
)

func main() {
	masterSet, err := utils.ParseCSVFile("./resources/master_set.csv")
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
	xyVals := make(models.XYLookupTable)
	xyVals.Populate(rawXYVals)
	rawCoeffs := models.MakeRawCoefficients(coefficientSet)
	dCoeffs := rawCoeffs.ConvertToMap()
	locationSevs := models.MakeLocationSeverities(locationsSevsSet)
	crimes := models.MarshalCrimes(masterSet)
	params := models.ScorerParameters{
		DConst:           70,
		DCoeffs:          dCoeffs,
		LocationSevConst: 20,
		LocationSevs:     locationSevs,
		DomesticConst:    5,
		ArrestedConst:    5,
		XYValues:         xyVals,
	}
	crimes.ScoreCrimes(&params)
	utils.WriteCSVToFile("./resources/output.csv", crimes.MakeCSVStr())
}
