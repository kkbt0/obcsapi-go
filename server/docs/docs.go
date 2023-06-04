// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "url": "https://gitee.com/kkbt/obcsapi-go"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/info": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "服务器信息与测试接口",
                "responses": {}
            }
        },
        "/ob/fv": {
            "post": {
                "security": [
                    {
                        "Token": []
                    }
                ],
                "description": "安卓软件 fv 悬浮球使用的 API 用于自定义任务的 图片、文字",
                "consumes": [
                    "text/plain",
                    "application/octet-stream"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Ob"
                ],
                "summary": "fv 悬浮球使用的 API",
                "parameters": [
                    {
                        "description": "fv payload 內容",
                        "name": "內容",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/ob/general": {
            "post": {
                "security": [
                    {
                        "Token": []
                    }
                ],
                "description": "通用 API 接口,添加 Memos",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Ob"
                ],
                "summary": "通用 API 接口 Memos",
                "parameters": [
                    {
                        "description": "MemosData",
                        "name": "json",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dao.MemosData"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/ob/general/{token}": {
            "post": {
                "description": "通用 API 接口,添加 Memos",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Ob"
                ],
                "summary": "通用 API 接口 (Memos Flomo Like API)",
                "parameters": [
                    {
                        "type": "string",
                        "description": "设定的 token 值",
                        "name": "token",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "MemosData",
                        "name": "json",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dao.MemosData"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/ob/generalall": {
            "get": {
                "security": [
                    {
                        "Token": []
                    }
                ],
                "description": "通用 API 接口，获取所有文件。需要配置声明允许使用该接口",
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Ob"
                ],
                "summary": "通用 API 接口 All",
                "parameters": [
                    {
                        "type": "string",
                        "description": "文件名，有路径，如 dir/text.md",
                        "name": "filekey",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            },
            "post": {
                "security": [
                    {
                        "Token": []
                    }
                ],
                "description": "通用 API 接口，覆盖修改或增添所有文件。需要配置声明允许使用该接口",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Ob"
                ],
                "summary": "通用 API 接口 All",
                "parameters": [
                    {
                        "description": "GeneralAllStruct",
                        "name": "json",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.GeneralAllStruct"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/ob/moonreader": {
            "post": {
                "security": [
                    {
                        "AuthorizationToken": []
                    }
                ],
                "description": "静读天下使用的 API，标注-设置-ReadWise 设置该路径和 token 值即可",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Ob"
                ],
                "summary": "静读天下使用的 API",
                "parameters": [
                    {
                        "description": "MoodReader",
                        "name": "划线和标注",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dao.MoodReader"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/ob/sr/webhook": {
            "post": {
                "security": [
                    {
                        "Token": []
                    }
                ],
                "description": "SimpRead 简悦 WebHook POST 简悦 WebHook 保存文章",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Ob"
                ],
                "summary": "简悦 WebHook 保存文章",
                "parameters": [
                    {
                        "description": "SimpRead 简悦 POST",
                        "name": "json",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.SimpReadWebHookStruct"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/ob/url": {
            "post": {
                "security": [
                    {
                        "Token": []
                    }
                ],
                "description": "裁剪网页",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Ob"
                ],
                "summary": "裁剪网页",
                "parameters": [
                    {
                        "description": "MemosData",
                        "name": "json",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dao.UrlStruct"
                        }
                    }
                ],
                "responses": {}
            }
        }
    },
    "definitions": {
        "dao.MemosData": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                }
            }
        },
        "dao.MoodReader": {
            "type": "object",
            "properties": {
                "highlights": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dao.MoodReaderHighlights"
                    }
                }
            }
        },
        "dao.MoodReaderHighlights": {
            "type": "object",
            "properties": {
                "author": {
                    "type": "string"
                },
                "note": {
                    "type": "string"
                },
                "text": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "dao.UrlStruct": {
            "type": "object",
            "properties": {
                "url": {
                    "type": "string"
                }
            }
        },
        "main.GeneralAllStruct": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "file_key": {
                    "type": "string",
                    "default": "test.md"
                },
                "mod": {
                    "type": "string",
                    "enum": [
                        "append",
                        "cover"
                    ]
                }
            }
        },
        "main.SimpReadWebHookStruct": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "desc": {
                    "type": "string"
                },
                "note": {
                    "type": "string"
                },
                "tags": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "AuthorizationToken": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        },
        "Token": {
            "type": "apiKey",
            "name": "Token",
            "in": "header"
        }
    },
    "externalDocs": {
        "url": "https://kkbt.gitee.io/obcsapi-go/"
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "v4.2.1 版本",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Obcsapi",
	Description:      "基于 Obsidian S3 存储， CouchDb ，本地存储和 WebDAV 的后端 API ,可借助 Obsidian 插件 Remotely-Save 插件，或者 Self-hosted LiveSync (ex:Obsidian-livesync) 插件 CouchDb 方式同步，保存消息到 Obsidian 库。该调试页面大部分仅提供对 Headers-Token 验证方式的支持，其他如 Query-token，Headers-Authorization 除了特殊的几个其他并不支持。可以使用 https://hoppscotch.io/ 或者 Postman 之类的工具，或者使用 VsCode 插件 REST Client ，使用 REST Client 可以在 https://gitee.com/kkbt/obcsapi-go/tree/master/http ，找到测试文件",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
