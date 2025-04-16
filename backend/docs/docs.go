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
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/admin": {
            "get": {
                "security": [
                    {
                        "JWT": []
                    }
                ],
                "description": "Whether user has admin rights",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "Check admin rights available",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/internal_api_get_admin.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    }
                }
            }
        },
        "/auth/login": {
            "post": {
                "description": "Generate ` + "`" + `access` + "`" + ` and ` + "`" + `refresh` + "`" + ` token pair. ` + "`" + `refreshToken` + "`" + ` sets in httpOnly cookie also.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Login by email and password",
                "parameters": [
                    {
                        "description": "existing user credentials",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_api_post_auth_login.Request"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/internal_api_post_auth_login.Response"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    }
                }
            }
        },
        "/auth/logout": {
            "post": {
                "description": "Remove ` + "`" + `refreshSession` + "`" + ` attached to ` + "`" + `refreshToken` + "`" + `. ` + "`" + `refreshToken` + "`" + ` can be passed in cookie",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Logout",
                "parameters": [
                    {
                        "description": "active refresh token in UUID RFC4122 format",
                        "name": "refreshToken",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/internal_api_post_auth_logout.Request"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    }
                }
            }
        },
        "/auth/refresh": {
            "post": {
                "description": "Refresh ` + "`" + `access` + "`" + ` and ` + "`" + `refresh` + "`" + ` token pair. ` + "`" + `refreshToken` + "`" + ` can be passed in cookie",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Refresh tokens",
                "parameters": [
                    {
                        "description": "active refresh token in UUID RFC4122 format",
                        "name": "refreshToken",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/internal_api_post_auth_refresh.Request"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/internal_api_post_auth_refresh.Response"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    }
                }
            }
        },
        "/balance/deposit": {
            "post": {
                "security": [
                    {
                        "JWT": []
                    }
                ],
                "description": "Add specified amount to authenticated user's balance",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "balance"
                ],
                "summary": "Deposit funds",
                "parameters": [
                    {
                        "description": "Deposit amount",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_api_post_balance_deposit.Request"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/internal_api_post_balance_deposit.Response"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    }
                }
            }
        },
        "/balance/withdraw": {
            "post": {
                "security": [
                    {
                        "JWT": []
                    }
                ],
                "description": "Subtract specified amount from authenticated user's balance",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "balance"
                ],
                "summary": "Withdraw funds",
                "parameters": [
                    {
                        "description": "Withdrawal amount",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_api_post_balance_withdraw.Request"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/internal_api_post_balance_withdraw.Response"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    }
                }
            }
        },
        "/health": {
            "get": {
                "description": "Whether REST-API alive or not",
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "health"
                ],
                "summary": "Check health",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    }
                }
            }
        },
        "/orders": {
            "post": {
                "security": [
                    {
                        "JWT": []
                    }
                ],
                "description": "Create order on behalf of authenticated user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "orders"
                ],
                "summary": "Create order",
                "parameters": [
                    {
                        "description": "order parameters",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_api_post_orders.Request"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/internal_api_post_orders.Response"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "echo.HTTPError": {
            "type": "object",
            "properties": {
                "message": {}
            }
        },
        "github_com_moevm_nosql1h25-writer_backend_internal_entity.SystemRoleType": {
            "type": "string",
            "enum": [
                "admin",
                "user"
            ],
            "x-enum-varnames": [
                "SystemRoleTypeAdmin",
                "SystemRoleTypeUser"
            ]
        },
        "internal_api_get_admin.Response": {
            "type": "object",
            "required": [
                "systemRole",
                "userId"
            ],
            "properties": {
                "systemRole": {
                    "allOf": [
                        {
                            "$ref": "#/definitions/github_com_moevm_nosql1h25-writer_backend_internal_entity.SystemRoleType"
                        }
                    ],
                    "example": "admin"
                },
                "userId": {
                    "type": "string",
                    "example": "5a2493c33c95a1281836eb6a"
                }
            }
        },
        "internal_api_post_auth_login.Request": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "format": "email",
                    "example": "test@gmail.com"
                },
                "password": {
                    "type": "string",
                    "maxLength": 72,
                    "minLength": 8,
                    "example": "Password123"
                }
            }
        },
        "internal_api_post_auth_login.Response": {
            "type": "object",
            "required": [
                "accessToken",
                "refreshToken"
            ],
            "properties": {
                "accessToken": {
                    "type": "string",
                    "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
                },
                "refreshToken": {
                    "type": "string",
                    "example": "289abe45-5920-4366-a12a-875ddb422ace"
                }
            }
        },
        "internal_api_post_auth_logout.Request": {
            "type": "object",
            "properties": {
                "refreshToken": {
                    "type": "string",
                    "example": "0e8f711e-b713-4869-b528-059a74311482"
                }
            }
        },
        "internal_api_post_auth_refresh.Request": {
            "type": "object",
            "properties": {
                "refreshToken": {
                    "type": "string",
                    "example": "0e8f711e-b713-4869-b528-059a74311482"
                }
            }
        },
        "internal_api_post_auth_refresh.Response": {
            "type": "object",
            "required": [
                "accessToken",
                "refreshToken"
            ],
            "properties": {
                "accessToken": {
                    "type": "string",
                    "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
                },
                "refreshToken": {
                    "type": "string",
                    "example": "289abe45-5920-4366-a12a-875ddb422ace"
                }
            }
        },
        "internal_api_post_balance_deposit.Request": {
            "type": "object",
            "required": [
                "amount"
            ],
            "properties": {
                "amount": {
                    "type": "integer",
                    "minimum": 1,
                    "example": 100
                }
            }
        },
        "internal_api_post_balance_deposit.Response": {
            "type": "object",
            "properties": {
                "newBalance": {
                    "type": "integer",
                    "example": 777
                }
            }
        },
        "internal_api_post_balance_withdraw.Request": {
            "type": "object",
            "required": [
                "amount"
            ],
            "properties": {
                "amount": {
                    "type": "integer",
                    "minimum": 1,
                    "example": 100
                }
            }
        },
        "internal_api_post_balance_withdraw.Response": {
            "type": "object",
            "properties": {
                "newBalance": {
                    "type": "integer",
                    "example": 111
                }
            }
        },
        "internal_api_post_orders.Request": {
            "type": "object",
            "required": [
                "comletionTime",
                "description",
                "title"
            ],
            "properties": {
                "comletionTime": {
                    "type": "integer",
                    "minimum": 3600000000000,
                    "example": 3600000000000
                },
                "cost": {
                    "type": "integer",
                    "minimum": 0,
                    "example": 500
                },
                "description": {
                    "type": "string",
                    "maxLength": 8192,
                    "minLength": 16,
                    "example": "Написать сценарий вот такой и такой"
                },
                "title": {
                    "type": "string",
                    "maxLength": 32,
                    "minLength": 3,
                    "example": "Сценарий"
                }
            }
        },
        "internal_api_post_orders.Response": {
            "type": "object",
            "required": [
                "id"
            ],
            "properties": {
                "id": {
                    "type": "string",
                    "example": "522bb79455449d881b004d27"
                }
            }
        }
    },
    "securityDefinitions": {
        "JWT": {
            "description": "JSON Web Token",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0.0",
	Host:             "localhost:80",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "Writer",
	Description:      "API for freelancer's site",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
