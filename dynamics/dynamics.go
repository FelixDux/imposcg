package dynamics
import (
    "github.com/FelixDux/imposcg/dynamics/impact"
    "github.com/FelixDux/imposcg/dynamics/forcingphase"
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

func NewChatterChecker(parameters parameters.Parameters, velocityThreshold float64, countThreshold uint) *ChatterChecker {
	sticking, _ := sticking.NewSticking(parameters)

	// TODO: handle err

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
				sticking: *sticking}
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

func DefaultChatterChecker(parameters parameters.Parameters) *ChatterChecker {
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
}

func ImpactMapGenerator(parameters parameters.Parameters) func (impact.Impact) *ImpactMap {
	initialiser, _ := motion.MotionBetweenImpactsInitialiser(parameters)

	// TODO: handle err

	return func(impact impact.Impact) *ImpactMap {return &ImpactMap{motion: *initialiser(impact), chatterChecker: *DefaultChatterChecker(parameters)}}
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

	return &IterationResult{longExcursions: longExcursions}
}

// Convenient overload
func (impactMap ImpactMap) IterateFromPoint(phi float64, v float64, numIterations uint) *IterationResult {
	t := impactMap.chatterChecker.sticking.Converter.TimeIntoCycle(phi)
	return impactMap.iterate(*impactMap.GenerateImpact(t, v), numIterations)
}

// Generate a singularity set
func (impactMap ImpactMap) singularity_set(numPoints uint) []impact.Impact {
	if numPoints == 0 {
		numPoints = 1
	}

	impacts := make([]impact.Impact, 0, numPoints)

	deltaTime := impactMap.chatterChecker.sticking.Converter.Period/float64(numPoints)

	startingTime := 0.0

	for i := uint(0); i < numPoints; i++ {
		impactResult := impactMap.apply(*impactMap.GenerateImpact(startingTime, 0))

		if impactResult.FoundImpact {
			impacts = append(impacts, impactResult.impact)
		}

		startingTime += deltaTime
	}

	return impacts
}