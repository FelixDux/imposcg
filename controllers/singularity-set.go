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

func singularitySetData(parameters *parameters.Parameters, numPoints uint) ([]impact.Impact, []impact.Impact, string) {
	impactMap, errMap := dynamics.NewImpactMap(*parameters)

	if errMap != nil {
		return nil, nil, errMap.Error()
	}

	singularity, dual := impactMap.SingularitySet(numPoints)
	
	return singularity, dual, ""
}

func singularitySetImage(parameters *parameters.Parameters, numPoints uint) string {
	singularity, _, errString := singularitySetData(parameters, numPoints)

	if singularity == nil {
		return errString
	} else {
		return charts.ImpactMapPlot(*parameters, singularity, 0.0, 0.0).Name()
	}
}

// PostSingularitySetImage godoc
// @Summary Return scatter plot of impacts which map to zero velocity impacts for a specified set of parameters
// @Description Return scatter plot of impacts which map to zero velocity impacts for a specified set of parameters
// @ID post-singularity-set-image
// @Accept x-www-form-urlencoded
// @Produce  png
// @Param frequency formData number true "Forcing frequency" minimum(0)
// @Param offset formData number true "Obstacle offset from origin"
// @Param r formData number true "Coefficient of restitution" minimum(0) maximum(1)
// @Param maxPeriods formData int false "Number of periods without an impact after which the algorithm will report 'long excursions'" default(100)
// @Param numPoints formData int false "Number of iterations of impact map" default(10000)
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

func AddSingularitySetControllers (r *gin.Engine) {
	iteration := r.Group("/api/singularity-set")
	// iteration.POST("/data",  PostSingularitySetData)
	iteration.POST("/image",  PostSingularitySetImage)
}