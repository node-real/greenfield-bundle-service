// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "description": "API for handling file bundling and querying objects in the Bundle Service.",
    "title": "Bundle Service API",
    "version": "1.0.0"
  },
  "host": "bundle-service.nodereal.io",
  "basePath": "/v1",
  "paths": {
    "/createBundle": {
      "post": {
        "description": "Initiates a new bundle, requiring details like bucket name and bundle name.\n",
        "consumes": [
          "application/json"
        ],
        "tags": [
          "Bundle"
        ],
        "summary": "Start a New Bundle",
        "operationId": "createBundle",
        "parameters": [
          {
            "type": "string",
            "description": "User's digital signature for authorization",
            "name": "X-Signature",
            "in": "header",
            "required": true
          },
          {
            "description": "Parameters for managing a bundle",
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object",
              "required": [
                "bucketName",
                "bundleName",
                "timestamp"
              ],
              "properties": {
                "bucketName": {
                  "description": "The name of the bucket",
                  "type": "string"
                },
                "bundleName": {
                  "description": "The name of the bundle to be managed",
                  "type": "string"
                },
                "timestamp": {
                  "description": "Timestamp of the request",
                  "type": "integer",
                  "format": "int64"
                }
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Successfully managed bundle"
          },
          "400": {
            "description": "Invalid request or parameters",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          },
          "500": {
            "description": "Internal server error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/finalizeBundle": {
      "post": {
        "description": "Completes the lifecycle of an existing bundle, requiring the bundle name for authorization.\n",
        "consumes": [
          "application/json"
        ],
        "tags": [
          "Bundle"
        ],
        "summary": "Finalize an Existing Bundle",
        "operationId": "finalizeBundle",
        "parameters": [
          {
            "type": "string",
            "description": "User's digital signature for authorization",
            "name": "X-Signature",
            "in": "header",
            "required": true
          },
          {
            "type": "integer",
            "format": "int64",
            "description": "Timestamp of the finalizeBundle request",
            "name": "timestamp",
            "in": "query",
            "required": true
          },
          {
            "description": "Parameters for managing a bundle",
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object",
              "required": [
                "bucketName",
                "bundleName",
                "timestamp"
              ],
              "properties": {
                "bucketName": {
                  "description": "The name of the bucket",
                  "type": "string"
                },
                "bundleName": {
                  "description": "The name of the bundle to be finalized",
                  "type": "string"
                },
                "timestamp": {
                  "description": "Timestamp of the request",
                  "type": "integer",
                  "format": "int64"
                }
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Successfully managed bundle"
          },
          "400": {
            "description": "Invalid request or parameters",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          },
          "500": {
            "description": "Internal server error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/setBundleRule": {
      "post": {
        "description": "Set new rules or replace old rules for bundling, including constraints like maximum size and number of files.\n",
        "consumes": [
          "application/json"
        ],
        "tags": [
          "Rule"
        ],
        "summary": "Set New Bundling Rules",
        "operationId": "setBundleRule",
        "parameters": [
          {
            "type": "string",
            "description": "User's digital signature for authorization",
            "name": "X-Signature",
            "in": "header",
            "required": true
          },
          {
            "description": "Bundle rule creation parameters",
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object",
              "required": [
                "bucketName",
                "maxBundleSize",
                "maxBundleFiles",
                "maxFinalizeTime",
                "timestamp"
              ],
              "properties": {
                "bucketName": {
                  "description": "Name of the bucket for which the rule applies",
                  "type": "string"
                },
                "maxBundleFiles": {
                  "description": "Maximum number of files in a bundle",
                  "type": "integer",
                  "format": "int64"
                },
                "maxBundleSize": {
                  "description": "Maximum size of a bundle in bytes",
                  "type": "integer",
                  "format": "int64"
                },
                "maxFinalizeTime": {
                  "description": "Maximum time in seconds before a bundle must be finalized",
                  "type": "integer",
                  "format": "int64"
                },
                "timestamp": {
                  "description": "Timestamp of the request",
                  "type": "integer",
                  "format": "int64"
                }
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Successfully added bundle rule"
          },
          "400": {
            "description": "Invalid request or parameters",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          },
          "500": {
            "description": "Internal server error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/uploadObject": {
      "post": {
        "description": "Allows users to upload a single file along with a signature for validation, and a timestamp.\n",
        "consumes": [
          "multipart/form-data"
        ],
        "tags": [
          "Bundle"
        ],
        "summary": "Upload a single object to a bundle",
        "operationId": "uploadObject",
        "parameters": [
          {
            "type": "string",
            "description": "User's digital signature for authentication",
            "name": "X-Signature",
            "in": "header",
            "required": true
          },
          {
            "type": "file",
            "description": "The file to be uploaded",
            "name": "file",
            "in": "formData",
            "required": true
          },
          {
            "type": "string",
            "description": "The name of the file to be uploaded",
            "name": "fileName",
            "in": "formData",
            "required": true
          },
          {
            "type": "string",
            "description": "User's signature for the file",
            "name": "signature",
            "in": "formData",
            "required": true
          },
          {
            "type": "integer",
            "format": "int64",
            "description": "Timestamp of the upload",
            "name": "timestamp",
            "in": "formData",
            "required": true
          },
          {
            "type": "string",
            "description": "Content type of the file",
            "name": "contentType",
            "in": "formData",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Successfully uploaded file",
            "schema": {
              "type": "object",
              "properties": {
                "bundleName": {
                  "description": "The name of the bundle where the file has been uploaded",
                  "type": "string"
                }
              }
            }
          },
          "400": {
            "description": "Invalid request or file format",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          },
          "500": {
            "description": "Internal server error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/view/{bucketName}/{bundleName}/{objectName}": {
      "get": {
        "description": "Fetches a specific object from a given bundle and returns it as a file.\n",
        "produces": [
          "application/octet-stream"
        ],
        "tags": [
          "Bundle"
        ],
        "summary": "Retrieve an object as a file from a bundle",
        "operationId": "bundleObject",
        "parameters": [
          {
            "type": "string",
            "description": "User's digital signature for authentication",
            "name": "X-Signature",
            "in": "header",
            "required": true
          },
          {
            "type": "string",
            "description": "The bucketName of the bundle",
            "name": "bucketName",
            "in": "path",
            "required": true
          },
          {
            "type": "string",
            "description": "The name of the bundle",
            "name": "bundleName",
            "in": "path",
            "required": true
          },
          {
            "type": "string",
            "description": "The name of the object within the bundle",
            "name": "objectName",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Successfully retrieved file",
            "schema": {
              "type": "file"
            }
          },
          "404": {
            "description": "Bundle or object not found",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          },
          "500": {
            "description": "Internal server error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "Error": {
      "type": "object",
      "properties": {
        "code": {
          "description": "HTTP error code",
          "type": "integer",
          "format": "int64",
          "example": "400/500"
        },
        "message": {
          "description": "Error message",
          "type": "string",
          "example": "Bad request/Internal server error"
        }
      }
    }
  }
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "description": "API for handling file bundling and querying objects in the Bundle Service.",
    "title": "Bundle Service API",
    "version": "1.0.0"
  },
  "host": "bundle-service.nodereal.io",
  "basePath": "/v1",
  "paths": {
    "/createBundle": {
      "post": {
        "description": "Initiates a new bundle, requiring details like bucket name and bundle name.\n",
        "consumes": [
          "application/json"
        ],
        "tags": [
          "Bundle"
        ],
        "summary": "Start a New Bundle",
        "operationId": "createBundle",
        "parameters": [
          {
            "type": "string",
            "description": "User's digital signature for authorization",
            "name": "X-Signature",
            "in": "header",
            "required": true
          },
          {
            "description": "Parameters for managing a bundle",
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object",
              "required": [
                "bucketName",
                "bundleName",
                "timestamp"
              ],
              "properties": {
                "bucketName": {
                  "description": "The name of the bucket",
                  "type": "string"
                },
                "bundleName": {
                  "description": "The name of the bundle to be managed",
                  "type": "string"
                },
                "timestamp": {
                  "description": "Timestamp of the request",
                  "type": "integer",
                  "format": "int64"
                }
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Successfully managed bundle"
          },
          "400": {
            "description": "Invalid request or parameters",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          },
          "500": {
            "description": "Internal server error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/finalizeBundle": {
      "post": {
        "description": "Completes the lifecycle of an existing bundle, requiring the bundle name for authorization.\n",
        "consumes": [
          "application/json"
        ],
        "tags": [
          "Bundle"
        ],
        "summary": "Finalize an Existing Bundle",
        "operationId": "finalizeBundle",
        "parameters": [
          {
            "type": "string",
            "description": "User's digital signature for authorization",
            "name": "X-Signature",
            "in": "header",
            "required": true
          },
          {
            "type": "integer",
            "format": "int64",
            "description": "Timestamp of the finalizeBundle request",
            "name": "timestamp",
            "in": "query",
            "required": true
          },
          {
            "description": "Parameters for managing a bundle",
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object",
              "required": [
                "bucketName",
                "bundleName",
                "timestamp"
              ],
              "properties": {
                "bucketName": {
                  "description": "The name of the bucket",
                  "type": "string"
                },
                "bundleName": {
                  "description": "The name of the bundle to be finalized",
                  "type": "string"
                },
                "timestamp": {
                  "description": "Timestamp of the request",
                  "type": "integer",
                  "format": "int64"
                }
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Successfully managed bundle"
          },
          "400": {
            "description": "Invalid request or parameters",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          },
          "500": {
            "description": "Internal server error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/setBundleRule": {
      "post": {
        "description": "Set new rules or replace old rules for bundling, including constraints like maximum size and number of files.\n",
        "consumes": [
          "application/json"
        ],
        "tags": [
          "Rule"
        ],
        "summary": "Set New Bundling Rules",
        "operationId": "setBundleRule",
        "parameters": [
          {
            "type": "string",
            "description": "User's digital signature for authorization",
            "name": "X-Signature",
            "in": "header",
            "required": true
          },
          {
            "description": "Bundle rule creation parameters",
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object",
              "required": [
                "bucketName",
                "maxBundleSize",
                "maxBundleFiles",
                "maxFinalizeTime",
                "timestamp"
              ],
              "properties": {
                "bucketName": {
                  "description": "Name of the bucket for which the rule applies",
                  "type": "string"
                },
                "maxBundleFiles": {
                  "description": "Maximum number of files in a bundle",
                  "type": "integer",
                  "format": "int64"
                },
                "maxBundleSize": {
                  "description": "Maximum size of a bundle in bytes",
                  "type": "integer",
                  "format": "int64"
                },
                "maxFinalizeTime": {
                  "description": "Maximum time in seconds before a bundle must be finalized",
                  "type": "integer",
                  "format": "int64"
                },
                "timestamp": {
                  "description": "Timestamp of the request",
                  "type": "integer",
                  "format": "int64"
                }
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Successfully added bundle rule"
          },
          "400": {
            "description": "Invalid request or parameters",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          },
          "500": {
            "description": "Internal server error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/uploadObject": {
      "post": {
        "description": "Allows users to upload a single file along with a signature for validation, and a timestamp.\n",
        "consumes": [
          "multipart/form-data"
        ],
        "tags": [
          "Bundle"
        ],
        "summary": "Upload a single object to a bundle",
        "operationId": "uploadObject",
        "parameters": [
          {
            "type": "string",
            "description": "User's digital signature for authentication",
            "name": "X-Signature",
            "in": "header",
            "required": true
          },
          {
            "type": "file",
            "description": "The file to be uploaded",
            "name": "file",
            "in": "formData",
            "required": true
          },
          {
            "type": "string",
            "description": "The name of the file to be uploaded",
            "name": "fileName",
            "in": "formData",
            "required": true
          },
          {
            "type": "string",
            "description": "User's signature for the file",
            "name": "signature",
            "in": "formData",
            "required": true
          },
          {
            "type": "integer",
            "format": "int64",
            "description": "Timestamp of the upload",
            "name": "timestamp",
            "in": "formData",
            "required": true
          },
          {
            "type": "string",
            "description": "Content type of the file",
            "name": "contentType",
            "in": "formData",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Successfully uploaded file",
            "schema": {
              "type": "object",
              "properties": {
                "bundleName": {
                  "description": "The name of the bundle where the file has been uploaded",
                  "type": "string"
                }
              }
            }
          },
          "400": {
            "description": "Invalid request or file format",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          },
          "500": {
            "description": "Internal server error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/view/{bucketName}/{bundleName}/{objectName}": {
      "get": {
        "description": "Fetches a specific object from a given bundle and returns it as a file.\n",
        "produces": [
          "application/octet-stream"
        ],
        "tags": [
          "Bundle"
        ],
        "summary": "Retrieve an object as a file from a bundle",
        "operationId": "bundleObject",
        "parameters": [
          {
            "type": "string",
            "description": "User's digital signature for authentication",
            "name": "X-Signature",
            "in": "header",
            "required": true
          },
          {
            "type": "string",
            "description": "The bucketName of the bundle",
            "name": "bucketName",
            "in": "path",
            "required": true
          },
          {
            "type": "string",
            "description": "The name of the bundle",
            "name": "bundleName",
            "in": "path",
            "required": true
          },
          {
            "type": "string",
            "description": "The name of the object within the bundle",
            "name": "objectName",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Successfully retrieved file",
            "schema": {
              "type": "file"
            }
          },
          "404": {
            "description": "Bundle or object not found",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          },
          "500": {
            "description": "Internal server error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "Error": {
      "type": "object",
      "properties": {
        "code": {
          "description": "HTTP error code",
          "type": "integer",
          "format": "int64",
          "example": "400/500"
        },
        "message": {
          "description": "Error message",
          "type": "string",
          "example": "Bad request/Internal server error"
        }
      }
    }
  }
}`))
}
