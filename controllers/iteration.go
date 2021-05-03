package controllers

import (
	"os"
	"log"
	"strings"
	"io"

	"github.com/FelixDux/imposcg/dynamics"
	"github.com/FelixDux/imposcg/dynamics/parameters"
	"github.com/FelixDux/imposcg/charts"
	"github.com/gin-gonic/gin"
	// "net/http"
)

func iterationData(parameters *parameters.Parameters, phi float64, v float64, numIterations uint) (*dynamics.IterationResult, string) {
	impactMap, errMap := dynamics.NewImpactMap(*parameters)

	if errMap != nil {
		return nil, errMap.Error()
	}
	
	return impactMap.IterateFromPoint(phi, v, numIterations), ""
}

func iterationImage(parameters *parameters.Parameters, phi float64, v float64, numIterations uint) string {
	data, errString := iterationData(parameters, phi, v, numIterations)

	if data == nil {
		return errString
	} else {
		return charts.ImpactMapPlot(*parameters, data.Impacts, phi, v).Name()
	}
}

func doScatter() string {

	params, errParams := parameters.NewParameters(4.85, 0.01, 0.8, 100)

	if len(errParams) > 0 {
	
		paramMessages := make([]string,len(errParams))
		for i, err := range(errParams) {

			paramMessages[i] = err.Error()

		}

		return strings.Join(paramMessages, "\n")
	}

	phi := 0.0
	v := 0.0
	
	return iterationImage(params, phi, v, 1000)
}

// GET /api/iteration/image
func GetIterationImage(c *gin.Context) {

	img, err := os.Open(doScatter())
    if err != nil {
        log.Print(err)
		c.JSON(500, err.Error())
    } else {
		defer img.Close()
		c.Writer.Header().Set("Content-Type", "image/png")
		io.Copy(c.Writer, img)
	}
}

// GET /api/iteration/data
func GetIterationData(c *gin.Context) {

	params, errParams := parameters.NewParameters(2.8, 0.0, 0.8, 100)

	if len(errParams) > 0 {
	
		paramMessages := make([]string,len(errParams))
		for i, err := range(errParams) {

			paramMessages[i] = err.Error()

		}

		c.JSON(500, strings.Join(paramMessages, "\n"))
		return
	}
	
	impactMap, errMap := dynamics.NewImpactMap(*params)

	if errMap != nil {
		c.JSON(500,  errMap.Error())
		return
	}

	phi := 0.0
	v := 0.0
	data := impactMap.IterateFromPoint(phi, v, 1000)

	c.JSON(200, gin.H{
		"message": data,
	})
}