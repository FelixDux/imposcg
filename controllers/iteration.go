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

)


func doScatter() string {

	params, errParams := parameters.NewParameters(2.8, 0.0, 0.8, 100)

	if len(errParams) > 0 {
	
		paramMessages := make([]string,len(errParams))
		for i, err := range(errParams) {

			paramMessages[i] = err.Error()

		}

		return strings.Join(paramMessages, "\n")
	}
	
	impactMap, errMap := dynamics.NewImpactMap(*params)

	if errMap != nil {
		return errMap.Error()
	}

	phi := 0.0
	v := 0.0
	data := impactMap.IterateFromPoint(phi, v, 1000)

	return charts.ImpactMapPlot(*params, data.Impacts, phi, v).Name()
}

// POST /api/iteration/image
func HandleScatter(c *gin.Context) {

	img, err := os.Open(doScatter())
    if err != nil {
        log.Fatal(err) // perhaps handle this nicer
    }
    defer img.Close()
    c.Writer.Header().Set("Content-Type", "image/png") // <-- set the content-type header
    io.Copy(c.Writer, img)
}

// POST /api/iteration/data
func HandleImpactMapData(c *gin.Context) {

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