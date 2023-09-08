package singleimpact

import (
    "github.com/FelixDux/imposcg/dynamics/forcingphase"
    "github.com/FelixDux/imposcg/dynamics/parameters"
	"math"
)

// 
// Response and stability curves for single-impact periodic orbits
//

type SingleImpactOrbitSpec struct {
	NumberOfForcingPeriods uint
	Parameters parameters.Parameters
}

type IntermediateParameters struct {
	Spec SingleImpactOrbitSpec
	Converter forcingphase.PhaseConverter
	Cn float64
	Sn float64
	SnCn float64
	Ro float64
	Divisor float64
	Discriminant float64
}

func NewIntermediateParameters(spec SingleImpactOrbitSpec) (*IntermediateParameters, error) {
	converter, err := forcingphase.NewPhaseConverter(spec.Parameters.ForcingFrequency)

	if (err != nil) {
		return nil, err
	} else {
		cn := math.Cos(converter.TimeIntoCycle(float64(spec.NumberOfForcingPeriods)))
		sn := math.Sin(converter.TimeIntoCycle(float64(spec.NumberOfForcingPeriods)))

		sncn := sn *(1.0 + spec.Parameters.CoefficientOfRestitution) / (1.0 - cn)

		ro := (1 - spec.Parameters.CoefficientOfRestitution)/spec.Parameters.ForcingFrequency

		divisor := math.Pow(sncn,2) + math.Pow(ro,2)

		discriminant := 4 * (math.Pow(sncn * spec.Parameters.CoefficientOfRestitution, 2) - (math.Pow(spec.Parameters.CoefficientOfRestitution, 2) - math.Pow(spec.Parameters.Gamma,2))*divisor)

		return &IntermediateParameters{Spec: spec, Converter: *converter, Cn: cn, Sn: sn, SnCn: sncn, Ro: ro, Divisor: divisor, Discriminant: discriminant}, nil
	}
}

func (parameters IntermediateParameters) PhaseForBranch(velocity float64) float64 {
	time := math.Asin(- 0.5 * velocity * parameters.Ro / parameters.Spec.Parameters.Gamma ) / parameters.Spec.Parameters.ForcingFrequency

	return parameters.Converter.TimeToPhase(time)
}

func (parameters IntermediateParameters) VelocityForBranch(upper bool) float64 {
	firstTerm := -2*parameters.SnCn*parameters.Spec.Parameters.CoefficientOfRestitution

	if (upper) {
		return (firstTerm + parameters.Discriminant) / parameters.Divisor
	} else {
		return (firstTerm - parameters.Discriminant) / parameters.Divisor
	}
}
