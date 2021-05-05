package controllers

import (
	"fmt"
	"strings"
	"strconv"
	"github.com/FelixDux/imposcg/dynamics/parameters"
	"github.com/gin-gonic/gin"
)

type ParametersInput struct {
	frequency float64
	offset float64
	r float64 
	maxPeriods uint
}

func (input ParametersInput) ParametersFromInput() (*parameters.Parameters, string) {
	params, errParams := parameters.NewParameters(input.frequency, input.offset, input.r, input.maxPeriods)

	if len(errParams) > 0 {
	
		paramMessages := make([]string,len(errParams))
		for i, err := range(errParams) {

			paramMessages[i] = err.Error()

		}

		return nil, strings.Join(paramMessages, "\n")
	} else {
		return params, ""
	}
}

func float64FromSource(c *gin.Context, name string, source func(string) string) (float64, string) {
	valueString := source(name)
	value, err := strconv.ParseFloat(valueString, 64)

	if err != nil {
		return 0.0, fmt.Sprintf("Invalid value %s for parameter %s", valueString, name)
	} else {
		return value, ""
	}
}

func uintFromSource(c *gin.Context, name string, source func(string) string) (uint, string) {
	valueString := source(name)
	value, err := strconv.ParseUint(valueString, 10, 32)

	if err != nil {
		return 0, fmt.Sprintf("Invalid value %s for parameter %s", valueString, name)
	} else {
		return uint(value), ""
	}
}

func ParametersFromSource(c *gin.Context, source func(string) string) (*parameters.Parameters, string) {
	input := ParametersInput{}
	errorStrings := make([]string, 0, 6)

	frequency, freqErrString := float64FromSource(c, "frequency", source)

	if freqErrString != "" {
		errorStrings = append(errorStrings, freqErrString)
	} else {
		input.frequency = frequency
	}

	offset, freqErrString := float64FromSource(c, "offset", source)

	if freqErrString != "" {
		errorStrings = append(errorStrings, freqErrString)
	} else {
		input.offset = offset
	}

	r, freqErrString := float64FromSource(c, "r", source)

	if freqErrString != "" {
		errorStrings = append(errorStrings, freqErrString)
	} else {
		input.r = r
	}

	maxPeriods, freqErrString := uintFromSource(c, "maxPeriods", source)

	if freqErrString != "" {
		errorStrings = append(errorStrings, freqErrString)
	} else {
		input.maxPeriods = maxPeriods
	}

	if len(errorStrings) > 0 {
		return nil, strings.Join(errorStrings, "\n")
	} else {
		return input.ParametersFromInput()
	}
}

func ParametersFromQueryString(c *gin.Context) (*parameters.Parameters, string) {
	return ParametersFromSource(c, c.Query)
}

func ParametersFromPost(c *gin.Context) (*parameters.Parameters, string) {
	return ParametersFromSource(c, c.Request.PostFormValue)
}

type IterationInputs struct {
	phi float64
	v float64
	numIterations uint
}

func IterationInputsFromSource(c *gin.Context, source func(string) string) (*IterationInputs, *parameters.Parameters, string) {
	input := IterationInputs{}
	errorStrings := make([]string, 0, 6)

	phi, freqErrString := float64FromSource(c, "phi", source)

	if freqErrString != "" {
		errorStrings = append(errorStrings, freqErrString)
	} else {
		input.phi = phi
	}

	v, freqErrString := float64FromSource(c, "v", source)

	if freqErrString != "" {
		errorStrings = append(errorStrings, freqErrString)
	} else {
		input.v = v
	}

	numIterations, freqErrString := uintFromSource(c, "numIterations", source)

	if freqErrString != "" {
		errorStrings = append(errorStrings, freqErrString)
	} else {
		input.numIterations = numIterations
	}

	if len(errorStrings) > 0 {
		return nil, nil, strings.Join(errorStrings, "\n")
	} else {
		parameters, paramErrors := ParametersFromSource(c, source)

		return &input, parameters, paramErrors
	}
}

func IterationInputsFromPost(c *gin.Context) (*IterationInputs, *parameters.Parameters, string) {
	return IterationInputsFromSource(c, c.Request.PostFormValue)
}