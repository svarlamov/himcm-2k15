package models

import (
	"fmt"
	"github.com/svarlamov/himcm2k15/utils"
)

type RawCoefficients struct {
	PrimaryDesc   []string
	SecondaryDesc []string
	Coefficient   []int64
}

func MakeRawCoefficients(rows [][]string) RawCoefficients {
	raw := RawCoefficients{
		PrimaryDesc:   make([]string, len(rows)),
		SecondaryDesc: make([]string, len(rows)),
		Coefficient:   make([]int64, len(rows)),
	}
	for i := 1; i < len(rows); i++ {
		raw.PrimaryDesc[i] = rows[i][0]
		raw.SecondaryDesc[i] = rows[i][1]
		raw.Coefficient[i] = utils.SToIP(rows[i][2])
	}
	return raw
}

func (rawData *RawCoefficients) ConvertToMap() CoefficientsMap {
	coeffsMap := make(CoefficientsMap)
	for i := 0; i < len(rawData.PrimaryDesc); i++ {
		coeffsMap[fmt.Sprintf("%s - %s", rawData.PrimaryDesc[i], rawData.SecondaryDesc[i])] = rawData.Coefficient[i]
	}
	return coeffsMap
}

type CoefficientsMap map[string]int64
