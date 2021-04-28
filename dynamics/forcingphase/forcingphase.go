package forcingphase

import (
    "math"
    "github.com/FelixDux/imposcg/dynamics/parameters"
)

type PhaseConverter struct {
    Period float64
}

func NewPhaseConverter(frequency float64) (*PhaseConverter, error) {
    if (frequency == 0) {
        return nil, parameters.ZeroForcingFrequencyError(frequency)
    }

    if (frequency < 0) {
        return nil, parameters.NegativeForcingFrequencyError(frequency)
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