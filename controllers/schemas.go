package controllers

import (
	// "log"
	"strings"
	"github.com/FelixDux/imposcg/dynamics/parameters"
	"github.com/gin-gonic/gin"
)

type ParametersInput struct {
	frequency float64 `json:"omega" binding:"required"`
	offset float64 `json:"sigma" binding:"required"`
	r float64 `json:"r" binding:"required"`
	maxPeriods uint `json:"maxPeriods" binding:"required"`
}

func (input ParametersInput) ParametersFromInput() (*parameters.Parameters, []error) {
	return parameters.NewParameters(input.frequency, input.offset, input.r, input.maxPeriods)
}

func ParametersFromPost(c *gin.Context) (*parameters.Parameters, string) {

  // Validate input
  var input ParametersInput
  if err := c.ShouldBindJSON(&input); err != nil {
    return nil, err.Error()
  } else {
	params, errParams := input.ParametersFromInput()

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
}

// func ParametersFromQueryString(c *gin.Context) (*parameters.Parameters, []error) {
// 	return parameters.NewParameters(c.Param("frequency"), c.Param("offset"), c.Param("r"), c.Param("maxPeriods"))
// }