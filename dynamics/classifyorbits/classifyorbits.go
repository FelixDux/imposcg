package classifyorbits

import (
    "github.com/FelixDux/imposcg/dynamics/forcingphase"
	"github.com/FelixDux/imposcg/dynamics"
	"fmt"
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

	return result
}


func (classifier OrbitClassifier) ClassifyForRange(numPhases uint, numVelocities uint, 
	minVelocity float64, maxVelocity float64, result chan OrbitClassificationResult) {

	deltaPhi := 1.0 / float64(numPhases+1)
	deltaV := (maxVelocity - minVelocity) / float64(numVelocities+1)

	phi := deltaPhi / 2.0
	v := deltaV / 2.0

	for i:=uint(0); i < numPhases; i++ {
		for j:=uint(0); j < numVelocities; j++ {
			result <- classifier.Classify(phi, v)

			phi += deltaPhi
			v += deltaV
		}
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