package controllers

// https://github.com/swaggo/gin-swagger to generate REST spec
// https://github.com/swaggo/swag#declarative-comments-format

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/FelixDux/imposcg/charts"
	"github.com/FelixDux/imposcg/dynamics"
	"github.com/FelixDux/imposcg/dynamics/parameters"
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

// PostIterationImage godoc
// @Summary Return scatter plot from iterating the impact map
// @Description Return scatter plot from iterating the impact map for a specified set of parameters
// @ID post-iteration-image
// @Accept x-www-form-urlencoded
// @Produce  png
// @Param frequency formData number true "Forcing frequency"
// @Param offset formData number true "Obstacle offset from origin"
// @Param r formData number true "Coefficient of restitution"
// @Param maxPeriods formData int true "Number of periods without an impact after which the algorithm will report 'long excursions'"
// @Success 200 {object} dynamics.IterationResult
// @Failure 400 {object} string "Invalid parameters"
// @Router /iteration/image/ [post]
func PostIterationImage(c *gin.Context) {
	parameters, errorString := ParametersFromPost(c)

	if parameters == nil {
        log.Print(errorString)
		c.JSON(400, fmt.Sprintf("Invalid parameters - %s", errorString))
	} else {
		img, err := os.Open(iterationImage(parameters, 0.0, 0.0, 1000))
		if err != nil {
			log.Print(err)
			c.JSON(400, fmt.Sprintf("Failed to complete iteration - %s", err.Error()))
		} else {
			defer img.Close()
			c.Writer.Header().Set("Content-Type", "image/png")
			io.Copy(c.Writer, img)
		}
	}
}

// PostIterationData godoc
// @Summary Return data from iterating the impact map
// @Description Return data from iterating the impact map for a specified set of parameters
// @ID post-iteration-data
// @Accept  x-www-form-urlencoded
// @Produce  json
// @Param frequency formData number true "Forcing frequency"
// @Param offset formData number true "Obstacle offset from origin"
// @Param r formData number true "Coefficient of restitution"
// @Param maxPeriods formData int true "Number of periods without an impact after which the algorithm will report 'long excursions'"
// @Success 200 {object} dynamics.IterationResult
// @Failure 400 {object} string "Invalid parameters"
// @Router /iteration/data/ [post]
func PostIterationData(c *gin.Context) {
	parameters, errorString := ParametersFromPost(c)

	if parameters == nil {
        log.Print(errorString)
		c.JSON(400, errorString)

		return
	}
	
	impactMap, errMap := dynamics.NewImpactMap(*parameters)

	if errMap != nil {
		c.JSON(400,  errMap.Error())
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
	iteration := r.Group("/api/iteration")
	iteration.POST("/data",  PostIterationData)
	iteration.POST("/image",  PostIterationImage)
}