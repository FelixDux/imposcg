package impact

import (
	"github.com/google/go-cmp/cmp"
    "github.com/FelixDux/imposcg/dynamics/forcingphase"
	"math"
)

// Each impact is uniquely specified by two parameters:
// `:phi`: the phase (time modulo and scaled by the forcing period) at
// which the impact occurs
// `:v`: the velocity of the impact, which cannot be negative
//
// In addition, we also record the actual time `:t`.
//
// Because `:phi` is periodic and `:v` non-negative, the surface on which
// impacts are defined is a half-cylinder. Whether a zero-velocity impact
// is physically meaningful depends on the value of `:phi` and on `sigma`
// the offset of the obstacle from the centre of motion.
//

type Impact struct {
	Phase float64
	Time float64
	Velocity float64
}

type Generator func(impactTime float64, impactVelocity float64) *Impact

func ImpactGenerator(phaseConverter forcingphase.PhaseConverter) func(impactTime float64, impactVelocity float64) *Impact {
	return func (impactTime float64, impactVelocity float64) *Impact  {
		return &Impact{Phase: phaseConverter.TimeToPhase(impactTime), Time: impactTime, Velocity: impactVelocity}
	}
}

func ImpactComparer(phaseTolerance, velocityTolerance float64) cmp.Option {
	return cmp.Comparer(func(x, y Impact) bool {
		// Compare phases
		if (math.Abs(x.Phase - y.Phase) >= phaseTolerance) {
			// Account for periodicity (i.e. 0 and 1 are the same)
			if math.Abs(x.Phase - y.Phase) < 1 - phaseTolerance {
				return false
			}
		}

		max_v := math.Max(x.Velocity, y.Velocity)

		if max_v==0 {
			max_v = 1
		}

		return (math.Abs((x.Velocity - y.Velocity)/max_v) < velocityTolerance);
	})
}

func (impact Impact) almostEqualAltOpt(other Impact, comparer cmp.Option) bool {
	return cmp.Equal(impact, other, comparer)
}

var defaultImpactComparer = ImpactComparer(1e-3, 1e-3)

func (impact Impact) almostEqual(other Impact) bool {
	return impact.almostEqualAltOpt(other, defaultImpactComparer)
}
