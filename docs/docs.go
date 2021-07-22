// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/iteration/data/": {
            "post": {
                "description": "Return data from iterating the impact map for a specified set of parameters",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Return data from iterating the impact map",
                "operationId": "post-iteration-data",
                "parameters": [
                    {
                        "minimum": 0,
                        "type": "number",
                        "description": "Forcing frequency",
                        "name": "frequency",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "number",
                        "description": "Obstacle offset from origin",
                        "name": "offset",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "maximum": 1,
                        "minimum": 0,
                        "type": "number",
                        "description": "Coefficient of restitution",
                        "name": "r",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "default": 100,
                        "description": "Number of periods without an impact after which the algorithm will report 'long excursions'",
                        "name": "maxPeriods",
                        "in": "formData"
                    },
                    {
                        "type": "number",
                        "default": 0,
                        "description": "Phase at initial impact",
                        "name": "phi",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "number",
                        "default": 0,
                        "description": "Velocity at initial impact",
                        "name": "v",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "default": 5000,
                        "description": "Number of iterations of impact map",
                        "name": "numIterations",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dynamics.IterationResult"
                        }
                    },
                    "400": {
                        "description": "Invalid parameters",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/iteration/image/": {
            "post": {
                "description": "Return scatter plot from iterating the impact map for a specified set of parameters",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "image/png"
                ],
                "summary": "Return scatter plot from iterating the impact map",
                "operationId": "post-iteration-image",
                "parameters": [
                    {
                        "minimum": 0,
                        "type": "number",
                        "description": "Forcing frequency",
                        "name": "frequency",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "number",
                        "description": "Obstacle offset from origin",
                        "name": "offset",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "maximum": 1,
                        "minimum": 0,
                        "type": "number",
                        "description": "Coefficient of restitution",
                        "name": "r",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "default": 100,
                        "description": "Number of periods without an impact after which the algorithm will report 'long excursions'",
                        "name": "maxPeriods",
                        "in": "formData"
                    },
                    {
                        "type": "number",
                        "default": 0,
                        "description": "Phase at initial impact",
                        "name": "phi",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "number",
                        "default": 0,
                        "description": "Velocity at initial impact",
                        "name": "v",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "default": 5000,
                        "description": "Number of iterations of impact map",
                        "name": "numIterations",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dynamics.IterationResult"
                        }
                    },
                    "400": {
                        "description": "Invalid parameters",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/singularity-set/data/": {
            "post": {
                "description": "Return impacts which map to and from zero velocity impacts for a specified set of parameters",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Return impacts which map to and from zero velocity impacts for a specified set of parameters",
                "operationId": "post-singularity-set-data",
                "parameters": [
                    {
                        "minimum": 0,
                        "type": "number",
                        "description": "Forcing frequency",
                        "name": "frequency",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "number",
                        "description": "Obstacle offset from origin",
                        "name": "offset",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "maximum": 1,
                        "minimum": 0,
                        "type": "number",
                        "description": "Coefficient of restitution",
                        "name": "r",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "default": 100,
                        "description": "Number of periods without an impact after which the algorithm will report 'long excursions'",
                        "name": "maxPeriods",
                        "in": "formData"
                    },
                    {
                        "type": "integer",
                        "default": 5000,
                        "description": "Number of impacts to map",
                        "name": "numPoints",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dynamics.IterationResult"
                        }
                    },
                    "400": {
                        "description": "Invalid parameters",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/singularity-set/image/": {
            "post": {
                "description": "Return scatter plot of impacts which map to and from zero velocity impacts for a specified set of parameters",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "image/png"
                ],
                "summary": "Return scatter plot of impacts which map to and from zero velocity impacts for a specified set of parameters",
                "operationId": "post-singularity-set-image",
                "parameters": [
                    {
                        "minimum": 0,
                        "type": "number",
                        "description": "Forcing frequency",
                        "name": "frequency",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "number",
                        "description": "Obstacle offset from origin",
                        "name": "offset",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "maximum": 1,
                        "minimum": 0,
                        "type": "number",
                        "description": "Coefficient of restitution",
                        "name": "r",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "default": 100,
                        "description": "Number of periods without an impact after which the algorithm will report 'long excursions'",
                        "name": "maxPeriods",
                        "in": "formData"
                    },
                    {
                        "type": "integer",
                        "default": 5000,
                        "description": "Number of impacts to map",
                        "name": "numPoints",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dynamics.IterationResult"
                        }
                    },
                    "400": {
                        "description": "Invalid parameters",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dynamics.IterationResult": {
            "type": "object",
            "properties": {
                "impacts": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/impact.Impact"
                    }
                }
            }
        },
        "impact.Impact": {
            "type": "object",
            "properties": {
                "phase": {
                    "type": "number"
                },
                "time": {
                    "type": "number"
                },
                "velocity": {
                    "type": "number"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "",
	BasePath:    "/api",
	Schemes:     []string{},
	Title:       "Impact Oscillator API",
	Description: "Analysis and simulation of a simple vibro-impact model developed in Go - principally as a learning exercise",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
