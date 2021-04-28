package parameters

import (
	"testing"
)

var ParameterErrorTests = [] struct {
	ForcingFrequency float64
	CoefficientOfRestitution float64
	ObstacleOffset float64
	MaximumPeriods uint 
	NumErrors uint
}{
	{2.8, 0.0, 0.1, 100, 0},
	{2.8, 0.0, 0.1, 0, 1},
	{-2.8, 0.8, 0.1, 100, 1},
	{2.8, 0.8, -0.1, 100, 0},
	{2.8, -0.5, 0.1, 100, 1},
	{0, 2.3, 0.1, 100, 2},
	{1, 1.2, -0.1, 0, 3},
}

func TestParameterErrors(t *testing.T) {
	for _, data := range ParameterErrorTests {
		params, errors := NewParameters(data.ForcingFrequency, data.ObstacleOffset, data.CoefficientOfRestitution, data.MaximumPeriods)

		if len(errors) != int(data.NumErrors) {
			t.Errorf("Expected %d errors for %+v, got %d", data.NumErrors, data, len(errors))
		}

		if data.NumErrors == 0 && params == nil {
			t.Errorf("Expected a Parameters struct for %+v", data)
		}
	}
}