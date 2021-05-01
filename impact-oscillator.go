package main

import (
	// "fmt"
	// "os"
	// "log"
	"strings"

	"github.com/FelixDux/imposcg/dynamics"
	"github.com/FelixDux/imposcg/dynamics/parameters"
	"github.com/FelixDux/imposcg/charts"
	"github.com/gin-gonic/gin"
	// "net/http"

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

func HandleScatter(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": doScatter(),
	})

	// //w http.ResponseWriter, r *http.Request
	// img, err := os.Open("example_t500.png")
    // if err != nil {
    //     log.Fatal(err) // perhaps handle this nicer
    // }
    // defer img.Close()
    // w.Header().Set("Content-Type", "image/png") // <-- set the content-type header
    // io.Copy(w, img)
}

func main() {
	r := gin.Default()
	r.GET("",  HandleScatter)
	r.Run() // listen and serve on 0.0.0.0:8080
}