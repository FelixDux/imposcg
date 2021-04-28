package forcingphase

import (
    "math"
    "fmt"
)

type PhaseConverter struct {
    Period float64
}

type ZeroForcingFrequencyError float64

func (f ZeroForcingFrequencyError) Error() string {
    return "Forcing frequency cannot be zero"
}

type NegativeForcingFrequencyError float64

func (f NegativeForcingFrequencyError) Error() string {
    return fmt.Sprintf("The model cannot handle negative forcing frequency %g", f)
}

func NewPhaseConverter(frequency float64) (*PhaseConverter, error) {
    if (frequency == 0) {
        return nil, ZeroForcingFrequencyError(frequency)
    }

    if (frequency < 0) {
        return nil, NegativeForcingFrequencyError(frequency)
    }

    period := 2.0 * math.Pi / frequency

    return &PhaseConverter{
        Period: period,
    }, nil
}

func (converter PhaseConverter) TimeToPhase(simtime float64) float64 {
    scaled_time := simtime / converter.Period

    return scaled_time - math.Floor(scaled_time)
}

func (converter PhaseConverter) TimeIntoCycle (phase float64) float64 {
    return phase * converter.Period
}

func (converter PhaseConverter) ForwardToPhase (starttime float64, phase float64) float64 {
    phase_change := phase - converter.TimeToPhase(starttime)

    if (phase_change < 0) {
        phase_change++
    }

    return starttime + converter.Period * phase_change
}

func (converter PhaseConverter) DifferenceInPeriods (starttime float64, endtime float64) int {
    return int(math.Round(math.Abs(endtime - starttime)/converter.Period))
}