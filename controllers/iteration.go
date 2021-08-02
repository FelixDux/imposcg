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
	"github.com/FelixDux/imposcg/dynamics/impact"
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
		return charts.ImpactMapPlot(*parameters, [][]impact.Impact{data.Impacts}, phi, v).Name()
	}
}

// PostIterationImage godoc
// @Summary Impact map
// @Description Return scatter plot from iterating the impact map for a specified set of parameters
// @ID post-iteration-image
// @Accept x-www-form-urlencoded
// @Produce  png
// @Param frequency formData number true "Forcing frequency" minimum(0) group("Dynamics")
// @Param offset formData number true "Obstacle offset from origin" group("Dynamics")
// @Param r formData number true "Coefficient of restitution" minimum(0) maximum(1) group("Dynamics")
// @Param maxPeriods formData int false "Number of periods without an impact after which the algorithm will report 'long excursions'" default(100) group("Control Parameters")
// @Param phi formData number true "Phase at initial impact" default(0.0) group("Initial Impact")
// @Param v formData number true "Velocity at initial impact" default(0.0) group("Initial Impact")
// @Param numIterations formData int false "Number of iterations of impact map" default(5000) group("Control Parameters")
// @Success 200 {object} dynamics.IterationResult
// @Failure 400 {object} string "Invalid parameters"
// @Router /iteration/image/ [post]
func PostIterationImage(c *gin.Context) {
	inputs, parameters, errorString := IterationInputsFromPost(c)

	if parameters == nil || inputs == nil {
        log.Print(errorString)
		c.JSON(400, errorString)
	} else {
		img, err := os.Open(iterationImage(parameters, inputs.phi, inputs.v, inputs.numIterations))
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
// @Summary Impact data
// @Description Return data from iterating the impact map for a specified set of parameters
// @ID post-iteration-data
// @Accept  x-www-form-urlencoded
// @Produce  json
// @Param frequency formData number true "Forcing frequency" minimum(0) group("Dynamics")
// @Param offset formData number true "Obstacle offset from origin" group("Dynamics")
// @Param r formData number true "Coefficient of restitution" minimum(0) maximum(1) group("Dynamics")
// @Param maxPeriods formData int false "Number of periods without an impact after which the algorithm will report 'long excursions'" default(100) group("Control Parameters")
// @Param phi formData number true "Phase at initial impact" default(0.0) group("Initial Impact")
// @Param v formData number true "Velocity at initial impact" default(0.0) group("Initial Impact")
// @Param numIterations formData int false "Number of iterations of impact map" default(5000) group("Control Parameters")
// @Success 200 {object} dynamics.IterationResult
// @Failure 400 {object} string "Invalid parameters"
// @Router /iteration/data/ [post]
func PostIterationData(c *gin.Context) {
	inputs, parameters, errorString := IterationInputsFromPost(c)

	if parameters == nil || inputs == nil {
        log.Print(errorString)
		c.JSON(400, errorString)

		return
	}
	
	impactMap, errMap := dynamics.NewImpactMap(*parameters)

	if errMap != nil {
		c.JSON(400,  errMap.Error())
		return
	}
	
	data := impactMap.IterateFromPoint(inputs.phi, inputs.v, inputs.numIterations)

	c.JSON(200,  data)
}

func AddIterationControllers (r *gin.Engine) {
	iteration := r.Group("/api/iteration")
	iteration.POST("/data",  PostIterationData)
	iteration.POST("/image",  PostIterationImage)
}