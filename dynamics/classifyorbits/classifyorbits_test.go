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

func TestMarshalClassifications(t *testing.T) {
	input := []OrbitClassificationResult{
		OrbitClassificationResult{Phi: 0.0, V: 0.0, Classification: OrbitClassification{ chaos: false, chatter: true, longExcursions: false, numImpacts: 1, numPeriods: 0}},
		OrbitClassificationResult{Phi: 0.5, V: 0.0, Classification: OrbitClassification{ chaos: true, chatter: false, longExcursions: false, numImpacts: 1, numPeriods: 0}},
		OrbitClassificationResult{Phi: 1.0, V: 0.0, Classification: OrbitClassification{ chaos: false, chatter: false, longExcursions: false, numImpacts: 3, numPeriods: 3}},

		OrbitClassificationResult{Phi: 0.0, V: 0.5, Classification: OrbitClassification{ chaos: false, chatter: true, longExcursions: false, numImpacts: 1, numPeriods: 0}},
		OrbitClassificationResult{Phi: 0.5, V: 0.5, Classification: OrbitClassification{ chaos: true, chatter: false, longExcursions: false, numImpacts: 1, numPeriods: 0}},
		OrbitClassificationResult{Phi: 1.0, V: 0.5, Classification: OrbitClassification{ chaos: false, chatter: false, longExcursions: false, numImpacts: 3, numPeriods: 3}},

		OrbitClassificationResult{Phi: 0.0, V: 1.0, Classification: OrbitClassification{ chaos: false, chatter: true, longExcursions: false, numImpacts: 1, numPeriods: 0}},
		OrbitClassificationResult{Phi: 0.5, V: 1.0, Classification: OrbitClassification{ chaos: true, chatter: false, longExcursions: false, numImpacts: 1, numPeriods: 0}},
		OrbitClassificationResult{Phi: 1.0, V: 1.0, Classification: OrbitClassification{ chaos: false, chatter: false, longExcursions: false, numImpacts: 3, numPeriods: 3}},
	}

	expected := map[string]float64{
		"(∞,∞)": 0.5,
		"chatter": 0.0,
		"(3,3)": 1.0,
	}

	actual := MarshalClassifications(&input)

	for label, impacts := range *actual {
		if phase, present := expected[label]; present {
			for _, impact := range impacts {
				if impact.Phase != phase {
					t.Errorf("Unexpected phase %v returned for label %s - expected %v", impact.Phase, label, phase)
				}
			}
		} else {
			t.Errorf("Label %s returned unexpectedly", label)
		}
	}

}
