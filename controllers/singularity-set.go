package controllers


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

type SingularitySetResult struct {
	Singularity []impact.Impact `json:"singularity"`
	Dual []impact.Impact `json:"dual"`
}

func singularitySetData(parameters *parameters.Parameters, numPoints uint) (*SingularitySetResult, string) {
	impactMap, errMap := dynamics.NewImpactMap(*parameters)

	if errMap != nil {
		return &SingularitySetResult{Singularity: nil, Dual: nil}, errMap.Error()
	}

	singularity, dual := impactMap.SingularitySet(numPoints)
	
	return &SingularitySetResult{Singularity: singularity, Dual: dual}, ""
}

func singularitySetImage(parameters *parameters.Parameters, numPoints uint) string {
	result, errString := singularitySetData(parameters, numPoints)

	if result.Singularity == nil || result.Dual == nil {
		return errString
	} else {
		return charts.ImpactMapPlot(*parameters, [][]impact.Impact{result.Singularity,result.Dual}, 0.0, 0.0).Name()
	}
}

// PostSingularitySetImage godoc
// @Summary Return scatter plot of impacts which map to and from zero velocity impacts for a specified set of parameters
// @Description Return scatter plot of impacts which map to and from zero velocity impacts for a specified set of parameters
// @ID post-singularity-set-image
// @Accept x-www-form-urlencoded
// @Produce  png
// @Param frequency formData number true "Forcing frequency" minimum(0)
// @Param offset formData number true "Obstacle offset from origin"
// @Param r formData number true "Coefficient of restitution" minimum(0) maximum(1)
// @Param maxPeriods formData int false "Number of periods without an impact after which the algorithm will report 'long excursions'" default(100)
// @Param numPoints formData int false "Number of impacts to map" default(5000)
// @Success 200 {object} dynamics.IterationResult
// @Failure 400 {object} string "Invalid parameters"
// @Router /singularity-set/image/ [post]
func PostSingularitySetImage(c *gin.Context) {
	numPoints, parameters, errorString := SingularitySetInputsFromPost(c)

	if parameters == nil || len(errorString) > 0 {
        log.Print(errorString)
		c.JSON(400, errorString)
	} else {
		img, err := os.Open(singularitySetImage(parameters, numPoints))
		if err != nil {
			log.Print(err)
			c.JSON(400, fmt.Sprintf("Failed to complete singularity set - %s", err.Error()))
		} else {
			defer img.Close()
			c.Writer.Header().Set("Content-Type", "image/png")
			io.Copy(c.Writer, img)
		}
	}
}

// PostSingularitySetData godoc
// @Summary Return impacts which map to and from zero velocity impacts for a specified set of parameters
// @Description Return impacts which map to and from zero velocity impacts for a specified set of parameters
// @ID post-singularity-set-data
// @Accept x-www-form-urlencoded
// @Produce json
// @Param frequency formData number true "Forcing frequency" minimum(0)
// @Param offset formData number true "Obstacle offset from origin"
// @Param r formData number true "Coefficient of restitution" minimum(0) maximum(1)
// @Param maxPeriods formData int false "Number of periods without an impact after which the algorithm will report 'long excursions'" default(100)
// @Param numPoints formData int false "Number of impacts to map" default(5000)
// @Success 200 {object} dynamics.IterationResult
// @Failure 400 {object} string "Invalid parameters"
// @Router /singularity-set/data/ [post]
func PostSingularitySetData(c *gin.Context) {
	numPoints, parameters, errorString := SingularitySetInputsFromPost(c)

	if parameters == nil || len(errorString) > 0 {
        log.Print(errorString)
		c.JSON(400, errorString)
	} else {
		result, errString := singularitySetData(parameters, numPoints)

		if result.Singularity == nil || result.Dual == nil {
			log.Print(errString)
			c.JSON(400, fmt.Sprintf("Failed to complete singularity set - %s", errString))
		} else {
			c.JSON(200, gin.H{"message": result,})
		}
	}
}

func AddSingularitySetControllers (r *gin.Engine) {
	iteration := r.Group("/api/singularity-set")
	iteration.POST("/data",  PostSingularitySetData)
	iteration.POST("/image",  PostSingularitySetImage)
}