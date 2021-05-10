package sticking

import (
	"testing"
    "github.com/FelixDux/imposcg/dynamics/parameters"
)

func TestStickingRegion(t *testing.T) {
	frequencies := [] float64 {2.8, 3.7, 4.0}

	r := 0.8

	offset := 0.0

	for _, frequency := range(frequencies) {
		params, _ := parameters.NewParameters(frequency, offset, r, 100)

		sticking, _ := NewSticking(*params)

		impactTime := 0.0

		if sticking.TimeSticks(impactTime) == false {
			t.Errorf("Impact time %f should be sticking for %+v but isn't", impactTime, params)
		}

		if !(sticking.releaseTime(impactTime) > impactTime) {
			t.Errorf("Release time for %+v should be greater then impact time of %f", params, impactTime)
		}

		if sticking.PhaseIn < sticking.PhaseOut {
			t.Errorf("Sticking interval (%f, %f) for %+v is the wrong way round", sticking.PhaseOut, sticking.PhaseIn, params)
		}

	}
}