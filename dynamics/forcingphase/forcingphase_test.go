package forcingphase

import (
	"math"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
    "github.com/FelixDux/imposcg/dynamics/parameters"
)

var frequencyErrorTests = [] struct {
    frequency float64
    errorType reflect.Type
}{
    {0, reflect.TypeOf(parameters.ZeroForcingFrequencyError(0))},
    {-3.4, reflect.TypeOf(parameters.NegativeForcingFrequencyError(0))},
    {4.3, reflect.TypeOf(nil)},
}

func TestFrequencyErrors(t *testing.T) {
    for _, tt := range frequencyErrorTests {
        _, gotErr := NewPhaseConverter(tt.frequency)

        gotErrType := reflect.TypeOf(gotErr)

        if gotErrType != tt.errorType {
            t.Errorf("Frequency of %g should cause %s not %s", tt.frequency, tt.errorType, gotErrType)
        }
    }
}

func TestConvertTimeToPhase(t *testing.T) {
	tm := 3.0
	f := math.Pi
    conv, _ := NewPhaseConverter(f)
	got := conv.TimeToPhase(tm)

	expect:= 0.5

	if got != expect {
		t.Errorf("ConvertTimeToPhase(%g)(%g) == %g, expected %g", f, tm, got, expect)
	}
}

func TestConvertTimeIntoCycle(t *testing.T) {
	phase := 0.80
	period := 1.25
	expect := 1.0
	f := 2.0*math.Pi/period
    conv, _ := NewPhaseConverter(f)

	got := conv.TimeIntoCycle(phase)

	if got != expect {
		t.Errorf("ConvertTimeIntoCycle(%g)(%g) == %g, expected %g", f, phase, got, expect)
	}
}

var ints = [] int {1, 2, 4, 5, 16}
var frequencies = [] float64 {4.89, 2.76}
var start_time = 0.02

const tol = 1e-6
var opt = cmp.Comparer(func(x, y float64) bool {
    return math.Abs(x-y) < tol
})

func TestShiftTimeInPeriods(t *testing.T) {
    for _, f := range frequencies {
        conv, _ := NewPhaseConverter(f)

        for _, i := range ints {
            time_shift := float64(i)*conv.Period

            shifted_time := time_shift + start_time

            new_time := conv.TimeIntoCycle(conv.TimeToPhase(shifted_time))

            if !cmp.Equal(new_time, start_time, opt) {
                t.Errorf("Converter with frequency %g does not convert consistently in both directions (start time %g, time shift %g, end time %g, %d periods)",
                        f, start_time, time_shift, new_time, i)
            }

            n := conv.DifferenceInPeriods(start_time, shifted_time)
            if i != n {
                t.Errorf("Converter with frequency %g does not return correct number of periods %g between times %g and %g: expected %d, got %d",
                        f, conv.Period, start_time, shifted_time, i, n)                
            }
        }
    }
}

func TestForwardToPhase(t *testing.T) {
    for _, f := range frequencies {
        conv, _ := NewPhaseConverter(f)

        for _, i := range ints {
            phase := 0.6
            small_time := 0.2
            big_time := 0.8

            time_delta := float64(i) * conv.Period

            new_small := conv.ForwardToPhase(time_delta + small_time, phase)
            new_small_phase := conv.TimeToPhase(new_small)

            if !cmp.Equal(phase, new_small_phase, opt) {
                t.Errorf("Phase converter with frequency %g runs forward time %g from time %g to phase %g, expected %g", f, time_delta, small_time, new_small_phase, phase)
            }

            new_big := conv.ForwardToPhase(time_delta + big_time, phase)
            new_big_phase := conv.TimeToPhase(new_big)

            if !cmp.Equal(phase, new_big_phase, opt) {
                t.Errorf("Phase converter with frequency %g runs forward time %g from time %g to phase %g, expected %g", f, time_delta, big_time, new_big_phase, phase)
            }
        }
    }
}