// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
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
        "/open/v1/oauth/login/{provider}": {
            "get": {
                "description": "前往第三方登录",
                "tags": [
                    "oauth"
                ],
                "summary": "前往第三方登录",
                "operationId": "oauth-login",
                "parameters": [
                    {
                        "type": "string",
                        "description": "登录成功后跳转网址",
                        "name": "nextURL",
                        "in": "query",
                        "required": true
                    },
                    {
                        "enum": [
                            "api",
                            "js"
                        ],
                        "type": "string",
                        "description": "同步登录方式(支持的值有: api, js)",
                        "name": "syncType",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "同步登录网址(syncType == js 时采用js方式应用此网址; syncType == api 时采用带加密参数后重定向跳转到此网址)",
                        "name": "loginURL",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "第三方网站标识(例如: qq, wechat 等)",
                        "name": "provider",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        },
        "/open/v1/oauth/providers": {
            "get": {
                "description": "列出支持的所有第三方登录提供商",
                "tags": [
                    "oauth"
                ],
                "summary": "第三方登录提供商列表",
                "operationId": "oauth-proviers",
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        },
        "/open/v1/payment/providers": {
            "get": {
                "description": "列出支持的所有支付网关提供商",
                "tags": [
                    "payment"
                ],
                "summary": "支付网关提供商列表",
                "operationId": "payment-proviers",
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        },
        "/open/v1/payment/{provider}": {
            "get": {
                "description": "跳转到付款界面",
                "tags": [
                    "payment"
                ],
                "summary": "前往付款",
                "operationId": "payment",
                "parameters": [
                    {
                        "type": "string",
                        "description": "appID",
                        "name": "appID",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "你的用户ID",
                        "name": "appUID",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "支付放弃后返回地址",
                        "name": "cancelURL",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "币种 支持的值：https://github.com/webx-top/payment/blob/master/config/constants.go",
                        "name": "currency",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "客户IP",
                        "name": "customerIP",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "设备 支持的值：https://github.com/webx-top/payment/blob/master/config/constants.go",
                        "name": "device",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "过期时间戳",
                        "name": "expiresTs",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "扩展信息",
                        "name": "extend",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "是否是虚拟商品",
                        "name": "isVirtualProduct",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "支付回调",
                        "name": "notifyURL",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "你的订单号",
                        "name": "outOrderNo",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "number",
                        "description": "价格",
                        "name": "price",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "商品ID",
                        "name": "productID",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "商品类型(自己定义)",
                        "name": "productType",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "支付成功后返回地址",
                        "name": "returnURL",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "订单主题(一般为商品名)",
                        "name": "subject",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "用于第四方支付时选择支付方式",
                        "name": "subtype",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "付款方式(alipay,wechat)",
                        "name": "type",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "支付方式标识(例如: alipay, wechat 等)",
                        "name": "provider",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        },
        "/open/v1/query/payment": {
            "get": {
                "description": "查询付款状态 (“你的订单号”和“平台订单号”必须至少提供一个)",
                "tags": [
                    "payment"
                ],
                "summary": "付款结果查询",
                "operationId": "payment-query",
                "parameters": [
                    {
                        "type": "string",
                        "description": "你的订单号",
                        "name": "outOrderNo",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "平台订单号",
                        "name": "orderNo",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        },
        "/open/v1/query/refund": {
            "get": {
                "description": "查询退款结果 (“你的退款单号”和“平台退款单号”必须至少提供一个)",
                "tags": [
                    "payment"
                ],
                "summary": "查询退款结果",
                "operationId": "payment-refund-query",
                "parameters": [
                    {
                        "type": "string",
                        "description": "你的退款单号",
                        "name": "outRefundNo",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "平台退款单号",
                        "name": "refundNo",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        },
        "/open/v1/refund": {
            "post": {
                "description": "发起退款",
                "tags": [
                    "payment"
                ],
                "summary": "发起退款",
                "operationId": "payment-refund",
                "parameters": [
                    {
                        "type": "string",
                        "description": "appID",
                        "name": "appID",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "通知接口网址",
                        "name": "notifyURL",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "平台订单号",
                        "name": "orderNo",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "你的订单号",
                        "name": "outOrderNo",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "number",
                        "description": "退款金额",
                        "name": "refundAmount",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "退款单号",
                        "name": "refundNo",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "退款原因",
                        "name": "refundReason",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        },
        "/open/v1/short_url/create": {
            "get": {
                "description": "创建短网址",
                "tags": [
                    "short_url"
                ],
                "summary": "创建短网址",
                "operationId": "short_url-create",
                "parameters": [
                    {
                        "type": "string",
                        "description": "长网址",
                        "name": "url",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/top.ResponseData"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "Data": {
                                            "$ref": "#/definitions/shorturl.ResponseAPICreate"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "shorturl.ResponseAPICreate": {
            "type": "object",
            "properties": {
                "qrcodeURL": {
                    "type": "string"
                },
                "shortURL": {
                    "type": "string"
                }
            }
        },
        "top.ResponseData": {
            "type": "object",
            "properties": {
                "Code": {
                    "type": "integer"
                },
                "Data": {
                    "type": "object"
                },
                "Info": {
                    "type": "string"
                },
                "State": {
                    "type": "string"
                },
                "URL": {
                    "type": "string"
                },
                "Zone": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        },
        "BasicAuth": {
            "type": "basic"
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
	Host:        "127.0.0.1:8080",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "开放平台接口",
	Description: "",
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
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
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
	swag.Register("swagger", &s{})
}
