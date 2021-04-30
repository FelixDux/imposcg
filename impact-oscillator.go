// Look here: https://github.com/golang-standards/project-layout

package main

import ("fmt"

"github.com/FelixDux/imposcg/dynamics/parameters"
"github.com/FelixDux/imposcg/dynamics")

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

	data := impactMap.IterateFromPoint(0.0, 0.0, 100)

	for _, point := range(data.Impacts) {
		fmt.Printf("%g\t%g\n", point.Phase, point.Velocity)
	}
}