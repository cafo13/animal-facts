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
        "contact": {
            "name": "Animal Facts API",
            "url": "https://animalfacts.app"
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
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/router.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Adding an animal fact",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "facts"
                ],
                "summary": "Add a new animal fact",
                "parameters": [
                    {
                        "description": "a new fact",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.Fact"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/types.Fact"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/router.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/router.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/fact/:id": {
            "get": {
                "description": "Getting an animal fact by id",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "facts"
                ],
                "summary": "Get animal fact by id",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/types.Fact"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/router.ErrorResponse"
                        }
                    }
                }
            },
            "put": {
                "description": "Updating an animal fact",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "facts"
                ],
                "summary": "Update an existing animal fact",
                "parameters": [
                    {
                        "description": "an updated fact",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/types.Fact"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/router.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/router.MessageResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/router.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Deleting an animal fact",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "facts"
                ],
                "summary": "Delete an animal fact",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/router.MessageResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/router.ErrorResponse"
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
        "router.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {}
            }
        },
        "router.MessageResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "types.Fact": {
            "type": "object",
            "properties": {
                "approved": {
                    "type": "boolean"
                },
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
	Version:          "0.0.1",
	Host:             "https://animalfacts.app",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Animal Facts API",
	Description:      "Awesome facts about animals.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
