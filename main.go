package main

import (
	"github.com/svarlamov/himcm2k15/utils"
	"github.com/svarlamov/himcm2k15/models"
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
	rawCoeffs := models.MakeRawCoefficients(coefficientSet)
	dCoeffs := rawCoeffs.ConvertToMap()
	locationSevs := models.MakeLocationSeverities(locationsSevsSet)
	crimes := models.MarshalCrimes(masterSet)
	params := models.ScorerParameters{
		DConst: 70,
		DCoeffs: dCoeffs,
		LocationSevConst: 20,
		LocationSevs: locationSevs,
		DomesticConst: 5,
		ArrestedConst: 5,
	}
	crimes.ScoreCrimes(&params)
	utils.WriteCSVToFile("./resources/output.csv", crimes.MakeCSVStr())
}