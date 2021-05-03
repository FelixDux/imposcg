package controllers

import (
	"os"
	"log"
	"io"

	"github.com/FelixDux/imposcg/dynamics"
	"github.com/FelixDux/imposcg/dynamics/parameters"
	"github.com/FelixDux/imposcg/charts"
	"github.com/gin-gonic/gin"
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

// GET /api/iteration/image
func GetIterationImage(c *gin.Context) {
	parameters, errorString := ParametersFromQueryString(c)

	if parameters == nil {
        log.Print(errorString)
		c.JSON(500, errorString)
	} else {
		img, err := os.Open(iterationImage(parameters, 0.0, 0.0, 1000))
		if err != nil {
			log.Print(err)
			c.JSON(500, err.Error())
		} else {
			defer img.Close()
			c.Writer.Header().Set("Content-Type", "image/png")
			io.Copy(c.Writer, img)
		}
	}
}

// GET /api/iteration/data
func GetIterationData(c *gin.Context) {
	parameters, errorString := ParametersFromQueryString(c)

	if parameters == nil {
        log.Print(errorString)
		c.JSON(500, errorString)

		return
	}
	
	impactMap, errMap := dynamics.NewImpactMap(*parameters)

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

func AddIterationControllers (r *gin.Engine) {
	r.GET("/api/iteration/data",  GetIterationData)
	r.GET("/api/iteration/image",  GetIterationImage)
}