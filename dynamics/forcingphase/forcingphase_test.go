package forcingphase

import "testing"

// Look here for table-driven tests https://dave.cheney.net/2013/06/09/writing-table-driven-tests-in-go

func TestConvertTimeToPhase(t *testing.T) {
	tm := 3.0
	p := 2.0
	got := ConvertTimeToPhase(p)(tm)

	expect:= 0.5

	if got != expect {
		t.Errorf("ConvertTimeToPhase(%f)(%f) == %f, expected %f", p, tm, got, expect)
	}
}

func TestConvertTimeIntoCycle(t *testing.T) {
	phase := 0.80
	period := 1.25
	expect := 1.0
	got := ConvertTimeIntoCycle(period)(phase)

	if got != expect {
		t.Errorf("ConvertTimeIntoCycle(%f)(%f) == %f, expected %f", period, phase, got, expect)
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

        @testset "Phase converter handles error conditions" for params in (
            Dict("frequency" => 0, "message" => "Forcing frequency cannot be zero"),
            Dict([("frequency", -1), ("message", "The model cannot handle negative forcing frequencies")]))
            begin
                @test_throws ErrorException(params["message"]) converter = PhaseConverter(params["frequency"])
            end
        end

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