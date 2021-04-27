package forcingphase

import "math"

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

func NewPhaseConverter(frequency float64) *PhaseConverter {
    period := 2.0 * math.Pi / frequency

    return &PhaseConverter{
        Period: period,
        TimeToPhase: ConvertTimeToPhase(period),
        TimeIntoCycle: ConvertTimeIntoCycle(period),
        ForwardToPhase: RunForwardToPhase(period),
        DifferenceInPeriods: GetDifferenceInPeriods(period)
    }
}

/*

struct PhaseConverter
    time_to_phase
    time_into_cycle
    forward_to_phase
    difference_in_periods
    period

    PhaseConverter(frequency::Frequency) = begin
        if frequency > 0
            period = 2 * pi / frequency
        elseif frequency == 0
            error("Forcing frequency cannot be zero")
        else
            error("The model cannot handle negative forcing frequencies")
        end

        time_to_phase = (time) -> convert_time_to_phase(period, time)
        time_into_cycle = (time) -> convert_time_into_cycle(period, time)
        forward_to_phase = (time, phase) -> convert_forward_to_phase(period, time, phase)
        difference_in_periods = (time1, time2) -> convert_difference_in_periods(period, time1, time2)

        new(time_to_phase, time_into_cycle, forward_to_phase, difference_in_periods, period)
    end
end
*/