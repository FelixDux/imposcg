package main

import (
	"github.com/FelixDux/imposcg/controllers"
	"github.com/gin-gonic/gin"

	"os"
	"strings"
	"log"
	"context"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/FelixDux/imposcg/docs"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
)

// @title Impact Oscillator API
// @version 1.0
// @description Analysis and simulation of a simple vibro-impact model developed in Go - principally as a learning exercise
// @host localhost:8080
// @BasePath /api

// Basic structure:
// / - SPA
// /swagger/*any - REST schema
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

var server *gin.Engine
var ginLambda *ginadapter.GinLambda

func EnvIsTruthy(env string) bool {
	value := strings.TrimSpace(os.Getenv(env))

	if len(value) > 0 {
		return true
	} else {
		return false
	}
}

func OnAwsLambda() bool {
	return EnvIsTruthy("AWS_LAMBDA_RUNTIME_API") && EnvIsTruthy("_LAMBDA_SERVER_PORT")
}

func setupServer() *gin.Engine {

	log.Printf("Gin cold start")

	r := gin.Default()

	controllers.AddIterationControllers(r)
	controllers.AddSingularitySetControllers(r)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

func init() {
	server = setupServer()

	if OnAwsLambda() {
		ginLambda = ginadapter.New(server)
	}
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	if OnAwsLambda() {
		lambda.Start(Handler)
	} else {
		server.Run()
	}
}