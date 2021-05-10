package sticking

import (
    "github.com/FelixDux/imposcg/dynamics/impact"
    "github.com/FelixDux/imposcg/dynamics/forcingphase"
    "github.com/FelixDux/imposcg/dynamics/parameters"
	"math"
)

type ReleaseImpact struct {
	NewImpact bool
	Impact impact.Impact
}

type Sticking struct {
	PhaseIn float64
	PhaseOut float64
	Converter forcingphase.PhaseConverter
	Generator impact.Generator
}

func NewSticking(parameters parameters.Parameters) (*Sticking, error) {
	var phaseIn float64
	var phaseOut float64

	if (1.0 <= parameters.ObstacleOffset) {
		// No sticking
		phaseIn = 0.0
		phaseOut = 0.0
	} else if -1.0 >= parameters.ObstacleOffset || parameters.ForcingFrequency == 0.0 {
		// Sticking for all phases
		phaseIn = 1.0
		phaseOut = 0.0
	} else { 
		converter, err := forcingphase.NewPhaseConverter(parameters.ForcingFrequency)

		if err == nil {

			// (OK to divide by.ForcingFrequency because zero case trapped above)
			angle := math.Acos(parameters.ObstacleOffset) 
			phase1 := converter.TimeToPhase(angle/parameters.ForcingFrequency) 
			phase2 := 1.0 - phase1

			if (math.Sin(angle) < 0.0) {
				phaseIn = phase1
				phaseOut = phase2
			} else {
				phaseIn = phase2
				phaseOut = phase1
			}

			return &Sticking{PhaseIn: phaseIn, PhaseOut: phaseOut, Converter: *converter, Generator: impact.ImpactGenerator(*converter)}, nil
		} else {
			return nil, err
		}
	}

	return &Sticking{PhaseIn: phaseIn, PhaseOut: phaseOut}, nil
}

func (sticking Sticking) never() bool {
	return sticking.PhaseIn == sticking.PhaseOut
}

func (sticking Sticking) always() bool {
	return sticking.PhaseIn == 1.0 && sticking.PhaseOut == 0.0
}

func (sticking Sticking) phaseSticks(phase float64) bool {
	if (sticking.never()) {return false}
	if (sticking.always()) {return true}

	return phase < sticking.PhaseOut || phase >= sticking.PhaseIn
}

func (sticking Sticking) TimeSticks(time float64) bool {
	return sticking.phaseSticks(sticking.Converter.TimeToPhase(time))
}

func (sticking Sticking) releaseTime(time float64) float64 {
	return sticking.Converter.ForwardToPhase(time, sticking.PhaseOut)
}

func (sticking Sticking) CheckImpact(impact impact.Impact) *ReleaseImpact {

	if (impact.Velocity == 0.0 && sticking.phaseSticks(impact.Phase) && !sticking.always()) {
		return &ReleaseImpact{NewImpact: true, Impact: *sticking.Generator(sticking.releaseTime(impact.Time), 0.0)}
	} else {
		return &ReleaseImpact{NewImpact: false, Impact: impact}
	}
}

