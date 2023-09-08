package singleimpact

import (
    "github.com/FelixDux/imposcg/dynamics/forcingphase"
	"github.com/FelixDux/imposcg/dynamics/impact"
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
	LeftIntercept float64
	RightIntercept float64
	SaddleNodePoint float64
}

func NewIntermediateParameters(spec SingleImpactOrbitSpec) (*IntermediateParameters, error) {
	converter, err := forcingphase.NewPhaseConverter(spec.Parameters.ForcingFrequency)

	if (err != nil) {
		return nil, err
	} else {
		cn := math.Cos(converter.TimeIntoCycle(float64(spec.NumberOfForcingPeriods)))
		sn := math.Sin(converter.TimeIntoCycle(float64(spec.NumberOfForcingPeriods)))

		sncn := sn *(1.0 + spec.Parameters.CoefficientOfRestitution) / (1.0 - cn)
		ro := (1 - spec.Parameters.CoefficientOfRestitution) / spec.Parameters.ForcingFrequency

		divisor := math.Pow(sncn,2) + math.Pow(ro,2)

		discriminant := 4 * (math.Pow(sncn * spec.Parameters.CoefficientOfRestitution, 2) - 
		(math.Pow(spec.Parameters.CoefficientOfRestitution, 2) - 
		math.Pow(spec.Parameters.Gamma,2))*divisor)

		rightIntercept := math.Abs(spec.Parameters.Gamma)

		saddleNodePoint := rightIntercept * math.Sqrt(1 + math.Pow(
			spec.Parameters.ForcingFrequency * sncn / (1 - spec.Parameters.CoefficientOfRestitution)))

		if (spec.Parameters.ForcingFrequency > 2 * float64(spec.NumberOfForcingPeriods)) {
			saddleNodePoint *= -1
		} else if (spec.Parameters.ForcingFrequency == 2 * float64(spec.NumberOfForcingPeriods)) {
			saddleNodePoint = -rightIntercept
		}

		return &IntermediateParameters{Spec: spec, Converter: *converter, Cn: cn, Sn: sn, 
			SnCn: sncn, Ro: ro, Divisor: divisor, Discriminant: discriminant, 
			LeftIntercept: -rightIntercept, RightIntercept: math.Abs(spec.Parameters.Gamma),
			SaddleNodePoint: saddleNodePoint}, nil
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

func (parameters IntermediateParameters) ImpactForBranch(upper bool) impact.SimpleImpact {
	velocity := parameters.VelocityForBranch(upper)

	return impact.SimpleImpact{Velocity: velocity, Phase: parameters.PhaseForBranch(velocity)}
}

func (parameters IntermediateParameters) IsPhysical(solution impact.SimpleImpact) bool {
	generator := impact.ImpactGenerator(parameters.Converter)

	nextImpact := generator(parameters.Converter.TimeIntoCycle(solution.Phase), solution.Velocity)

	return nextImpact.AlmostEqual(impact.Impact{Time: solution.Phase, Phase: solution.Phase, 
		Velocity: solution.Velocity})
}

func (parameters IntermediateParameters) PeriodDoublingOffset(velocity float64) float64 {
	return velocity * ( (1 - parameters.Spec.Parameters.CoefficientOfRestitution - 
		2*parameters.Spec.Parameters.CoefficientOfRestitution*parameters.Cn)/parameters.Sn +
		parameters.SnCn / (2 * parameters.Spec.Parameters.Gamma) ) / math.Pow(parameters.Spec.Parameters.ForcingFrequency, 2)
}
