// Look here: https://github.com/golang-standards/project-layout

package main

import ("fmt"

"github.com/FelixDux/imposcg/dynamics/forcingphase")

func main() {
	conv, _ := forcingphase.NewPhaseConverter(3)
	fmt.Println(conv.TimeToPhase(2))
}