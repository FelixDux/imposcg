package parameters

//  parameters which define a 1-d impact oscillator system with no damping between impacts

import ("fmt")

type Parameters struct {
	ForcingFrequency float64
	CoefficientOfRestitution float64
	ObstacleOffset float64
	MaximumPeriods uint // maximum forcing periods to detect impact
}

type ZeroForcingFrequencyError float64

func (f ZeroForcingFrequencyError) Error() string {
    return "Forcing frequency cannot be zero"
}

type NegativeForcingFrequencyError float64

func (f NegativeForcingFrequencyError) Error() string {
    return fmt.Sprintf("The model cannot handle negative forcing frequency %g", f)
}

type ResonantForcingFrequencyError float64

func (f ResonantForcingFrequencyError) Error() string {
    return "A forcing frequency of 1 is a resonant case with unbounded solutions"
}

type LargeCoefficientOfRestitutionError float64

func (r LargeCoefficientOfRestitutionError) Error() string {
    return fmt.Sprintf("A coefficient of restitution of %g > 1 will generate unbounded solutions", r)
}

type NegativeCoefficientOfRestitutionError float64

func (r NegativeCoefficientOfRestitutionError) Error() string {
    return fmt.Sprintf("A negative coefficient of restitution %g will generate unphysical solutions", r)
}

type ZeroMaxmimumPeriodsError uint

func (n ZeroMaxmimumPeriodsError) Error() string {
	return "maximum forcing periods to detect impact must be > 0"
}

func NewParameters(frequency float64, offset float64, r float64, maxPeriods uint) (*Parameters, []error) {
	errorList := make([]error, 6)

	i := 0

	if frequency == 0.0 {
		errorList[i] = ZeroForcingFrequencyError(frequency)
		i++
	}

	if 0.0 > frequency {
		errorList[i] = NegativeForcingFrequencyError(frequency)
		i++
	}

	if frequency == 1.0 {
		errorList[i] = ResonantForcingFrequencyError(frequency)
		i++
	}

	if 1.0 < r {
		errorList[i] = LargeCoefficientOfRestitutionError(r)
		i++
	}

	if 0.0 > r {
		errorList[i] = NegativeCoefficientOfRestitutionError(r)
		i++
	}

	if maxPeriods == 0 {
		errorList[i] = ZeroMaxmimumPeriodsError(maxPeriods)
		i++
	}

	errorList = errorList[:i]

	if i > 0 {
		return nil, errorList
	} else {
		return &Parameters{ForcingFrequency: frequency, ObstacleOffset: offset, CoefficientOfRestitution: r, MaximumPeriods: maxPeriods}, errorList
	}
}
