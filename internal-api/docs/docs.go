// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {
            "name": "MIT",
            "url": "https://github.com/cafo13/animal-facts/blob/main/LICENSE"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/facts": {
            "post": {
                "description": "create a new fact",
                "produces": [
                    "application/json"
                ],
                "summary": "create fact",
                "parameters": [
                    {
                        "description": "fact",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.CreateUpdateFact"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/api.CreateFactResult"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResult"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResult"
                        }
                    }
                }
            }
        },
        "/facts/:id": {
            "put": {
                "description": "update an existing fact",
                "produces": [
                    "application/json"
                ],
                "summary": "update fact",
                "parameters": [
                    {
                        "description": "fact",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.CreateUpdateFact"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "fact updated",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResult"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResult"
                        }
                    }
                }
            },
            "delete": {
                "description": "delete an existing fact",
                "produces": [
                    "application/json"
                ],
                "summary": "delete fact",
                "responses": {
                    "200": {
                        "description": "fact deleted",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResult"
                        }
                    }
                }
            }
        },
        "/facts/:id/approve": {
            "post": {
                "description": "approve an existing fact, so that it gets available in the public API",
                "produces": [
                    "application/json"
                ],
                "summary": "approve fact",
                "responses": {
                    "200": {
                        "description": "fact approved",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResult"
                        }
                    }
                }
            }
        },
        "/facts/:id/unapprove": {
            "post": {
                "description": "unapprove an existing fact, so that it is no longer available in the public API",
                "produces": [
                    "application/json"
                ],
                "summary": "unapprove fact",
                "responses": {
                    "200": {
                        "description": "fact unapproved",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResult"
                        }
                    }
                }
            }
        },
        "/facts/all": {
            "get": {
                "description": "gets all facts (approved and unapproved) from the database",
                "produces": [
                    "application/json"
                ],
                "summary": "gets all facts",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "array",
                                "items": {
                                    "$ref": "#/definitions/repository.Fact"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResult"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.CreateFactResult": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                }
            }
        },
        "api.CreateUpdateFact": {
            "type": "object",
            "properties": {
                "fact": {
                    "type": "string"
                },
                "source": {
                    "type": "string"
                }
            }
        },
        "api.ErrorResult": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "repository.Fact": {
            "type": "object",
            "properties": {
                "approved": {
                    "type": "boolean"
                },
                "createdAt": {
                    "type": "string"
                },
                "createdBy": {
                    "type": "string"
                },
                "fact": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "source": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                },
                "updatedBy": {
                    "type": "string"
                }
            }
        }
    },
    "externalDocs": {
        "description": "OpenAPI",
        "url": "https://swagger.io/resources/open-api/"
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.0.4",
	Host:             "https://animal-facts-internal.cafo.dev",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Animal Facts Internal API",
	Description:      "This API provides facts about animals.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
