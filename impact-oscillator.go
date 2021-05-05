package main

import (
	"github.com/FelixDux/imposcg/controllers"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/FelixDux/imposcg/docs"
)

// @title Impact Oscillator API
// @version 1.0
// @description Analysis and simulation of a simple vibro-impact model developed in Go - principally as a learning exercise
// @host localhost:8080
// @BasePath /api

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

func setupServer() *gin.Engine {

	r := gin.Default()

	controllers.AddIterationControllers(r)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

func main() {
	setupServer().Run() // listen and serve on 0.0.0.0:8080
}