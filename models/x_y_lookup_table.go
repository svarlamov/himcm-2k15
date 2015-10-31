package models

import (
	"github.com/svarlamov/himcm2k15/utils"
	"fmt"
)

type XYLookupTable map[int64][2]int64

func (table XYLookupTable) Populate(raw [][]string) {
	for ind, row := range raw {
		if ind == 0 {
			continue
		}
		table[utils.SToIP(row[0])] = [2]int64{utils.SToIP(row[1]), utils.SToIP(row[2])}
	}
	fmt.Println(len(table))
}
