package main

import (
	"github.com/FelixDux/imposcg/controllers"
	"github.com/FelixDux/imposcg/docs"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/static"

	"context"
	"log"
	"os"
	"strings"

	_ "github.com/FelixDux/imposcg/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
)

// https://github.com/swaggo/swag/issues/367
// https://github.com/swaggo/swag#how-to-use-it-with-gin

// @title Impact Oscillator
// @version 1.0
// @description Analysis and simulation of a simple vibro-impact model developed in Go - principally as a learning exercise
// @BasePath /api

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

func EnvOrDefault(env string, defaultValue string) string {
	value := strings.TrimSpace(os.Getenv(env))

	if len(value) > 0 {
		return value
	} else {
		return defaultValue
	}
}

func SwagHost() string {
	return EnvOrDefault("SWAG_HOST", "")
}

func LogEnv(env string) {
	log.Printf("%s = %s", env, os.Getenv(env))
}

func OnAwsLambda() bool {
	return EnvIsTruthy("_LAMBDA_SERVER_PORT")
}

func setupServer() *gin.Engine {

	log.Printf("Gin cold start")

	r := gin.Default()

	controllers.AddIterationControllers(r)
	controllers.AddSingularitySetControllers(r)
	controllers.AddDOAControllers(r)
	controllers.AddParameterInfoControllers(r)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Serve frontend static files
	r.Use(static.Serve("/", static.LocalFile("./static", true)))
	r.Static("/js/", "./static/js/")

	return r
}

func SetupSwagInfo() {
	docs.SwaggerInfo.Host = SwagHost()
}

func init() {
	SetupSwagInfo()
	
	server = setupServer()

	if OnAwsLambda() {
		log.Printf("Creating lambda adapter")
		ginLambda = ginadapter.New(server)
	} else {
		LogEnv("_AWS_LAMBDA_RUNTIME_API")
		LogEnv("_LAMBDA_SERVER_PORT")
	}
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	log.Printf("Handling lambda request: %v+", req)
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	if OnAwsLambda() {
		log.Printf("Starting lambda adapter")
		lambda.Start(Handler)
	} else {
		log.Printf("Starting Gin server")
		server.Run()
	}
}