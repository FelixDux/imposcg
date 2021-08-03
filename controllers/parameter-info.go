package controllers

// https://github.com/swaggo/gin-swagger to generate REST spec
// https://github.com/swaggo/swag#declarative-comments-format

import (
	"github.com/gin-gonic/gin"
)

type ParameterProperty struct {
	Parameter string
	Property string
}

type ParameterInfo struct {
	Symbols []ParameterProperty
}

func (info ParameterInfo) Add(parameter string, property string) *ParameterInfo {
	info.Symbols = append(info.Symbols, ParameterProperty{Parameter: parameter, Property: property})

	return &info
}

// GetParameterSymbols godoc
// @Summary Parameter symbols
// @Description Greek symbols to be used for rendering specified parameters
// @ID get-parameter-symbols
// @Produce  json
// @Success 200 {object} controllers.ParameterInfo
// @Router /parameter-info/symbols/ [get]
func GetParameterSymbols(c *gin.Context) {
	info := &ParameterInfo{}

	symbols := map[string]string{
		"frequency": "ω",
		"offset": "σ",
		"phi": "φ",}

	for k, v := range symbols {
		info = info.Add(k, v)
	}

	c.JSON(200, info)
}

// GetParameterGroups godoc
// @Summary Parameter groups
// @Description Groups for displaying related parameters
// @ID get-parameter-groups
// @Produce  json
// @Success 200 {object} controllers.ParameterInfo
// @Router /parameter-info/groups/ [get]
func GetParameterGroups(c *gin.Context) {
	info := &ParameterInfo{}

	const sysParams = "System parameters"
	const ctlParams = "Control parameters"
	const initialImpact = "Initial impact"

	groups := map[string]string{
		"frequency": sysParams,
		"offset": sysParams,
		"r": sysParams,
		"phi": initialImpact,
		"maxPeriods": ctlParams,
		"numIterations": ctlParams,
		"numPoints": ctlParams,
	}

	for k, v := range groups {
		info = info.Add(k, v)
	}

	c.JSON(200, info)
}

func AddParameterInfoControllers(r *gin.Engine) {
	parameterInfo := r.Group("/api/parameter-info")
	parameterInfo.GET("/symbols", GetParameterSymbols)
}