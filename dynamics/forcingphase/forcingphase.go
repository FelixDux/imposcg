package forcingphase

import (
    "math"
    "fmt"
)

func ConvertTimeToPhase(period float64) func(float64) float64 {
    return func(simtime float64) float64 {
        scaled_time := simtime / period

        return scaled_time - math.Floor(scaled_time)
    }
}

func ConvertTimeIntoCycle(period float64) func(float64) float64 {
    return func(phase float64) float64 {
        return phase * period
    }
}

func RunForwardToPhase(period float64) func(float64, float64) float64 {
    return func(starttime float64, phase float64) float64 {
        phase_change := phase - ConvertTimeToPhase(period)(starttime)

        if (phase_change < 0) {
            phase_change++
        }

        return starttime + period * phase_change
    }
}

func GetDifferenceInPeriods(period float64) func(float64, float64) int {
    return func (starttime float64, endtime float64) int {
        return int(math.Round(math.Abs(endtime - starttime)/period))
    }
}

type PhaseConverter struct {
    Period float64

    TimeToPhase func(float64) float64
    TimeIntoCycle func(float64) float64
    ForwardToPhase func(float64, float64) float64
    DifferenceInPeriods func(float64, float64) int
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
        TimeToPhase: ConvertTimeToPhase(period),
        TimeIntoCycle: ConvertTimeIntoCycle(period),
        ForwardToPhase: RunForwardToPhase(period),
        DifferenceInPeriods: GetDifferenceInPeriods(period),
    }, nil
}
