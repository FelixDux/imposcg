package controllers

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/FelixDux/imposcg/charts"
	"github.com/FelixDux/imposcg/dynamics"
	"github.com/FelixDux/imposcg/dynamics/classifyorbits"
	"github.com/FelixDux/imposcg/dynamics/impact"
	"github.com/FelixDux/imposcg/dynamics/parameters"
	"github.com/gin-gonic/gin"
)

func GenerateDOAData(parameters *parameters.Parameters, inputs *DOAInputs) (*map[string][]impact.SimpleImpact, error) {
	impactMap, err := dynamics.NewImpactMap(*parameters)

	if err != nil {
		return nil, err
	}

	classifier := classifyorbits.NewOrbitClassifier(impactMap, inputs.numIterations)

	classifications := classifier.BuildClassification(inputs.numPhases, inputs.numVelocities, inputs.maxVelocity)

	return classifyorbits.MarshalClassifications(&classifications), nil
}

// PostDOAData godoc
// @Summary Domain of attraction data
// @Description Return domains of attraction for the impact map for a specified set of parameters
// @ID post-doa-data
// @Accept  x-www-form-urlencoded
// @Produce  json
// @Param frequency formData number true "Forcing frequency" minimum(0)
// @Param offset formData number true "Obstacle offset from origin"
// @Param r formData number true "Coefficient of restitution" minimum(0) maximum(1)
// @Param maxPeriods formData int false "Number of periods without an impact after which the algorithm will report 'long excursions'" default(100)
// @Param numIterations formData int false "Number of iterations of impact map" default(500)
// @Param maxVelocity formData number true "Upper limit of impact velocity range for DOA plot" default(4.0)
// @Param numPhases formData int false "Size of grid along the φ-axis" default(100)
// @Param numVelocities formData int false "Size of grid along the v-axis" default(100)
// @Success 200 {object} dynamics.IterationResult
// @Failure 400 {object} string "Invalid parameters"
// @Router /doa/data/ [post]
func PostDOAData(c *gin.Context) {
	inputs, parameters, errorString := DOAInputsFromPost(c)

	if parameters == nil || inputs == nil {
        log.Print(errorString)
		c.JSON(400, errorString)

		return
	}
	
	data, err := GenerateDOAData(parameters, inputs)

	if err != nil {
		c.JSON(400,  err.Error())
		return
	}

	c.JSON(200,  data)
}

// PostDOAImage godoc
// @Summary Domain of attraction plot
// @Description Plot domains of attraction for the impact map for a specified set of parameters
// @ID post-doa-image
// @Accept  x-www-form-urlencoded
// @Produce  png
// @Param frequency formData number true "Forcing frequency" minimum(0)
// @Param offset formData number true "Obstacle offset from origin"
// @Param r formData number true "Coefficient of restitution" minimum(0) maximum(1)
// @Param maxPeriods formData int false "Number of periods without an impact after which the algorithm will report 'long excursions'" default(100)
// @Param numIterations formData int false "Number of iterations of impact map" default(500)
// @Param maxVelocity formData number true "Upper limit of impact velocity range for DOA plot" default(4.0)
// @Param numPhases formData int false "Size of grid along the φ-axis" default(100)
// @Param numVelocities formData int false "Size of grid along the v-axis" default(100)
// @Success 200 {object} dynamics.IterationResult
// @Failure 400 {object} string "Invalid parameters"
// @Router /doa/image/ [post]
func PostDOAImage(c *gin.Context) {
	inputs, parameters, errorString := DOAInputsFromPost(c)

	if parameters == nil || inputs == nil {
        log.Print(errorString)
		c.JSON(400, errorString)

		return
	}
	
	data, err := GenerateDOAData(parameters, inputs)

	if err != nil {
		c.JSON(400,  err.Error())
		return
	}

	plot := charts.DOAPlot(*parameters, data)

	img, err := os.Open(plot.Name())
	if err != nil {
		log.Print(err)
		c.JSON(400, fmt.Sprintf("Failed to complete DAO plot - %s", err.Error()))
	} else {
		defer img.Close()
		c.Writer.Header().Set("Content-Type", "image/png")
		io.Copy(c.Writer, img)
	}

	c.JSON(200,  data)
}

func AddDOAControllers (r *gin.Engine) {
	doa := r.Group("/api/doa")
	doa.POST("/data",  PostDOAData)
	doa.POST("/image", PostDOAImage)
}