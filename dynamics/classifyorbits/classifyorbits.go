package classifyorbits

import (
    // "github.com/FelixDux/imposcg/dynamics/impact"
    "github.com/FelixDux/imposcg/dynamics/forcingphase"
    // "github.com/FelixDux/imposcg/dynamics/parameters"
	"github.com/FelixDux/imposcg/dynamics"
	// "math"
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
	if classification.chaos {
		return "(∞,∞)"
	} else if classification.longExcursions {
		return "long excursions"
	} else if classification.chatter {
		return "chatter"
	} else {
		return fmt.Sprintf("(%d,%d)", classification.numImpacts, classification.numPeriods)
	}
}

func NewOrbitClassifier(converter forcingphase.PhaseConverter) func(*dynamics.IterationResult) OrbitClassification {
	return func(iterationResult *dynamics.IterationResult) OrbitClassification {
		result := NewClassification()

		if iterationResult.LongExcursions {
			result.longExcursions = true
		} else {
			// make the last impact the comparator
			// reverse iterate through impacts
			// a zero velocity implies chatter
			// if find an impact which is 'equal' to the comparator and the difference in periods > 0
			// then we have found a periodic orbit of (number of impacts, difference of periods)
		}

		return *result
	}
}