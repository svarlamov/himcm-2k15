package main

import (
	"github.com/svarlamov/himcm2k15/utils"
	"github.com/svarlamov/himcm2k15/models"
	"fmt"
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
	rawCoeffs := models.MakeRawCoefficients(coefficientSet)
	dCoeffs := rawCoeffs.ConvertToMap()
	crimes := models.MarshalCrimes(masterSet)
	params := models.ScorerParameters{
		DConst: 1,
		DCoeffs: dCoeffs,
	}
	crimes.ScoreCrimes(&params)
	fmt.Println(crimes.MakeCSVStr())
}