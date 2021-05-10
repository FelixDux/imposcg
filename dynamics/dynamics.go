package dynamics
import (
    "github.com/FelixDux/imposcg/dynamics/impact"
    "github.com/FelixDux/imposcg/dynamics/parameters"
    "github.com/FelixDux/imposcg/dynamics/sticking"
    "github.com/FelixDux/imposcg/dynamics/motion"
	"math"
)

type ChatterResult struct {
	isChatter bool
	accumulationImpact impact.Impact
}

type ChatterChecker struct {
	
	// Detects and numerically approximates 'Chatter', which is when an infinite sequence of impact.Impacts accumulates 
	// in a finite time on a 'sticking' impact.Impact. It is the analogue in this system to a real-world situation in 
	// which the mass judders against the stop. To handle it numerically it is necessary to detect when it is 
	// happening and then extrapolate forward to the accumulation point.
		
		velocityThreshold float64
		countThreshold uint
		sticking sticking.Sticking
		accumulationTime func (impact.Impact) float64 
		canChatter bool

		impactCount uint
}

func NewChatterChecker(parameters parameters.Parameters, velocityThreshold float64, countThreshold uint) (*ChatterChecker, error) {
	sticking, err := sticking.NewSticking(parameters)

	if err != nil {
		return nil, err
	}

	canChatter := true
	accumulationTime := func (impact impact.Impact) float64 {return impact.Time}

	if parameters.CoefficientOfRestitution < 1.0 && parameters.CoefficientOfRestitution >=0.0 { 
		accumulationTime = func (impact impact.Impact) float64 {
			return impact.Time - 2*impact.Velocity / (1-parameters.CoefficientOfRestitution) /
				(math.Cos(impact.Time * parameters.ForcingFrequency) - parameters.ObstacleOffset)
		}
	} else {
		canChatter = false
	}

	return &ChatterChecker{
				velocityThreshold: velocityThreshold,
				countThreshold: countThreshold,
				impactCount: 0,
				canChatter: canChatter,
				accumulationTime: accumulationTime,
				sticking: *sticking}, nil
}

func (checker ChatterChecker) Check(impact impact.Impact) *ChatterResult {
	if checker.canChatter && impact.Velocity < checker.velocityThreshold {
		checker.impactCount++
		if (checker.impactCount > checker.countThreshold) {
			checker.impactCount = 0
			newTime := checker.accumulationTime(impact)

			if (checker.sticking.TimeSticks(newTime)) {
				return &ChatterResult{isChatter: true, accumulationImpact: *checker.sticking.Generator(newTime, 0)}
			}
		}
	}
	
	return &ChatterResult{isChatter: false, accumulationImpact: impact}
}

func DefaultChatterChecker(parameters parameters.Parameters) (*ChatterChecker, error) {
	return NewChatterChecker(parameters, 0.05, 10)
}

type IterationResult struct 
{
	Impacts []impact.Impact

	longExcursions bool
}

type ImpactResult struct 
{
	impact impact.Impact

	FoundImpact bool
}


type ImpactMap struct {
	
	// Transformation of the impact surface (an infinite half cylinder parametrised by phase and velocity)
	// which maps impacts to impacts
		
	motion motion.MotionBetweenImpacts
	chatterChecker ChatterChecker
	dualMaker func(impact.Impact) impact.Impact
}

func NewImpactMap(parameters parameters.Parameters) (*ImpactMap, error) {
	motion, errMotion := motion.NewMotionBetweenImpacts(parameters)

	if errMotion != nil {
		return nil, errMotion
	}

	chatterChecker, errChatter := DefaultChatterChecker(parameters)

	if errChatter != nil {
		return nil, errChatter
	}

	return &ImpactMap{motion: *motion, chatterChecker: *chatterChecker, dualMaker: impact.DualMaker(parameters.CoefficientOfRestitution)}, nil
}

func (impactMap ImpactMap) GenerateImpact(Time float64, Velocity float64) *impact.Impact {
	return impactMap.chatterChecker.sticking.Generator(Time, Velocity)
}

// Apply the map to an impact
 func (impactMap ImpactMap) apply(impact impact.Impact) *ImpactResult {
	trajectory := *impactMap.motion.NextImpact(impact)

	stateAtImpact := trajectory.Motion[len(trajectory.Motion)-1]

	return &ImpactResult{impact: *impactMap.GenerateImpact(stateAtImpact.Time, stateAtImpact.Velocity), FoundImpact: trajectory.FoundImpact}
}

// Iterate the map 
func (impactMap ImpactMap) iterate(initialImpact impact.Impact, numIterations uint) *IterationResult {

	longExcursions := false

	trajectory := make([]impact.Impact, 0, numIterations)

	trajectory = append(trajectory, initialImpact)

	for i := uint(0); i < numIterations; i++ {
		next_impact := impactMap.apply(trajectory[len(trajectory)-1])

		trajectory = append(trajectory, next_impact.impact)

		if !next_impact.FoundImpact {
			longExcursions = true
		}

		// Now check for chatter
		chatterResult := impactMap.chatterChecker.Check(trajectory[len(trajectory)-1])

		if (chatterResult.isChatter) {
			trajectory = append(trajectory, chatterResult.accumulationImpact)
		}
	}

	return &IterationResult{longExcursions: longExcursions, Impacts: trajectory}
}

// Convenient overload
func (impactMap ImpactMap) IterateFromPoint(phi float64, v float64, numIterations uint) *IterationResult {
	t := impactMap.chatterChecker.sticking.Converter.TimeIntoCycle(phi)
	return impactMap.iterate(*impactMap.GenerateImpact(t, v), numIterations)
}

// Generate a singularity set
func (impactMap ImpactMap) SingularitySet(numPoints uint) ([]impact.Impact, []impact.Impact) {
	if numPoints == 0 {
		numPoints = 1
	}

	singularitySet := make([]impact.Impact, 0, numPoints)
	dual := make([]impact.Impact, 0, numPoints)

	converter := impactMap.chatterChecker.sticking.Converter

	startingTime := converter.Period * impactMap.chatterChecker.sticking.PhaseOut
	endingTime := converter.Period * impactMap.chatterChecker.sticking.PhaseIn

	deltaTime := (endingTime - startingTime)/float64(numPoints)

	for i := uint(0); i < numPoints; i++ {
		impactResult := impactMap.apply(*impactMap.GenerateImpact(startingTime, 0.0))

		if impactResult.FoundImpact {
			dual = append(dual, impactResult.impact)

			singularitySet = append(singularitySet, impactMap.dualMaker(impactResult.impact))
		}

		startingTime += deltaTime
	}

	return singularitySet, dual
}
