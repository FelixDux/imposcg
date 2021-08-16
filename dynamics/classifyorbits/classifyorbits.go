package classifyorbits

import (
	"fmt"
	"runtime"

	"github.com/FelixDux/imposcg/dynamics"
	"github.com/FelixDux/imposcg/dynamics/forcingphase"
	"github.com/FelixDux/imposcg/dynamics/impact"
)

type OrbitClassification struct {
	chaos bool
	chatter bool
	longExcursions bool
	numImpacts uint
	numPeriods uint
}

func NewClassification() *OrbitClassification {
	return &OrbitClassification{ chaos: true, chatter: false, longExcursions: false, numImpacts: 1, numPeriods: 0}
}

func (classification OrbitClassification) label () string {
	if classification.longExcursions {
		return "long excursions"
	} else if classification.chaos {
		return "(∞,∞)"
	} else if classification.chatter {
		return "chatter"
	} else {
		return fmt.Sprintf("(%d,%d)", classification.numImpacts, classification.numPeriods)
	}
}

type OrbitClassifier struct {
	impactMap *dynamics.ImpactMap
	numIterations uint
	classifier func(*dynamics.IterationResult) OrbitClassification
}

type OrbitClassificationResult struct {
	Phi float64
	V float64
	Classification OrbitClassification
}

func (classifier OrbitClassifier) Classify(phi float64, velocity float64) OrbitClassificationResult {
	iterationResult := classifier.impactMap.IterateFromPoint(phi, velocity, classifier.numIterations)

	return OrbitClassificationResult{Phi: phi, V: velocity, Classification: classifier.classifier(iterationResult)}
}


func (classifier OrbitClassifier) BuildClassification(numPhases uint, numVelocities uint, maxVelocity float64) [] OrbitClassificationResult {

	result := make([]OrbitClassificationResult, numPhases*numVelocities)

	numRanges := runtime.GOMAXPROCS(0)-1

	if numRanges < 1 {
		numRanges = 1
	}

	channels := make([]chan OrbitClassificationResult, numRanges)

	resultsCounters := make([]uint, numRanges)

	deltaMaxVelocity := maxVelocity / float64(numRanges)

	velocitiesPerRange := numVelocities / uint(numRanges)

	minVelocity := float64(0)

	maxVelocityForRange := deltaMaxVelocity

	velocitiesRemaining := numVelocities

	for i := 0; i < numRanges; i++ {
		if i == numRanges-1 {
			velocitiesPerRange = velocitiesRemaining
			maxVelocityForRange = maxVelocity
		} else {
			velocitiesRemaining -= velocitiesPerRange
		}

		resultsCounters[i] = velocitiesPerRange*numPhases

		channels[i] = make(chan OrbitClassificationResult, velocitiesPerRange*numPhases)
		
		go classifier.ClassifyForRange(numPhases, velocitiesPerRange, minVelocity, maxVelocityForRange,
			channels[i])

		minVelocity = maxVelocityForRange
		maxVelocityForRange += deltaMaxVelocity
	}

	resultCount := uint(0)

	numResults := numPhases*numVelocities

	for resultCount < numResults {

		for i := 0; i < numRanges && resultCount < numResults; i++ {
			if resultsCounters[i] > 0 {
				result[resultCount] = <- channels[i]

				resultCount++

				resultsCounters[i]--
			}
		}
	}

	return result
}


func (classifier OrbitClassifier) ClassifyForRange(numPhases uint, numVelocities uint, 
	minVelocity float64, maxVelocity float64, result chan<- OrbitClassificationResult) {

	deltaPhi := 1.0 / float64(numPhases+1)
	deltaV := (maxVelocity - minVelocity) / float64(numVelocities+1)

	phi := deltaPhi / 2.0

	for i:=uint(0); i < numPhases; i++ {
		v := deltaV / 2.0
		for j:=uint(0); j < numVelocities; j++ {
			result <- classifier.Classify(phi, v)
			v += deltaV
		}
		phi += deltaPhi
	}

}

func NewOrbitClassifier(impactMap *dynamics.ImpactMap, numIterations uint) *OrbitClassifier {
	return &OrbitClassifier{impactMap: impactMap, numIterations: numIterations, 
	classifier: NewOrbitClassifierFunction(*impactMap.Converter()),}
}

func NewOrbitClassifierFunction(converter forcingphase.PhaseConverter) func(*dynamics.IterationResult) OrbitClassification {
	return func(iterationResult *dynamics.IterationResult) OrbitClassification {
		result := NewClassification()

		if iterationResult.LongExcursions {
			result.longExcursions = true
		} else {

			lastIdx := len(iterationResult.Impacts)-1

			// make the last impact the comparator
			comparator := iterationResult.Impacts[lastIdx]

			// reverse iterate through impacts
			for i := lastIdx; lastIdx >= 0 && result.chaos; lastIdx-- {
				impact := iterationResult.Impacts[i]
				result.numImpacts++

				// a zero velocity implies chatter
				if impact.Velocity == 0 {
					result.chatter = true
				}

				// if find an impact which is 'equal' to the comparator and the difference in periods > 0
				// then we have found a periodic orbit of (number of impacts, difference of periods)
				if impact.AlmostEqual(comparator) {
					result.numPeriods = uint(converter.DifferenceInPeriods(comparator.Time, impact.Time))

					if result.numPeriods > 0 {
						result.chaos = false
					}
				}
			}
		}

		return *result
	}
}

func MarshalClassifications(classifications * [] OrbitClassificationResult) *map[string][]impact.SimpleImpact {
	result := map[string][]impact.SimpleImpact{}

	for _, classification := range *classifications {
		label := classification.Classification.label()

		result[label] = append(result[label], impact.SimpleImpact{Phase: classification.Phi, Velocity: classification.V})
	}

	return &result
}