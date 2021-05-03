package controllers

import (
	// "log"
	"fmt"
	"strings"
	"strconv"
	"github.com/FelixDux/imposcg/dynamics/parameters"
	"github.com/gin-gonic/gin"
)

type ParametersInput struct {
	frequency float64 `json:"omega" binding:"required"`
	offset float64 `json:"sigma" binding:"required"`
	r float64 `json:"r" binding:"required"`
	maxPeriods uint `json:"maxPeriods" binding:"required"`
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

func ParametersFromPost(c *gin.Context) (*parameters.Parameters, string) {

  // Validate input
  var input ParametersInput
  if err := c.ShouldBindJSON(&input); err != nil {
    return nil, err.Error()
  } else {
	return input.ParametersFromInput()
  }
}

func float64FromQueryString(c *gin.Context, name string) (float64, string) {
	valueString := c.Query(name)
	value, err := strconv.ParseFloat(valueString, 64)

	if err != nil {
		return 0.0, fmt.Sprintf("Invalid value %s for parameter %s", valueString, name)
	} else {
		return value, ""
	}
}

func uintFromQueryString(c *gin.Context, name string) (uint, string) {
	valueString := c.Query(name)
	value, err := strconv.ParseUint(valueString, 10, 32)

	if err != nil {
		return 0, fmt.Sprintf("Invalid value %s for parameter %s", valueString, name)
	} else {
		return uint(value), ""
	}
}

func ParametersFromQueryString(c *gin.Context) (*parameters.Parameters, string) {
	input := ParametersInput{}
	errorStrings := make([]string, 0, 6)

	frequency, freqErrString := float64FromQueryString(c, "frequency")

	if freqErrString != "" {
		errorStrings = append(errorStrings, freqErrString)
	} else {
		input.frequency = frequency
	}

	offset, freqErrString := float64FromQueryString(c, "offset")

	if freqErrString != "" {
		errorStrings = append(errorStrings, freqErrString)
	} else {
		input.offset = offset
	}

	r, freqErrString := float64FromQueryString(c, "r")

	if freqErrString != "" {
		errorStrings = append(errorStrings, freqErrString)
	} else {
		input.r = r
	}

	maxPeriods, freqErrString := uintFromQueryString(c, "maxPeriods")

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