package motion

import (
    "github.com/FelixDux/imposcg/dynamics/impact"
    "github.com/FelixDux/imposcg/dynamics/forcingphase"
    "github.com/FelixDux/imposcg/dynamics/parameters"
    "github.com/FelixDux/imposcg/dynamics/sticking"
	"math"
)
//
// Time evolution of the system from one impact to the next
//

type StateOfMotion struct {
	// 	State and phase variables for the motion between impacts
	Time float64
	Displacement float64
	Velocity float64
}

type MotionAtTime struct {	
	// Coefficients for time evolution of the system from one impact to the next 
	Parameters parameters.Parameters
	ImpactTime float64
	CosCoefficient float64
	SinCoefficient float64
	LongExcursion func(float64) bool
}

func LongExcursionChecker(MaximumPeriods uint, converter *forcingphase.PhaseConverter, impactTime float64) func(float64) bool {
	return func(time float64) bool {return time-impactTime > float64(MaximumPeriods) * converter.Period}
}

func NewMotionAtTime(parameters parameters.Parameters, converter *forcingphase.PhaseConverter, impact impact.Impact) *MotionAtTime {
	cosCoefficient := parameters.ObstacleOffset - parameters.Gamma * math.Cos(parameters.ForcingFrequency*impact.Time);
	
	sinCoefficient := -(parameters.CoefficientOfRestitution * impact.Velocity) + parameters.ForcingFrequency * parameters.Gamma * math.Sin(parameters.ForcingFrequency*impact.Time)

	return &MotionAtTime{
		Parameters: parameters, 
		ImpactTime: impact.Time, 
		CosCoefficient: cosCoefficient, 
		SinCoefficient: sinCoefficient, 
		LongExcursion: LongExcursionChecker(parameters.MaximumPeriods, converter, impact.Time)}
}

type SearchParameters struct {
	InitialStepSize float64
	MinimumStepSize float64
}

func DefaultSearchParameters() SearchParameters {
	return SearchParameters{InitialStepSize: 0.1, MinimumStepSize: 0.000001}
}

type MotionBetweenImpacts struct {
//
// Generates a trajectory from one impact to the next
//

	motion MotionAtTime
	sticking sticking.Sticking
	search SearchParameters
	offset float64

}

func MotionBetweenImpactsInitialiser(parameters parameters.Parameters) (func(impact.Impact) *MotionBetweenImpacts, error) {
	sticking, err := sticking.NewSticking(parameters)

	if err != nil {
		return nil, err
	} else {
		converter := &sticking.Converter
		return func(impact impact.Impact) *MotionBetweenImpacts {
			return &MotionBetweenImpacts{motion: *NewMotionAtTime(parameters, converter, impact), sticking: *sticking, search: DefaultSearchParameters(), offset: parameters.ObstacleOffset}}, nil
	}
}

func (motion MotionAtTime) State(time float64) StateOfMotion {
	lambda := time - motion.ImpactTime

	cosLambda := math.Cos(lambda)
	sinLambda := math.Sin(lambda)

	return StateOfMotion {Time: time,
		Displacement: motion.CosCoefficient * cosLambda + motion.SinCoefficient * sinLambda + 
			motion.Parameters.Gamma * math.Cos( time * motion.Parameters.ForcingFrequency),
		Velocity: motion.SinCoefficient * cosLambda - motion.CosCoefficient * sinLambda - 
			motion.Parameters.ForcingFrequency * motion.Parameters.Gamma * math.Sin(time * motion.Parameters.ForcingFrequency) }
}


type NextImpactResult struct {
	Motion []StateOfMotion

	FoundImpact bool
}

func (result NextImpactResult) Grow(state StateOfMotion) *NextImpactResult {
	result.Motion = append(result.Motion, state)

	return &result
}

func (motion MotionBetweenImpacts) NewNextImpactResult(cap uint, impact impact.Impact) *NextImpactResult {

	trajectory := make([]StateOfMotion, 0, cap)
	
	trajectory = append(trajectory, StateOfMotion {Time: impact.Time, Displacement: motion.offset, Velocity: impact.Velocity})
	
	releaseImpact := motion.sticking.CheckImpact(impact)
	
	if (releaseImpact.NewImpact) {
		trajectory = append(trajectory, StateOfMotion{
			Time: releaseImpact.Impact.Time, 
			Displacement: motion.offset, 
			Velocity: releaseImpact.Impact.Velocity})
	}

	return &NextImpactResult{Motion: trajectory, FoundImpact: false}
}

func (motion MotionBetweenImpacts) DefaultNextImpactResult(impact impact.Impact) *NextImpactResult {
	return motion.NewNextImpactResult(300, impact)
}

func (motion MotionBetweenImpacts) NextImpact(impact impact.Impact) *NextImpactResult {

	result := motion.DefaultNextImpactResult(impact)

	result.FoundImpact = true

	stepSize := motion.search.InitialStepSize

	currentTime := result.Motion[len(result.Motion) - 1].Time

	for math.Abs(stepSize) > motion.search.MinimumStepSize && result.FoundImpact {
		currentTime += stepSize

		currentState := motion.motion.State(currentTime)

		// Update step size - this is the bisection search algorithm
		if currentState.Displacement < motion.offset {
			// only record the state if it is physical
			// (i.e. non-penetrating)
			result = result.Grow(currentState)

			if stepSize < 0.0 {
				stepSize *= -0.5
			}
		} else if currentState.Displacement > motion.offset {
			if (stepSize > 0) {
				stepSize *= -0.5;
			}
		} else {
			result = result.Grow(currentState)
			stepSize = 0;
		}

		if (motion.motion.LongExcursion(currentTime)) {
			result.FoundImpact = false;
		}
	}
	
	return result
}

