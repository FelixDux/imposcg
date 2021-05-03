package main

import (
	"github.com/FelixDux/imposcg/controllers"
	"github.com/gin-gonic/gin"

)

// Basic structure:
// / - SPA
// /api - REST schema
// /api/iteration/data
// /api/iteration/image
// /api/singularity-set/data
// /api/singularity-set/image
// /api/doa/data
// /api/doa/image
// /api/offset-response/data
// /api/offset-response/image
// /api/frequency-response/data
// /api/frequency-response/image

func main() {
	r := gin.Default()
	r.GET("/api/iteration/data",  controllers.HandleImpactMapData)
	r.GET("/api/iteration/image",  controllers.HandleScatter)
	r.Run() // listen and serve on 0.0.0.0:8080
}