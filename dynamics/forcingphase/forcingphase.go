package forcingphase

import "math"

func TimeToPhase(period float64, simtime float64) float64 {
	scaled_time := simtime / period

	return scaled_time - math.Floor(scaled_time)
}

func TimeIntoCycle(period float64, phase float64) float64 {
	return phase * period
}

func ForwardToPhase(period float64, starttime float64, phase float64) float64 {
    phase_change := phase - TimeToPhase(period, starttime)

    if (phase_change < 0) {
        phase_change++
	}

    return starttime + period * phase_change
}

/*

convert_difference_in_periods(period::Time, time1::Time, time2::Time)::Int = Int(round(abs(time1 - time2)/period))

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