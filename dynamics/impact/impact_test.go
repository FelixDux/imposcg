package impact

import (
	"testing"

	"github.com/FelixDux/imposcg/dynamics/forcingphase"
)

func TestImpactEqualityDefault(t *testing.T) {
	converter, _ := forcingphase.NewPhaseConverter(2.0)

	generator := ImpactGenerator(*converter)

	checkEqual := func (x, y Impact, expected bool) {
		shoulder := func() string {
			if expected {
				return "should"
			} else {
				return "should not"
			}
		}

		if x.almostEqual(y) != expected {
			t.Errorf("Default comparer %s treat impacts %+v and %+v as equal", shoulder(), x, y)
		}
	}

	impact1 := *generator(0, 1)
	impact2 := *generator(0.0001, 1.0001)
	impact3 := *generator(0.3, 0.2)
	impact4 := *generator(1.0000001*converter.Period, 1)
	impact5 := *generator(0.99999*converter.Period, 1)

	checkEqual(impact1, impact2, true)
	checkEqual(impact3, impact2, false)
	checkEqual(impact1, impact4, true)
	checkEqual(impact1, impact5, true)
	
}