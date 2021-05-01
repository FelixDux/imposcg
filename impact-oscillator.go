package main

import (
	"fmt"

	"github.com/FelixDux/imposcg/dynamics"
	"github.com/FelixDux/imposcg/dynamics/parameters"
	"github.com/FelixDux/imposcg/charts"
)

func main() {
	fmt.Println("Impact Oscillator")

	params, errParams := parameters.NewParameters(2.8, 0.0, 0.8, 100)

	if len(errParams) > 0 {
		for _, err := range(errParams) {

			fmt.Println(err.Error())

		}

		return
	}
	
	impactMap, errMap := dynamics.NewImpactMap(*params)

	if errMap != nil {
		fmt.Println(errMap.Error())

		return
	}

	phi := 0.0
	v := 0.0
	data := impactMap.IterateFromPoint(phi, v, 1000)

	fmt.Println(charts.ImpactMapPlot(*params, data.Impacts, phi, v).Name())

}