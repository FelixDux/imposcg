package classifyorbits

import (
	"testing"
	// "reflect"
)

func TestClassificationLabels(t *testing.T) {

	classificationLabelTests := map[string]OrbitClassification{
		"(∞,∞)": {chaos: true, chatter: false, longExcursions: false, numImpacts: 99, numPeriods: 12},
		"(99,12)": {chaos: false, chatter: false, longExcursions: false, numImpacts: 99, numPeriods: 12},
		"chatter": {chaos: false, chatter: true, longExcursions: false, numImpacts: 99, numPeriods: 12},
		"long excursions": {chaos: true, chatter: false, longExcursions: true, numImpacts: 99, numPeriods: 12},
	}

    for expected, classification := range classificationLabelTests {
        actual := classification.label()

        if actual != expected {
            t.Errorf("Label for %+v should be %s not %s", classification, expected, actual)
        }
    }
}
