// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "https://animalfact.app/terms",
        "contact": {
            "name": "Animal Facts API",
            "url": "https://animalfacts.app/support",
            "email": "support@animalfacts.app"
        },
        "license": {
            "name": "MIT",
            "url": "https://github.com/cafo13/animal-facts/blob/main/LICENSE"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/fact": {
            "get": {
                "description": "Getting a random animal fact",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "facts"
                ],
                "summary": "Get random animal fact",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/types.Fact"
                        }
                    }
                }
            }
        },
        "/fact/:id": {
            "get": {
                "description": "Getting an animal fact by ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "facts"
                ],
                "summary": "Get animal fact by ID",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/types.Fact"
                        }
                    }
                }
            }
        },
        "/healthy": {
            "get": {
                "description": "Checking the health of the API",
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "general"
                ],
                "summary": "Get health status",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "types.Fact": {
            "type": "object",
            "properties": {
                "category": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "source": {
                    "type": "string"
                },
                "text": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "https://animalfacts.app",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Animal Facts API",
	Description:      "Get awesome facts about animals.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
