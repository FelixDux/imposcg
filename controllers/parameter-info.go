package controllers

// https://github.com/swaggo/gin-swagger to generate REST spec
// https://github.com/swaggo/swag#declarative-comments-format

import (
	"github.com/gin-gonic/gin"
)

type ParameterSymbol struct {
	Parameter string
	Symbol string
}

type ParameterInfo struct {
	Symbols []ParameterSymbol
}

func (info ParameterInfo) Add(parameter string, symbol string) *ParameterInfo {
	info.Symbols = append(info.Symbols, ParameterSymbol{Parameter: parameter, Symbol: symbol})

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

	c.JSON(200, gin.H{
		"message": info,
	})
}

func AddParameterInfoControllers(r *gin.Engine) {
	parameterInfo := r.Group("/api/parameter-info")
	parameterInfo.GET("/symbols", GetParameterSymbols)
}