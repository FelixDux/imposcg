package forcingphase

import ("testing"
        "reflect"
        "math")

// Look here for table-driven tests https://dave.cheney.net/2013/06/09/writing-table-driven-tests-in-go

var frequencyErrorTests = [] struct {
    frequency float64
    errorType reflect.Type
}{
    {0, reflect.TypeOf(ZeroForcingFrequencyError(0))},
    {-3.4, reflect.TypeOf(NegativeForcingFrequencyError(0))},
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

/*
module TestDynamics
    using Test
    using Dynamics

    @testset "Phase conversion" begin
        ints = (1, 2, 4, 5, 16)
        frequencies = (4.89, 2.76)
        start_time = 0.02;

        @testset "Phase converter converts consistently in both directions" for i in ints, f in frequencies
            begin
                converter = PhaseConverter(f)

                new_time = i * converter.period + start_time |> converter.time_to_phase |> converter.time_into_cycle

                @test isapprox(new_time, start_time)
                
            end
        end

        @testset "Phase converter returns correct number of periods between times" for i in ints, f in frequencies
            begin
                converter = PhaseConverter(f)

                new_time = i * converter.period + start_time 

                @test i == converter.difference_in_periods(start_time, new_time)
                
            end
        end

        @testset "Phase converter runs forward to specified phase correctly" for i in ints, f in frequencies[2:end-1]
            begin
                converter = PhaseConverter(f)
                phase = 0.6
                small_time = 0.2
                big_time = 0.8

                time_delta = i * converter.period

                new_small = converter.forward_to_phase(time_delta + small_time, phase)

                @test isapprox(converter.time_to_phase(new_small), phase)
                @test new_small > small_time
                @test isapprox(new_small - time_delta, phase*converter.period)

                new_big = converter.forward_to_phase(time_delta + big_time, phase)

                @test isapprox(converter.time_to_phase(new_big), phase)
                @test new_big > big_time
                @test isapprox(new_big - time_delta, phase*converter.period)
            end
        end
    end
end
*/