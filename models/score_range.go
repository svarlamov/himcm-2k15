package models
import "github.com/svarlamov/himcm2k15/utils"

type ScoreRange struct {
	Upper int64
	Lower int64
}

type ScoreRanges []ScoreRange

func (scoreRanges ScoreRanges) Populate(raw [][]string) {
	for ind, row := range raw {
		if ind == 0 {
			continue
		}
		scoreRanges[int(utils.SToIP(row[0]))].Upper = utils.SToIP(row[1])
		scoreRanges[int(utils.SToIP(row[0]))].Lower = utils.SToIP(row[2])
	}
}