// Look here: https://github.com/golang-standards/project-layout

package main

import ("fmt"

"github.com/FelixDux/imposcg/dynamics/forcingphase")

func main() {
	fmt.Println(forcingphase.ConvertTimeToPhase(3)(2))
}